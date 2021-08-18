// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tsdb

import (
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/cumulativecache"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

// for testing
var (
	newReplicaSequenceFunc = newReplicaSequence
	newIntervalSegmentFunc = newIntervalSegment
	newKVStoreFunc         = kv.NewStore
	newIndexDBFunc         = indexdb.NewIndexDatabase
	newMemoryDBFunc        = memdb.NewMemoryDatabase
)

var (
	shardScope                 = linmetric.NewScope("lindb.tsdb.shard")
	badMetricsVec              = shardScope.NewDeltaCounterVec("bad_metrics", "db", "shard")
	outOfRangeMetricsVec       = shardScope.NewDeltaCounterVec("metrics_out_of_range", "db", "shard")
	writeMetricsVec            = shardScope.NewDeltaCounterVec("write_metrics", "db", "shard")
	writeMetricFailuresVec     = shardScope.NewDeltaCounterVec("write_metric_failures", "db", "shard")
	writeFieldsVec             = shardScope.NewDeltaCounterVec("write_fields", "db", "shard")
	cumulativeTransformedVec   = shardScope.NewDeltaCounterVec("cumulative_transformed", "db", "shard")
	cumulativeUnTransformedVec = shardScope.NewDeltaCounterVec("cumulative_untransformed", "db", "shard")
	escapedFieldNameVec        = shardScope.NewDeltaCounterVec("escaped_fields", "db", "shard")
	memdbTotalSizeVec          = shardScope.NewGaugeVec("memdb_total_size", "db", "shard")
	memdbNumberVec             = shardScope.NewGaugeVec("memdb_number", "db", "shard")
	memFlushTimerVec           = shardScope.Scope("memdb_flush_duration").NewDeltaHistogramVec("db", "shard")
	indexFlushTimerVec         = shardScope.Scope("indexdb_flush_duration").NewDeltaHistogramVec("db", "shard")
)

const (
	replicaDir       = "replica"
	segmentDir       = "segment"
	indexParentDir   = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
	metaDir          = "meta"
	tempDir          = "temp"
)

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// DatabaseName returns the database name
	DatabaseName() string
	// ShardID returns the shard id
	ShardID() int32
	// CurrentInterval returns current interval for metric write.
	CurrentInterval() timeutil.Interval
	// ShardInfo returns the unique shard info
	ShardInfo() string
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// GetOrCreateMemoryDatabase makes sure that a memory database will always be returned by given family time.
	GetOrCreateMemoryDatabase(familyTime int64) (memdb.MemoryDatabase, error)
	// IndexDatabase returns the index-database
	IndexDatabase() indexdb.IndexDatabase
	// Write writes the metric-point into memory-database.
	Write(metric *protoMetricsV1.Metric) error
	// GetOrCreateSequence gets the replica sequence by given remote peer if exist, else creates a new sequence
	GetOrCreateSequence(replicaPeer string) (replication.Sequence, error)
	// MemDBTotalSize returns the total size of mutable and immutable memdb
	MemDBTotalSize() int64
	// Flush flushes index and memory data to disk
	Flush() error
	// NeedFlush checks if shard need to flush memory data
	NeedFlush() bool
	// IsFlushing checks if this shard is in flushing
	IsFlushing() bool
	// initIndexDatabase initializes index database
	initIndexDatabase() error
	// Closer releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer
	// DataFilter filters the data based on condition
	flow.DataFilter
}

type shardMetrics struct {
	badMetrics              *linmetric.BoundDeltaCounter
	outOfRangeMetrics       *linmetric.BoundDeltaCounter
	writeMetrics            *linmetric.BoundDeltaCounter
	writeMetricFailures     *linmetric.BoundDeltaCounter
	writeFields             *linmetric.BoundDeltaCounter
	cumulativeTransformed   *linmetric.BoundDeltaCounter
	cumulativeUnTransformed *linmetric.BoundDeltaCounter
	escapedFields           *linmetric.BoundDeltaCounter
	memdbTotalSize          *linmetric.BoundGauge
	memdbNumber             *linmetric.BoundGauge
	memFlushTimer           *linmetric.BoundDeltaHistogram
	indexFlushTimer         *linmetric.BoundDeltaHistogram
}

func newShardMetrics(dbName string, shardID int32) *shardMetrics {
	shardIDStr := strconv.Itoa(int(shardID))
	return &shardMetrics{
		badMetrics:              badMetricsVec.WithTagValues(dbName, shardIDStr),
		outOfRangeMetrics:       outOfRangeMetricsVec.WithTagValues(dbName, shardIDStr),
		writeMetrics:            writeMetricsVec.WithTagValues(dbName, shardIDStr),
		writeMetricFailures:     writeMetricFailuresVec.WithTagValues(dbName, shardIDStr),
		writeFields:             writeFieldsVec.WithTagValues(dbName, shardIDStr),
		cumulativeTransformed:   cumulativeTransformedVec.WithTagValues(dbName, shardIDStr),
		cumulativeUnTransformed: cumulativeUnTransformedVec.WithTagValues(dbName, shardIDStr),
		escapedFields:           escapedFieldNameVec.WithTagValues(dbName, shardIDStr),
		memdbTotalSize:          memdbTotalSizeVec.WithTagValues(dbName, shardIDStr),
		memdbNumber:             memdbNumberVec.WithTagValues(dbName, shardIDStr),
		memFlushTimer:           memFlushTimerVec.WithTagValues(dbName, shardIDStr),
		indexFlushTimer:         indexFlushTimerVec.WithTagValues(dbName, shardIDStr),
	}
}

// shard implements Shard interface
// directory tree:
//    xx/shard/1/ (path)
//    xx/shard/1/replica
//    xx/shard/1/temp/123213123131 // time of ns
//    xx/shard/1/meta/
//    xx/shard/1/index/inverted/
//    xx/shard/1/data/20191012/
//    xx/shard/1/data/20191013/
type shard struct {
	databaseName string
	id           int32
	path         string
	option       option.DatabaseOption
	sequence     ReplicaSequence

	mutex    sync.Mutex     // mutex for update families
	families familyMemDBSet // memory database for each family time

	indexDB  indexdb.IndexDatabase
	metadata metadb.Metadata
	// write accept time range
	interval timeutil.Interval
	ahead    timeutil.Interval
	behind   timeutil.Interval
	// segments keeps all interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	segments       map[timeutil.IntervalType]IntervalSegment
	segment        IntervalSegment // smallest interval for writing data
	isFlushing     atomic.Bool     // restrict flusher concurrency
	flushCondition sync.WaitGroup  // flush condition

	indexStore     kv.Store  // kv stores
	forwardFamily  kv.Family // forward store
	invertedFamily kv.Family // inverted store
	logger         *logger.Logger
	metrics        shardMetrics

	// cumulative field value-> delta cache
	once4Cache      sync.Once
	cumulativeCache *cumulativecache.Cache
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(
	db Database,
	shardID int32,
	shardPath string,
	option option.DatabaseOption,
) (Shard, error) {
	var err error
	if err = option.Validate(); err != nil {
		return nil, fmt.Errorf("engine option is invalid, err: %s", err)
	}
	var interval timeutil.Interval
	_ = interval.ValueOf(option.Interval)

	if err := mkDirIfNotExist(shardPath); err != nil {
		return nil, err
	}
	replicaSequence, err := newReplicaSequenceFunc(filepath.Join(shardPath, replicaDir))
	if err != nil {
		return nil, err
	}
	createdShard := &shard{
		databaseName: db.Name(),
		id:           shardID,
		path:         shardPath,
		option:       option,
		sequence:     replicaSequence,
		families:     *newFamilyMemDBSet(),
		metadata:     db.Metadata(),
		interval:     interval,
		segments:     make(map[timeutil.IntervalType]IntervalSegment),
		isFlushing:   *atomic.NewBool(false),
		metrics:      *newShardMetrics(db.Name(), shardID),
		logger:       logger.GetLogger("tsdb", "Shard"),
	}
	// new segment for writing
	createdShard.segment, err = newIntervalSegmentFunc(
		interval,
		filepath.Join(shardPath, segmentDir, interval.Type().String()))

	if err != nil {
		return nil, err
	}
	_ = createdShard.ahead.ValueOf(option.Ahead)
	_ = createdShard.behind.ValueOf(option.Behind)
	// add writing segment into segment list
	createdShard.segments[interval.Type()] = createdShard.segment

	defer func() {
		if err == nil {
			return
		}
		if err = createdShard.Close(); err != nil {
			engineLogger.Error("close shard error when create shard fail",
				logger.Int32("shardID", createdShard.id),
				logger.String("database", createdShard.databaseName),
				logger.String("shard", createdShard.path), logger.Error(err))
		}
	}()
	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	// add shard into global shard manager
	GetShardManager().AddShard(createdShard)
	return createdShard, nil
}

func (s *shard) getCache() *cumulativecache.Cache {
	s.once4Cache.Do(func() {
		s.cumulativeCache = cumulativecache.NewCache(
			32,
			time.Minute, // max pull duration is 1 minute
			time.Second*30,
			shardScope.Scope("cumulative_caches"),
		)
	})
	return s.cumulativeCache
}

// DatabaseName returns the database name
func (s *shard) DatabaseName() string {
	return s.databaseName
}

// ShardID returns the shard id
func (s *shard) ShardID() int32 {
	return s.id
}

// ShardInfo returns the unique shard info
func (s *shard) ShardInfo() string {
	return s.path
}

// CurrentInterval returns current interval for metric  write.
func (s *shard) CurrentInterval() timeutil.Interval {
	return s.interval
}

func (s *shard) GetOrCreateSequence(replicaPeer string) (replication.Sequence, error) {
	return s.sequence.getOrCreateSequence(replicaPeer)
}

func (s *shard) IndexDatabase() indexdb.IndexDatabase {
	return s.indexDB
}

func (s *shard) GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily {
	segment, ok := s.segments[intervalType]
	if ok {
		return segment.getDataFamilies(timeRange)
	}
	return nil
}

// GetOrCreateMemoryDatabase returns memory database by given family time.
func (s *shard) GetOrCreateMemoryDatabase(familyTime int64) (memdb.MemoryDatabase, error) {
	db, exist := s.families.GetMutableFamily(familyTime)
	if exist {
		return db, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	// double check
	db, exist = s.families.GetMutableFamily(familyTime)
	if exist {
		return db, nil
	}
	newDB, err := s.createMemoryDatabase(familyTime)
	if err != nil {
		return nil, err
	}
	s.families.InsertFamily(familyTime, newDB)
	return newDB, nil
}

// Filter filters the data based on metric/time range/seriesIDs,
// if finds data then returns the flow.FilterResultSet, else returns nil
func (s *shard) Filter(
	metricID uint32,
	seriesIDs *roaring.Bitmap,
	timeRange timeutil.TimeRange,
	fields field.Metas,
) (rs []flow.FilterResultSet, err error) {
	entries := s.families.Entries()
	for idx := range entries {
		// check family time if in query time range
		if timeRange.Contains(entries[idx].familyTime) {
			resultSet, err := entries[idx].memDB.Filter(metricID, seriesIDs, timeRange, fields)
			if err != nil {
				return nil, err
			}
			rs = append(rs, resultSet...)
		}
	}
	return
}

func (s *shard) FindMemoryDatabase() (rs []memdb.MemoryDatabase) {
	entries := s.families.Entries()
	for idx := range entries {
		rs = append(rs, entries[idx].memDB)
	}
	return rs
}

func (s *shard) validateMetric(metric *protoMetricsV1.Metric) (isCumulative bool, err error) {
	if metric == nil {
		return isCumulative, constants.ErrMetricPBNilMetric
	}
	if len(metric.Name) == 0 {
		return isCumulative, constants.ErrMetricPBEmptyMetricName
	}
	// empty field
	if len(metric.SimpleFields) == 0 && metric.CompoundField == nil {
		return isCumulative, constants.ErrMetricPBEmptyField
	}
	timestamp := metric.Timestamp
	now := fasttime.UnixMilliseconds()
	// check metric timestamp if in acceptable time range
	if (s.behind.Int64() > 0 && timestamp < now-s.behind.Int64()) ||
		(s.ahead.Int64() > 0 && timestamp > now+s.ahead.Int64()) {
		s.metrics.outOfRangeMetrics.Incr()
		return isCumulative, constants.ErrMetricOutOfTimeRange
	}
	// validate empty tags
	if len(metric.Tags) > 0 {
		for idx := range metric.Tags {
			// nil tag
			if metric.Tags[idx] == nil {
				return isCumulative, constants.ErrMetricEmptyTagKeyValue
			}
			// empty key value
			if metric.Tags[idx].Key == "" || metric.Tags[idx].Value == "" {
				return isCumulative, constants.ErrMetricEmptyTagKeyValue
			}
		}
	}

	// check simple fields
	for idx := range metric.SimpleFields {
		// nil value
		if metric.SimpleFields[idx] == nil {
			return isCumulative, constants.ErrBadMetricPBFormat
		}
		// field-name empty
		if metric.SimpleFields[idx].Name == "" {
			return isCumulative, constants.ErrMetricEmptyFieldName
		}
		// check sanitize
		if field.HistogramConverter.NeedToSanitize(metric.SimpleFields[idx].Name) {
			s.metrics.escapedFields.Incr()
			metric.SimpleFields[idx].Name = field.HistogramConverter.Sanitize(metric.SimpleFields[idx].Name)
		}
		// field type unspecified
		switch metric.SimpleFields[idx].Type {
		case protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED:
			return false, constants.ErrBadMetricPBFormat
		case protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM:
			isCumulative = true
		}
		v := metric.SimpleFields[idx].Value
		if math.IsNaN(v) {
			return isCumulative, constants.ErrMetricNanField
		}
		if math.IsInf(v, 0) {
			return isCumulative, constants.ErrMetricInfField
		}
	}
	// no more compound field
	if metric.CompoundField == nil {
		return isCumulative, nil
	}
	// compound field-type unspecified
	switch metric.CompoundField.Type {
	case protoMetricsV1.CompoundFieldType_COMPOUND_UNSPECIFIED:
		return isCumulative, constants.ErrBadMetricPBFormat
	case protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM:
		isCumulative = true
	}
	// value length zero or length not match
	if len(metric.CompoundField.Values) != len(metric.CompoundField.ExplicitBounds) ||
		len(metric.CompoundField.Values) <= 2 {
		return isCumulative, constants.ErrBadMetricPBFormat
	}
	// ensure compound field value > 0
	if (metric.CompoundField.Max < 0) ||
		metric.CompoundField.Min < 0 ||
		metric.CompoundField.Sum < 0 ||
		metric.CompoundField.Count < 0 {
		return isCumulative, constants.ErrBadMetricPBFormat
	}

	for idx := 0; idx < len(metric.CompoundField.Values); idx++ {
		// ensure value > 0
		if metric.CompoundField.Values[idx] < 0 || metric.CompoundField.ExplicitBounds[idx] < 0 {
			return isCumulative, constants.ErrBadMetricPBFormat
		}
		// ensure explicate bounds increase progressively
		if idx >= 1 && metric.CompoundField.ExplicitBounds[idx] < metric.CompoundField.ExplicitBounds[idx-1] {
			return isCumulative, constants.ErrBadMetricPBFormat
		}
		// ensure last bound is +Inf
		if idx == len(metric.CompoundField.ExplicitBounds)-1 && !math.IsInf(metric.CompoundField.ExplicitBounds[idx], 1) {
			return isCumulative, constants.ErrBadMetricPBFormat
		}
	}
	return isCumulative, nil
}

func (s *shard) howManyFieldsWillWrite(metric *protoMetricsV1.Metric) int {
	var count = len(metric.SimpleFields)
	if metric.CompoundField == nil {
		return count
	}
	// min, max is a feature in lindb's field
	count += len(metric.CompoundField.Values)
	if metric.CompoundField.Min > 0 {
		count++
	}
	if metric.CompoundField.Max > 0 {
		count++
	}
	// assume that all compound field will contains sum/count field
	count += 2
	return count
}

func (s *shard) lookupMetricMeta(metric *protoMetricsV1.Metric) (*memdb.MetricPoint, error) {
	ns := metric.Namespace
	if len(ns) == 0 {
		ns = constants.DefaultNamespace
	}

	metricID, err := s.metadata.MetadataDatabase().GenMetricID(ns, metric.Name)
	if err != nil {
		s.metrics.writeMetricFailures.Incr()
		return nil, err
	}
	var seriesID uint32
	isCreated := false
	if len(metric.Tags) == 0 {
		// if metric without tags, uses default series id(0)
		seriesID = constants.SeriesIDWithoutTags
	} else {
		seriesID, isCreated, err = s.indexDB.GetOrCreateSeriesID(metricID, metric.TagsHash)
		if err != nil {
			s.metrics.writeMetricFailures.Incr()
			return nil, err
		}
	}
	if isCreated {
		// if series id is new, need build inverted index
		s.indexDB.BuildInvertIndex(ns, metric.Name, metric.Tags, seriesID)
	}

	fieldsCount := s.howManyFieldsWillWrite(metric)
	var mm = memdb.MetricPoint{
		MetricID: metricID,
		SeriesID: seriesID,
		FieldIDs: make([]field.ID, fieldsCount)[:0],
		Proto:    metric,
	}

	for idx := range metric.SimpleFields {
		var fieldType field.Type
		switch metric.SimpleFields[idx].Type {
		case protoMetricsV1.SimpleFieldType_DELTA_SUM, protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM:
			fieldType = field.SumField
		case protoMetricsV1.SimpleFieldType_GAUGE:
			fieldType = field.GaugeField
		}
		fieldID, err := s.metadata.MetadataDatabase().GenFieldID(
			ns, metric.Name, field.Name(metric.SimpleFields[idx].Name), fieldType)
		if err != nil {
			return nil, err
		}
		mm.FieldIDs = append(mm.FieldIDs, fieldID)
	}
	if metric.CompoundField == nil {
		return &mm, nil
	}
	// min
	if metric.CompoundField.Min > 0 {
		minFieldID, err := s.metadata.MetadataDatabase().GenFieldID(
			ns, metric.Name, field.Name(field.HistogramConverter.MinFieldName), field.MinField)
		if err != nil {
			return nil, err
		}
		mm.FieldIDs = append(mm.FieldIDs, minFieldID)
	}
	// max
	if metric.CompoundField.Max > 0 {
		maxFieldID, err := s.metadata.MetadataDatabase().GenFieldID(
			ns, metric.Name, field.Name(field.HistogramConverter.MaxFieldName), field.MaxField)
		if err != nil {
			return nil, err
		}
		mm.FieldIDs = append(mm.FieldIDs, maxFieldID)
	}
	// sum
	sumFieldID, err := s.metadata.MetadataDatabase().GenFieldID(
		ns, metric.Name, field.Name(field.HistogramConverter.SumFieldName), field.SumField)
	if err != nil {
		return nil, err
	}
	mm.FieldIDs = append(mm.FieldIDs, sumFieldID)
	// count
	countFieldID, err := s.metadata.MetadataDatabase().GenFieldID(
		ns, metric.Name, field.Name(field.HistogramConverter.CountFieldName), field.SumField)
	if err != nil {
		return nil, err
	}
	mm.FieldIDs = append(mm.FieldIDs, countFieldID)
	// values
	for idx := range metric.CompoundField.ExplicitBounds {
		histogramFieldID, err := s.metadata.MetadataDatabase().GenFieldID(
			ns, metric.Name,
			field.Name(field.HistogramConverter.BucketName(metric.CompoundField.ExplicitBounds[idx])),
			field.HistogramField)
		if err != nil {
			return nil, err
		}
		mm.FieldIDs = append(mm.FieldIDs, histogramFieldID)
	}
	return &mm, nil
}

// Write writes the metric-point into memory-database.
func (s *shard) Write(metric *protoMetricsV1.Metric) (err error) {
	isCumulative, err := s.validateMetric(metric)
	if err != nil {
		s.metrics.badMetrics.Incr()
		return err
	}
	timestamp := metric.Timestamp
	point, err := s.lookupMetricMeta(metric)
	if err != nil {
		return err
	}

	// calculate family start time and slot index
	intervalCalc := s.interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)              // day
	family := intervalCalc.CalcFamily(timestamp, segmentTime)           // hours
	familyTime := intervalCalc.CalcFamilyStartTime(segmentTime, family) // family timestamp
	db, err := s.GetOrCreateMemoryDatabase(familyTime)
	if err != nil {
		s.metrics.writeMetricFailures.Incr()
		return err
	}

	point.SlotIndex = uint16(intervalCalc.CalcSlot(timestamp, familyTime, s.interval.Int64())) // slot offset of family
	if isCumulative {
		if updated := s.getCache().CumulativePointToDelta(point); updated {
			s.metrics.cumulativeTransformed.Incr()
		} else {
			s.metrics.cumulativeUnTransformed.Incr()
		}
	}

	db.AcquireWrite()
	// write metric point into memory db
	err = db.Write(point)
	db.CompleteWrite()

	if err == nil {
		s.metrics.writeMetrics.Incr()
		s.metrics.writeFields.Add(float64(len(point.FieldIDs)))
	} else {
		s.metrics.writeMetricFailures.Incr()
	}
	// if memdb size is above threshold, it will be put into immutable list
	s.validateMemDBSize(familyTime, db)
	return err
}

func (s *shard) Close() error {
	// wait previous flush job completed
	s.flushCondition.Wait()

	GetShardManager().RemoveShard(s)
	if s.indexDB != nil {
		if err := s.indexDB.Close(); err != nil {
			return err
		}
	}
	if s.indexStore != nil {
		if err := s.indexStore.Close(); err != nil {
			return err
		}
	}
	for _, entry := range s.families.Entries() {
		if err := s.flushMemoryDatabase(entry.memDB); err != nil {
			return err
		}
	}
	s.ackReplicaSeq()
	return s.sequence.Close()
}

// IsFlushing checks if this shard is in flushing
func (s *shard) IsFlushing() bool { return s.isFlushing.Load() }

func (s *shard) validateMemDBSize(familyTime int64, m memdb.MemoryDatabase) {
	// memory usage lower than threshold
	maxMemDBSize := int64(config.GlobalStorageConfig().TSDB.MaxMemDBSize)
	if m.MemSize() < maxMemDBSize {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.families.SetFamilyImmutable(familyTime)

	s.logger.Info("memdb is above memory threshold, switch to immutable",
		logger.Int32("shardID", s.id),
		logger.String("database", s.databaseName),
		logger.Int64("familyTime", familyTime),
		logger.String("uptime", m.Uptime().String()),
		logger.String("memdb-size", ltoml.Size(m.MemSize()).String()),
		logger.Int64("max-memdb-size", maxMemDBSize),
	)
}

func (s *shard) tryEvictMutable() {
	// fast path, there is no expired mutable memdb
	ttl := config.GlobalStorageConfig().TSDB.MutableMemDBTTL.Duration()
	mutable := s.families.MutableEntries()
	for _, entry := range mutable {
		if entry.memDB.Uptime() > ttl {
			goto MoveMutable
		}
	}
	return

MoveMutable:
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, entry := range s.families.MutableEntries() {
		if entry.memDB.Uptime() > ttl {
			s.families.SetFamilyImmutable(entry.familyTime)
			s.logger.Info("switch a expired mutable memdb to immutable",
				logger.Int32("shardID", s.id),
				logger.String("database", s.databaseName),
				logger.String("uptime", entry.memDB.Uptime().String()),
				logger.String("mutable-memdb-ttl", ttl.String()),
			)
		}
	}
}

func (s *shard) MemDBTotalSize() int64 {
	return s.families.TotalSize()
}

// NeedFlush checks if shard need to flush memory data
func (s *shard) NeedFlush() bool {
	if s.IsFlushing() {
		return false
	}
	s.metrics.memdbNumber.Update(float64(len(s.families.Entries())))
	s.metrics.memdbTotalSize.Update(float64(s.MemDBTotalSize()))

	s.tryEvictMutable()

	cfg := config.GlobalStorageConfig()
	// too many memdbs
	number := len(s.families.Entries())
	if number > cfg.TSDB.MaxMemDBNumber {
		s.logger.Info("number of memdb is above threshold, waiting for flush",
			logger.Int32("shardID", s.id),
			logger.String("database", s.databaseName),
			logger.Int32("memdb-number", int32(number)),
			logger.Int32("max-memdb-number", int32(cfg.TSDB.MaxMemDBNumber)),
		)
		return true
	}
	// total size too much
	totalSize := s.families.TotalSize()
	if totalSize > int64(cfg.TSDB.MaxMemDBTotalSize) {
		s.logger.Info("total size of memdb is above threshold, waiting for flush",
			logger.Int32("shardID", s.id),
			logger.String("database", s.databaseName),
			logger.Int64("memdb-total-size", totalSize),
			logger.Int64("max-memdb-total-size", int64(cfg.TSDB.MaxMemDBTotalSize)),
		)
		return true
	}
	return false
}

// Flush flushes index and memory data to disk
func (s *shard) Flush() (err error) {
	// another flush process is running
	if !s.isFlushing.CAS(false, true) {
		return nil
	}
	// 1. mark flush job doing
	s.flushCondition.Add(1)

	defer func() {
		//TODO add commit kv meta after ack successfully
		// mark flush job complete, notify
		s.flushCondition.Done()
		s.isFlushing.Store(false)
	}()

	startTime := time.Now()
	//FIXME stone1100
	// index flush
	if s.indexDB != nil {
		if err = s.indexDB.Flush(); err != nil {
			s.logger.Error("failed to flush indexDB ",
				logger.Int32("shardID", s.id),
				logger.String("database", s.databaseName),
				logger.Error(err))
			return err
		}
		s.logger.Info("flush indexDB successfully",
			logger.Int32("shardID", s.id),
			logger.String("database", s.databaseName),
		)
		s.metrics.indexFlushTimer.UpdateSince(startTime)
	}

	var waitingFlushMemDB memdb.MemoryDatabase
	immutable := s.families.ImmutableEntries()
	// flush first immutable memdb
	if len(immutable) > 0 {
		waitingFlushMemDB = immutable[0].memDB
	} else {
		s.mutex.Lock()
		// force picks a mutable memdb from memory
		if evictedMutable := s.families.SetLargestMutableMemDBImmutable(); evictedMutable {
			waitingFlushMemDB = s.families.ImmutableEntries()[0].memDB
			s.logger.Info("forcefully switch a memdb to immutable for flushing",
				logger.Int32("shardID", s.id),
				logger.String("database", s.databaseName),
				logger.Int64("familyTime", waitingFlushMemDB.FamilyTime()),
				logger.Int64("memDBSize", waitingFlushMemDB.MemSize()),
			)
		}
		s.mutex.Unlock()
	}
	if waitingFlushMemDB == nil {
		s.logger.Warn("there is no memdb to flush", logger.Int32("shardID", s.id))
		return nil
	}

	startTime = time.Now()
	if err := s.flushMemoryDatabase(waitingFlushMemDB); err != nil {
		s.logger.Error("failed to flush memdb",
			logger.Int32("shardID", s.id),
			logger.String("database", s.databaseName),
			logger.Int64("familyTime", waitingFlushMemDB.FamilyTime()),
			logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
		return err
	}
	// flush success, remove it from the immutable list
	s.mutex.Lock()
	s.families.RemoveHeadImmutable()
	s.mutex.Unlock()

	endTime := time.Now()
	s.logger.Info("flush memdb successfully",
		logger.Int32("shardID", s.id),
		logger.String("database", s.databaseName),
		logger.String("flush-duration", endTime.Sub(startTime).String()),
		logger.Int64("familyTime", waitingFlushMemDB.FamilyTime()),
		logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
	s.metrics.memFlushTimer.UpdateDuration(endTime.Sub(startTime))

	//FIXME(stone1100) commit replica sequence
	s.ackReplicaSeq()
	return nil
}

// initIndexDatabase initializes the index database
func (s *shard) initIndexDatabase() error {
	var err error
	storeOption := kv.DefaultStoreOption(filepath.Join(s.path, indexParentDir))
	s.indexStore, err = newKVStoreFunc(storeOption.Path, storeOption)
	if err != nil {
		return err
	}
	s.forwardFamily, err = s.indexStore.CreateFamily(
		forwardIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(invertedindex.SeriesForwardMerger)})
	if err != nil {
		return err
	}
	s.invertedFamily, err = s.indexStore.CreateFamily(
		invertedIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(invertedindex.SeriesInvertedMerger)})
	if err != nil {
		return err
	}
	s.indexDB, err = newIndexDBFunc(
		context.TODO(),
		filepath.Join(s.path, metaDir),
		s.metadata, s.forwardFamily,
		s.invertedFamily)
	if err != nil {
		return err
	}
	return nil
}

// createMemoryDatabase creates a new memory database for writing data points
func (s *shard) createMemoryDatabase(familyTime int64) (memdb.MemoryDatabase, error) {
	return newMemoryDBFunc(memdb.MemoryDatabaseCfg{
		FamilyTime: familyTime,
		Name:       s.databaseName,
		TempPath:   filepath.Join(s.path, filepath.Join(tempDir, fmt.Sprintf("%d", timeutil.Now()))),
	})
}

// flushMemoryDatabase flushes memory database to disk kv store
func (s *shard) flushMemoryDatabase(memDB memdb.MemoryDatabase) error {
	startTime := time.Now()
	defer s.metrics.memFlushTimer.UpdateSince(startTime)

	segmentName := s.interval.Calculator().GetSegment(memDB.FamilyTime())
	segment, err := s.segment.GetOrCreateSegment(segmentName)
	if err != nil {
		return err
	}
	thisDataFamily, err := segment.GetDataFamily(memDB.FamilyTime())
	if err != nil {
		return err
	}
	// flush family data
	if err := memDB.FlushFamilyTo(
		metricsdata.NewFlusher(thisDataFamily.Family().NewFlusher())); err != nil {
		return err
	}
	if err := memDB.Close(); err != nil {
		return err
	}
	return nil
}

// ackReplicaSeq commits the replica sequence
// NOTICE: if fail, maybe data will write duplicate if system restart
func (s *shard) ackReplicaSeq() {
	allHeads := s.sequence.getAllHeads()
	if err := s.sequence.ack(allHeads); err != nil {
		engineLogger.Error("ack replica sequence error", logger.String("shard", s.path), logger.Error(err))
	}
}
