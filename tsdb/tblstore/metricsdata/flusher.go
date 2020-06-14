package metricsdata

import (
	"hash/crc32"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsdata

// Flusher is a wrapper of kv.Builder, provides ability to flush a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: series entry
// Level3: compressed field data
//
// flush step:
// 1. flush field metas of metric level
// 2. flush field store of one series
// 3. flush series id
// 4. flush metric data include field metadata and all series ids data
type Flusher interface {
	// FlushFieldMetas writes the meta info a field
	FlushFieldMetas(fieldMetas field.Metas)
	// FlushField writes a compressed field data to writer.
	FlushField(data []byte)
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries(seriesID uint32)
	// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
	FlushMetric(metricID uint32, start, end uint16) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error
	// GetFieldMetas returns current metric's field metas
	GetFieldMetas() field.Metas
}

// flusher implements Flusher.
type flusher struct {
	kvFlusher kv.Flusher

	writer     *stream.BufferWriter
	fieldMetas field.Metas

	// context for building field entry
	seriesWriter  *stream.BufferWriter
	fieldOffsets  *encoding.FixedOffsetEncoder
	seriesHasData bool

	seriesIDs           *roaring.Bitmap
	highOffsets         *encoding.FixedOffsetEncoder // high value of series ids
	lowOffsets          *encoding.FixedOffsetEncoder // low container of series ids
	highKey             uint16
	seriesCountOfBucket int
}

// NewFlusher returns a new Flusher,
// interval is used to calculate the time-range of field data slots.`
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	return &flusher{
		kvFlusher: kvFlusher,
		// metric block context
		writer:       stream.NewBufferWriter(nil),
		fieldOffsets: encoding.NewFixedOffsetEncoder(),

		seriesWriter: stream.NewBufferWriter(nil),

		seriesIDs:   roaring.New(),
		lowOffsets:  encoding.NewFixedOffsetEncoder(),
		highOffsets: encoding.NewFixedOffsetEncoder(),
	}
}

// FlushFieldMetas writes the field-meta of the metric
func (w *flusher) FlushFieldMetas(fieldMetas field.Metas) {
	w.fieldMetas = fieldMetas
}

// FlushField writes a compressed field data to writer.
func (w *flusher) FlushField(data []byte) {
	hasData := len(data) > 0
	w.seriesHasData = true
	if w.fieldMetas.Len() == 1 {
		if hasData {
			// if metric only has one field, just writes field data
			w.seriesWriter.PutBytes(data) // write field data
		}
	} else {
		if hasData {
			pos := w.seriesWriter.Len()   // field start position
			w.seriesWriter.PutBytes(data) // write field data
			w.fieldOffsets.Add(pos)       // add field start position
		} else {
			w.fieldOffsets.Add(encoding.EmptyOffset) // write empty offset, need align fields order with metric level
		}
	}
}

// FlushSeries writes a full series, this will be called after writing all fields of this entry.
// 1. only one field: series data = field data
// 2. mutli-fields: series data = field offsets + fields data
func (w *flusher) FlushSeries(seriesID uint32) {
	if !w.seriesHasData {
		// if not field data, needn't flush series data
		return
	}
	defer func() {
		w.seriesHasData = false
		w.seriesWriter.Reset()
		w.fieldOffsets.Reset()
	}()

	highKey := encoding.HighBits(seriesID)
	if highKey != w.highKey {
		// flush data by diff high key
		w.flushSeriesBucket()
		w.highKey = highKey // set high key, for next container storage
	}

	pos := w.writer.Len() // field data offset/field offset block start position
	if w.fieldMetas.Len() > 1 {
		// metric has mutli-fields, need write field offsets
		// write field offsets into offset block of series level
		w.writer.PutBytes(w.fieldOffsets.MarshalBinary())
	}

	// write series's field data
	seriesData, _ := w.seriesWriter.Bytes()
	w.writer.PutBytes(seriesData)

	w.lowOffsets.Add(pos) // add field offset's position

	// add series id into metric's index block
	w.seriesIDs.Add(seriesID)
	w.seriesCountOfBucket++
}

// flushSeriesBucket flushes a suit series data in one container(roaring.Bitmap)
func (w *flusher) flushSeriesBucket() {
	if w.seriesCountOfBucket == 0 {
		// if no series data in bucket, return it
		return
	}

	defer func() {
		w.lowOffsets.Reset()
		w.seriesCountOfBucket = 0
	}()

	pos := w.writer.Len() // low container's start position
	// write low offsets into offset block of high container
	w.writer.PutBytes(w.lowOffsets.MarshalBinary())
	w.highOffsets.Add(pos)
}

// reset resets the context for flushing metric block
func (w *flusher) reset() {
	w.writer.Reset()

	w.fieldOffsets.Reset()
	w.lowOffsets.Reset()
	w.highOffsets.Reset()
	w.highKey = 0
	w.seriesHasData = false
	w.seriesWriter.Reset()
	w.seriesIDs.Clear()

	w.fieldMetas = w.fieldMetas[:0]
}

// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
func (w *flusher) FlushMetric(metricID uint32, start, end uint16) error {
	defer w.reset()

	if w.seriesIDs.IsEmpty() {
		// if metric hasn't series ids
		return nil
	}

	// check if has pending series bucket not flush
	w.flushSeriesBucket()

	// write fields-meta
	fieldsMetaPos := w.writer.Len()
	// write fields count
	w.writer.PutByte(byte(len(w.fieldMetas)))
	// write field-id, field-type list
	for _, fm := range w.fieldMetas {
		// write field-id
		w.writer.PutByte(byte(fm.ID))
		// write field-type
		w.writer.PutByte(byte(fm.Type))
	}
	// write series ids bitmap
	seriesIDsBlock, err := encoding.BitmapMarshal(w.seriesIDs)
	if err != nil {
		return err
	}
	seriesIDsPos := w.writer.Len()
	w.writer.PutBytes(seriesIDsBlock)
	// write high offsets
	offsetPos := w.writer.Len()
	w.writer.PutBytes(w.highOffsets.MarshalBinary())

	//////////////////////////////////////////////////
	// build footer (field meta's offset+series ids' offset+high level offsets+crc32 checksum)
	// (2 bytes + 2 bytes +4 bytes + 4 bytes + 4 bytes + 4 bytes)
	//////////////////////////////////////////////////
	// write time range of metric level
	w.writer.PutUInt16(start)
	w.writer.PutUInt16(end)
	// write field metas' start position
	w.writer.PutUint32(uint32(fieldsMetaPos))
	// write series ids' start position
	w.writer.PutUint32(uint32(seriesIDsPos))
	// write offset block start position
	w.writer.PutUint32(uint32(offsetPos))
	// write CRC32 checksum
	data, _ := w.writer.Bytes()
	w.writer.PutUint32(crc32.ChecksumIEEE(data))
	// real flush process
	data, _ = w.writer.Bytes()
	return w.kvFlusher.Add(metricID, data)
}

// Commit adds the footer and then closes the kv builder, this will be called after writing all metric-blocks.
func (w *flusher) Commit() error {
	return w.kvFlusher.Commit()
}

// GetFieldMetas returns current metric's field metas
func (w *flusher) GetFieldMetas() field.Metas {
	return w.fieldMetas
}
