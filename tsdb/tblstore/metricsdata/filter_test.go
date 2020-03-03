package metricsdata

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series/field"
)

func TestFileFilterResultSet_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl)

	rs := newFileFilterResultSet(1, field.Metas{}, nil, reader)
	reader.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	rs.Load(nil, []field.ID{1}, 0, nil)
}

func TestMetricsDataFilter_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockReader(ctrl)
	filter := NewFilter(10, nil, []Reader{reader})
	// case 1: field not found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	rs, err := filter.Filter([]field.ID{1, 30}, roaring.BitmapOf(1, 2, 3))
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 2: series ids found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	reader.EXPECT().GetSeriesIDs().Return(roaring.BitmapOf(10, 200))
	rs, err = filter.Filter([]field.ID{2, 30}, roaring.BitmapOf(1, 2, 3))
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 3: data found
	reader.EXPECT().GetFields().Return(field.Metas{{ID: 2}, {ID: 20}})
	reader.EXPECT().GetSeriesIDs().Return(roaring.BitmapOf(10, 200))
	rs, err = filter.Filter([]field.ID{2, 30}, roaring.BitmapOf(1, 200, 3))
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.EqualValues(t, roaring.BitmapOf(200).ToArray(), rs[0].SeriesIDs().ToArray())
}
