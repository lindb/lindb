package memdb

import (
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/timeutil"
	pb "github.com/eleme/lindb/rpc/proto/field"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

//go:generate mockgen -source ./field_store.go -destination=./field_store_mock_test.go -package memdb

// fStoreINTF abstracts a field-store
type fStoreINTF interface {
	// getFieldType returns the field-type
	getFieldType() field.Type
	// write writes the metric's field with writeContext
	write(f *pb.Field, writeCtx writeContext)
	// flushFieldTo flushes field data of the specific familyTime
	// return false if there is no data related of familyTime
	flushFieldTo(tableFlusher metrictbl.TableFlusher, familyTime int64) (flushed bool)
	// timeRange returns the start-time and end-time of fStore's data
	// ok means data is available
	timeRange(interval int64) (timeRange timeutil.TimeRange, ok bool)
}

// todo:@codingcrush, replace segments with slice
// fieldStore holds the relation of familyStartTime and segmentStore.
type fieldStore struct {
	fieldType field.Type           // sum, gauge, min, max
	fieldID   uint16               // default 0
	segments  map[int64]sStoreINTF // familyTime => segment store
}

// newFieldStore returns a new fieldStore.
func newFieldStore(fieldID uint16, fieldType field.Type) fStoreINTF {
	return &fieldStore{
		fieldID:   fieldID,
		fieldType: fieldType,
		segments:  make(map[int64]sStoreINTF),
	}
}

// getFieldType returns field type for current field store
func (fs *fieldStore) getFieldType() field.Type {
	return fs.fieldType
}

func (fs *fieldStore) write(f *pb.Field, writeCtx writeContext) {
	switch fields := f.Field.(type) {
	case *pb.Field_Sum:
		sStore, exist := fs.segments[writeCtx.familyTime]
		if !exist {
			//TODO ???
			sStore = newSimpleFieldStore(field.GetAggFunc(field.Sum))
			fs.segments[writeCtx.familyTime] = sStore
		}
		sStore.writeFloat(fields.Sum, writeCtx)
	default:
		memDBLogger.Warn("convert field error, unknown field type")
	}
}

// flushFieldTo flushes segments' data to writer and reset the segments-map.
func (fs *fieldStore) flushFieldTo(tableFlusher metrictbl.TableFlusher, familyTime int64) (flushed bool) {
	ss, ok := fs.segments[familyTime]
	if !ok {
		return false
	}
	delete(fs.segments, familyTime)

	data, startSlot, endSlot, err := ss.bytes()

	if err != nil {
		memDBLogger.Error("read segment data error:", logger.Error(err))
		return false
	}
	tableFlusher.FlushField(fs.fieldID, fs.fieldType, data, startSlot, endSlot)
	return true
}

func (fs *fieldStore) timeRange(interval int64) (timeRange timeutil.TimeRange, ok bool) {
	for familyTime, sStore := range fs.segments {
		startSlot, endSlot, err := sStore.slotRange()
		if err != nil {
			continue
		}
		ok = true
		startTime := familyTime + int64(startSlot)*interval
		endTime := familyTime + int64(endSlot)*interval
		if timeRange.Start == 0 || startTime < timeRange.Start {
			timeRange.Start = startTime
		}
		if timeRange.End == 0 || timeRange.End < endTime {
			timeRange.End = endTime
		}
	}
	return
}
