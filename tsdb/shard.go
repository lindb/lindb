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

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
	"github.com/lindb/lindb/tsdb/tblstore/tagindex"
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
	ShardID() models.ShardID
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
	// WriteRows writes metric rows with same family in batch
	WriteRows(familyTime int64, rows []metric.StorageRow) error
	// GetOrCreateSequence gets the replica sequence by given remote peer if exist, else creates a new sequence
	GetOrCreateSequence(replicaPeer string) (queue.Sequence, error)
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
	id           models.ShardID
	path         string
	option       option.DatabaseOption
	sequence     ReplicaSequence

	mutex    sync.Mutex     // mutex for update families
	families familyMemDBSet // memory database for each family time

	indexDB  indexdb.IndexDatabase
	metadata metadb.Metadata
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
		writeBatches        *linmetric.BoundCounter
		writeMetrics        *linmetric.BoundCounter
		writeMetricFailures *linmetric.BoundCounter
		writeFields         *linmetric.BoundCounter
		memdbTotalSize      *linmetric.BoundGauge
		memdbNumber         *linmetric.BoundGauge
		memFlushTimer       *linmetric.BoundHistogram
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
		logger:       logger.GetLogger("tsdb", "Shard"),
	}
	// initialize metrics
	shardIDStr := strconv.Itoa(int(shardID))
	createdShard.statistics.writeBatches = writeBatchesVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.writeMetrics = writeMetricsVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.writeMetricFailures = writeMetricFailuresVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.writeFields = writeFieldsVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.memdbTotalSize = memdbTotalSizeVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.memdbNumber = memdbNumberVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.memFlushTimer = memFlushTimerVec.WithTagValues(db.Name(), shardIDStr)
	createdShard.statistics.indexFlushTimer = indexFlushTimerVec.WithTagValues(db.Name(), shardIDStr)

	// new segment for writing
	createdShard.segment, err = newIntervalSegmentFunc(
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
	// add shard into global shard manager
	GetShardManager().AddShard(createdShard)
	return createdShard, nil
}

// DatabaseName returns the database name
func (s *shard) DatabaseName() string { return s.databaseName }

// ShardID returns the shard id.
func (s *shard) ShardID() models.ShardID { return s.id }

// ShardInfo returns the unique shard info.
func (s *shard) ShardInfo() string { return s.path }

// CurrentInterval returns current interval for metric  write.
func (s *shard) CurrentInterval() timeutil.Interval { return s.interval }

func (s *shard) GetOrCreateSequence(replicaPeer string) (queue.Sequence, error) {
	return s.sequence.getOrCreateSequence(replicaPeer)
}

func (s *shard) IndexDatabase() indexdb.IndexDatabase { return s.indexDB }

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
		familyStartTime := entries[idx].familyTime
		familyEndTime := s.interval.Calculator().CalcFamilyEndTime(familyStartTime)
		if !timeRange.Overlap(timeutil.TimeRange{Start: familyStartTime, End: familyEndTime}) {
			continue
		}
		resultSet, err := entries[idx].memDB.Filter(metricID, seriesIDs, timeRange, fields)
		if err != nil {
			return nil, err
		}
		rs = append(rs, resultSet...)
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

func (s *shard) WriteRows(familyTime int64, rows []metric.StorageRow) error {
	defer s.statistics.writeBatches.Incr()

	intervalCalc := s.interval.Calculator()

	for idx := range rows {
		if err := s.lookupRowMeta(&rows[idx]); err != nil {
			s.logger.Error("failed to lookup meta of row", logger.Error(err))
			continue
		}
		rows[idx].SlotIndex = uint16(intervalCalc.CalcSlot(
			rows[idx].Timestamp(),
			familyTime,
			s.interval.Int64()),
		)
	}
	db, err := s.GetOrCreateMemoryDatabase(familyTime)
	if err != nil {
		// all rows are dropped
		s.statistics.writeMetricFailures.Add(float64(len(rows)))
		return err
	}
	db.AcquireWrite()
	defer db.CompleteWrite()

	releaseFunc := db.WithLock()
	defer releaseFunc()

	for idx := range rows {
		if !rows[idx].Writable {
			s.statistics.writeMetricFailures.Incr()
			continue
		}
		if err = db.WriteRow(&rows[idx]); err == nil {
			s.statistics.writeMetrics.Incr()
			s.statistics.writeFields.Add(float64(len(rows[idx].FieldIDs)))
		} else {
			s.statistics.writeMetricFailures.Incr()
			s.logger.Error("failed writing row", logger.Error(err))
		}
	}
	// if memdb size is above threshold, it will be put into immutable list
	s.validateMemDBSize(familyTime, db)
	return nil
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
		logger.Any("shardID", s.id),
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
				logger.Any("shardID", s.id),
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
	s.statistics.memdbNumber.Update(float64(len(s.families.Entries())))
	s.statistics.memdbTotalSize.Update(float64(s.MemDBTotalSize()))

	s.tryEvictMutable()

	cfg := config.GlobalStorageConfig()
	// too many memdbs
	number := len(s.families.Entries())
	if number > cfg.TSDB.MaxMemDBNumber {
		s.logger.Info("number of memdb is above threshold, waiting for flush",
			logger.Any("shardID", s.id),
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
			logger.Any("shardID", s.id),
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
				logger.Any("shardID", s.id),
				logger.String("database", s.databaseName),
				logger.Int64("familyTime", waitingFlushMemDB.FamilyTime()),
				logger.Int64("memDBSize", waitingFlushMemDB.MemSize()),
			)
		}
		s.mutex.Unlock()
	}
	if waitingFlushMemDB == nil {
		s.logger.Warn("there is no memdb to flush", logger.Any("shardID", s.id))
		return nil
	}

	startTime = time.Now()
	if err := s.flushMemoryDatabase(waitingFlushMemDB); err != nil {
		s.logger.Error("failed to flush memdb",
			logger.Any("shardID", s.id),
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
		logger.Any("shardID", s.id),
		logger.String("database", s.databaseName),
		logger.String("flush-duration", endTime.Sub(startTime).String()),
		logger.Int64("familyTime", waitingFlushMemDB.FamilyTime()),
		logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
	s.statistics.memFlushTimer.UpdateDuration(endTime.Sub(startTime))

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
	defer s.statistics.memFlushTimer.UpdateSince(startTime)

	segmentName := s.interval.Calculator().GetSegment(memDB.FamilyTime())
	segment, err := s.segment.GetOrCreateSegment(segmentName)
	if err != nil {
		return err
	}
	thisDataFamily, err := segment.GetDataFamily(memDB.FamilyTime())
	if err != nil {
		return err
	}
	dataFlusher, err := metricsdata.NewFlusher(thisDataFamily.Family().NewFlusher())
	if err != nil {
		return err
	}
	// flush family data
	if err := memDB.FlushFamilyTo(dataFlusher); err != nil {
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
