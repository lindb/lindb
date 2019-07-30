package memdb

import (
	"strconv"
	"testing"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/timeutil"
	pb "github.com/eleme/lindb/rpc/proto/field"
	"github.com/eleme/lindb/tsdb/index"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newTimeSeriesStore(t *testing.T) {
	tStore := newTimeSeriesStore(100, 100)
	assert.NotNil(t, tStore)
	assert.True(t, tStore.isNoData())
	assert.False(t, tStore.isExpired())

	_, ok := tStore.timeRange()
	assert.False(t, ok)
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

	mockFStore := NewMockfStoreINTF(ctrl)
	mockFStore.EXPECT().write(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockFStore.EXPECT().getFieldType().Return(field.SumField).AnyTimes()
	mockFStore.EXPECT().getFieldName().Return("sum").AnyTimes()
	// get existed fStore
	tStore.insertFStore(mockFStore)
	writeCtx := writeContext{
		metricID:   1,
		generator:  makeMockIDGenerator(ctrl),
		blockStore: newBlockStore(30)}
	err := tStore.write(
		&pb.Metric{
			Fields: []*pb.Field{
				{Name: "sum", Field: &pb.Field_Sum{}},
				{Name: "min", Field: &pb.Field_Min{}},
				{Name: "unknown", Field: nil}},
		}, writeCtx)
	assert.Nil(t, err)
	assert.False(t, tStoreInterface.isNoData())
	// error field type
	err = tStore.write(&pb.Metric{Fields: []*pb.Field{{Name: "sum", Field: &pb.Field_Min{}}}}, writeCtx)
	assert.Equal(t, models.ErrWrongFieldType, err)

	// create new fStore
	err = tStore.write(&pb.Metric{Fields: []*pb.Field{{Name: "sum1", Field: &pb.Field_Sum{}}}}, writeCtx)
	assert.Nil(t, err)
	// too many fields
	for i := 0; i < 3000; i++ {
		tStore.insertFStore(newFieldStore(strconv.Itoa(i), uint16(i), field.SumField))
	}
	err = tStore.write(&pb.Metric{Fields: []*pb.Field{{Name: "sum2", Field: &pb.Field_Sum{}}}}, writeCtx)
	assert.Equal(t, models.ErrTooManyFields, err)
}

func Test_tStore_GenFieldID_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)
	// mock id generator
	mockGen := index.NewMockIDGenerator(ctrl)
	mockGen.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		uint16(0), models.ErrWrongFieldType).AnyTimes()
	// error field type from generator
	tStore.fStoreNodes = nil
	err := tStore.write(&pb.Metric{Fields: []*pb.Field{{Name: "field1", Field: &pb.Field_Sum{}}}}, writeContext{
		metricID:   1,
		generator:  mockGen,
		blockStore: newBlockStore(30)})
	assert.Equal(t, models.ErrWrongFieldType, err)
}

func Test_tStore_afterWrite(t *testing.T) {
	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)

	writeCtx := writeContext{
		timeInterval: 10 * 1000,
		slotIndex:    40,
		familyTime:   timeutil.Now() / 3600 / 1000 * 3600 * 1000}
	tStore.afterWrite(writeCtx)
	assert.True(t, tStore.hasData)
	assert.Equal(t, tStore.startDelta, tStore.endDelta)
	timeRange, _ := tStore.timeRange()
	assert.Equal(t, timeRange.Start, timeRange.End)

	writeCtx.slotIndex = 380
	tStore.afterWrite(writeCtx)
	timeRange, _ = tStore.timeRange()
	assert.True(t, timeRange.Start < timeRange.End)
	assert.True(t, timeRange.End > timeutil.Now())
}

func Test_tStore_flushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore(100, 100)
	tStore := tStoreInterface.(*timeSeriesStore)

	mockTF := makeMockTableFlusher(ctrl)

	familyTime := timeutil.Now() / 3600 / 1000 * 3600 * 1000
	// has data
	mockFStore1 := NewMockfStoreINTF(ctrl)
	mockFStore1.EXPECT().getFieldName().Return("1").AnyTimes()
	mockFStore1.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore1.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*60, End: familyTime + 1000*120}, true).AnyTimes()
	mockFStore2 := NewMockfStoreINTF(ctrl)
	mockFStore2.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore2.EXPECT().getFieldName().Return("2").AnyTimes()
	mockFStore2.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{
		Start: familyTime + 1000*70, End: familyTime + 1000*130}, true).AnyTimes()
	mockFStore3 := NewMockfStoreINTF(ctrl)
	mockFStore3.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mockFStore3.EXPECT().timeRange(gomock.Any()).Return(
		timeutil.TimeRange{Start: 100, End: 200}, false).AnyTimes()
	mockFStore3.EXPECT().getFieldName().Return("3").AnyTimes()

	tStore.insertFStore(mockFStore1)
	tStore.insertFStore(mockFStore2)
	tStore.insertFStore(mockFStore3)
	assert.True(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
	assert.False(t, tStoreInterface.isNoData())
	timeRange, _ := tStoreInterface.timeRange()
	assert.Equal(t, int64(70), (timeRange.End-timeRange.Start)/1000)

	// flush error
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)

	assert.False(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))

	// no-data
	mockFStore4 := NewMockfStoreINTF(ctrl)
	mockFStore4.EXPECT().flushFieldTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockFStore4.EXPECT().timeRange(gomock.Any()).Return(timeutil.TimeRange{Start: 0, End: 0}, false).AnyTimes()
	mockFStore4.EXPECT().getFieldName().Return("4").AnyTimes()
	tStore.fStoreNodes = nil
	tStore.insertFStore(mockFStore3)
	tStore.insertFStore(mockFStore4)
	assert.True(t, tStore.flushSeriesTo(mockTF, flushContext{timeInterval: 10 * 1000}))
	assert.True(t, tStoreInterface.isNoData())
}
