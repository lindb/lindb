package memdb

import (
	"testing"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newTimeSeriesStore(t *testing.T) {
	tStore := newTimeSeriesStore()
	assert.NotNil(t, tStore)
	assert.True(t, tStore.IsNoData())
	assert.False(t, tStore.IsExpired())
}

func Test_tStore_expired(t *testing.T) {
	tStore := newTimeSeriesStore()
	time.Sleep(time.Millisecond * 1)
	assert.False(t, tStore.IsExpired())

	seriesTTL.Store(time.Nanosecond)
	time.Sleep(time.Millisecond * 1)
	assert.True(t, tStore.IsExpired())
}

func Test_tStore_write_sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock fieldID getter
	mockFieldIDGetter := NewMockmStoreFieldIDGetter(ctrl)
	mockFieldIDGetter.EXPECT().GetFieldIDOrGenerate(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(uint16(1), nil).AnyTimes()
	// mock field-store
	mockFStore := NewMockfStoreINTF(ctrl)
	mockFStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(1).AnyTimes()
	mockFStore.EXPECT().GetFieldID().Return(uint16(1)).AnyTimes()
	// get existed fStore
	_, err := tStore.Write(
		&pb.Metric{
			Fields: []*pb.Field{
				{Name: "sum", Field: &pb.Field_Sum{Sum: &pb.Sum{
					Value: 1.0,
				}}},
				{Name: "unknown", Field: nil}},
		}, writeContext{
			metricID:            1,
			blockStore:          newBlockStore(30),
			mStoreFieldIDGetter: mockFieldIDGetter})
	assert.Nil(t, err)
	assert.False(t, tStoreInterface.IsNoData())
	fStore, ok := tStoreInterface.GetFStore(uint16(1))
	assert.True(t, ok)
	assert.NotNil(t, fStore)

	// insert test
	tStore.insertFStore(newFieldStore(3))
	tStore.insertFStore(newFieldStore(2))
}

func Test_tStore_write_gauge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock fieldID getter
	mockFieldIDGetter := NewMockmStoreFieldIDGetter(ctrl)
	mockFieldIDGetter.EXPECT().GetFieldIDOrGenerate(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(uint16(1), nil).AnyTimes()
	// mock field-store
	mockFStore := NewMockfStoreINTF(ctrl)
	mockFStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(1).AnyTimes()
	mockFStore.EXPECT().GetFieldID().Return(uint16(1)).AnyTimes()
	// get existed fStore
	_, err := tStore.Write(
		&pb.Metric{
			Fields: []*pb.Field{
				{Name: "gauge", Field: &pb.Field_Gauge{Gauge: &pb.Gauge{
					Value: 1.0,
				}}},
				{Name: "unknown", Field: nil}},
		}, writeContext{
			metricID:            1,
			blockStore:          newBlockStore(30),
			mStoreFieldIDGetter: mockFieldIDGetter})
	assert.NoError(t, err)
	assert.False(t, tStoreInterface.IsNoData())
	fStore, ok := tStoreInterface.GetFStore(uint16(1))
	assert.True(t, ok)
	assert.NotNil(t, fStore)
}

func Test_tStore_GenFieldID_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock id generator
	mockGetter := NewMockmStoreFieldIDGetter(ctrl)
	mockGetter.EXPECT().GetFieldIDOrGenerate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uint16(1), series.ErrWrongFieldType).AnyTimes()
	// error field type from generator
	tStore.fStoreNodes = nil
	_, err := tStore.Write(
		&pb.Metric{Fields: []*pb.Field{{Name: "field1", Field: &pb.Field_Sum{}}}}, writeContext{
			metricID:            1,
			blockStore:          newBlockStore(30),
			mStoreFieldIDGetter: mockGetter})
	assert.Equal(t, series.ErrWrongFieldType, err)
}

func Test_tStore_flushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)
	assert.Equal(t, emptyTimeSeriesStoreSize, tStore.MemSize())

	mockTF := makeMockDataFlusher(ctrl)

	familyTime := timeutil.Now() / 3600 / 1000 * 3600 * 1000
	// has data
	mockFStore1 := NewMockfStoreINTF(ctrl)
	mockFStore1.EXPECT().SegmentsCount().Return(1).AnyTimes()
	mockFStore1.EXPECT().GetFieldID().Return(uint16(1)).AnyTimes()
	mockFStore1.EXPECT().MemSize().Return(emptyFieldStoreSize).AnyTimes()
	mockFStore1.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()).Return(100).AnyTimes()
	mockFStore1.EXPECT().TimeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*60, End: familyTime + 1000*120}, true).AnyTimes()
	mockFStore2 := NewMockfStoreINTF(ctrl)
	mockFStore2.EXPECT().SegmentsCount().Return(1).AnyTimes()
	mockFStore2.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()).Return(100).AnyTimes()
	mockFStore2.EXPECT().GetFieldID().Return(uint16(2)).AnyTimes()
	mockFStore2.EXPECT().MemSize().Return(emptyFieldStoreSize).AnyTimes()
	mockFStore2.EXPECT().TimeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*70, End: familyTime + 1000*130}, true).AnyTimes()
	mockFStore3 := NewMockfStoreINTF(ctrl)
	mockFStore3.EXPECT().SegmentsCount().Return(1).AnyTimes()
	mockFStore3.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()).Return(0).AnyTimes()
	mockFStore3.EXPECT().MemSize().Return(emptyFieldStoreSize).AnyTimes()
	mockFStore3.EXPECT().TimeRange(gomock.Any()).Return(
		timeutil.TimeRange{Start: 100, End: 200}, false).AnyTimes()
	mockFStore3.EXPECT().GetFieldID().Return(uint16(3)).AnyTimes()

	tStore.insertFStore(mockFStore1)
	tStore.insertFStore(mockFStore2)
	tStore.insertFStore(mockFStore3)
	assert.NotEqual(t, emptyTimeSeriesStoreSize, tStore.MemSize())
	assert.NotZero(t, tStore.FlushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
	assert.False(t, tStoreInterface.IsNoData())

	// flush error
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)

	assert.Zero(t, tStore.FlushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))

	// no-data
	mockFStore4 := NewMockfStoreINTF(ctrl)
	mockFStore4.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()).Return(10).AnyTimes()
	mockFStore4.EXPECT().TimeRange(gomock.Any()).Return(timeutil.TimeRange{Start: 0, End: 0}, false).AnyTimes()
	mockFStore4.EXPECT().GetFieldID().Return(uint16(4)).AnyTimes()
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)
	tStore.insertFStore(mockFStore4)
	assert.NotZero(t, tStore.FlushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
}
