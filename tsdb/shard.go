package tsdb

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/replication"
	pb "github.com/lindb/lindb/rpc/proto/field"
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
	// ShardInfo returns the unique shard info
	ShardInfo() string
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// MemoryDatabase returns memory database
	MemoryDatabase() memdb.MemoryDatabase
	// IndexDatabase returns the index-database
	IndexDatabase() indexdb.IndexDatabase
	// Write writes the metric-point into memory-database.
	Write(metric *pb.Metric) error
	// GetOrCreateSequence gets the replica sequence by given remote peer if exist, else creates a new sequence
	GetOrCreateSequence(replicaPeer string) (replication.Sequence, error)
	// Close releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer
	// Flush flushes index and memory data to disk
	Flush() error
	// NeedFlush checks if shard need to flush memory data
	NeedFlush() bool
	// IsFlushing checks if this shard is in flushing
	IsFlushing() bool
	// initIndexDatabase initializes index database
	initIndexDatabase() error
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

	mutable   memdb.MemoryDatabase // current accept user write data points
	immutable memdb.MemoryDatabase // need flush data to disk persist

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

	rwMutex sync.RWMutex
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
		metadata:     db.Metadata(),
		interval:     interval,
		segments:     make(map[timeutil.IntervalType]IntervalSegment),
		isFlushing:   *atomic.NewBool(false),
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
		if err != nil {
			if err := createdShard.Close(); err != nil {
				engineLogger.Error("close shard error when create shard fail",
					logger.String("shard", createdShard.path), logger.Error(err))
			}
		}
	}()
	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	memDB, err := createdShard.createMemoryDatabase()
	if err != nil {
		return nil, err
	}
	createdShard.mutable = memDB
	// add shard into global shard manager
	GetShardManager().AddShard(createdShard)
	return createdShard, nil
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

// MemoryDatabase returns memory database
func (s *shard) MemoryDatabase() memdb.MemoryDatabase {
	var memDB memdb.MemoryDatabase
	s.rwMutex.RLock()
	memDB = s.mutable
	s.rwMutex.RUnlock()
	return memDB
}

// Write writes the metric-point into memory-database.
func (s *shard) Write(metric *pb.Metric) error {
	if metric == nil {
		return constants.ErrNilMetric
	}
	if len(metric.Name) == 0 {
		return constants.ErrEmptyMetricName
	}
	if len(metric.Fields) == 0 {
		return constants.ErrEmptyField
	}
	timestamp := metric.Timestamp
	now := timeutil.Now()

	// check metric timestamp if in acceptable time range
	if (s.behind.Int64() > 0 && timestamp < now-s.behind.Int64()) ||
		(s.ahead.Int64() > 0 && timestamp > now+s.ahead.Int64()) {
		return nil
	}
	ns := metric.Namespace
	if len(ns) == 0 {
		ns = constants.DefaultNamespace
	}
	metricID, err := s.metadata.MetadataDatabase().GenMetricID(ns, metric.Name)
	if err != nil {
		return err
	}
	var seriesID uint32
	isCreated := false
	if len(metric.Tags) == 0 {
		// if metric without tags, uses default series id(0)
		seriesID = constants.SeriesIDWithoutTags
	} else {
		seriesID, isCreated, err = s.indexDB.GetOrCreateSeriesID(metricID, metric.TagsHash)
		if err != nil {
			return err
		}
	}

	if isCreated {
		// if series id is new, need build inverted index
		s.indexDB.BuildInvertIndex(ns, metric.Name, metric.Tags, seriesID)
	}
	db := s.MemoryDatabase()

	// mark writing data
	db.AcquireWrite()
	// set write completed
	defer db.CompleteWrite()
	// write metric point into memory db
	return db.Write(ns, metric.Name, metricID, seriesID, metric.Timestamp, metric.Fields)
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
	if s.mutable != nil {
		if err := s.flushMemoryDatabase(s.mutable); err != nil {
			return err
		}
	}
	if s.immutable != nil {
		if err := s.flushMemoryDatabase(s.immutable); err != nil {
			return err
		}
	}
	s.ackReplicaSeq()
	return nil
}

// IsFlushing checks if this shard is in flushing
func (s *shard) IsFlushing() bool { return s.isFlushing.Load() }

// NeedFlush checks if shard need to flush memory data
func (s *shard) NeedFlush() bool {
	if s.IsFlushing() {
		return false
	}
	if s.hasImmutable() {
		return false
	}

	memDB := s.MemoryDatabase()
	//TODO add time threshold???
	return memDB.MemSize() > constants.ShardMemoryUsedThreshold
}

// Flush flushes index and memory data to disk
func (s *shard) Flush() (err error) {
	// another flush process is running
	if !s.isFlushing.CAS(false, true) {
		return nil
	}
	// 1. swap memory database
	s.swapMemoryDatabase()
	// 2. mark flush job doing
	s.flushCondition.Add(1)

	defer func() {
		//TODO add commit kv meta after ack successfully
		// mark flush job complete, notify
		s.flushCondition.Done()
		s.isFlushing.Store(false)
	}()

	//FIXME stone1100
	// index flush
	if s.indexDB != nil {
		if err = s.indexDB.Flush(); err != nil {
			return err
		}
	}

	// flush immutable, if exist
	// maybe not exist, when swap fail
	if s.immutable != nil {
		if err := s.flushMemoryDatabase(s.immutable); err != nil {
			return err
		}
		// after flush success, mark immutable as nil
		s.rwMutex.Lock()
		s.immutable = nil
		s.rwMutex.Unlock()
	}
	// finally, commit replica sequence
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
		s.databaseName,
		filepath.Join(s.path, metaDir),
		s.metadata, s.forwardFamily,
		s.invertedFamily)
	if err != nil {
		return err
	}
	return nil
}

// hasImmutable checks if has immutable memory database
func (s *shard) hasImmutable() bool {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	return s.immutable != nil
}

// swapMemoryDatabase swaps mutable/immutable memory database
func (s *shard) swapMemoryDatabase() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	// if immutable is not nil, cannot do swap
	if s.immutable != nil {
		return
	}

	memDB, err := s.createMemoryDatabase()
	if err != nil {
		engineLogger.Error("create new memory database error when swap",
			logger.String("shard", s.path), logger.Error(err))
		return
	}
	s.immutable = s.mutable // mark old memory database is immutable
	s.mutable = memDB       // create new  memory database as mutable
}

// createMemoryDatabase creates a new memory database for writing data points
func (s *shard) createMemoryDatabase() (memdb.MemoryDatabase, error) {
	return newMemoryDBFunc(memdb.MemoryDatabaseCfg{
		Interval: s.interval,
		Metadata: s.metadata,
		TempPath: filepath.Join(s.path, filepath.Join(tempDir, fmt.Sprintf("%d", timeutil.Now()))),
	})
}

// flushMemoryDatabase flushes memory database to disk kv store
func (s *shard) flushMemoryDatabase(memDB memdb.MemoryDatabase) error {
	for _, familyTime := range memDB.Families() {
		segmentName := s.interval.Calculator().GetSegment(familyTime)
		segment, err := s.segment.GetOrCreateSegment(segmentName)
		if err != nil {
			return err
		}
		thisDataFamily, err := segment.GetDataFamily(familyTime)
		if err != nil {
			continue
		}
		// flush family data
		if err := memDB.FlushFamilyTo(
			metricsdata.NewFlusher(thisDataFamily.Family().NewFlusher()), familyTime); err != nil {
			return err
		}
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
