package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestGroupByAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gIt := series.NewMockGroupedIterator(ctrl)
	sIt := series.NewMockIterator(ctrl)
	fIt := series.NewMockFieldIterator(ctrl)

	now, _ := timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	agg := NewGroupingAggregator(timeutil.OneSecond, &timeutil.TimeRange{
		Start: now,
		End:   now + 3*timeutil.OneHour,
	},
		AggregatorSpecs{
			NewAggregatorSpec("b", field.SumField),
			NewAggregatorSpec("a", field.SumField),
		})

	gomock.InOrder(
		gIt.EXPECT().Tags().Return(map[string]string{"host": "1.1.1.1"}),
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(sIt),
		// series it
		sIt.EXPECT().FieldName().Return("a"),
		sIt.EXPECT().HasNext().Return(true),
		sIt.EXPECT().Next().Return(familyTime, fIt),
		fIt.EXPECT().HasNext().Return(false),
		// series it
		sIt.EXPECT().HasNext().Return(true),
		sIt.EXPECT().Next().Return(familyTime, nil),
		sIt.EXPECT().HasNext().Return(false),
		// series it
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(sIt),
		sIt.EXPECT().FieldName().Return("c"),

		gIt.EXPECT().HasNext().Return(false),
	)
	agg.Aggregate(gIt)
	rs := agg.ResultSet()
	assert.Equal(t, 1, len(rs))

	gomock.InOrder(
		gIt.EXPECT().Tags().Return(map[string]string{"host": "1.1.1.2"}),
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(sIt),
		// series it
		sIt.EXPECT().FieldName().Return("a"),
		sIt.EXPECT().HasNext().Return(true),
		sIt.EXPECT().Next().Return(familyTime, fIt),
		fIt.EXPECT().HasNext().Return(false),
		// series it
		sIt.EXPECT().HasNext().Return(true),
		sIt.EXPECT().Next().Return(familyTime, nil),
		sIt.EXPECT().HasNext().Return(false),
		// series it
		gIt.EXPECT().HasNext().Return(true),
		gIt.EXPECT().Next().Return(sIt),
		sIt.EXPECT().FieldName().Return("c"),

		gIt.EXPECT().HasNext().Return(false),
	)
	agg.Aggregate(gIt)

	rs = agg.ResultSet()
	assert.Equal(t, 2, len(rs))

	agg = NewGroupingAggregator(timeutil.OneSecond, &timeutil.TimeRange{
		Start: now,
		End:   now + 3*timeutil.OneHour,
	},
		AggregatorSpecs{})
	rs = agg.ResultSet()
	assert.Nil(t, rs)

}
