package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

type mockScanWorker struct {
	events []*series.FieldEvent
}

func (w *mockScanWorker) Emit(event *series.FieldEvent) {
	w.events = append(w.events, event)
}
func (w *mockScanWorker) Complete(seriesID uint32) {}
func (w *mockScanWorker) Close()                   {}

func TestFieldStore_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bs := newBlockStore(10)
	calc := interval.GetCalculator(interval.Day)

	now, _ := timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	familyTime, _ := timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	tStore := newTimeSeriesStore(100)
	ts := tStore.(*timeSeriesStore)

	fStore := newFieldStore(10)
	sCtx := &series.ScanContext{TimeRange: timeutil.TimeRange{
		Start: now - 100,
		End:   now + 1000,
	}, IntervalCalc: calc, Interval: 10000}
	fieldMeta := &fieldMeta{fieldID: 1, fieldName: "f1", fieldType: field.SumField}
	// no data
	fStore.Scan(sCtx, series.Version(10), uint32(10), fieldMeta, ts)

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

	// time range not match
	now, _ = timeutil.ParseTimestamp("20190702 20:10:48", "20060102 15:04:05")
	fStore.Scan(&series.ScanContext{TimeRange: timeutil.TimeRange{
		Start: now - 100,
		End:   now + 1000,
	}, IntervalCalc: calc, Interval: 10000}, series.Version(10), uint32(10), fieldMeta, ts)

	// found it
	now, _ = timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	worker := &mockScanWorker{}
	fStore.Scan(&series.ScanContext{
		TimeRange: timeutil.TimeRange{
			Start: now - 100,
			End:   now + 1000,
		},
		IntervalCalc: calc,
		Worker:       worker,
		Interval:     10000,
	}, series.Version(10), uint32(10), fieldMeta, ts)

	assert.Equal(t, 1, len(worker.events))
	it := worker.events[0].FieldIt
	_, err := it.Bytes()
	assert.NotNil(t, err)
	assert.Equal(t, familyTime, it.SegmentStartTime())
	assert.Equal(t, uint16(1), it.FieldID())
	assert.Equal(t, "f1", it.FieldName())
	assert.Equal(t, field.SumField, it.FieldType())
	assert.True(t, it.HasNext())
	pIt := it.Next()
	assert.True(t, pIt.HasNext())
	slot, val := pIt.Next()
	assert.Equal(t, 20, slot)
	assert.Equal(t, 1.0, val)
	assert.False(t, pIt.HasNext())
	assert.False(t, it.HasNext())

	// write data
	for i := 0; i < 10; i++ {
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
	}
	// found it
	now, _ = timeutil.ParseTimestamp("20190702 19:10:48", "20060102 15:04:05")
	worker = &mockScanWorker{}
	fStore.Scan(&series.ScanContext{
		TimeRange: timeutil.TimeRange{
			Start: now - 100,
			End:   now + 1000,
		},
		IntervalCalc: calc,
		Worker:       worker,
		Interval:     10000,
	}, series.Version(10), uint32(10), fieldMeta, ts)

	it = worker.events[0].FieldIt
	assert.True(t, it.HasNext())
	pIt = it.Next()
	assert.True(t, pIt.HasNext())
	slot, val = pIt.Next()
	assert.Equal(t, 20, slot)
	assert.Equal(t, 11.0, val)
	assert.False(t, pIt.HasNext())
	assert.False(t, it.HasNext())
}
