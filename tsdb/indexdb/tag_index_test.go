package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

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
