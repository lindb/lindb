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
	"github.com/lindb/lindb/series"
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
	_, err = newSegment("20190903", int64(10000), interval.Day, filepath.Join(segPath, "20190903"))
	if err != nil {
		t.Fatal(err)
	}
	// cannot re-open kv-store
	s, err = newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	assert.Nil(t, s)
	assert.NotNil(t, err)

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

func TestIntervalSegment_getDataFamilies(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	segment1, _ := s.GetOrCreateSegment("20190902")
	now, _ := timeutil.ParseTimestamp("20190902 19:10:48", "20060102 15:04:05")
	_, _ = segment1.GetDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190902 20:10:48", "20060102 15:04:05")
	_, _ = segment1.GetDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190902 22:10:48", "20060102 15:04:05")
	_, _ = segment1.GetDataFamily(now)
	segment2, _ := s.GetOrCreateSegment("20190904")
	now, _ = timeutil.ParseTimestamp("20190904 22:10:48", "20060102 15:04:05")
	_, _ = segment2.GetDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190904 20:10:48", "20060102 15:04:05")
	_, _ = segment2.GetDataFamily(now)

	start, _ := timeutil.ParseTimestamp("20190901 20:10:48", "20060102 15:04:05")
	end, _ := timeutil.ParseTimestamp("20190901 22:10:48", "20060102 15:04:05")
	segments := s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 0, len(segments))

	start, _ = timeutil.ParseTimestamp("20190905 20:10:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190905 22:10:48", "20060102 15:04:05")
	segments = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 0, len(segments))

	start, _ = timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190905 22:10:48", "20060102 15:04:05")
	segments = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 5, len(segments))

	start, _ = timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190902 20:40:48", "20060102 15:04:05")
	segments = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 2, len(segments))

	start, _ = timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190904 19:40:48", "20060102 15:04:05")
	segments = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 3, len(segments))

	start, _ = timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190902 19:40:48", "20060102 15:04:05")
	segments = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 1, len(segments))
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
	familyEndTime, _ := timeutil.ParseTimestamp("20190904 20:00:00", "20060102 15:04:05")
	assert.Equal(t, timeutil.TimeRange{
		Start: familyBaseTime,
		End:   familyEndTime - 1,
	}, dataFamily.TimeRange())
	dataFamily1, _ := seg.GetDataFamily(now)
	assert.Equal(t, dataFamily, dataFamily1)

	// segment not match
	now, _ = timeutil.ParseTimestamp("20190903 19:10:48", "20060102 15:04:05")
	dataFamily, err = seg.GetDataFamily(now)
	assert.Nil(t, dataFamily)
	assert.NotNil(t, err)
	now, _ = timeutil.ParseTimestamp("20190905 19:10:48", "20060102 15:04:05")
	dataFamily, err = seg.GetDataFamily(now)
	assert.Nil(t, dataFamily)
	assert.NotNil(t, err)

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
	s, err := newSegment("20190904", int64(10000), interval.Day, testPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s)
	s, err = newSegment("20190904", int64(10000), interval.Day, testPath)
	assert.NotNil(t, err)
	assert.Nil(t, s)
}
