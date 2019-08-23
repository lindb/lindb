package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./segment.go -destination=./segment_mock.go -package=tsdb -self_package=github.com/lindb/lindb/tsdb

// IntervalSegment represents a interval segment, there are some segments in a shard.
type IntervalSegment interface {
	// GetOrCreateSegment creates new segment if not exist, if exist return it
	GetOrCreateSegment(segmentName string) (Segment, error)
	// GetSegments returns segment list by time range, return nil if not match
	GetSegments(timeRange timeutil.TimeRange) []Segment
	// Close closes interval segment, release resource
	Close()
}

// intervalSegment implements IntervalSegment interface
type intervalSegment struct {
	path         string
	interval     int64
	intervalType interval.Type

	segments sync.Map

	mutex sync.Mutex
}

// newIntervalSegment create interval segment based on interval/type/path etc.
func newIntervalSegment(interval int64, intervalType interval.Type, path string) (IntervalSegment, error) {
	if err := fileutil.MkDirIfNotExist(path); err != nil {
		return nil, err
	}
	intervalSegment := &intervalSegment{
		path:         path,
		interval:     interval,
		intervalType: intervalType,
	}

	// load segments if exist
	//TODO too many kv store load???
	segmentNames, err := fileutil.ListDir(path)
	if err != nil {
		//TODO return error????
		return nil, err
	}
	for _, segmentName := range segmentNames {
		seg, err := newSegment(segmentName, intervalType, filepath.Join(path, segmentName))
		if err != nil {
			return nil, fmt.Errorf("create segmenet error:%s", err)
		}
		intervalSegment.segments.Store(segmentName, seg)
	}

	return intervalSegment, nil
}

// GetOrCreateSegment creates new segment if not exist, if exist return it
func (s *intervalSegment) GetOrCreateSegment(segmentName string) (Segment, error) {
	segment := s.getSegment(segmentName)
	if segment == nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// double check, make sure only create segment once
		segment = s.getSegment(segmentName)
		if segment == nil {
			seg, err := newSegment(segmentName, s.intervalType, filepath.Join(s.path, segmentName))
			if err != nil {
				return nil, fmt.Errorf("create segmenet error:%s", err)
			}
			s.segments.Store(segmentName, seg)
			return seg, nil
		}
	}
	return segment, nil
}

// GetSegments returns segment list by time range, return nil if not match
func (s *intervalSegment) GetSegments(timeRange timeutil.TimeRange) []Segment {
	calc, err := interval.GetCalculator(s.intervalType)
	if err != nil {
		return nil
	}

	var segments []Segment
	start := calc.CalcSegmentTime(timeRange.Start)
	end := calc.CalcSegmentTime(timeRange.End)
	s.segments.Range(func(k, v interface{}) bool {
		segment, ok := v.(Segment)
		if ok {
			baseTime := segment.BaseTime()
			if baseTime >= start && baseTime <= end {
				segments = append(segments, segment)
			}
		}
		return true
	})
	return segments
}

// Close closes interval segment, release resource
func (s *intervalSegment) Close() {
	s.segments.Range(func(k, v interface{}) bool {
		seg, ok := v.(Segment)
		if ok {
			seg.Close()
		}
		return true
	})
}

// getSegment returns segment by name, if not exist return nil
func (s *intervalSegment) getSegment(segmentName string) Segment {
	segment, _ := s.segments.Load(segmentName)
	seg, ok := segment.(Segment)
	if ok {
		return seg
	}
	return nil
}

// Segment represents a time based segment, there are some segments in a interval segment.
// A segment use k/v store for storing time series data.
type Segment interface {
	// BaseTime returns segment base time
	BaseTime() int64
	// GetDataFamilies returns data family list by time range, return nil if not match
	GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// Close closes segment, include kv store
	Close()
}

// segment implements Segment interface
type segment struct {
	baseTime     int64
	kvStore      kv.Store
	intervalType interval.Type
	//TODO
	// families     map[int64]kv.Family

	logger *logger.Logger
}

// newSegment returns segment, segment is wrapper of kv store
func newSegment(segmentName string, intervalType interval.Type, path string) (Segment, error) {
	kvStore, err := kv.NewStore(segmentName, kv.StoreOption{Path: path})
	if err != nil {
		return nil, fmt.Errorf("create  kv store for segment error:%s", err)
	}
	// parse base time from segment name
	calc, err := interval.GetCalculator(intervalType)
	if err != nil {
		return nil, err
	}
	baseTime, err := calc.ParseSegmentTime(segmentName)
	if err != nil {
		return nil, fmt.Errorf("parse segment[%s] base time error", path)
	}

	return &segment{
		baseTime:     baseTime,
		kvStore:      kvStore,
		intervalType: intervalType,
		logger:       logger.GetLogger("tsdb", "Segment"),
	}, nil
}

// BaseTime returns segment base time
func (s *segment) BaseTime() int64 {
	return s.baseTime
}
func (s *segment) GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	//TODO need impl
	return nil
}

// Close closes segment, include kv store
func (s *segment) Close() {
	if err := s.kvStore.Close(); err != nil {
		s.logger.Error("close kv store error", logger.Error(err))
	}
}
