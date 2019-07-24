package tsdb

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/eleme/lindb/tsdb/memdb"

	pb "github.com/eleme/lindb/rpc/proto/field"

	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/timeutil"
	"github.com/eleme/lindb/pkg/util"
)

const segmentPath = "segment"

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// GetSegments returns segment list by interval type and time range, return nil if not match
	GetSegments(intervalType interval.Type, timeRange timeutil.TimeRange) []Segment
	// Write writes the metric-point into memory-database.
	Write(metric *pb.Metric) error
	// Close releases shard's resource, such as flush data, spawned goroutines etc.
	Close()
}

// shard implements Shard interface
type shard struct {
	id     int32
	path   string
	option option.ShardOption
	memDB  memdb.MemoryDatabase

	segment IntervalSegment // smallest interval for writing data

	// segments keeps all interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	segments map[interval.Type]IntervalSegment
	cancel   context.CancelFunc
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(shardID int32, path string, option option.ShardOption) (Shard, error) {
	if option.Interval <= 0 {
		return nil, fmt.Errorf("interval cannot be negative")
	}
	if _, err := interval.GetCalculator(option.IntervalType); err != nil {
		return nil, fmt.Errorf("interval type[%d] not define", option.IntervalType)
	}
	if err := util.MkDirIfNotExist(path); err != nil {
		return nil, err
	}

	// new segment for writing
	segment, err := newIntervalSegment(option.Interval,
		option.IntervalType,
		filepath.Join(path, segmentPath, option.IntervalType.String()))
	if err != nil {
		return nil, err
	}
	var memDB memdb.MemoryDatabase
	ctx, cancel := context.WithCancel(context.Background())
	memDB, err = memdb.NewMemoryDatabase(ctx, option.TimeWindow, int64(option.Interval), option.IntervalType)
	if err != nil {
		//if create memory database error, cancel background context
		cancel()
		return nil, err
	}
	shard := &shard{
		id:       shardID,
		path:     path,
		option:   option,
		memDB:    memDB,
		segment:  segment,
		segments: make(map[interval.Type]IntervalSegment),
		cancel:   cancel,
	}
	// add writing segment into segment list
	shard.segments[option.IntervalType] = segment
	return shard, nil
}

// GetSegments returns segment list by interval type and time range, return nil if not match
func (s *shard) GetSegments(intervalType interval.Type, timeRange timeutil.TimeRange) []Segment {
	segment, ok := s.segments[intervalType]
	if ok {
		return segment.GetSegments(timeRange)
	}
	return nil
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
	if (s.option.Behind > 0 && timestamp < now-s.option.Behind) ||
		(s.option.Ahead > 0 && timestamp > now+s.option.Ahead) {
		return nil
	}

	// write metric point into memory db
	return s.memDB.Write(metric)
}

// Close closes the memDatabase and spawned goroutines.
func (s *shard) Close() {
	s.cancel()
}
