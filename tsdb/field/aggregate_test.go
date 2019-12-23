package field

import (
	"math"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

func TestAggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Aggregate(field.Unknown, nil, nil, nil)

	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	fAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg})

	tsd := encoding.GetTSDDecoder()
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.Bytes()
	pAgg.EXPECT().Aggregate(10, 10.0) // agg value
	Aggregate(field.SumField, fAgg, tsd, data)
}
