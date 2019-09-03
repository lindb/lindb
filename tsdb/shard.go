package tsdb

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/memdb"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

const segmentPath = "segment"

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType interval.Type, timeRange timeutil.TimeRange) []DataFamily
	// GetSeriesIDsFilter returns series index for searching series(tags),
	// using this filter for filtering data in kv store.
	GetSeriesIDsFilter() series.Filter
	// GetMetaGetter returns tags meta getter
	GetMetaGetter() series.MetaGetter
	// GetMemoryDatabase returns memory database
	GetMemoryDatabase() memdb.MemoryDatabase
	// Write writes the metric-point into memory-database.
	Write(metric *pb.Metric) error
	// Close releases shard's resource, such as flush data, spawned goroutines etc.
	Close()
}

// shard implements Shard interface
type shard struct {
	id       int32
	path     string
	option   option.EngineOption
	memDB    memdb.MemoryDatabase
	indexDB  diskdb.IndexDatabase
	interval int64

	//TODO codingcrush add kv store for data storage

	// write accept time range
	ahead  int64
	behind int64

	segment IntervalSegment // smallest interval for writing data

	// segments keeps all interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	segments map[interval.Type]IntervalSegment
	cancel   context.CancelFunc
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(shardID int32, path string, index Index, option option.EngineOption) (Shard, error) {
	if err := option.Validation(); err != nil {
		return nil, fmt.Errorf("engine option is invalid, err:%s", err)
	}
	intervalVal, _ := timeutil.ParseInterval(option.Interval)
	intervalType := interval.CalcIntervalType(intervalVal)
	if err := fileutil.MkDirIfNotExist(path); err != nil {
		return nil, err
	}

	indexDB, err := index.CreateIndexDatabase(shardID)
	if err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error:%s", shardID, err)
	}

	// new segment for writing
	segment, err := newIntervalSegment(intervalVal, intervalType,
		filepath.Join(path, segmentPath, string(intervalType)))
	if err != nil {
		return nil, err
	}
	var memDB memdb.MemoryDatabase
	ctx, cancel := context.WithCancel(context.Background())
	memDB = memdb.NewMemoryDatabase(ctx, memdb.MemoryDatabaseCfg{
		TimeWindow:    option.TimeWindow,
		IntervalValue: intervalVal,
		IntervalType:  intervalType,
		Generator:     index.GetIDSequencer(),
	})
	shard := &shard{
		id:       shardID,
		path:     path,
		option:   option,
		memDB:    memDB,
		interval: intervalVal,
		indexDB:  indexDB,
		segment:  segment,
		segments: make(map[interval.Type]IntervalSegment),
		cancel:   cancel,
	}
	shard.ahead, _ = timeutil.ParseInterval(option.Ahead)
	shard.behind, _ = timeutil.ParseInterval(option.Behind)
	// add writing segment into segment list
	shard.segments[intervalType] = segment
	return shard, nil
}

// GetDataFamilies returns data family list by interval type and time range, return nil if not match
func (s *shard) GetDataFamilies(intervalType interval.Type, timeRange timeutil.TimeRange) []DataFamily {
	segment, ok := s.segments[intervalType]
	if ok {
		return segment.getDataFamilies(timeRange)
	}
	return nil
}

// GetSeriesIDsFilter returns series index for searching series(tags)
func (s *shard) GetSeriesIDsFilter() series.Filter {
	return s.indexDB
}

// GetMetaGetter returns tags meta getter
func (s *shard) GetMetaGetter() series.MetaGetter {
	return s.indexDB
}

// GetMemoryDatabase returns memory database
func (s *shard) GetMemoryDatabase() memdb.MemoryDatabase {
	return s.memDB
}

// Write writes the metric-point into memory-database.
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

// Close closes the memDatabase and spawned goroutines.
func (s *shard) Close() {
	s.cancel()
}
