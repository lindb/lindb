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
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/tagindex"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

// for testing
var (
	newIntervalSegmentFunc = newIntervalSegment
	newKVStoreFunc         = kv.NewStore
	newIndexDBFunc         = indexdb.NewIndexDatabase
	newMemoryDBFunc        = memdb.NewMemoryDatabase
)

var (
	shardScope             = linmetric.NewScope("lindb.tsdb.shard")
	writeMetricFailuresVec = shardScope.NewCounterVec("write_metric_failures", "db", "shard")
	writeBatchesVec        = shardScope.NewCounterVec("write_batches", "db", "shard")
	writeMetricsVec        = shardScope.NewCounterVec("write_metrics", "db", "shard")
	writeFieldsVec         = shardScope.NewCounterVec("write_fields", "db", "shard")
	memdbTotalSizeVec      = shardScope.NewGaugeVec("memdb_total_size", "db", "shard")
	memdbNumberVec         = shardScope.NewGaugeVec("memdb_number", "db", "shard")
	memFlushTimerVec       = shardScope.Scope("memdb_flush_duration").NewHistogramVec("db", "shard")
	indexFlushTimerVec     = shardScope.Scope("indexdb_flush_duration").NewHistogramVec("db", "shard")
)

const (
	segmentDir       = "segment"
	indexParentDir   = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
	metaDir          = "meta"
	bufferDir        = "buffer"
)

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// DatabaseName returns the database name.
	DatabaseName() string
	// ShardID returns the shard id.
	ShardID() models.ShardID
	// Path returns shard's storage directory.
	Path() string
	// CurrentInterval returns current interval for metric write.
	CurrentInterval() timeutil.Interval
	// Indicator returns the unique shard info.
	Indicator() string
	// GetOrCrateDataFamily returns data family, if not exist create a new data family.
	GetOrCrateDataFamily(familyTime int64) (DataFamily, error)
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// IndexDatabase returns the index-database
	IndexDatabase() indexdb.IndexDatabase
	BufferManager() memdb.BufferManager
	// WriteRows writes metric rows with same family in batch
	WriteRows(rows []metric.StorageRow) error

	Flush() error
	// initIndexDatabase initializes index database
	initIndexDatabase() error
	// Closer releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer
}

// shard implements Shard interface
// directory tree:
//    xx/shard/1/ (path)
//    xx/shard/1/buffer/123213123131 // time of ns
//    xx/shard/1/meta/
//    xx/shard/1/index/inverted/
//    xx/shard/1/data/20191012/
//    xx/shard/1/data/20191013/
type shard struct {
	databaseName string
	id           models.ShardID
	path         string
	option       option.DatabaseOption

	bufferMgr memdb.BufferManager
	indexDB   indexdb.IndexDatabase
	metadata  metadb.Metadata
	// write accept time range
	interval timeutil.Interval
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

	statistics struct {
		writeMetricFailures *linmetric.BoundCounter
		indexFlushTimer     *linmetric.BoundHistogram
	}
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(
	db Database,
	shardID models.ShardID,
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
	createdShard := &shard{
		databaseName: db.Name(),
		id:           shardID,
		path:         shardPath,
		option:       option,
		metadata:     db.Metadata(),
		bufferMgr:    memdb.NewBufferManager(filepath.Join(shardPath, bufferDir)),
		interval:     interval,
		segments:     make(map[timeutil.IntervalType]IntervalSegment),
		isFlushing:   *atomic.NewBool(false),
		logger:       logger.GetLogger("tsdb", "Shard"),
	}
	//try cleanup history dirty write buffer
	createdShard.bufferMgr.Cleanup()

	// initialize metrics
	shardIDStr := strconv.Itoa(int(shardID))
	createdShard.statistics.writeMetricFailures = writeMetricFailuresVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.indexFlushTimer = indexFlushTimerVec.WithTagValues(db.Name(), shardIDStr)

	// new segment for writing
	createdShard.segment, err = newIntervalSegmentFunc(
		createdShard,
		interval,
		filepath.Join(shardPath, segmentDir, interval.Type().String()))

	if err != nil {
		return nil, err
	}
	// add writing segment into segment list
	createdShard.segments[interval.Type()] = createdShard.segment

	defer func() {
		if err == nil {
			return
		}
		if err = createdShard.Close(); err != nil {
			engineLogger.Error("close shard error when create shard fail",
				logger.Any("shardID", createdShard.id),
				logger.String("database", createdShard.databaseName),
				logger.String("shard", createdShard.path), logger.Error(err))
		}
	}()
	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	return createdShard, nil
}

// DatabaseName returns the database name
func (s *shard) DatabaseName() string { return s.databaseName }

// ShardID returns the shard id.
func (s *shard) ShardID() models.ShardID { return s.id }

func (s *shard) Path() string {
	return s.path
}

// ShardInfo returns the unique shard info.
func (s *shard) Indicator() string { return s.path }

// CurrentInterval returns current interval for metric  write.
func (s *shard) CurrentInterval() timeutil.Interval { return s.interval }

func (s *shard) IndexDatabase() indexdb.IndexDatabase { return s.indexDB }

func (s *shard) BufferManager() memdb.BufferManager {
	return s.bufferMgr
}

func (s *shard) GetOrCrateDataFamily(familyTime int64) (DataFamily, error) {
	segmentName := s.interval.Calculator().GetSegment(familyTime)
	segment, err := s.segment.GetOrCreateSegment(segmentName)
	if err != nil {
		return nil, err
	}
	family, err := segment.GetOrCreateDataFamily(familyTime)
	if err != nil {
		return nil, err
	}
	return family, nil
}

func (s *shard) GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily {
	segment, ok := s.segments[intervalType]
	if ok {
		return segment.getDataFamilies(timeRange)
	}
	return nil
}

func (s *shard) lookupRowMeta(row *metric.StorageRow) (err error) {
	namespace := constants.DefaultNamespace
	metricName := string(row.Name())

	if len(row.NameSpace()) > 0 {
		namespace = string(row.NameSpace())
	}

	row.MetricID, err = s.metadata.MetadataDatabase().GenMetricID(namespace, metricName)
	if err != nil {
		s.statistics.writeMetricFailures.Incr()
		return err
	}
	var isCreated bool
	if row.TagsLen() == 0 {
		// if metric without tags, uses default series id(0)
		row.SeriesID = constants.SeriesIDWithoutTags
	} else {
		row.SeriesID, isCreated, err = s.indexDB.GetOrCreateSeriesID(row.MetricID, row.TagsHash())
		if err != nil {
			s.statistics.writeMetricFailures.Incr()
			return err
		}
	}
	if isCreated {
		// if series id is new, need build inverted index
		s.indexDB.BuildInvertIndex(
			namespace,
			metricName,
			row.NewKeyValueIterator(),
			row.SeriesID)
	}
	// set field id
	simpleFieldItr := row.NewSimpleFieldIterator()
	var fieldID field.ID
	for simpleFieldItr.HasNext() {
		if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
			namespace, metricName,
			simpleFieldItr.NextName(),
			simpleFieldItr.NextType()); err != nil {
			return err
		}
		row.FieldIDs = append(row.FieldIDs, fieldID)
	}

	compoundFieldItr, ok := row.NewCompoundFieldIterator()
	if !ok {
		goto Done
	}
	// min
	if compoundFieldItr.Min() > 0 {
		if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
			namespace, metricName, compoundFieldItr.HistogramMinFieldName(), field.MinField); err != nil {
			return err
		}
		row.FieldIDs = append(row.FieldIDs, fieldID)
	}
	// max
	if compoundFieldItr.Max() > 0 {
		if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
			namespace, metricName, compoundFieldItr.HistogramMaxFieldName(), field.MaxField); err != nil {
			return err
		}
		row.FieldIDs = append(row.FieldIDs, fieldID)
	}
	// sum
	if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
		namespace, metricName, compoundFieldItr.HistogramSumFieldName(), field.SumField); err != nil {
		return err
	}
	row.FieldIDs = append(row.FieldIDs, fieldID)
	// count
	if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
		namespace, metricName, compoundFieldItr.HistogramCountFieldName(), field.SumField); err != nil {
		return err
	}
	row.FieldIDs = append(row.FieldIDs, fieldID)
	// explicit bounds
	for compoundFieldItr.HasNextBucket() {
		if fieldID, err = s.metadata.MetadataDatabase().GenFieldID(
			namespace, metricName,
			compoundFieldItr.BucketName(), field.HistogramField); err != nil {
			return err
		}
		row.FieldIDs = append(row.FieldIDs, fieldID)
	}

Done:
	row.Writable = true
	return nil
}

func (s *shard) WriteRows(rows []metric.StorageRow) error {
	for idx := range rows {
		if err := s.lookupRowMeta(&rows[idx]); err != nil {
			s.logger.Error("failed to lookup meta of row", logger.Error(err))
			continue
		}
	}
	return nil
}

func (s *shard) Close() error {
	// wait previous flush job completed
	s.flushCondition.Wait()

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
	// close segment/flush family data
	s.segment.Close()
	return nil
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
				logger.Any("shardID", s.id),
				logger.String("database", s.databaseName),
				logger.Error(err))
			return err
		}
		s.logger.Info("flush indexDB successfully",
			logger.Any("shardID", s.id),
			logger.String("database", s.databaseName),
		)
		s.statistics.indexFlushTimer.UpdateSince(startTime)
	}

	//FIXME(stone1100) commit replica sequence
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
			Merger:           string(tagindex.SeriesForwardMerger)})
	if err != nil {
		return err
	}
	s.invertedFamily, err = s.indexStore.CreateFamily(
		invertedIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(tagindex.SeriesInvertedMerger)})
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
