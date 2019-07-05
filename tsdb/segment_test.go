package tsdb

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/util"
)

var segPath = filepath.Join(testPath, shardPath, "1", segmentPath, interval.Day.String())

func TestNewIntervalSegment(t *testing.T) {
	defer util.RemoveDir(testPath)
	s, err := newIntervalSegment(time.Second*10, interval.Day, segPath)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.True(t, util.Exist(segPath))
}

func TestNewSegment(t *testing.T) {
	defer util.RemoveDir(testPath)
	s, _ := newIntervalSegment(time.Second*10, interval.Day, segPath)

	seg, err := s.GetOrCreateSegment("20190702")
	assert.Nil(t, err)
	assert.NotNil(t, seg)
	assert.True(t, util.Exist(filepath.Join(segPath, "20190702")))

	s.Close()

	s, _ = newIntervalSegment(time.Second*10, interval.Day, segPath)

	seg1, ok := s.(*intervalSegment)
	if ok {
		seg = seg1.getSegment("20190702")
		assert.Nil(t, err)
		assert.NotNil(t, seg)
		assert.True(t, util.Exist(filepath.Join(segPath, "20190702")))
	} else {
		t.Fail()
	}
}

func TestGetSegmentsByTimeRange(t *testing.T) {
	defer util.RemoveDir(testPath)
	s, _ := newIntervalSegment(time.Second*10, interval.Day, segPath)
	s.GetOrCreateSegment("20190705")
	t2, _ := util.ParseTimestamp("20190705", "20060102")
	segments := s.GetSegments(models.TimeRange{Start: t2, End: t2 + 60*60*1000})
	fmt.Println(len(segments))
	assert.Equal(t, 1, len(segments))

	segments = s.GetSegments(models.TimeRange{Start: t2 + 50*1000, End: t2 + 60*60*1000})
	assert.Equal(t, 1, len(segments))
}
