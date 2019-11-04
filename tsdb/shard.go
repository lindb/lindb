package tsdb

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

const (
	segmentDir       = "segment"
	indexParDir      = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
)

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType interval.Type, timeRange timeutil.TimeRange) []DataFamily
	// MemoryDatabase returns memory database
	MemoryDatabase() memdb.MemoryDatabase
	// IndexDatabase returns the index-database
	IndexDatabase() indexdb.IndexDatabase
	// Write writes the metric-point into memory-database.
	Write(metric *pb.Metric) error
	// Close releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer

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
//    xx/shard/1/index/forward/
//    xx/shard/1/index/inverted/
//    xx/shard/1/data/20191012/
//    xx/shard/1/data/20191013/
type shard struct {
	id          int32
	path        string
	interval    int64
	option      option.DatabaseOption
	memDB       memdb.MemoryDatabase
	indexDB     indexdb.IndexDatabase
	idSequencer metadb.IDSequencer
	// write accept time range
	ahead  int64
	behind int64

	// segments keeps all interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	segments   map[interval.Type]IntervalSegment
	segment    IntervalSegment    // smallest interval for writing data
	cancel     context.CancelFunc // cancel function
	indexStore kv.Store           // kv stores
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(
	shardID int32,
	shardPath string,
	idSequencer metadb.IDSequencer,
	option option.DatabaseOption,
) (
	s Shard,
	err error,
) {
	if err = option.Validation(); err != nil {
		return nil, fmt.Errorf("engine option is invalid, err: %s", err)
	}
	intervalVal, _ := timeutil.ParseInterval(option.Interval)
	intervalType := interval.CalcIntervalType(intervalVal)
	if err := fileutil.MkDirIfNotExist(shardPath); err != nil {
		return nil, err
	}
	createdShard := &shard{
		id:          shardID,
		path:        shardPath,
		option:      option,
		interval:    intervalVal,
		idSequencer: idSequencer,
		segments:    make(map[interval.Type]IntervalSegment),
	}
	// new segment for writing
	createdShard.segment, err = newIntervalSegment(
		intervalVal,
		intervalType,
		filepath.Join(shardPath, segmentDir, string(intervalType)))
	if err != nil {
		return nil, err
	}
	createdShard.ahead, _ = timeutil.ParseInterval(option.Ahead)
	createdShard.behind, _ = timeutil.ParseInterval(option.Behind)
	// add writing segment into segment list
	createdShard.segments[intervalType] = createdShard.segment

	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	var ctx context.Context
	ctx, createdShard.cancel = context.WithCancel(context.Background())
	createdShard.memDB = memdb.NewMemoryDatabase(ctx, memdb.MemoryDatabaseCfg{
		TimeWindow:    option.TimeWindow,
		IntervalValue: intervalVal,
		IntervalType:  intervalType,
		Generator:     idSequencer,
	})
	return createdShard, nil
}

func (s *shard) IndexDatabase() indexdb.IndexDatabase {
	return s.indexDB
}

func (s *shard) GetDataFamilies(intervalType interval.Type, timeRange timeutil.TimeRange) []DataFamily {
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
	if (s.behind > 0 && timestamp < now-s.behind) ||
		(s.ahead > 0 && timestamp > now+s.ahead) {
		return nil
	}
	// write metric point into memory db
	return s.memDB.Write(metric)
}

func (s *shard) Close() error {
	defer s.cancel()
	return s.indexStore.Close()
}

func (s *shard) initIndexDatabase() error {
	storeOption := kv.DefaultStoreOption(filepath.Join(s.path, indexParDir))
	indexStore, err := kv.NewStore(storeOption.Path, storeOption)
	if err != nil {
		return err
	}

	invertedFamily, err := indexStore.CreateFamily(
		forwardIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           invertedIndexMerger})
	if err != nil {
		return err
	}
	forwardFamily, err := indexStore.CreateFamily(
		invertedIndexDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           forwardIndexMerger})
	if err != nil {
		return err
	}
	s.indexDB = indexdb.NewIndexDatabase(s.idSequencer, invertedFamily, forwardFamily)
	s.indexStore = indexStore
	return nil
}

func (s *shard) MemoryFilter() series.Filter         { return s.memDB }
func (s *shard) IndexFilter() series.Filter          { return s.indexDB }
func (s *shard) MemoryMetaGetter() series.MetaGetter { return s.memDB }
func (s *shard) IndexMetaGetter() series.MetaGetter  { return s.indexDB }
