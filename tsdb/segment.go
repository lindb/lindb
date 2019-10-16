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
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./segment.go -destination=./segment_mock.go -package=tsdb

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
	path         string
	interval     int64
	intervalType interval.Type
	calc         interval.Calculator

	segments sync.Map

	mutex sync.Mutex
}

// newIntervalSegment create interval segment based on interval/type/path etc.
func newIntervalSegment(intervalVal int64, intervalType interval.Type, path string) (segment IntervalSegment, err error) {
	calc := interval.GetCalculator(intervalType)
	if err = fileutil.MkDirIfNotExist(path); err != nil {
		return segment, err
	}
	intervalSegment := &intervalSegment{
		path:         path,
		interval:     intervalVal,
		intervalType: intervalType,
		calc:         calc,
	}

	defer func() {
		// if create or init segment fail, need close segment
		if err != nil {
			intervalSegment.Close()
		}
	}()

	// load segments if exist
	//TODO too many kv store load???
	segmentNames, err := fileutil.ListDir(path)
	if err != nil {
		return segment, err
	}
	for _, segmentName := range segmentNames {
		seg, err := newSegment(segmentName, intervalSegment.interval, intervalType, filepath.Join(path, segmentName))
		if err != nil {
			err = fmt.Errorf("create segmenet error:%s", err)
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
	segment := s.getSegment(segmentName)
	if segment == nil {
		// double check, make sure only create segment once
		s.mutex.Lock()
		defer s.mutex.Unlock()
		segment = s.getSegment(segmentName)
		if segment == nil {
			seg, err := newSegment(segmentName, s.interval, s.intervalType, filepath.Join(s.path, segmentName))
			if err != nil {
				return nil, fmt.Errorf("create segmenet error:%s", err)
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
	segmentQueryTimeRange := &timeutil.TimeRange{
		Start: s.calc.CalcSegmentTime(timeRange.Start), // need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
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
	// GetDataFamily returns the data family based on timestamp
	GetDataFamily(timestamp int64) (DataFamily, error)
	// Close closes segment, include kv store
	Close()

	// getDataFamilies returns data family list by time range, return nil if not match
	getDataFamilies(timeRange timeutil.TimeRange) []DataFamily
}

// segment implements Segment interface
type segment struct {
	baseTime     int64
	kvStore      kv.Store
	interval     int64
	intervalType interval.Type
	calc         interval.Calculator
	families     sync.Map

	mutex sync.Mutex

	logger *logger.Logger
}

// newSegment returns segment, segment is wrapper of kv store
func newSegment(segmentName string, intervalVal int64, intervalType interval.Type, path string) (Segment, error) {
	// parse base time from segment name
	calc := interval.GetCalculator(intervalType)
	baseTime, err := calc.ParseSegmentTime(segmentName)
	if err != nil {
		return nil, fmt.Errorf("parse segment[%s] base time error", path)
	}
	kvStore, err := kv.NewStore(segmentName, kv.DefaultStoreOption(path))
	if err != nil {
		return nil, fmt.Errorf("create  kv store for segment error:%s", err)
	}

	return &segment{
		baseTime:     baseTime,
		kvStore:      kvStore,
		interval:     intervalVal,
		intervalType: intervalType,
		calc:         calc,
		logger:       logger.GetLogger("tsdb", "Segment"),
	}, nil
}

// BaseTime returns segment base time
func (s *segment) BaseTime() int64 {
	return s.baseTime
}

// GetDataFamilies returns data family list by time range, return nil if not match
func (s *segment) getDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	familyQueryTimeRange := timeutil.TimeRange{
		Start: s.calc.CalcFamilyStartTime(s.baseTime, s.calc.CalcFamily(timeRange.Start, s.baseTime)),
		End:   s.calc.CalcFamilyStartTime(s.baseTime, s.calc.CalcFamily(timeRange.End, s.baseTime)),
	}
	s.families.Range(func(k, v interface{}) bool {
		family, ok := v.(DataFamily)
		if ok {
			timeRange := family.TimeRange()
			if familyQueryTimeRange.Overlap(&timeRange) {
				result = append(result, family)
			}
		}
		return true
	})
	return result
}

// GetDataFamily returns the data family based on timestamp
func (s *segment) GetDataFamily(timestamp int64) (DataFamily, error) {
	segmentTime := s.calc.CalcSegmentTime(timestamp)
	if segmentTime != s.baseTime {
		return nil, fmt.Errorf("segment base time not match")
	}

	familyTime := s.calc.CalcFamily(timestamp, s.baseTime)
	family, ok := s.families.Load(familyTime)
	if !ok {
		// double check
		s.mutex.Lock()
		defer s.mutex.Unlock()
		family, ok = s.families.Load(familyTime)
		if !ok {
			//FIXME codingcrush need impl merger
			familyOption := kv.FamilyOption{
				CompactThreshold: 0,
				Merger:           "mock_merger",
			}
			// create kv family
			f, err := s.kvStore.CreateFamily(fmt.Sprintf("%d", familyTime), familyOption)
			if err != nil {
				return nil, fmt.Errorf("create data family error:%s", err)
			}
			// create data family
			familyStartTime := s.calc.CalcFamilyStartTime(s.baseTime, familyTime)
			dataFamily := newDataFamily(s.interval, timeutil.TimeRange{
				Start: familyStartTime,
				End:   s.calc.CalcFamilyEndTime(familyStartTime),
			}, f)
			s.families.Store(familyTime, dataFamily)
			return dataFamily, nil
		}
	}
	f, ok := family.(DataFamily)
	if !ok {
		return nil, series.ErrNotFound
	}
	return f, nil
}

// Close closes segment, include kv store
func (s *segment) Close() {
	if err := s.kvStore.Close(); err != nil {
		s.logger.Error("close kv store error", logger.Error(err))
	}
}
