package query

import (
	"math"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

func TestAggregate_Sum(t *testing.T) {
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

func TestAggregate_Summary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	sumAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	countAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	fAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{sumAgg, countAgg})

	fieldWriter := stream.NewBufferWriter(nil)
	fieldIDs := collections.NewBitArray(nil)
	offsets := encoding.NewFixedOffsetEncoder()

	tsd := encoding.GetTSDDecoder()
	sum := encoding.NewTSDEncoder(10)
	sum.AppendTime(bit.One)
	sum.AppendValue(math.Float64bits(10.0))
	data, _ := sum.Bytes()
	offset := fieldWriter.Len()
	offsets.Add(offset)
	fieldWriter.PutBytes(data)
	fieldIDs.SetBit(uint16(0))

	count := encoding.NewTSDEncoder(10)
	count.AppendTime(bit.One)
	count.AppendValue(encoding.ZigZagEncode(5))
	data, _ = count.Bytes()
	offset = fieldWriter.Len()
	offsets.Add(offset)
	fieldWriter.PutBytes(data)
	fieldIDs.SetBit(uint16(1))

	writer := stream.NewBufferWriter(nil)
	writer.PutBytes(fieldIDs.Bytes())
	writer.PutBytes(offsets.MarshalBinary())
	d, _ := fieldWriter.Bytes()
	writer.PutBytes(d)
	fieldData, _ := writer.Bytes()

	sumAgg.EXPECT().Aggregate(10, 10.0)  // agg value
	countAgg.EXPECT().Aggregate(10, 5.0) // agg value
	Aggregate(field.SummaryField, fAgg, tsd, fieldData)
}
