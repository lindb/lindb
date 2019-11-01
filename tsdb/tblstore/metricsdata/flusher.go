package metricsdata

import (
	"hash/crc32"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsdata

// Flusher is a wrapper of kv.Builder, provides ability to flush a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: version entry
// Level3: series entry
// Level4: compressed field data
type Flusher interface {
	// FlushFieldMetas writes the meta info a field
	FlushFieldMetas([]field.Meta)
	// FlushField writes a compressed field data to writer.
	FlushField(fieldID uint16, data []byte, startSlot, endSlot int)
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries(seriesID uint32)
	// FlushVersion writes a version of the metric
	FlushVersion(version series.Version)
	// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
	FlushMetric(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error
}

// NewFlusher returns a new Flusher,
// interval is used to calculate the time-range of field data slots.`
func NewFlusher(kvFlusher kv.Flusher, interval int64) Flusher {
	return &flusher{
		interval:  interval,
		kvFlusher: kvFlusher,
		// metric block context
		writer: stream.NewBufferWriter(nil),
		// version entry context
		seriesOffsets: encoding.NewDeltaBitPackingEncoder(),
		seriesIDs:     roaring.New(),
		// series entry context
		fieldsData: make(map[uint16][]byte),
		bitArray:   collections.NewBitArray(nil)}
}

// flusher implements Flusher.
type flusher struct {
	interval  int64
	kvFlusher kv.Flusher

	writer *stream.BufferWriter
	// context for building metric block
	versionBlocks []struct {
		length  int            // length of flushed version blocks
		version series.Version // flushed version
	}
	fieldMetas []field.Meta
	// context for building version block
	versionStartPos  int   // start position of writer
	versionStartTime int64 // startTime of all data
	versionEndTime   int64 // endTime of all data
	seriesOffsets    *encoding.DeltaBitPackingEncoder
	seriesIDs        *roaring.Bitmap
	// context for building series entry
	fieldsData      map[uint16][]byte
	bitArray        *collections.BitArray
	seriesStartTime int64 // startTime of the series
	seriesEndTime   int64 // endTime of the series
}

// FlushFieldMetas writes the field-meta of the metric
func (w *flusher) FlushFieldMetas(fieldMetas []field.Meta) {
	w.fieldMetas = fieldMetas
}

// FlushField writes a compressed field data to writer.
func (w *flusher) FlushField(fieldID uint16, data []byte, startSlot, endSlot int) {
	startTime := int64(startSlot) * w.interval
	endTime := int64(endSlot) * w.interval

	// collect min-startTime and max-endTime of the the version.
	if w.versionStartTime > startTime || w.versionStartTime == 0 {
		w.versionStartTime = startTime
	}
	if w.versionEndTime < endTime {
		w.versionEndTime = endTime
	}
	// collect min-startTime and min-endTime of the series
	if w.seriesStartTime > startTime || w.seriesStartTime == 0 {
		w.seriesStartTime = startTime
	}
	if w.seriesEndTime < endTime {
		w.seriesEndTime = endTime
	}
	// record mapping of fieldID and field-data
	w.fieldsData[fieldID] = data
}

func (w *flusher) ResetSeriesContext() {
	for fieldID := range w.fieldsData {
		delete(w.fieldsData, fieldID)
	}
	w.bitArray.Reset(nil)
	w.seriesStartTime = 0
	w.seriesEndTime = 0
}

// FlushSeries writes a full series, this will be called after writing all fields of this entry.
func (w *flusher) FlushSeries(seriesID uint32) {
	defer w.ResetSeriesContext()

	seriesEntryStartPos := w.writer.Len() - w.versionStartPos
	w.seriesOffsets.Add(int32(seriesEntryStartPos))
	w.seriesIDs.Add(seriesID)

	// Fields Info Block
	// write start-time
	w.writer.PutVarint64(w.seriesStartTime)
	// write end-time
	w.writer.PutVarint64(w.seriesEndTime)
	// build and write bit-array
	for idx, fm := range w.fieldMetas {
		if _, ok := w.fieldsData[fm.ID]; !ok {
			continue
		}
		w.bitArray.SetBit(uint16(idx))
	}
	w.writer.PutBytes(w.bitArray.Bytes())
	// write data length
	for _, fm := range w.fieldMetas {
		if data, ok := w.fieldsData[fm.ID]; ok {
			w.writer.PutUvarint64(uint64(len(data)))
		}
	}

	// Fields Data Block
	// write fields data
	for _, fm := range w.fieldMetas {
		if data, ok := w.fieldsData[fm.ID]; ok {
			w.writer.PutBytes(data)
		}
	}
}

func (w *flusher) ResetVersionContext() {
	w.versionStartTime = 0
	w.versionEndTime = 0
	w.seriesOffsets.Reset()
	w.seriesIDs.Clear()
}

// FlushVersion writes a version of the metric
func (w *flusher) FlushVersion(version series.Version) {
	defer w.ResetVersionContext()

	// write series offset
	seriesOffsetPos := w.writer.Len() - w.versionStartPos
	w.writer.PutBytes(w.seriesOffsets.Bytes())

	// write series bitmap
	w.seriesIDs.RunOptimize()
	seriesBitmapPos := w.writer.Len() - w.versionStartPos
	data, _ := w.seriesIDs.MarshalBinary()
	w.writer.PutBytes(data)

	// write fields-meta
	fieldsMetaPos := w.writer.Len() - w.versionStartPos
	// write start-time, end-time of this version
	w.writer.PutVarint64(w.versionStartTime)
	w.writer.PutVarint64(w.versionEndTime)
	// write fields count
	w.writer.PutUvarint64(uint64(len(w.fieldMetas)))
	// write field-id, field-type list
	for _, fm := range w.fieldMetas {
		// write field-id
		w.writer.PutUInt16(fm.ID)
		// write field-type
		w.writer.PutByte(byte(fm.Type))
		// write field-name
		w.writer.PutUvarint64(uint64(len(fm.Name)))
		w.writer.PutBytes([]byte(fm.Name))
	}
	// write footer, length: 4+4+4
	w.writer.PutUint32(uint32(seriesOffsetPos))
	w.writer.PutUint32(uint32(seriesBitmapPos))
	w.writer.PutUint32(uint32(fieldsMetaPos))
	// record version length
	w.versionBlocks = append(w.versionBlocks, struct {
		length  int
		version series.Version
	}{
		length:  w.writer.Len() - w.versionStartPos,
		version: version,
	})
	w.versionStartPos = w.writer.Len()
}

// Reset resets the context for flushing metric block
func (w *flusher) Reset() {
	w.writer.Reset()
	w.versionBlocks = w.versionBlocks[:0]
	w.fieldMetas = w.fieldMetas[:0]
	w.versionStartPos = 0
}

// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
func (w *flusher) FlushMetric(metricID uint32) error {
	defer w.Reset()
	// no version was flushed before
	if len(w.versionBlocks) == 0 {
		return nil
	}
	//////////////////////////////////////////////////
	// build Version Offsets Block
	//////////////////////////////////////////////////
	// start position of the offsets block
	posOfVersionOffsets := w.writer.Len()
	// write versions count
	w.writer.PutUvarint64(uint64(len(w.versionBlocks)))
	// write all versions and version lengths
	for _, versionBlock := range w.versionBlocks {
		// write version
		w.writer.PutInt64(versionBlock.version.Int64())
		// write version block length
		w.writer.PutUvarint64(uint64(versionBlock.length))
	}
	//////////////////////////////////////////////////
	// build Footer
	//////////////////////////////////////////////////
	// write position of the offsets block
	w.writer.PutUint32(uint32(posOfVersionOffsets))
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
