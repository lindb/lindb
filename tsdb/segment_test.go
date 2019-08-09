package tsdb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
)

var segPath = filepath.Join(testPath, shardPath, "1", segmentPath, string(interval.Day))

func TestNewIntervalSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, err := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.True(t, fileutil.Exist(segPath))
}

func TestNewSegment(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	s, _ := newIntervalSegment(int64(time.Second*10), interval.Day, segPath)

	seg, err := s.GetOrCreateSegment("20190702")
	assert.Nil(t, err)
	assert.NotNil(t, seg)
	assert.True(t, fileutil.Exist(filepath.Join(segPath, "20190702")))

	s.Close()

	s, _ = newIntervalSegment(int64(time.Second*10), interval.Day, segPath)

	seg1, ok := s.(*intervalSegment)
	if ok {
		seg = seg1.getSegment("20190702")
		assert.Nil(t, err)
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

}
