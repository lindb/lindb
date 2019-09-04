package tblstore

import (
	"hash/crc32"
	"sync/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/tsdb/field"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./metrics_data_flusher.go -destination=./metrics_data_flusher_mock.go -package tblstore

// MetricsDataFlusher is a wrapper of kv.Builder, provides ability to build a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: TSEntry
// Level3: compressed field data
type MetricsDataFlusher interface {
	// FlushFieldMeta writes the meta info a field
	FlushFieldMeta(fieldID uint16, fieldType field.Type)
	// FlushField writes a compressed field data to writer.
	FlushField(fieldID uint16, data []byte, startSlot, endSlot int)
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries(seriesID uint32)
	// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
	FlushMetric(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error
}

// NewMetricsDataFlusher returns a new MetricsDataFlusher,
// interval is used to calculate the time-range of field data slots.`
func NewMetricsDataFlusher(flusher kv.Flusher, interval int64) MetricsDataFlusher {
	return &metricsDataFlusher{
		interval:     interval,
		flusher:      flusher,
		blockBuilder: newBlockBuilder(),
		entryBuilder: newSeriesEntryBuilder()}
}

// metricsDataFlusher implements MetricsDataFlusher.
type metricsDataFlusher struct {
	interval     int64
	flusher      kv.Flusher
	blockBuilder *blockBuilder
	entryBuilder *entryBuilder
}

// FlushFieldMeta writes the meta info a field
func (w *metricsDataFlusher) FlushFieldMeta(fieldID uint16, fieldType field.Type) {
	w.blockBuilder.appendFieldMeta(fieldID, fieldType)
}

// FlushField writes a compressed field data to writer.
func (w *metricsDataFlusher) FlushField(fieldID uint16, data []byte, startSlot, endSlot int) {
	startTime := int64(startSlot) * w.interval
	endTime := int64(endSlot) * w.interval

	w.blockBuilder.addStartEndTime(startTime, endTime)
	w.entryBuilder.addField(fieldID, data, startTime, endTime)
}

// FlushSeries writes a full series, this will be called after writing all fields of this entry.
func (w *metricsDataFlusher) FlushSeries(seriesID uint32) {
	w.blockBuilder.addSeries(seriesID, w.entryBuilder.bytes(w.blockBuilder.metaFieldsID))
	w.entryBuilder.reset()
}

// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
func (w *metricsDataFlusher) FlushMetric(metricID uint32) error {
	if err := w.blockBuilder.finish(); err != nil {
		return err
	}
	if err := w.flusher.Add(metricID, w.blockBuilder.bytes()); err != nil {
		return err
	}
	w.blockBuilder.reset()
	return nil
}

// Commit adds the footer and then closes the kv builder, this will be called after writing all metric-blocks.
func (w *metricsDataFlusher) Commit() error {
	return w.flusher.Commit()
}

// blockBuilder builds a metric-block containing multi TSEntry in order.
type blockBuilder struct {
	minStartTime    int64                 // startTime of fields-meta
	maxEndTime      int64                 // endTime of fields-meta
	metaFieldsID    []uint16              // fieldID list of fields-meta
	metaFieldsIDMap map[uint16]int        // set
	metaFieldsType  map[uint16]field.Type // field-id -> fieldType
	writer          *stream.BufferWriter
	offset          *encoding.DeltaBitPackingEncoder
	keys            *roaring.Bitmap
}

// newBlockBuilder returns a new metricBlockBuilder.
func newBlockBuilder() *blockBuilder {
	return &blockBuilder{
		metaFieldsIDMap: make(map[uint16]int),
		metaFieldsType:  make(map[uint16]field.Type),
		keys:            roaring.New(),
		writer:          stream.NewBufferWriter(nil),
		offset:          encoding.NewDeltaBitPackingEncoder()}
}

// addSeries puts tsID and metric-block data in memory.
func (blockBuilder *blockBuilder) addSeries(seriesID uint32, data []byte) {
	offset := blockBuilder.writer.Len()
	blockBuilder.offset.Add(int32(offset))
	blockBuilder.keys.Add(seriesID)
	blockBuilder.writer.PutBytes(data)
}

// appendFieldMeta updates min start-time and max end-time.
func (blockBuilder *blockBuilder) addStartEndTime(startTime, endTime int64) {
	// collect min-startTime and max-endTime of one entry.
	atomic.CompareAndSwapInt64(&blockBuilder.minStartTime, 0, startTime)
	if blockBuilder.minStartTime > startTime {
		blockBuilder.minStartTime = startTime
	}
	if blockBuilder.maxEndTime < endTime {
		blockBuilder.maxEndTime = endTime
	}
}

// appendFieldMeta builds a field-id list in order.
func (blockBuilder *blockBuilder) appendFieldMeta(fieldID uint16, fieldType field.Type) {
	if _, ok := blockBuilder.metaFieldsIDMap[fieldID]; ok {
		return
	}
	blockBuilder.metaFieldsIDMap[fieldID] = len(blockBuilder.metaFieldsID)
	blockBuilder.metaFieldsType[fieldID] = fieldType
	blockBuilder.metaFieldsID = append(blockBuilder.metaFieldsID, fieldID)
}

// reset resets the buffer, keys and offset.
func (blockBuilder *blockBuilder) reset() {
	blockBuilder.minStartTime = 0
	blockBuilder.maxEndTime = 0
	blockBuilder.metaFieldsID = blockBuilder.metaFieldsID[:0]
	blockBuilder.metaFieldsIDMap = make(map[uint16]int)
	blockBuilder.metaFieldsType = make(map[uint16]field.Type)
	blockBuilder.writer.Reset()
	blockBuilder.keys.Clear()
	blockBuilder.offset.Reset()
}

// finish writes keys, offset.
func (blockBuilder *blockBuilder) finish() error {
	// write offset
	posOfOffset := blockBuilder.writer.Len()
	offset := blockBuilder.offset.Bytes()
	blockBuilder.writer.PutBytes(offset)

	// write keys
	blockBuilder.keys.RunOptimize()
	keys, err := blockBuilder.keys.MarshalBinary()
	if err != nil {
		return err
	}
	posOfKeys := blockBuilder.writer.Len()
	blockBuilder.writer.PutBytes(keys)

	// write fields-meta
	posOfMeta := blockBuilder.writer.Len()
	// write start-time, end-time
	blockBuilder.writer.PutUvarint64(uint64(blockBuilder.minStartTime))
	blockBuilder.writer.PutUvarint64(uint64(blockBuilder.maxEndTime))
	// write fields count
	blockBuilder.writer.PutUvarint64(uint64(len(blockBuilder.metaFieldsID)))
	// write field-id, field-type list
	for _, fieldID := range blockBuilder.metaFieldsID {
		// write field-id
		blockBuilder.writer.PutUInt16(fieldID)
		// write field-type
		fieldType := blockBuilder.metaFieldsType[fieldID]
		blockBuilder.writer.PutByte(byte(fieldType))
	}
	// write footer, length: 4+4+4+4
	blockBuilder.writer.PutUint32(uint32(posOfOffset))
	blockBuilder.writer.PutUint32(uint32(posOfKeys))
	blockBuilder.writer.PutUint32(uint32(posOfMeta))
	// write crc32
	blockBuilder.writer.PutUint32(crc32.ChecksumIEEE(blockBuilder.bytes()))
	return nil
}

// bytes returns a slice of the underlying data written before.
func (blockBuilder *blockBuilder) bytes() []byte {
	data, _ := blockBuilder.writer.Bytes()
	return data
}

// Level4: builder
// entryBuilder builds a series containing different fields.
type entryBuilder struct {
	dataWriter   *stream.BufferWriter
	lenWriter    *stream.BufferWriter
	fieldsData   map[uint16][]byte
	bitArray     *collections.BitArray
	minStartTime int64 // startTime of fields-meta
	maxEndTime   int64 // endTime of fields-meta
}

// newSeriesEntryBuilder returns a new entryBuilder, default first 2 byte is the column count.
func newSeriesEntryBuilder() *entryBuilder {
	bitArray, _ := collections.NewBitArray(nil)
	return &entryBuilder{
		bitArray:   bitArray,
		dataWriter: stream.NewBufferWriter(nil),
		lenWriter:  stream.NewBufferWriter(nil),
		fieldsData: make(map[uint16][]byte)}
}

// addField puts fieldName and data into buffer in memory.
func (entryBuilder *entryBuilder) addField(fieldID uint16, data []byte, startTime, endTime int64) {
	atomic.CompareAndSwapInt64(&entryBuilder.minStartTime, 0, startTime)
	if entryBuilder.minStartTime > startTime {
		entryBuilder.minStartTime = startTime
	}
	if entryBuilder.maxEndTime < endTime {
		entryBuilder.maxEndTime = endTime
	}
	entryBuilder.fieldsData[fieldID] = data
}

// bytes builds a slice of the underlying fields-info data written before.
func (entryBuilder *entryBuilder) bytes(metaFieldsID []uint16) []byte {
	entryBuilder.dataWriter.Reset()
	entryBuilder.lenWriter.Reset()

	if len(metaFieldsID) == 0 {
		return nil
	}

	var existedFieldsID []uint16
	// write start-time
	entryBuilder.dataWriter.PutUvarint64(uint64(entryBuilder.minStartTime))
	// write end-time
	entryBuilder.dataWriter.PutUvarint64(uint64(entryBuilder.maxEndTime))
	// build bit-array
	for idx, fieldID := range metaFieldsID {
		if _, ok := entryBuilder.fieldsData[fieldID]; !ok {
			continue
		}
		existedFieldsID = append(existedFieldsID, fieldID)
		entryBuilder.bitArray.SetBit(uint16(idx))
	}
	// write bit-array length
	entryBuilder.dataWriter.PutUvarint64(uint64(entryBuilder.bitArray.Len()))
	// write bit-array
	entryBuilder.dataWriter.PutBytes(entryBuilder.bitArray.Bytes())
	// write variant length in order of fields in fields-meta
	for _, fieldID := range existedFieldsID {
		theData := entryBuilder.fieldsData[fieldID]
		entryBuilder.lenWriter.PutUvarint64(uint64(len(theData)))
	}
	lenData, _ := entryBuilder.lenWriter.Bytes()
	entryBuilder.dataWriter.PutBytes(lenData)
	// write data in order of fields in fields-meta
	for _, fieldID := range existedFieldsID {
		theData := entryBuilder.fieldsData[fieldID]
		entryBuilder.dataWriter.PutBytes(theData)
	}
	data, _ := entryBuilder.dataWriter.Bytes()
	return data
}

// reset resets the inner buffer and time-range.
func (entryBuilder *entryBuilder) reset() {
	entryBuilder.dataWriter.Reset()
	entryBuilder.lenWriter.Reset()
	entryBuilder.bitArray.Reset()
	entryBuilder.fieldsData = make(map[uint16][]byte)
	entryBuilder.minStartTime = 0
	entryBuilder.maxEndTime = 0
}
