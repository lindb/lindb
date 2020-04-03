package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./interval_segment.go -destination=./interval_segment_mock.go -package=tsdb

// IntervalSegment represents a interval segment, there are some segments in a shard.
type IntervalSegment interface {
	// GetOrCreateSegment creates new segment if not exist, if exist return it
	GetOrCreateSegment(segmentName string) (Segment, error)
	// getDataFamilies returns data family list by time range, return nil if not match
	getDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// Close closes interval segment, release resource
	Close()
}

// intervalSegment implements IntervalSegment interface
type intervalSegment struct {
	path     string
	interval timeutil.Interval
	segments sync.Map

	mutex sync.Mutex
}

// newIntervalSegment create interval segment based on interval/type/path etc.
func newIntervalSegment(
	interval timeutil.Interval,
	path string,
) (
	segment IntervalSegment,
	err error,
) {
	if err = mkDirIfNotExist(path); err != nil {
		return segment, err
	}
	intervalSegment := &intervalSegment{
		path:     path,
		interval: interval,
	}

	defer func() {
		// if create or init segment fail, need close segment
		if err != nil {
			intervalSegment.Close()
		}
	}()

	// load segments if exist
	//TODO too many kv store load???
	segmentNames, err := listDir(path)
	if err != nil {
		return segment, err
	}
	for _, segmentName := range segmentNames {
		seg, err := newSegment(segmentName, intervalSegment.interval, filepath.Join(path, segmentName))
		if err != nil {
			err = fmt.Errorf("create segmenet error: %s", err)
			return segment, err
		}
		intervalSegment.segments.Store(segmentName, seg)
	}

	// set segment
	segment = intervalSegment
	return segment, err
}

// GetOrCreateSegment creates new segment if not exist, if exist return it
func (s *intervalSegment) GetOrCreateSegment(segmentName string) (Segment, error) {
	segment, ok := s.getSegment(segmentName)
	if !ok {
		// double check, make sure only create segment once
		s.mutex.Lock()
		defer s.mutex.Unlock()
		segment, ok = s.getSegment(segmentName)
		if !ok {
			seg, err := newSegment(segmentName, s.interval, filepath.Join(s.path, segmentName))
			if err != nil {
				return nil, fmt.Errorf("create segmenet error: %s", err)
			}
			s.segments.Store(segmentName, seg)
			return seg, nil
		}
	}
	return segment, nil
}

// getDataFamilies returns data family list by time range, return nil if not match
func (s *intervalSegment) getDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	intervalCalc := s.interval.Calculator()
	segmentQueryTimeRange := &timeutil.TimeRange{
		Start: intervalCalc.CalcSegmentTime(timeRange.Start), // need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		End:   timeRange.End,
	}
	s.segments.Range(func(k, v interface{}) bool {
		segment, ok := v.(Segment)
		if ok {
			baseTime := segment.BaseTime()
			if segmentQueryTimeRange.Contains(baseTime) {
				familyQueryTimeRange := segmentQueryTimeRange.Intersect(&timeRange)
				families := segment.getDataFamilies(*familyQueryTimeRange)
				if len(families) > 0 {
					result = append(result, families...)
				}
			}
		}
		return true
	})
	return result
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

// getSegment returns segment by name
func (s *intervalSegment) getSegment(segmentName string) (Segment, bool) {
	segment, _ := s.segments.Load(segmentName)
	seg, ok := segment.(Segment)
	return seg, ok
}
