package indexdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestInvertedIndex_buildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)

	seriesIDs, err := index.FindSeriesIDsByExpr(1, &stmt.EqualsExpr{Key: "host", Value: "1.1.1.1"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)

	seriesIDs, err = index.FindSeriesIDsByExpr(2, &stmt.EqualsExpr{Key: "zone", Value: "bj"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(2), seriesIDs)

	seriesIDs, err = index.FindSeriesIDsByExpr(4, &stmt.EqualsExpr{Key: "zone", Value: "bj"})
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)

	seriesIDs, err = index.GetSeriesIDsForTag(4)
	assert.Error(t, err)
	assert.Nil(t, seriesIDs)

	seriesIDs, err = index.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
}

func TestInvertedIndex_SuggestTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	index := prepareInvertedIndex(ctrl)
	assert.Nil(t, index.SuggestTagValues(5, "", 10))
	assert.Len(t, index.SuggestTagValues(2, "sh", 10), 1)
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

func prepareInvertedIndex(ctrl *gomock.Controller) InvertedIndex {
	generator := metadb.NewMockIDGenerator(ctrl)
	generator.EXPECT().GenTagKeyID(gomock.Any(), "host").Return(uint32(1)).AnyTimes()
	generator.EXPECT().GenTagKeyID(gomock.Any(), "zone").Return(uint32(2)).AnyTimes()
	index := newInvertedIndex(generator)
	index.buildInvertIndex(10, map[string]string{
		"host": "1.1.1.1",
		"zone": "sh",
	}, 1)
	index.buildInvertIndex(10, map[string]string{
		"host": "1.1.1.1",
		"zone": "bj",
	}, 2)
	return index
}
