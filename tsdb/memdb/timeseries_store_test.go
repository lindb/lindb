package memdb

import (
	"testing"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/series"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newTimeSeriesStore(t *testing.T) {
	tStore := newTimeSeriesStore(100, 100)
	assert.NotNil(t, tStore)
	assert.True(t, tStore.isNoData())
	assert.False(t, tStore.isExpired())
	assert.Equal(t, uint64(100), tStore.getHash())
}

func Test_tStore_expired(t *testing.T) {
	tStore := newTimeSeriesStore(100, 100)
	time.Sleep(time.Millisecond * 1)
	assert.False(t, tStore.isExpired())
	setTagsIDTTL(1)
	time.Sleep(time.Millisecond * 1)
	assert.True(t, tStore.isExpired())
}

func Test_tStore_write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock fieldID getter
	mockFieldIDGetter := NewMockmStoreFieldIDGetter(ctrl)
	mockFieldIDGetter.EXPECT().getFieldIDOrGenerate(gomock.Any(),
		gomock.Any(), gomock.Any()).Return(uint16(1), nil).AnyTimes()
	// mock field-store
	mockFStore := NewMockfStoreINTF(ctrl)
	mockFStore.EXPECT().write(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockFStore.EXPECT().getFieldID().Return(uint16(1)).AnyTimes()
	// get existed fStore
	err := tStore.write(
		&pb.Metric{
			Fields: []*pb.Field{
				{Name: "sum", Field: &pb.Field_Sum{}},
				{Name: "min", Field: &pb.Field_Min{}},
				{Name: "unknown", Field: nil}},
		}, writeContext{
			metricID:            1,
			blockStore:          newBlockStore(30),
			mStoreFieldIDGetter: mockFieldIDGetter})
	assert.Nil(t, err)
	assert.False(t, tStoreInterface.isNoData())

	// insert test
	tStore.insertFStore(newFieldStore(3))
	tStore.insertFStore(newFieldStore(2))
}

func Test_tStore_GenFieldID_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock id generator
	mockGetter := NewMockmStoreFieldIDGetter(ctrl)
	mockGetter.EXPECT().getFieldIDOrGenerate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uint16(1), series.ErrWrongFieldType).AnyTimes()
	// error field type from generator
	tStore.fStoreNodes = nil
	err := tStore.write(&pb.Metric{Fields: []*pb.Field{{Name: "field1", Field: &pb.Field_Sum{}}}}, writeContext{
		metricID:            1,
		blockStore:          newBlockStore(30),
		mStoreFieldIDGetter: mockGetter})
	assert.Equal(t, series.ErrWrongFieldType, err)
}

func Test_tStore_afterWrite(t *testing.T) {
	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)

	tStore.afterWrite()
	assert.True(t, tStore.hasData)
}

func Test_tStore_flushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)

	mockTF := makeMockDataFlusher(ctrl)

	familyTime := timeutil.Now() / 3600 / 1000 * 3600 * 1000
	// has data
	mockFStore1 := NewMockfStoreINTF(ctrl)
	mockFStore1.EXPECT().getFieldID().Return(uint16(1)).AnyTimes()
	mockFStore1.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore1.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*60, End: familyTime + 1000*120}, true).AnyTimes()
	mockFStore2 := NewMockfStoreINTF(ctrl)
	mockFStore2.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore2.EXPECT().getFieldID().Return(uint16(2)).AnyTimes()
	mockFStore2.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*70, End: familyTime + 1000*130}, true).AnyTimes()
	mockFStore3 := NewMockfStoreINTF(ctrl)
	mockFStore3.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mockFStore3.EXPECT().timeRange(gomock.Any()).Return(
		timeutil.TimeRange{Start: 100, End: 200}, false).AnyTimes()
	mockFStore3.EXPECT().getFieldID().Return(uint16(3)).AnyTimes()

	tStore.insertFStore(mockFStore1)
	tStore.insertFStore(mockFStore2)
	tStore.insertFStore(mockFStore3)
	assert.True(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
	assert.False(t, tStoreInterface.isNoData())

	// flush error
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)

	assert.False(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))

	// no-data
	mockFStore4 := NewMockfStoreINTF(ctrl)
	mockFStore4.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore4.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{Start: 0, End: 0}, false).AnyTimes()
	mockFStore4.EXPECT().getFieldID().Return(uint16(4)).AnyTimes()
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)
	tStore.insertFStore(mockFStore4)
	assert.True(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
	assert.True(t, tStoreInterface.isNoData())
}
