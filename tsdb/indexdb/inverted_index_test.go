package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/tsdb/metadb"
)

func TestInvertedIndex_buildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)

	// case 1: get series ids by tag value ids
	seriesIDs, err := index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(2, roaring.BitmapOf(2))
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(2), seriesIDs)

	// case 2: tag key is not exist
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(4, roaring.BitmapOf(1))
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)

	// case 3: tag value ids is not exist
	seriesIDs, err = index.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(10, 20))
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), seriesIDs)
	// case 4: tag key not exist
	seriesIDs, err = index.GetSeriesIDsForTag(4)
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)
	// case 5: get all series ids under tag key
	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)

}

func TestInvertedIndex_GetGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)
	ctx, err := index.GetGroupingContext([]uint32{3, 4})
	assert.Error(t, err)
	assert.Nil(t, ctx)

	ctx, err = index.GetGroupingContext([]uint32{1, 2})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	ctx, err = index.GetGroupingContext([]uint32{1})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
}

func TestInvertedIndex_FlushInvertedIndexTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)
	assert.Panics(t, func() {
		_ = index.FlushInvertedIndexTo(nil)
	})
}

func prepareInvertedIndex(ctrl *gomock.Controller) InvertedIndex {
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	tagMetadata := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMetadata).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(1), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone").Return(uint32(2), nil).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), "zone_err").Return(uint32(0), fmt.Errorf("err")).AnyTimes()
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.1").Return(uint32(1), nil).Times(2)
	tagMetadata.EXPECT().GenTagValueID(uint32(1), "1.1.1.5").Return(uint32(0), fmt.Errorf("err"))
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "sh").Return(uint32(1), nil)
	tagMetadata.EXPECT().GenTagValueID(uint32(2), "bj").Return(uint32(2), nil)
	index := newInvertedIndex(metadata)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host": "1.1.1.1",
		"zone": "sh",
	}, 1)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host": "1.1.1.1",
		"zone": "bj",
	}, 2)
	index.buildInvertIndex("ns", "name", map[string]string{
		"host":     "1.1.1.5",
		"zone_err": "bj",
	}, 3)
	return index
}
