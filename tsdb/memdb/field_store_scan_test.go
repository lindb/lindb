package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
)

type mockScanWorker struct {
	events []series.ScanEvent
}

func (w *mockScanWorker) Emit(event series.ScanEvent) {
	w.events = append(w.events, event)
}
func (w *mockScanWorker) Close() {}

func TestFieldStore_Scan(t *testing.T) {
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
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	gomock.InOrder(
		agg.EXPECT().GetAggregator(familyTime).Return(fieldAgg, true),
		fieldAgg.EXPECT().GetAllAggregates().Return([]aggregation.PrimitiveAggregator{pAgg}),
		pAgg.EXPECT().Aggregate(20, 1.0).Return(false),
	)
	fStore.scan(agg, sCtx)
}
