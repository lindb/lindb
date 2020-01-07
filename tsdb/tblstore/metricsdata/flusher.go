package metricsdata

import (
	"fmt"
	"hash/crc32"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsdata

// Flusher is a wrapper of kv.Builder, provides ability to flush a metric-table file to disk.
// The layout is available in `tsdb/doc.go`
// Level1: metric-block
// Level2: version entry
// Level3: series entry
// Level4: compressed field data
//
// flush step:
// 1. flush field store of one series
// 2. flush series bucket data based on one container of roaring bitmap
// 3. flush series bucket info such as series data's offsets
// 4. flush metric's version data
// 5. flush metric data include field metadata and all version data
type Flusher interface {
	// FlushFieldMetas writes the meta info a field
	FlushFieldMetas(fieldMetas []field.Meta)
	FlushPrimitiveField(pFieldID uint16, data []byte)
	// FlushField writes a compressed field data to writer.
	FlushField(fieldID uint16)
	// FlushSeries writes a full series, this will be called after writing all fields of this entry.
	FlushSeries()
	// FlushSeriesBucket writes a suit series data in one container(roaring.Bitmap),
	// this will be called after writing a suit series data.
	FlushSeriesBucket()
	// FlushVersion writes a version of the metric
	FlushVersion(version series.Version, seriesIDs *roaring.Bitmap)
	// FlushMetric writes a full metric-block, this will be called after writing all entries of this metric.
	FlushMetric(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric-blocks.
	Commit() error

	GetFieldMeta(fieldID uint16) (field.Meta, bool)
}

// NewFlusher returns a new Flusher,
// interval is used to calculate the time-range of field data slots.`
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	return &flusher{
		kvFlusher: kvFlusher,
		// metric block context
		writer: stream.NewBufferWriter(nil),
		// version entry context
		seriesOffsets:       encoding.NewFixedOffsetEncoder(),
		seriesBucketOffsets: encoding.NewFixedOffsetEncoder(),
		// series entry context
		fieldsData:            make(map[uint16]*fieldData),
		primitiveFieldsData:   make(map[uint16][]byte),
		primitiveFieldOffsets: encoding.NewFixedOffsetEncoder(),
		fieldWriter:           stream.NewBufferWriter(nil),
		bitArray:              collections.NewBitArray(nil),
		bitArray1:             collections.NewBitArray(nil),
	}
}

type fieldData struct {
	fieldIDs, offset, data []byte
}

func (f *fieldData) Len() int {
	return len(f.fieldIDs) + len(f.offset) + len(f.data)
}

// flusher implements Flusher.
type flusher struct {
	kvFlusher kv.Flusher

	writer *stream.BufferWriter
	// context for building metric block
	versionBlocks []struct {
		length  int            // length of flushed version blocks
		version series.Version // flushed version
	}
	fieldMetas field.Metas
	// context for building version block
	versionStartPos     int // start position of writer
	seriesOffsets       *encoding.FixedOffsetEncoder
	seriesBucketOffsets *encoding.FixedOffsetEncoder

	// context for building field entry
	primitiveFieldsData   map[uint16][]byte
	primitiveFieldOffsets *encoding.FixedOffsetEncoder
	fieldWriter           *stream.BufferWriter
	bitArray1             *collections.BitArray

	// context for building series entry
	fieldsData map[uint16]*fieldData
	bitArray   *collections.BitArray
}

// FlushFieldMetas writes the field-meta of the metric
func (w *flusher) FlushFieldMetas(fieldMetas []field.Meta) {
	w.fieldMetas = fieldMetas
}

func (w *flusher) FlushPrimitiveField(pFieldID uint16, data []byte) {
	// record mapping of primitive field id and field-data
	w.primitiveFieldsData[pFieldID] = data
}

func (w *flusher) resetFieldContext() {
	for fieldID := range w.primitiveFieldsData {
		delete(w.primitiveFieldsData, fieldID)
	}
	w.fieldWriter.Reset()
	w.primitiveFieldOffsets.Reset()
	w.bitArray1.Reset(nil)
}

// FlushField writes a compressed field data to writer.
func (w *flusher) FlushField(fieldID uint16) {
	defer w.resetFieldContext()
	fieldMeta, ok := w.fieldMetas.GetFromID(fieldID)
	if !ok {
		return
	}
	fieldType := fieldMeta.Type
	primitiveFields := fieldType.GetSchema().GetAllPrimitiveFields()
	switch fieldType {
	case field.SummaryField: //complex field
		for _, pFieldID := range primitiveFields {
			offset := w.fieldWriter.Len()
			data, ok := w.primitiveFieldsData[pFieldID]
			if !ok {
				continue
			}
			w.fieldWriter.PutBytes(data)
			w.primitiveFieldOffsets.Add(offset)
			w.bitArray1.SetBit(pFieldID)
		}
		// ignore err
		data, _ := w.fieldWriter.Bytes()
		w.setFieldData(fieldID, w.bitArray1.Bytes(), w.primitiveFieldOffsets.MarshalBinary(), data)
	default: //simple field
		// record mapping of fieldID and field-data
		w.setFieldData(fieldID, nil, nil, w.primitiveFieldsData[primitiveFields[0]])
	}
}

func (w *flusher) setFieldData(fieldID uint16, fieldIDs, offset, data []byte) {
	fieldValue, ok := w.fieldsData[fieldID]
	if !ok {
		fieldValue = &fieldData{}
		w.fieldsData[fieldID] = fieldValue
	}
	fieldValue.fieldIDs = fieldIDs
	fieldValue.offset = offset
	fieldValue.data = data
}

// ResetSeriesContext resets the series context for reuse
func (w *flusher) ResetSeriesContext() {
	for fieldID := range w.fieldsData {
		delete(w.fieldsData, fieldID)
	}
	w.bitArray.Reset(nil)
}

// FlushSeries writes a full series, this will be called after writing all fields of this entry.
func (w *flusher) FlushSeries() {
	defer w.ResetSeriesContext()

	seriesEntryStartPos := w.writer.Len() - w.versionStartPos
	w.seriesOffsets.Add(seriesEntryStartPos)

	// Fields Info Block
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
			w.writer.PutUvarint32(uint32(data.Len()))
		}
	}

	// Fields Data Block
	// write fields data
	for _, fm := range w.fieldMetas {
		if data, ok := w.fieldsData[fm.ID]; ok {
			if len(data.fieldIDs) > 0 {
				w.writer.PutBytes(data.fieldIDs)
			}
			if len(data.offset) > 0 {
				w.writer.PutBytes(data.offset)
			}
			if len(data.data) > 0 {
				w.writer.PutBytes(data.data)
			}
		}
	}
}

// FlushSeriesBucket flushes a suit series data in one container(roaring.Bitmap)
func (w *flusher) FlushSeriesBucket() {
	defer w.seriesOffsets.Reset()

	// write series bucket offset
	seriesBucketOffsetPos := w.writer.Len() - w.versionStartPos
	w.writer.PutBytes(w.seriesOffsets.MarshalBinary())
	w.seriesBucketOffsets.Add(seriesBucketOffsetPos)
}

// ResetVersionContext resets the version context for reuse
func (w *flusher) ResetVersionContext() {
	w.seriesBucketOffsets.Reset()
}

// FlushVersion writes a version of the metric
func (w *flusher) FlushVersion(version series.Version, seriesIDs *roaring.Bitmap) {
	defer w.ResetVersionContext()

	// write series bitmap
	seriesIDs.RunOptimize()
	seriesPos := w.writer.Len() - w.versionStartPos
	data, _ := seriesIDs.MarshalBinary()
	w.writer.PutBytes(data)
	seriesBucketPos := w.writer.Len() - w.versionStartPos
	// write bitmap container offsets
	w.writer.PutBytes(w.seriesBucketOffsets.MarshalBinary())

	// write fields-meta
	fieldsMetaPos := w.writer.Len() - w.versionStartPos
	// write fields count
	w.writer.PutUInt16(uint16(len(w.fieldMetas)))
	// write field-id, field-type list
	for _, fm := range w.fieldMetas {
		// write field-id
		w.writer.PutUInt16(fm.ID)
		// write field-type
		w.writer.PutByte(byte(fm.Type))
	}
	// write footer, length: 4+4+4
	w.writer.PutUint32(uint32(seriesPos))       // series bitmap position
	w.writer.PutUint32(uint32(seriesBucketPos)) // bucket offset position
	w.writer.PutUint32(uint32(fieldsMetaPos))   // field metadata position
	// record version length
	w.versionBlocks = append(w.versionBlocks, struct {
		length  int
		version series.Version
	}{
		length:  w.writer.Len() - w.versionStartPos,
		version: version,
	})
	// next version start pos
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
	fmt.Printf("flush data metricID:%d\n", metricID)

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

func (w *flusher) GetFieldMeta(fieldID uint16) (field.Meta, bool) {
	return w.fieldMetas.GetFromID(fieldID)
}
