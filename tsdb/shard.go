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
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

// for testing
var (
	newReplicaSequenceFunc = newReplicaSequence
	newIntervalSegmentFunc = newIntervalSegment
	newKVStoreFunc         = kv.NewStore
	newIndexDBFunc         = indexdb.NewIndexDatabase
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
	// Flush index and memory data to disk
	Flush() error
	// IsFlushing checks if this shard is in flushing
	IsFlushing() bool
	// initIndexDatabase initializes index database
	initIndexDatabase() error
}

// shard implements Shard interface
// directory tree:
//    xx/shard/1/ (path)
//    xx/shard/1/replica
//    xx/shard/1/temp/123213123131 // time of ms
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

	mutable memdb.MemoryDatabase

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

	writing atomic.Bool
	rwMutex sync.RWMutex
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(
	db Database,
	shardID int32,
	shardPath string,
	option option.DatabaseOption,
) (
	s Shard,
	err error,
) {
	if err = option.Validate(); err != nil {
		return nil, fmt.Errorf("engine option is invalid, err: %s", err)
	}
	var interval timeutil.Interval
	_ = interval.ValueOf(option.Interval)

	if err := mkdirFunc(shardPath); err != nil {
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

	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	createdShard.mutable = memdb.NewMemoryDatabase(memdb.MemoryDatabaseCfg{
		Interval: interval,
		Metadata: createdShard.metadata,
		TempPath: filepath.Join(shardPath, tempDir),
	})
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

func (s *shard) MemoryDatabase() memdb.MemoryDatabase {
	return s.mutable
}

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
	seriesID, isCreated, err := s.indexDB.GetOrCreateSeriesID(metricID, metric.TagsHash)
	if err != nil {
		return err
	}

	if isCreated {
		// if series id is new, need build inverted index
		s.indexDB.BuildInvertIndex(ns, metric.Name, metric.Tags, seriesID)
	}

	var db memdb.MemoryDatabase
	s.rwMutex.Lock()
	s.writing.Store(true)
	db = s.mutable
	s.rwMutex.Unlock()

	// set write completed
	defer s.writing.Store(false)
	// write metric point into memory db
	return db.Write(ns, metric.Name, metricID, seriesID, metric.Timestamp, metric.Fields)
}

func (s *shard) Close() error {
	if err := s.indexDB.Close(); err != nil {
		return err
	}
	if err := s.Flush(); err != nil {
		return err
	}
	return s.indexStore.Close()
}

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
			Merger:           invertedIndexMerger})
	if err != nil {
		return err
	}
	s.invertedFamily, err = s.indexStore.CreateFamily(
		invertedIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           invertedIndexMerger})
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

func (s *shard) IsFlushing() bool { return s.isFlushing.Load() }

func (s *shard) Flush() (err error) {
	// another flush process is running
	if !s.isFlushing.CAS(false, true) {
		return nil
	}
	// mark flush job doing
	s.flushCondition.Add(1)
	allHeads := s.sequence.getAllHeads()

	defer func() {
		// if ack fail, maybe data will write duplicate if system restart
		if err := s.sequence.ack(allHeads); err != nil {
			engineLogger.Error("ack replica sequence error", logger.Error(err))
		}
		//TODO add commit kv meta after ack successfully

		// mark flush job complete, notify
		s.flushCondition.Done()
		s.isFlushing.Store(false)
	}()

	//FIXME stone1100
	// index flush
	//if err = s.memDB.FlushInvertedIndexTo(
	//	invertedindex.NewFlusher(s.invertedFamily.NewFlusher())); err != nil {
	//	return err
	//}

	for _, familyTime := range s.mutable.Families() {
		segmentName := s.interval.Calculator().GetSegment(familyTime)
		segment, err := s.segment.GetOrCreateSegment(segmentName)
		if err != nil {
			return err
		}
		thisDataFamily, err := segment.GetDataFamily(familyTime)
		if err != nil {
			continue
		}
		if err := s.mutable.FlushFamilyTo(
			metricsdata.NewFlusher(thisDataFamily.Family().NewFlusher()), familyTime); err != nil {
			return err
		}
	}
	return nil
}
