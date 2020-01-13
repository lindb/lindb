package memdb

import (
	"math"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./segment_store.go -destination=./segment_store_mock_test.go -package memdb

// for testing
var (
	encodeFunc = encoding.NewTSDEncoder
)

const (
	emptySimpleFieldStoreSize = 8 + // familyTime
		2 + // start time slot
		8 + // aggFunc
		8 // block pointer
)

// sStoreINTF represents segment-store,
// which abstracts a store for storing field data based on family start time
type sStoreINTF interface {
	GetFamilyTime() int64

	SlotRange() (startSlot, endSlot uint16)

	// FlushFieldTo flushes segment's data to writer
	FlushFieldTo(
		tableFlusher metricsdata.Flusher,
		fieldMeta field.Meta,
		flushCtx flushContext,
	) (
		flushedSize int,
	)

	// Write writes the metric's field with writeContext
	// 1) block is nil, create new block
	// 2) in current time window, if has old value need do rollup
	Write(
		fieldType field.Type,
		f *pb.Field,
		writeCtx writeContext,
	) (writtenSize int)

	// CheckAndCompact checks current time window's block for storing field data based on slot time and value type.
	// returns memSize that is compress data's length.
	// if slot time out of current time window, need compress time window then create new one
	CheckAndCompact(fieldType field.Type, writeCtx writeContext) (memSize int)

	// MemSize returns the segment store memory size
	MemSize() int

	// scan scans segment store data based on query time range
	scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext)
}

// isInCurrentTimeWindow returns the current time slot is in current time window
func isInCurrentTimeWindow(start, end uint16, current uint16) bool {
	return current >= start && current <= end
}

// getTimeSlotRange returns the final time slot range based on start/end
func getTimeSlotRange(startSlot1, endSlot1 uint16, startSlot2, endSlot2 uint16) (start, end uint16) {
	start = startSlot1
	end = endSlot1
	if end < endSlot2 {
		end = endSlot2
	}
	if start > startSlot2 {
		start = startSlot2
	}
	return
}

// compactInt compress block data
func compact(aggFunc field.AggFunc, tsd *encoding.TSDDecoder, block block,
	startTime, startSlot, endSlot uint16, withTimeRange bool,
) (compress []byte, err error) {
	if block == nil && tsd == nil {
		return
	}
	encode := encodeFunc(startSlot)
	for i := startSlot; i <= endSlot; i++ {
		newValue, hasNewValue := getCurrentFloatValue(block, startTime, i)
		oldValue, hasOldValue := getOldFloatValue(tsd, i)
		switch {
		case hasNewValue && !hasOldValue:
			// just compress current block value with pos
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(newValue))
		case hasNewValue && hasOldValue:
			// merge and compress
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(aggFunc.AggregateFloat(newValue, oldValue)))
		case !hasNewValue && hasOldValue:
			// compress old value
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(oldValue))
		default:
			// append empty value
			encode.AppendTime(bit.Zero)
		}
	}
	if withTimeRange {
		return encode.Bytes()
	}
	// get compress data without time slot range
	return encode.BytesWithoutTime()
}

func getOldFloatValue(tsd *encoding.TSDDecoder, timeSlot uint16) (value float64, hasValue bool) {
	if tsd == nil {
		return
	}
	if !tsd.HasValueWithSlot(timeSlot) {
		return
	}
	hasValue = true
	value = math.Float64frombits(tsd.Value())
	return
}

func getCurrentFloatValue(block block, startTime uint16, timeSlot uint16) (value float64, hasValue bool) {
	if block == nil {
		return
	}
	if timeSlot < startTime || timeSlot > startTime+block.getSize() {
		return
	}
	hasValue = true
	value = block.getFloatValue(timeSlot - startTime)
	return
}
