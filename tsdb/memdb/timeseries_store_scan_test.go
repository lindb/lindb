package memdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

func TestTimeSeriesStore_scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	mCtx := &memScanContext{
		fieldIDs:   []uint16{1, 2, 3},
		fieldCount: 3,
	}
	// no data
	tStore.scan(mCtx)

	// mock fieldID getter
	mockFieldIDGetter := NewMockmStoreFieldIDGetter(ctrl)
	for i := 0; i < 10; i++ {
		mockFieldIDGetter.EXPECT().GetFieldIDOrGenerate(gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Return(uint16(i+10), nil)
		_, err := tStore.Write(
			&pb.Metric{
				Fields: []*pb.Field{
					{Name: fmt.Sprintf("sum-%d", i), Field: &pb.Field_Sum{Sum: &pb.Sum{
						Value: 1.0,
					}}},
					{Name: "unknown", Field: nil}},
			}, writeContext{
				metricID:            1,
				familyTime:          10,
				blockStore:          newBlockStore(30),
				mStoreFieldIDGetter: mockFieldIDGetter})
		assert.NoError(t, err)
	}
	// find data
	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	fieldsAgg := aggregation.FieldAggregates{sAgg, sAgg}
	gomock.InOrder(
		sAgg.EXPECT().GetAggregator(int64(10)).Return(nil, false).MaxTimes(2),
	)
	mCtx = &memScanContext{
		fieldIDs:    []uint16{12, 13},
		aggregators: fieldsAgg,
		fieldCount:  2,
	}
	tStore.scan(mCtx)

	// not match field
	mCtx = &memScanContext{
		fieldIDs:    []uint16{2, 3},
		aggregators: fieldsAgg,
		fieldCount:  2,
	}
	tStore.scan(mCtx)
}
