package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

func TestFieldStore_simple_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	agg := aggregation.NewMockSeriesAggregator(ctrl)

	bs := newBlockStore(10)

	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")

	fStore := newFieldStore(10)
	sCtx := &memScanContext{}
	// no data
	fStore.scan(agg, sCtx)

	// write data
	fStore.Write(
		field.SumField,
		&pb.Field{
			Name: "f1",
			Field: &pb.Field_Sum{Sum: &pb.Sum{
				Value: 1.0,
			}}},
		writeContext{
			blockStore: bs,
			familyTime: familyTime,
			slotIndex:  20,
			metricID:   uint32(10),
		})

	fieldAgg := aggregation.NewMockFieldAggregator(ctrl)
	agg.EXPECT().GetFieldType().Return(field.SumField)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	gomock.InOrder(
		agg.EXPECT().GetAggregator(familyTime).Return(fieldAgg, true),
		fieldAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg}),
		pAgg.EXPECT().Aggregate(20, 1.0).Return(false),
	)
	fStore.scan(agg, sCtx)
}
func TestFieldStore_complex_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	agg := aggregation.NewMockSeriesAggregator(ctrl)

	bs := newBlockStore(10)

	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")

	fStore := newFieldStore(10)
	sCtx := &memScanContext{}
	// no data
	fStore.scan(agg, sCtx)

	// write data
	fStore.Write(
		field.SummaryField,
		&pb.Field{
			Name: "f1",
			Field: &pb.Field_Summary{Summary: &pb.Summary{
				Sum:   10.0,
				Count: 2.0,
			}}},
		writeContext{
			blockStore: bs,
			familyTime: familyTime,
			slotIndex:  20,
			metricID:   uint32(10),
		})

	fieldAgg := aggregation.NewMockFieldAggregator(ctrl)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	gomock.InOrder(
		agg.EXPECT().GetAggregator(familyTime).Return(fieldAgg, true),
		fieldAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg}),
		agg.EXPECT().GetFieldType().Return(field.SummaryField),
		pAgg.EXPECT().FieldID().Return(uint16(2)),
		//pAgg.EXPECT().Aggregate(20, 2.0).Return(false),//FIXME stone1100
	)
	fStore.scan(agg, sCtx)
}
