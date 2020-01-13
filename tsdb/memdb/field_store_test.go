package memdb

import (
	"sort"
	"testing"

	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getMockSStore(ctrl *gomock.Controller, familyTime int64) *MocksStoreINTF {
	mockSStore := NewMocksStoreINTF(ctrl)
	mockSStore.EXPECT().GetFamilyTime().Return(familyTime).AnyTimes()
	mockSStore.EXPECT().MemSize().Return(emptySimpleFieldStoreSize).AnyTimes()
	return mockSStore
}

func Test_newFieldStore(t *testing.T) {
	fStore := newFieldStore(1)
	assert.NotNil(t, fStore)
	assert.Equal(t, uint16(1), fStore.GetFieldID())
	timeRange, ok := fStore.TimeRange(10)
	assert.False(t, ok)
	assert.Equal(t, int64(0), timeRange.Start)
	assert.Equal(t, int64(0), timeRange.End)
}

func Test_fStore_write(t *testing.T) {
	fStore := newFieldStore(10)
	theFieldStore := fStore.(*fieldStore)
	writeCtx := writeContext{familyTime: 15, blockStore: newBlockStore(30)}

	//unknown field
	theFieldStore.Write(field.Unknown, &pb.Field{Name: "unknown"}, writeCtx)
	// sum field
	theFieldStore.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 1.0,
		},
	}}, writeCtx)
}

func Test_fStore_timeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore(10)
	theFieldStore := fStore.(*fieldStore)

	mockSStore1 := getMockSStore(ctrl, 1564300800000)
	mockSStore1.EXPECT().SlotRange().Return(uint16(1), uint16(10)).AnyTimes()
	mockSStore2 := getMockSStore(ctrl, 1564304400000)
	mockSStore2.EXPECT().SlotRange().Return(uint16(3), uint16(5)).AnyTimes()
	mockSStore3 := getMockSStore(ctrl, 1564297200000)
	mockSStore3.EXPECT().SlotRange().Return(uint16(6), uint16(13)).AnyTimes()
	mockSStore4 := getMockSStore(ctrl, 1564308000000)
	mockSStore4.EXPECT().SlotRange().Return(uint16(4), uint16(14)).AnyTimes()

	theFieldStore.insertSStore(mockSStore1)
	timeRange, ok := theFieldStore.TimeRange(10 * 1000)
	assert.Equal(t, int64(1564300810000), timeRange.Start)
	assert.Equal(t, int64(1564300900000), timeRange.End)
	assert.True(t, ok)

	theFieldStore.insertSStore(mockSStore2)
	theFieldStore.insertSStore(mockSStore4)
	timeRange, ok = theFieldStore.TimeRange(10 * 1000)
	assert.Equal(t, int64(1564300810000), timeRange.Start)
	assert.Equal(t, int64(1564308140000), timeRange.End)
	assert.True(t, ok)
}

func Test_fStore_flushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fStore := newFieldStore(10)
	theFieldStore := fStore.(*fieldStore)

	mockTF := makeMockDataFlusher(ctrl)
	mockTF.EXPECT().GetFieldMeta(gomock.Any()).Return(field.Meta{}, false)
	mockTF.EXPECT().GetFieldMeta(gomock.Any()).Return(field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	}, true).MaxTimes(2)
	mockSStore1 := getMockSStore(ctrl, 1564304400000)
	mockSStore2 := getMockSStore(ctrl, 1564308000000)
	mockSStore2.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(0)

	theFieldStore.insertSStore(mockSStore1)
	theFieldStore.insertSStore(mockSStore2)

	assert.Len(t, theFieldStore.sStoreNodes, 2)
	// familyTime not exist
	assert.Zero(t, theFieldStore.FlushFieldTo(mockTF, flushContext{familyTime: 1564297200000}))
	assert.Len(t, theFieldStore.sStoreNodes, 2)
	// mock error
	assert.Zero(t, theFieldStore.FlushFieldTo(mockTF, flushContext{familyTime: 1564304400000}))
	assert.Len(t, theFieldStore.sStoreNodes, 1)
	// mock ok
	assert.NotZero(t, theFieldStore.FlushFieldTo(mockTF, flushContext{familyTime: 1564308000000}))
	assert.Len(t, theFieldStore.sStoreNodes, 0)
}

func Test_fStore_removeSStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsINTF := newFieldStore(1)
	assert.Equal(t, emptyFieldStoreSize, fsINTF.MemSize())
	fs := fsINTF.(*fieldStore)
	// segments empty
	fs.removeSStore(0)
	fs.removeSStore(1)

	// assert sorted
	fs.insertSStore(getMockSStore(ctrl, 1))
	fs.insertSStore(getMockSStore(ctrl, 9))
	fs.insertSStore(getMockSStore(ctrl, 2))
	fs.insertSStore(getMockSStore(ctrl, 3))
	fs.insertSStore(getMockSStore(ctrl, 7))
	fs.insertSStore(getMockSStore(ctrl, 5))
	assert.NotEqual(t, emptyFieldStoreSize, fsINTF.MemSize())
	assert.True(t, sort.IsSorted(fs.sStoreNodes))
	// remove greater
	fs.removeSStore(10)
	// remove not exist
	fs.removeSStore(8)
	// remove smaller
	fs.removeSStore(0)
	// remove existed
	fs.removeSStore(9)
	fs.removeSStore(1)
	fs.removeSStore(3)
	fs.removeSStore(4)
	fs.removeSStore(2)
	fs.removeSStore(7)
}
