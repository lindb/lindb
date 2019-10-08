package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
)

func TestGroupedIterator_HasNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sIt1 := NewMockSeriesAggregator(ctrl)
	sIt2 := NewMockSeriesAggregator(ctrl)
	fIt := series.NewMockIterator(ctrl)
	tags := map[string]string{"host": "1.1.1.1"}
	it := newGroupedIterator(tags, FieldAggregates{sIt1, sIt2})
	gomock.InOrder(
		sIt1.EXPECT().ResultSet().Return(fIt),
		sIt2.EXPECT().ResultSet().Return(fIt),
	)
	assert.Equal(t, tags, it.Tags())
	assert.True(t, it.HasNext())
	assert.Equal(t, fIt, it.Next())
	assert.True(t, it.HasNext())
	assert.Equal(t, fIt, it.Next())
	assert.False(t, it.HasNext())
}
