package tsdb

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/replication"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/forwardindex"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

const (
	replicaDir       = "replica"
	segmentDir       = "segment"
	indexParDir      = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
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

	MemoryFilter() series.Filter
	IndexFilter() series.Filter
	MemoryMetaGetter() series.MetaGetter
	IndexMetaGetter() series.MetaGetter
	// initIndexDatabase initializes index database
	initIndexDatabase() error
}

// shard implements Shard interface
// directory tree:
//    xx/shard/1/ (path)
//    xx/shard/1/replica
//    xx/shard/1/index/forward/
//    xx/shard/1/index/inverted/
//    xx/shard/1/data/20191012/
//    xx/shard/1/data/20191013/
type shard struct {
	databaseName string
	id           int32
	path         string
	option       option.DatabaseOption
	sequence     *replicaSequence
	memDB        memdb.MemoryDatabase
	indexDB      indexdb.IndexDatabase
	idSequencer  metadb.IDSequencer
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

	cancel         context.CancelFunc // cancel function
	indexStore     kv.Store           // kv stores
	invertedFamily kv.Family
	forwardFamily  kv.Family
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(
	databaseName string,
	shardID int32,
	shardPath string,
	idSequencer metadb.IDSequencer,
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

	if err := fileutil.MkDirIfNotExist(shardPath); err != nil {
		return nil, err
	}
	replicaSequence, err := newReplicaSequence(filepath.Join(shardPath, replicaDir))
	if err != nil {
		return nil, err
	}
	createdShard := &shard{
		databaseName: databaseName,
		id:           shardID,
		path:         shardPath,
		option:       option,
		sequence:     replicaSequence,
		interval:     interval,
		idSequencer:  idSequencer,
		segments:     make(map[timeutil.IntervalType]IntervalSegment),
		isFlushing:   *atomic.NewBool(false),
	}
	// new segment for writing
	createdShard.segment, err = newIntervalSegment(
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
	var ctx context.Context
	ctx, createdShard.cancel = context.WithCancel(context.Background())
	createdShard.memDB = memdb.NewMemoryDatabase(ctx, memdb.MemoryDatabaseCfg{
		TimeWindow: option.TimeWindow,
		Interval:   interval,
		Generator:  idSequencer,
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
	return s.memDB
}

func (s *shard) Write(metric *pb.Metric) error {
	if metric == nil {
		return fmt.Errorf("metric is nil")
	}
	if metric.Fields == nil {
		return fmt.Errorf("fields is nil")
	}
	timestamp := metric.Timestamp
	now := timeutil.Now()

	// check metric timestamp if in acceptable time range
	if (s.behind.Int64() > 0 && timestamp < now-s.behind.Int64()) ||
		(s.ahead.Int64() > 0 && timestamp > now+s.ahead.Int64()) {
		return nil
	}

	// if doing flush job, need wait flush job completed
	if s.isFlushing.Load() {
		s.flushCondition.Wait()
	}

	// write metric point into memory db
	return s.memDB.Write(metric)
}

func (s *shard) Close() error {
	if err := s.Flush(); err != nil {
		return err
	}
	defer s.cancel()
	return s.indexStore.Close()
}

func (s *shard) initIndexDatabase() error {
	var err error
	storeOption := kv.DefaultStoreOption(filepath.Join(s.path, indexParDir))
	s.indexStore, err = kv.NewStore(storeOption.Path, storeOption)
	if err != nil {
		return err
	}
	s.invertedFamily, err = s.indexStore.CreateFamily(
		forwardIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           invertedIndexMerger})
	if err != nil {
		return err
	}
	s.forwardFamily, err = s.indexStore.CreateFamily(
		invertedIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           forwardIndexMerger})
	if err != nil {
		return err
	}
	s.indexDB = indexdb.NewIndexDatabase(s.idSequencer, s.invertedFamily, s.forwardFamily)
	return nil
}

func (s *shard) MemoryFilter() series.Filter         { return s.memDB }
func (s *shard) IndexFilter() series.Filter          { return s.indexDB }
func (s *shard) MemoryMetaGetter() series.MetaGetter { return s.memDB }
func (s *shard) IndexMetaGetter() series.MetaGetter  { return s.indexDB }
func (s *shard) IsFlushing() bool                    { return s.isFlushing.Load() }

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

	if err = s.memDB.FlushForwardIndexTo(
		forwardindex.NewFlusher(s.forwardFamily.NewFlusher())); err != nil {
		return err
	}
	if err = s.memDB.FlushInvertedIndexTo(
		invertedindex.NewFlusher(s.invertedFamily.NewFlusher())); err != nil {
		return err
	}

	for _, familyTime := range s.memDB.Families() {
		segmentName := s.interval.Calculator().GetSegment(familyTime)
		segment, err := s.segment.GetOrCreateSegment(segmentName)
		if err != nil {
			return err
		}
		thisDataFamily, err := segment.GetDataFamily(familyTime)
		if err != nil {
			continue
		}
		if err := s.memDB.FlushFamilyTo(
			metricsdata.NewFlusher(thisDataFamily.Family().NewFlusher()), familyTime); err != nil {
			return err
		}
	}
	return nil
}
