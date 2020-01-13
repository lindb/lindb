package memdb

import (
	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

// for testing
var (
	newTSDStreamWriter = encoding.NewTSDStreamWriter
)

const (
	emptyComplexFieldStoreSize = 8 + // familyTime
		2 + //start time
		8 + // block pointer
		8 // compress pointer
)

type complexFieldStore struct {
	familyTime int64
	startTime  uint16

	blocks   map[uint16]block
	compress []byte
}

func newComplexFieldStore(familyTime int64) sStoreINTF {
	return &complexFieldStore{
		familyTime: familyTime,
		blocks:     make(map[uint16]block),
	}
}

func (fs *complexFieldStore) Write(
	fieldType field.Type,
	f *pb.Field,
	writeCtx writeContext,
) (writtenSize int) {
	switch fieldType {
	case field.SummaryField:
		writtenSize += fs.write(fieldType, 1, f.GetSummary().Sum, writeCtx)
		writtenSize += fs.write(fieldType, 2, f.GetSummary().Count, writeCtx)
	case field.HistogramField:
		//FIXME stone1100
	}

	return
}

func (fs *complexFieldStore) write(fieldType field.Type, pFieldID uint16, value float64, writeCtx writeContext) (writtenSize int) {
	block, ok := fs.blocks[pFieldID]
	current := writeCtx.slotIndex
	if !ok {
		block = writeCtx.blockStore.allocFloatBlock()
		fs.startTime = current
		block.setFloatValue(0, value)
		fs.blocks[pFieldID] = block
		return block.memsize()
	}
	pos := current - fs.startTime
	if block.hasValue(pos) {
		// do rollup using agg func
		aggFunc := fieldType.GetSchema().GetAggFunc(pFieldID)
		block.setFloatValue(pos, aggFunc.AggregateFloat(block.getFloatValue(pos), value))
	} else {
		block.setFloatValue(pos, value)
	}
	return
}

func (fs *complexFieldStore) CheckAndCompact(fieldType field.Type, writeCtx writeContext) (memSize int) {
	if len(fs.blocks) == 0 {
		return
	}
	current := writeCtx.slotIndex
	if !isInCurrentTimeWindow(fs.startTime, fs.startTime+writeCtx.blockStore.timeWindow-1, current) {
		// if current slot time out of current time window, need compress block data, start new time window
		start, end := fs.SlotRange()
		size := len(fs.compress)
		compress := fs.compact(fieldType, start, end)
		fs.compress = compress

		// reset block values
		for _, block := range fs.blocks {
			block.reset()
		}
		// !!!IMPORTANT: must reset start time
		fs.startTime = current
		return len(fs.compress) - size
	}
	return
}

func (fs *complexFieldStore) compact(fieldType field.Type, start, end uint16) []byte {
	size := len(fs.compress)
	fieldSchema := fieldType.GetSchema()
	pFieldIDs := fieldSchema.GetAllPrimitiveFields()
	tsdWriter := newTSDStreamWriter(start, end)
	var reader encoding.TSDStreamReader
	var tsd *encoding.TSDDecoder
	var fieldData *encoding.TSDDecoder
	var fieldID uint16
	next := true
	if size > 0 {
		reader = encoding.NewTSDStreamReader(fs.compress)
		defer reader.Close()
	}

	for _, pFieldID := range pFieldIDs {
		aggFunc := fieldSchema.GetAggFunc(pFieldID)
		block := fs.blocks[pFieldID]
		//FIXME stone1100 need test for field ids skip
		if reader != nil && next && reader.HasNext() {
			fieldID, tsd = reader.Next()
			if fieldID == pFieldID {
				fieldData = tsd
				next = true
			} else {
				next = false
			}
		}
		data, err := compact(aggFunc, fieldData, block, fs.startTime, start, end, false)
		fieldData = nil
		if err != nil {
			memDBLogger.Error("compact block data error, data will lost", logger.Error(err))
		}
		if len(data) > 0 {
			tsdWriter.WriteField(pFieldID, data)
		}
	}
	compress, err := tsdWriter.Bytes()
	if err != nil {
		memDBLogger.Error("compact complex fields data error, data will lost", logger.Error(err))
	}
	return compress
}

func (fs *complexFieldStore) GetFamilyTime() int64 {
	return fs.familyTime
}

func (fs *complexFieldStore) SlotRange() (
	startSlot,
	endSlot uint16) {

	endSlot = 0
	startSlot = fs.startTime
	for _, block := range fs.blocks {
		end := startSlot + block.getSize()
		if end > endSlot {
			endSlot = end
		}
	}
	if len(fs.compress) == 0 {
		return
	}
	start, end := encoding.DecodeTSDTime(fs.compress)
	return getTimeSlotRange(start, end, startSlot, endSlot)
}

func (fs *complexFieldStore) FlushFieldTo(
	tableFlusher metricsdata.Flusher,
	fieldMeta field.Meta,
	flushCtx flushContext,
) (
	flushedSize int,
) {
	if len(fs.blocks) == 0 {
		return
	}
	data := fs.compact(fieldMeta.Type, flushCtx.start, flushCtx.end)
	if data == nil {
		return
	}
	fs.compress = data
	tableFlusher.FlushPrimitiveField(uint16(0), fs.compress)

	return fs.MemSize()
}

func (fs *complexFieldStore) MemSize() int {
	if len(fs.blocks) == 0 {
		return emptyComplexFieldStoreSize
	}
	size := emptyComplexFieldStoreSize
	for _, block := range fs.blocks {
		size += 2 + // pField id size
			block.memsize() // block mem size
	}
	return size
}

// load loads block data, then aggregates the data
func (fs *complexFieldStore) load(
	fieldType field.Type,
	startSlot, endSlot uint16,
	agg []aggregation.PrimitiveAggregator,
	memScanCtx *memScanContext,
) {
	//FIXME stone100
}
