package memdb

import (
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

// singleFieldStore represents single field store
type simpleFieldStore struct {
	familyTime int64
	startTime  uint16

	block    block
	compress []byte
}

// newSingleFieldStore returns a new segment store for simple field store
func newSimpleFieldStore(familyTime int64) sStoreINTF {
	return &simpleFieldStore{
		familyTime: familyTime,
	}
}

func (fs *simpleFieldStore) CheckAndCompact(fieldType field.Type, writeCtx writeContext) (memSize int) {
	if fs.block == nil {
		return
	}
	current := writeCtx.slotIndex
	if !isInCurrentTimeWindow(fs.startTime, fs.startTime+writeCtx.blockStore.timeWindow-1, current) {
		// if current slot time out of current time window, need compress block data, start new time window
		size := len(fs.compress)
		start, end := fs.SlotRange()
		data := fs.compact(fieldType, start, end)

		fs.compress = data
		fs.block.reset()
		// !!!IMPORTANT: must reset start time
		fs.startTime = current
		return len(fs.compress) - size
	}
	return
}

func (fs *simpleFieldStore) Write(
	fieldType field.Type,
	f *pb.Field,
	writeCtx writeContext,
) (writtenSize int) {
	current := writeCtx.slotIndex
	value := fs.getFieldValue(fieldType, f)
	if fs.block == nil {
		fs.block = writeCtx.blockStore.allocFloatBlock()
		fs.startTime = current
		fs.block.setFloatValue(0, value)
		return fs.block.memsize()
	}
	pos := current - fs.startTime
	if fs.block.hasValue(pos) {
		// do rollup using agg func
		aggFunc := fieldType.GetSchema().GetAggFunc(field.SimpleFieldPFieldID)
		fs.block.setFloatValue(pos, aggFunc.AggregateFloat(fs.block.getFloatValue(pos), value))
	} else {
		fs.block.setFloatValue(pos, value)
	}
	return
}

func (fs *simpleFieldStore) getFieldValue(fieldType field.Type, f *pb.Field) float64 {
	switch fieldType {
	case field.SumField:
		return f.GetSum().Value
	case field.MinField:
		return f.GetMin().Value
	case field.MaxField:
		return f.GetMax().Value
	case field.GaugeField:
		return f.GetGauge().Value
	default:
		return 0
	}
}

func (fs *simpleFieldStore) MemSize() int {
	return emptySimpleFieldStoreSize
}

func (fs *simpleFieldStore) GetFamilyTime() int64 {
	return fs.familyTime
}

func (fs *simpleFieldStore) FlushFieldTo(
	tableFlusher metricsdata.Flusher,
	fieldMeta field.Meta,
	flushCtx flushContext,
) (
	flushedSize int,
) {
	size := len(fs.compress)
	data := fs.compact(fieldMeta.Type, flushCtx.start, flushCtx.end)
	if data == nil {
		return
	}
	fs.compress = data
	tableFlusher.FlushPrimitiveField(field.SimpleFieldPFieldID, fs.compress)

	return fs.MemSize() + size
}

func (fs *simpleFieldStore) compact(fieldType field.Type, start, end uint16) []byte {
	aggFunc := fieldType.GetSchema().GetAggFunc(field.SimpleFieldPFieldID)
	size := len(fs.compress)
	var tsd *encoding.TSDDecoder
	if size > 0 {
		// calc new start/end based on old compress values
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	data, err := compact(aggFunc, tsd, fs.block, fs.startTime, start, end, true)
	if err != nil {
		memDBLogger.Error("compact simple segment store data err", logger.Error(err))
	}
	return data
}

func (fs *simpleFieldStore) SlotRange() (startSlot, endSlot uint16) {
	startSlot = fs.startTime
	endSlot = fs.startTime + fs.block.getSize()
	if len(fs.compress) == 0 {
		return
	}
	start, end := encoding.DecodeTSDTime(fs.compress)
	return getTimeSlotRange(start, end, startSlot, endSlot)
}

// load loads block data, then aggregates the data
func (fs *simpleFieldStore) load(
	fieldType field.Type,
	startSlot, endSlot uint16,
	agg []aggregation.PrimitiveAggregator,
	memScanCtx *memScanContext,
) {
	hasOld := len(fs.compress) > 0
	aggFunc := fieldType.GetSchema().GetAggFunc(field.SimpleFieldPFieldID)

	var tsd *encoding.TSDDecoder
	if hasOld {
		// calc new start/end based on old compress values
		tsd = memScanCtx.tsd
		tsd.Reset(fs.compress)
	}
	completed := false
	value := 0.0
	for i := startSlot; i <= endSlot; i++ {
		newValue, hasNewValue := getCurrentFloatValue(fs.block, fs.startTime, i)
		oldValue, hasOldValue := getOldFloatValue(tsd, i)

		switch {
		case hasNewValue && !hasOldValue:
			// get value from new block buffer
			value = newValue
		case hasNewValue && hasOldValue:
			// merge data from new and old
			value = aggFunc.AggregateFloat(newValue, oldValue)
		case !hasNewValue && hasOldValue:
			// get old value from compress data
			value = oldValue
		}
		if hasNewValue || hasOldValue {
			//FIXME stone1100

			idx := int(i)
			for _, a := range agg {
				// aggregate the data
				completed = a.Aggregate(idx, value)
			}
		}
		if completed {
			return
		}
	}
}
