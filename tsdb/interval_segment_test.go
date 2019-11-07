package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestNewIntervalSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, err := newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s)
	assert.True(t, fileutil.Exist(segPath))
	s.Close()

	// create fail
	_, err = newSegment(
		"20190903",
		timeutil.Interval(timeutil.OneSecond*10),
		filepath.Join(segPath, "20190903"))
	if err != nil {
		t.Fatal(err)
	}
	// cannot re-open kv-store
	s, err = newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
	assert.Nil(t, s)
	assert.NotNil(t, err)

}

func TestIntervalSegment_GetOrCreateSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
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

	s, _ = newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)

	s1, ok := s.(*intervalSegment)
	if ok {
		seg, ok = s1.getSegment("20190702")
		assert.NotNil(t, seg)
		assert.True(t, ok)
		assert.True(t, fileutil.Exist(filepath.Join(segPath, "20190702")))
	} else {
		t.Fail()
	}
}

func TestIntervalSegment_getDataFamilies(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(timeutil.Interval(timeutil.OneSecond*10), segPath)
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
