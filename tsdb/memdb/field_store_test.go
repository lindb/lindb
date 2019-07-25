package memdb

import (
	"fmt"
	"testing"

	"github.com/eleme/lindb/pkg/field"
	pb "github.com/eleme/lindb/rpc/proto/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newFieldStore(t *testing.T) {
	fStore := newFieldStore(10, field.SumField)
	assert.NotNil(t, fStore)
	assert.Equal(t, fStore.getFieldType(), field.SumField)
	timeRange, ok := fStore.timeRange(10)
	assert.False(t, ok)
	assert.Equal(t, int64(0), timeRange.Start)
	assert.Equal(t, int64(0), timeRange.End)
}

func Test_fStore_write(t *testing.T) {
	fStore := newFieldStore(10, field.SumField)
	theFieldStore := fStore.(*fieldStore)
	writeCtx := writeContext{familyTime: 15, blockStore: newBlockStore(30)}

	// unknown field
	theFieldStore.write(&pb.Field{Name: "unknown"}, writeCtx)
	// sum field
	theFieldStore.write(&pb.Field{Name: "sum", Field: &pb.Field_Sum{}}, writeCtx)
	theFieldStore.write(&pb.Field{Name: "sum", Field: &pb.Field_Sum{}}, writeCtx)
}

func Test_fStore_timeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore(10, field.SumField)
	theFieldStore := fStore.(*fieldStore)

	mockSStore1 := NewMocksStoreINTF(ctrl)
	mockSStore1.EXPECT().slotRange().Return(1, 10, nil).AnyTimes()
	mockSStore2 := NewMocksStoreINTF(ctrl)
	mockSStore2.EXPECT().slotRange().Return(3, 5, nil).AnyTimes()
	mockSStore3 := NewMocksStoreINTF(ctrl)
	mockSStore3.EXPECT().slotRange().Return(6, 13, fmt.Errorf("error")).AnyTimes()
	mockSStore4 := NewMocksStoreINTF(ctrl)
	mockSStore4.EXPECT().slotRange().Return(4, 14, nil).AnyTimes()

	// error case
	theFieldStore.segments[1564297200000] = mockSStore3
	timeRange, ok := theFieldStore.timeRange(10 * 1000)
	assert.Equal(t, int64(0), timeRange.Start)
	assert.Equal(t, int64(0), timeRange.End)
	assert.False(t, ok)

	theFieldStore.segments[1564300800000] = mockSStore1
	timeRange, ok = theFieldStore.timeRange(10 * 1000)
	assert.Equal(t, int64(1564300810000), timeRange.Start)
	assert.Equal(t, int64(1564300900000), timeRange.End)
	assert.True(t, ok)

	theFieldStore.segments[1564304400000] = mockSStore2
	theFieldStore.segments[1564308000000] = mockSStore4
	timeRange, ok = theFieldStore.timeRange(10 * 1000)
	assert.Equal(t, int64(1564300810000), timeRange.Start)
	assert.Equal(t, int64(1564308140000), timeRange.End)
	assert.True(t, ok)
}

func Test_fStore_flushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore(10, field.SumField)
	theFieldStore := fStore.(*fieldStore)

	mockTF := makeMockTableFlusher(ctrl)
	mockSStore1 := NewMocksStoreINTF(ctrl)
	mockSStore1.EXPECT().bytes().Return(nil, 0, 0, fmt.Errorf("error"))
	mockSStore2 := NewMocksStoreINTF(ctrl)
	mockSStore2.EXPECT().bytes().Return(nil, 1, 212, nil)

	theFieldStore.segments[1564304400000] = mockSStore1
	theFieldStore.segments[1564308000000] = mockSStore2

	assert.Len(t, theFieldStore.segments, 2)
	// familyTime not exist
	assert.False(t, theFieldStore.flushFieldTo(mockTF, 1564297200000))
	assert.Len(t, theFieldStore.segments, 2)
	// mock error
	assert.False(t, theFieldStore.flushFieldTo(mockTF, 1564304400000))
	assert.Len(t, theFieldStore.segments, 1)
	// mock ok
	assert.True(t, theFieldStore.flushFieldTo(mockTF, 1564308000000))
	assert.Len(t, theFieldStore.segments, 0)
}
