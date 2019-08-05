package metrictbl

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sync/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/field"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metrictbl

// TableFlusher is a wrapper of kv.Builder, provides ability to build a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: TSEntry
// Level3: compressed field data
type TableFlusher interface {
	// FlushField writes a compressed field data to writer.
	FlushField(fieldID uint16, fieldType field.Type, data []byte, startSlot, endSlot int)
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries(seriesID uint32)
	// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
	FlushMetric(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error
}

// NewTableFlusher returns a new TableWriter, interval is used to calculate the time-range of field data slots.`
func NewTableFlusher(flusher kv.Flusher, interval int64) TableFlusher {
	return &tableFlusher{
		interval:     interval,
		flusher:      flusher,
		blockBuilder: newBlockBuilder(),
		entryBuilder: newSeriesEntryBuilder()}
}

// tableFlusher implements TableWriter.
type tableFlusher struct {
	interval     int64
	flusher      kv.Flusher
	blockBuilder *blockBuilder
	entryBuilder *entryBuilder
}

// FlushField writes a compressed field data to writer.
func (w *tableFlusher) FlushField(fieldID uint16, fieldType field.Type, data []byte, startSlot, endSlot int) {
	startTime := int64(startSlot) * w.interval
	endTime := int64(endSlot) * w.interval

	w.blockBuilder.appendFieldMeta(fieldID, fieldType, startTime, endTime)
	w.entryBuilder.addField(fieldID, data, startTime, endTime)
}

// FlushSeries writes a full series, this will be called after writing all fields of this entry.
func (w *tableFlusher) FlushSeries(seriesID uint32) {
	w.blockBuilder.addSeries(seriesID, w.entryBuilder.bytes(w.blockBuilder.metaFieldsID))
	w.entryBuilder.reset()
}

// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
func (w *tableFlusher) FlushMetric(metricID uint32) error {
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
func (w *tableFlusher) Commit() error {
	return w.flusher.Commit()
}

// blockBuilder builds a metric-block containing multi TSEntry in order.
type blockBuilder struct {
	minStartTime    int64                 // startTime of fields-meta
	maxEndTime      int64                 // endTime of fields-meta
	metaFieldsID    []uint16              // fieldID list of fields-meta
	metaFieldsIDMap map[uint16]int        // set
	metaFieldsType  map[uint16]field.Type // field-id -> fieldType
	buf             bytes.Buffer
	offset          *encoding.DeltaBitPackingEncoder
	keys            *roaring.Bitmap
}

// newBlockBuilder returns a new metricBlockBuilder.
func newBlockBuilder() *blockBuilder {
	return &blockBuilder{
		metaFieldsIDMap: make(map[uint16]int),
		metaFieldsType:  make(map[uint16]field.Type),
		keys:            roaring.New(),
		offset:          encoding.NewDeltaBitPackingEncoder()}
}

// addSeries puts tsID and metric-block data in memory.
func (blockBuilder *blockBuilder) addSeries(seriesID uint32, data []byte) {
	offset := blockBuilder.buf.Len()
	blockBuilder.offset.Add(int32(offset))
	blockBuilder.keys.Add(seriesID)
	_, _ = blockBuilder.buf.Write(data)
}

// appendFieldMeta builds a field-id list, min start-time and max end-time.
func (blockBuilder *blockBuilder) appendFieldMeta(fieldID uint16, fieldType field.Type, startTime, endTime int64) {
	// collect min-startTime and max-endTime of one entry.
	atomic.CompareAndSwapInt64(&blockBuilder.minStartTime, 0, startTime)
	if blockBuilder.minStartTime > startTime {
		blockBuilder.minStartTime = startTime
	}
	if blockBuilder.maxEndTime < endTime {
		blockBuilder.maxEndTime = endTime
	}
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
	blockBuilder.buf.Reset()
	blockBuilder.keys.Clear()
	blockBuilder.offset = encoding.NewDeltaBitPackingEncoder()
}

// finish writes keys, offset.
func (blockBuilder *blockBuilder) finish() error {
	// write offset
	posOfOffset := blockBuilder.buf.Len()
	offset, err := blockBuilder.offset.Bytes()
	if err != nil {
		return err
	}
	_, _ = blockBuilder.buf.Write(offset)

	// write keys
	blockBuilder.keys.RunOptimize()
	keys, err := blockBuilder.keys.MarshalBinary()
	if err != nil {
		return err
	}
	posOfKeys := blockBuilder.buf.Len()
	_, _ = blockBuilder.buf.Write(keys)

	// write fields-meta
	posOfMeta := blockBuilder.buf.Len()
	var buf [16]byte
	// write start-time, end-time
	size := binary.PutUvarint(buf[:], uint64(blockBuilder.minStartTime))
	_, _ = blockBuilder.buf.Write(buf[:size])
	size = binary.PutUvarint(buf[:], uint64(blockBuilder.maxEndTime))
	_, _ = blockBuilder.buf.Write(buf[:size])
	// write fields count
	size = binary.PutUvarint(buf[:], uint64(len(blockBuilder.metaFieldsID)))
	blockBuilder.buf.Write(buf[:size])
	// write field-id, field-type list
	for _, fieldID := range blockBuilder.metaFieldsID {
		// write field-id
		binary.BigEndian.PutUint16(buf[:], fieldID)
		// write field-type
		fieldType := blockBuilder.metaFieldsType[fieldID]
		buf[2] = byte(fieldType)
		blockBuilder.buf.Write(buf[:3])
	}
	// write footer, length: 4+4+4+4
	binary.BigEndian.PutUint32(buf[:4], uint32(posOfOffset))
	binary.BigEndian.PutUint32(buf[4:8], uint32(posOfKeys))
	binary.BigEndian.PutUint32(buf[8:12], uint32(posOfMeta))
	_, _ = blockBuilder.buf.Write(buf[:12])
	// write crc32
	h := crc32.ChecksumIEEE(blockBuilder.buf.Bytes())
	binary.BigEndian.PutUint32(buf[12:16], h)
	_, _ = blockBuilder.buf.Write(buf[12:16])

	return nil
}

// bytes returns a slice of the underlying data written before.
func (blockBuilder *blockBuilder) bytes() []byte {
	return blockBuilder.buf.Bytes()
}

// Level4: builder
// entryBuilder builds a series containing different fields.
type entryBuilder struct {
	dataBuf      bytes.Buffer
	lenBuf       bytes.Buffer
	fieldsData   map[uint16][]byte
	bitArray     *bitArray
	minStartTime int64 // startTime of fields-meta
	maxEndTime   int64 // endTime of fields-meta
}

// newSeriesEntryBuilder returns a new entryBuilder, default first 2 byte is the column count.
func newSeriesEntryBuilder() *entryBuilder {
	return &entryBuilder{
		bitArray:   &bitArray{},
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
	entryBuilder.dataBuf.Reset()
	entryBuilder.lenBuf.Reset()

	if len(metaFieldsID) == 0 {
		return nil
	}

	var buf [8]byte
	var existedFieldsID []uint16
	// write start-time
	size1 := binary.PutUvarint(buf[:], uint64(entryBuilder.minStartTime))
	entryBuilder.dataBuf.Write(buf[:size1])
	// write end-time
	size2 := binary.PutUvarint(buf[:], uint64(entryBuilder.maxEndTime))
	entryBuilder.dataBuf.Write(buf[:size2])
	// build bit-array
	for idx, fieldID := range metaFieldsID {
		if _, ok := entryBuilder.fieldsData[fieldID]; !ok {
			continue
		}
		existedFieldsID = append(existedFieldsID, fieldID)
		entryBuilder.bitArray.setBit(uint16(idx))
	}
	// write bit-array length
	size3 := binary.PutUvarint(buf[:], uint64(entryBuilder.bitArray.getLen()))
	entryBuilder.dataBuf.Write(buf[:size3])
	// write bit-array
	entryBuilder.dataBuf.Write(entryBuilder.bitArray.payload)
	// write variant length in order of fields in fields-meta
	for _, fieldID := range existedFieldsID {
		theData := entryBuilder.fieldsData[fieldID]
		size4 := binary.PutUvarint(buf[:], uint64(len(theData)))
		entryBuilder.lenBuf.Write(buf[:size4])
	}
	entryBuilder.dataBuf.Write(entryBuilder.lenBuf.Bytes())
	// write data in order of fields in fields-meta
	for _, fieldID := range existedFieldsID {
		theData := entryBuilder.fieldsData[fieldID]
		entryBuilder.dataBuf.Write(theData)
	}
	return entryBuilder.dataBuf.Bytes()
}

// reset resets the inner buffer and time-range.
func (entryBuilder *entryBuilder) reset() {
	entryBuilder.dataBuf.Reset()
	entryBuilder.bitArray.reset()
	entryBuilder.fieldsData = make(map[uint16][]byte)
	entryBuilder.minStartTime = 0
	entryBuilder.maxEndTime = 0
}
