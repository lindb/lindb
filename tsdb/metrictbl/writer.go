package metrictbl

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sync/atomic"

	"github.com/eleme/lindb/kv"
	"github.com/eleme/lindb/pkg/encoding"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./writer.go -destination=./writer_mock.go -package metrictbl

// TableWriter is a wrapper of kv.Builder, provides ability to build a metric-table file to disk.
// Level1: metric-block
// Level2: TSEntry
// Level3: compressed field data
type TableWriter interface {
	// WriteField writes a compressed field data to writer.
	WriteField(fieldID uint32, data []byte, startSlot, endSlot int)
	// WriteTSEntry writes a full tsEntry, this will be called after writing all fields of this entry.
	WriteTSEntry(tsID uint32)
	// WriteMetricBlock writes a full metric-block, this will be called after writing all entries of this metric.
	WriteMetricBlock(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error
}

// NewTableWriter returns a new TableWriter, interval is used to calculate the time-range of field data slots.`
func NewTableWriter(flusher kv.Flusher, interval int64) TableWriter {
	return newTableWriter(flusher, interval)
}

// newTableWriter returns a new newTableWriter.
func newTableWriter(flusher kv.Flusher, interval int64) *tableWriter {
	return &tableWriter{
		interval:     interval,
		flusher:      flusher,
		blockBuilder: newBlockBuilder(),
		entryBuilder: newTSEntryBuilder()}
}

// tableWriter implements TableWriter.
type tableWriter struct {
	interval     int64
	flusher      kv.Flusher
	blockBuilder *blockBuilder
	entryBuilder *entryBuilder
}

// WriteField writes a compressed field data to writer.
func (w *tableWriter) WriteField(fieldID uint32, data []byte, startSlot, endSlot int) {
	startTime := int64(startSlot) * w.interval
	endTime := int64(endSlot) * w.interval

	w.blockBuilder.appendFieldMeta(fieldID, startTime, endTime)
	w.entryBuilder.addField(fieldID, data, startTime, endTime)
}

// WriteTSEntry writes a full tsEntry, this will be called after writing all fields of this entry.
func (w *tableWriter) WriteTSEntry(tsID uint32) {
	w.blockBuilder.addTSEntry(tsID, w.entryBuilder.bytes(w.blockBuilder.metaFieldsID))
	w.entryBuilder.reset()
}

// WriteMetricBlock writes a full metric-block, this will be called after writing all entries of this metric.
func (w *tableWriter) WriteMetricBlock(metricID uint32) error {
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
func (w *tableWriter) Commit() error {
	return w.flusher.Commit()
}

// blockBuilder builds a metric-block containing multi TSEntry in order.
type blockBuilder struct {
	minStartTime    int64          // startTime of fields-meta
	maxEndTime      int64          // endTime of fields-meta
	metaFieldsID    []uint32       // fieldID list of fields-meta
	metaFieldsIDMap map[uint32]int // set
	buf             bytes.Buffer
	offset          *encoding.DeltaBitPackingEncoder
	keys            *roaring.Bitmap
}

// newBlockBuilder returns a new metricBlockBuilder.
func newBlockBuilder() *blockBuilder {
	return &blockBuilder{
		metaFieldsIDMap: make(map[uint32]int),
		keys:            roaring.New(),
		offset:          encoding.NewDeltaBitPackingEncoder()}
}

// addTSEntry puts tsID and metric-block data in memory.
func (blockBuilder *blockBuilder) addTSEntry(tsID uint32, data []byte) {
	offset := blockBuilder.buf.Len()
	blockBuilder.offset.Add(int32(offset))
	blockBuilder.keys.Add(tsID)
	_, _ = blockBuilder.buf.Write(data)
}

// appendFieldMeta builds a field-id list, min start-time and max end-time.
func (blockBuilder *blockBuilder) appendFieldMeta(fieldID uint32, startTime, endTime int64) {
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
	blockBuilder.metaFieldsID = append(blockBuilder.metaFieldsID, fieldID)
}

// reset resets the buffer, keys and offset.
func (blockBuilder *blockBuilder) reset() {
	blockBuilder.minStartTime = 0
	blockBuilder.maxEndTime = 0
	blockBuilder.metaFieldsID = blockBuilder.metaFieldsID[:0]
	blockBuilder.metaFieldsIDMap = make(map[uint32]int)

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
	// write field-id list
	for _, fieldID := range blockBuilder.metaFieldsID {
		binary.BigEndian.PutUint32(buf[:], fieldID)
		blockBuilder.buf.Write(buf[:4])
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
// entryBuilder builds a tsEntry containing different fields.
type entryBuilder struct {
	dataBuf      bytes.Buffer
	lenBuf       bytes.Buffer
	fieldsData   map[uint32][]byte
	bitArray     *bitArray
	minStartTime int64 // startTime of fields-meta
	maxEndTime   int64 // endTime of fields-meta
}

// newTSEntryBuilder returns a new tSEntryBuilder, default first 2 byte is the column count.
func newTSEntryBuilder() *entryBuilder {
	return &entryBuilder{
		bitArray:   &bitArray{},
		fieldsData: make(map[uint32][]byte)}
}

// addField puts fieldName and data into buffer in memory.
func (entryBuilder *entryBuilder) addField(fieldID uint32, data []byte, startTime, endTime int64) {
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
func (entryBuilder *entryBuilder) bytes(metaFieldsID []uint32) []byte {
	entryBuilder.dataBuf.Reset()
	entryBuilder.lenBuf.Reset()

	if len(metaFieldsID) == 0 {
		return nil
	}

	var buf [8]byte
	var existedFieldsID []uint32
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
	entryBuilder.fieldsData = make(map[uint32][]byte)
	entryBuilder.minStartTime = 0
	entryBuilder.maxEndTime = 0
}
