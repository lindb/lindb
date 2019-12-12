package tsdb

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./segment.go -destination=./segment_mock.go -package=tsdb

// for testing
var (
	newStore = kv.NewStore
)

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
	baseTime int64
	kvStore  kv.Store
	interval timeutil.Interval
	families sync.Map

	mutex sync.Mutex

	logger *logger.Logger
}

// newSegment returns segment, segment is wrapper of kv store
func newSegment(
	segmentName string,
	interval timeutil.Interval,
	path string,
) (
	Segment,
	error,
) {
	// parse base time from segment name
	calc := interval.Calculator()
	baseTime, err := calc.ParseSegmentTime(segmentName)
	if err != nil {
		return nil, fmt.Errorf("parse segment[%s] base time error", path)
	}
	kvStore, err := newStore(segmentName, kv.DefaultStoreOption(path))
	if err != nil {
		return nil, fmt.Errorf("create kv store for segment error:%s", err)
	}
	familyNames := kvStore.ListFamilyNames()
	s := &segment{
		baseTime: baseTime,
		kvStore:  kvStore,
		interval: interval,
		logger:   logger.GetLogger("tsdb", "Segment"),
	}
	for _, familyName := range familyNames {
		familyTime, err := strconv.Atoi(familyName)
		if err != nil {
			return nil, fmt.Errorf("load data family error:%s", err)
		}
		_ = s.initDataFamily(familyTime, kvStore.GetFamily(familyName))
	}
	return s, nil
}

// BaseTime returns segment base time
func (s *segment) BaseTime() int64 {
	return s.baseTime
}

// GetDataFamilies returns data family list by time range, return nil if not match
func (s *segment) getDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	calc := s.interval.Calculator()

	familyQueryTimeRange := timeutil.TimeRange{
		Start: calc.CalcFamilyStartTime(s.baseTime, calc.CalcFamily(timeRange.Start, s.baseTime)),
		End:   calc.CalcFamilyStartTime(s.baseTime, calc.CalcFamily(timeRange.End, s.baseTime)),
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
	calc := s.interval.Calculator()

	segmentTime := calc.CalcSegmentTime(timestamp)
	if segmentTime != s.baseTime {
		return nil, fmt.Errorf("segment base time not match")
	}

	familyTime := calc.CalcFamily(timestamp, s.baseTime)
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
				Merger:           nopMerger,
			}
			// create kv family
			f, err := s.kvStore.CreateFamily(fmt.Sprintf("%d", familyTime), familyOption)
			if err != nil {
				return nil, fmt.Errorf("create data family error: %s", err)
			}
			dataFamily := s.initDataFamily(familyTime, f)
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

func (s *segment) initDataFamily(familyTime int, family kv.Family) DataFamily {
	calc := s.interval.Calculator()
	// create data family
	familyStartTime := calc.CalcFamilyStartTime(s.baseTime, familyTime)
	dataFamily := newDataFamily(s.interval, timeutil.TimeRange{
		Start: familyStartTime,
		End:   calc.CalcFamilyEndTime(familyStartTime),
	}, family)
	s.families.Store(familyTime, dataFamily)
	return dataFamily
}
