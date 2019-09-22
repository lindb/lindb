package aggregation

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestGroupByAggregator_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now, _ := timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	agg := NewGroupByAggregator(timeutil.OneSecond, &timeutil.TimeRange{
		Start: now,
		End:   now + 3*timeutil.OneHour,
	}, true,
		AggregatorSpecs{
			NewAggregatorSpec("b", field.SumField),
			NewAggregatorSpec("a", field.SumField),
		})
	agg.ResultSet()
}
