package tsdb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

var segPath = filepath.Join(testPath, shardDir, "2", segmentDir, timeutil.Day.String())

func TestSegment_Close(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
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
	s, _ := newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
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
	s, err := newSegment("20190904", timeutil.Interval(timeutil.OneSecond*10), testPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s)
	s, err = newSegment("20190904", timeutil.Interval(timeutil.OneSecond*10), testPath)
	assert.NotNil(t, err)
	assert.Nil(t, s)
}
