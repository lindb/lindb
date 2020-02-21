package indexdb

import (
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

func TestTagIndex_GetGroupingScanner(t *testing.T) {
	index := prepareTagIdx()
	// case 1: series ids not match
	scanners, err := index.GetGroupingScanner(roaring.BitmapOf(1000, 2000))
	assert.NoError(t, err)
	assert.Nil(t, scanners)
	// case 2: get scanner
	scanners, err = index.GetGroupingScanner(roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Len(t, scanners, 1)
	container, tagValueIDs := scanners[0].GetSeriesAndTagValue(0)
	assert.Equal(t, 8, container.GetCardinality())
	assert.Len(t, tagValueIDs, 8)
	container, tagValueIDs = scanners[0].GetSeriesAndTagValue(1)
	assert.Nil(t, container)
	assert.Nil(t, tagValueIDs)
}

func TestTagIndex_buildInvertedIndex(t *testing.T) {
	index := newTagIndex()
	index.buildInvertedIndex(2, 1)
	index.buildInvertedIndex(2, 3)
	index.buildInvertedIndex(2, 2)
	index.buildInvertedIndex(1, 1)
	index.buildInvertedIndex(1, 2)
	values := index.getValues()
	seriesIDs, ok := values.Get(1)
	assert.True(t, ok)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	seriesIDs, ok = values.Get(1)
	assert.True(t, ok)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), index.getAllSeriesIDs())
}

func TestTagIndex_getSeriesIDsByTagValueIDs(t *testing.T) {
	tagIndex := prepareTagIdx()
	// tag-value not exist
	assert.Equal(t, roaring.New(), tagIndex.getSeriesIDsByTagValueIDs(roaring.BitmapOf(40, 50, 30)))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(4), tagIndex.getSeriesIDsByTagValueIDs(roaring.BitmapOf(4)))
}

func TestTagIndex_getAllSeriesIDs(t *testing.T) {
	tagIndex := prepareTagIdx()
	assert.Equal(t, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8), tagIndex.getAllSeriesIDs())
}

func TestTagIndex_flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIndex := prepareTagIdx()
	forward := invertedindex.NewMockForwardFlusher(ctrl)
	inverted := invertedindex.NewMockInvertedFlusher(ctrl)
	// case 1: flush forward err
	forward.EXPECT().FlushForwardIndex(gomock.Any()).AnyTimes()
	forward.EXPECT().FlushTagKeyID(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err := tagIndex.flush(12, forward, inverted)
	assert.Error(t, err)
	// case 2: flush tag level series ids err
	forward.EXPECT().FlushTagKeyID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	inverted.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = tagIndex.flush(12, forward, inverted)
	assert.Error(t, err)
	// case 3: flush tag value series ids
	inverted.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(nil)
	inverted.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = tagIndex.flush(13, forward, inverted)
	assert.Error(t, err)
	// case 4: flush tag key err
	inverted.EXPECT().FlushInvertedIndex(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	inverted.EXPECT().FlushTagKeyID(gomock.Any()).Return(fmt.Errorf("err"))
	err = tagIndex.flush(14, forward, inverted)
	assert.Error(t, err)
	// case 3: flush tag key err
	inverted.EXPECT().FlushTagKeyID(gomock.Any()).Return(nil)
	err = tagIndex.flush(14, forward, inverted)
	assert.NoError(t, err)
}

func prepareTagIdx() TagIndex {
	tagIndex := newTagIndex()
	tagIndex.buildInvertedIndex(1, 1)
	tagIndex.buildInvertedIndex(2, 2)
	tagIndex.buildInvertedIndex(3, 3)
	tagIndex.buildInvertedIndex(4, 4)
	tagIndex.buildInvertedIndex(5, 5)
	tagIndex.buildInvertedIndex(6, 6)
	tagIndex.buildInvertedIndex(7, 7)
	tagIndex.buildInvertedIndex(8, 8)
	return tagIndex
}

type groupingScanner struct {
	forward *ForwardStore
}

func (g *groupingScanner) GetSeriesAndTagValue(highKey uint16) (roaring.Container, []uint32) {
	idx := g.forward.Keys().GetContainerIndex(highKey)
	if idx == -1 {
		return nil, nil
	}
	return g.forward.Keys().GetContainerAtIndex(idx), g.forward.Values()[idx]
}

func newScanner() series.GroupingScanner {
	return &groupingScanner{
		forward: NewForwardStore(),
	}
}

func BenchmarkForwardStore_Grouping(b *testing.B) {
	hosts := newScanner()
	disks := newScanner()
	partitions := newScanner()
	h := hosts.(*groupingScanner)
	d := disks.(*groupingScanner)
	p := partitions.(*groupingScanner)
	id := uint32(1)
	count := uint32(40000)
	for i := uint32(1); i <= count; i++ {
		for j := uint32(1); j <= 4; j++ {
			for k := uint32(1); k <= 20; k++ {
				id++
				h.forward.Put(id, i)
				d.forward.Put(id, j)
				p.forward.Put(id, k)
			}
		}
	}

	fmt.Println(id)
	seriesIDs := roaring.New()
	seriesIDs.AddRange(0, uint64(1000000))
	keys := seriesIDs.GetHighKeys()
	// test single group by tag keys
	scanners := make(map[uint32][]series.GroupingScanner)
	scanners[1] = []series.GroupingScanner{partitions}
	scanners[2] = []series.GroupingScanner{hosts}
	ctx := query.NewGroupContext([]uint32{1, 2}, scanners)

	now := timeutil.Now()
	var wait sync.WaitGroup
	var c atomic.Int32
	for idx, key := range keys {
		container := seriesIDs.GetContainerAtIndex(idx)
		//i += container.GetCardinality()
		k := key
		wait.Add(1)
		go func() {
			rs := ctx.BuildGroup(k, container)
			c.Add(int32(len(rs)))
			wait.Done()
		}()
		//assert.Len(t, rs, 80)
	}
	wait.Wait()
	fmt.Println(c.Load())
	fmt.Println(timeutil.Now() - now)
}
