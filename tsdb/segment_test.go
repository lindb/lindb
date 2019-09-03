package tsdb

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/series"
)

var segPath = filepath.Join(testPath, shardPath, "1", segmentPath, string(interval.Day))

func TestNewIntervalSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, err := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s)
	assert.True(t, fileutil.Exist(segPath))
	s.Close()

	// create fail
	_ = fileutil.MkDirIfNotExist(filepath.Join(segPath, "20190903"))
	s, err = newIntervalSegment(int64(time.Second*10), interval.Unknown, segPath)
	assert.Nil(t, s)
	assert.NotNil(t, err)
}

func TestIntervalSegment_GetOrCreateSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	seg, err := s.GetOrCreateSegment("20190702")
	assert.Nil(t, err)
	assert.NotNil(t, seg)
	assert.True(t, fileutil.Exist(filepath.Join(segPath, "20190702")))

	seg1, err1 := s.GetOrCreateSegment("20190702")
	if err1 != nil {
		t.Fatal(err1)
	}
	assert.Equal(t, seg, seg1)

	// test create fail
	seg, err = s.GetOrCreateSegment("201907-a")
	assert.Nil(t, seg)
	assert.NotNil(t, err)

	s.Close()

	s, _ = newIntervalSegment(int64(time.Second*10), interval.Day, segPath)

	s1, ok := s.(*intervalSegment)
	if ok {
		seg = s1.getSegment("20190702")
		assert.NotNil(t, seg)
		assert.True(t, fileutil.Exist(filepath.Join(segPath, "20190702")))
	} else {
		t.Fail()
	}
}

func TestGetSegmentsByTimeRange(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	_, _ = s.GetOrCreateSegment("20190702")
	t2, _ := timeutil.ParseTimestamp("20190702", "20060102")
	segments := s.GetSegments(timeutil.TimeRange{Start: t2, End: t2 + 60*60*1000})
	assert.Equal(t, 1, len(segments))

	segments = s.GetSegments(timeutil.TimeRange{Start: t2 + 50*1000, End: t2 + 60*60*1000})
	assert.Equal(t, 1, len(segments))

	t2, _ = timeutil.ParseTimestamp("20190701", "20060102")
	segments = s.GetSegments(timeutil.TimeRange{Start: t2, End: t2 + 25*60*60*1000})
	assert.Equal(t, 1, len(segments))

	seg := s.(*intervalSegment)
	seg.intervalType = interval.Unknown
	t2, _ = timeutil.ParseTimestamp("20190701", "20060102")
	segments = s.GetSegments(timeutil.TimeRange{Start: t2, End: t2 + 25*60*60*1000})
	assert.Nil(t, segments)
}

func TestSegment_Close(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	seg, _ := s.GetOrCreateSegment("20190702")
	seg1 := seg.(*segment)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := kv.NewMockStore(ctrl)
	seg1.kvStore = store
	store.EXPECT().Close().Return(fmt.Errorf("err"))
	seg.Close()
}

func TestSegment_GetDataFamily(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	seg, _ := s.GetOrCreateSegment("20190904")
	now, _ := timeutil.ParseTimestamp("20190904 19:10:48", "20060102 15:04:05")
	familyBaseTime, _ := timeutil.ParseTimestamp("20190904 19:00:00", "20060102 15:04:05")
	assert.NotNil(t, seg)
	dataFamily, err := seg.GetDataFamily(now)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, familyBaseTime, dataFamily.BaseTime())
	dataFamily1, _ := seg.GetDataFamily(now)
	assert.Equal(t, dataFamily, dataFamily1)

	// wrong data family type
	wrongTime, _ := timeutil.ParseTimestamp("20190904 23:10:48", "20060102 15:04:05")
	seg1 := seg.(*segment)
	seg1.families.Store(23, "err data family")
	result, err := seg.GetDataFamily(wrongTime)
	assert.Equal(t, series.ErrNotFound, err)
	assert.Nil(t, result)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := kv.NewMockStore(ctrl)
	seg1.kvStore = store
	wrongTime, _ = timeutil.ParseTimestamp("20190904 11:10:48", "20060102 15:04:05")
	store.EXPECT().CreateFamily("11", gomock.Any()).Return(nil, fmt.Errorf("err"))
	dataFamily, err = seg.GetDataFamily(wrongTime)
	assert.NotNil(t, err)
	assert.Nil(t, dataFamily)
}

func TestSegment_New(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, err := newSegment("20190904", interval.Day, testPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s)
	s, err = newSegment("20190904", interval.Day, testPath)
	assert.NotNil(t, err)
	assert.Nil(t, s)
}
