package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/eleme/lindb/tsdb/memdb"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/util"
)

const segmentPath = "segment"

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// GetSegments returns segment list by interval type and time range, return nil if not match
	GetSegments(intervalType interval.Type, timeRange models.TimeRange) []Segment
	// Write writes the metric-point into memory-database.
	Write(point models.Point) error
	// Close releases shard's resource, such as flush data, spawned goroutines etc.
	Close()
}

// shard implements Shard interface
type shard struct {
	id     int32
	path   string
	option option.ShardOption
	memDb  memdb.MemoryDatabase

	segment IntervalSegment // smallest interval for writing data

	// segments keeps all interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	segments map[interval.Type]IntervalSegment

	cancel func()

	intervalCalc interval.Calculator
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if fail.
func newShard(shardID int32, path string, option option.ShardOption) (Shard, error) {
	if option.Interval <= 0 {
		return nil, fmt.Errorf("interval cannot be negative")
	}
	if interval.GetCalculator(option.IntervalType) == nil {
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
	ctx, cancel := context.WithCancel(context.Background())
	shard := &shard{
		id:       shardID,
		path:     path,
		option:   option,
		memDb:    memdb.NewMemoryDatabase(ctx),
		segment:  segment,
		segments: make(map[interval.Type]IntervalSegment),
		cancel:   cancel,
	}
	// add writing segment into segment list
	shard.segments[option.IntervalType] = segment
	return shard, nil
}

// GetSegments returns segment list by interval type and time range, return nil if not match
func (s *shard) GetSegments(intervalType interval.Type, timeRange models.TimeRange) []Segment {
	segment, ok := s.segments[intervalType]
	if ok {
		return segment.GetSegments(timeRange)
	}
	return nil
}

// Write writes the metric-point into memory-database.
func (s *shard) Write(point models.Point) error {
	timestamp := point.Timestamp()
	now := time.Now().Unix()

	if timestamp < now-s.option.Behind {
		return nil
	}

	if timestamp > now+s.option.Ahead {
		return nil
	}
	//use family base time for memory store
	segmentTime := s.intervalCalc.CalFamilyBaseTime(point.Timestamp())
	//slot time for ts data store
	slotTime := s.intervalCalc.CalSlot(point.Timestamp())

	return s.memDb.Write(point, segmentTime, slotTime)
}

// Close closes the memDatabase and spawned goroutines.
func (s *shard) Close() {
	s.cancel()
}
