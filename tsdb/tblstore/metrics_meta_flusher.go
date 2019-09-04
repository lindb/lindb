package tblstore

import (
	"bytes"
	"math"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/tsdb/field"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./metrics_meta_flusher.go -destination=./metrics_meta_flusher_mock.go -package tblstore

var (
	nameIDIndexFlusherLogger = logger.GetLogger("tsdb", "IndexTableFlusher")
)

// MetricsNameIDFlusher is a wrapper of kv.Builder, provides the ability to store metricNames and metricIDs to disk.
// The layout is available in `tsdb/doc.go`(Metric NameID Table)
type MetricsNameIDFlusher interface {
	// FlushMetricsNS flushes a mapping relation of metric-name and metric-ID of a namespace to kv table.
	// NameSpace is a concept for multi-tenancy, each value is a isolated index.
	FlushMetricsNS(nsID uint32, data []byte, metricIDSeq, tagIDSeq uint32) error
	// Commit closes the writer, this will be called after writing all namespaces.
	Commit() error
}

// MetricsMetaFlusher is a wrapper of kv.Builder, provides the ability to store meta info of a metricID.
// The layout is available in `tsdb/doc.go`(Metric Meta Table)
type MetricsMetaFlusher interface {
	// FlushTagKeyID flushes the relation of tagKey and tagID to buffer
	FlushTagKeyID(tagKey string, tagID uint32)
	// FlushFieldID flushes the relation of fieldName and fieldID to buffer
	FlushFieldID(fieldName string, fieldType field.Type, fieldID uint16)
	// FlushMetricsMeta flushes meta info above to the underlying kv table
	FlushMetricMeta(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric meta info.
	Commit() error
}

// metricsNameIDFlusher implements MetricsNameIDFlusher
type metricsNameIDFlusher struct {
	flusher kv.Flusher
}

//NewMetricsNameIDFlusher returns a new MetricsNameIDFlusher
func NewMetricsNameIDFlusher(flusher kv.Flusher) MetricsNameIDFlusher {
	return &metricsNameIDFlusher{flusher: flusher}
}

// FlushMetricsNS flushes a mapping relation of metric-name and metric-ID to the underlying kv table.
func (f *metricsNameIDFlusher) FlushMetricsNS(nsID uint32, data []byte, metricIDSeq, tagIDSeq uint32) error {
	writer := stream.NewBufferWriter(bytes.NewBuffer(data))
	// write metricIDSeq
	writer.PutUint32(metricIDSeq)
	// write tagIDSeq
	writer.PutUint32(tagIDSeq)
	block, _ := writer.Bytes()
	return f.flusher.Add(nsID, block)
}

// Commit closes the writer, this will be called after writing all namespaces.
func (f *metricsNameIDFlusher) Commit() error { return f.flusher.Commit() }

// metricsMetaFlusher implements MetricsMetaFlusher
type metricsMetaFlusher struct {
	flusher        kv.Flusher
	valueBufWriter *stream.BufferWriter
	tagsBufWriter  *stream.BufferWriter
	fieldBufWriter *stream.BufferWriter
}

// NewMetricsMetaFlusher returns a new MetricsMetaFlusher
func NewMetricsMetaFlusher(flusher kv.Flusher) MetricsMetaFlusher {
	return &metricsMetaFlusher{
		flusher:        flusher,
		valueBufWriter: stream.NewBufferWriter(nil),
		tagsBufWriter:  stream.NewBufferWriter(nil),
		fieldBufWriter: stream.NewBufferWriter(nil)}
}

// FlushTagKeyID flushes the relation of tagKey and tagID to buffer
func (f *metricsMetaFlusher) FlushTagKeyID(tagKey string, tagID uint32) {
	if tagKey == "" {
		return
	}
	if len(tagKey) > math.MaxUint8 {
		nameIDIndexFlusherLogger.Error("tagKey too long", zap.Int("length", len(tagKey)))
	}
	// write tagKey
	f.tagsBufWriter.PutByte(byte(len(tagKey)))
	f.tagsBufWriter.PutBytes([]byte(tagKey))
	// write tagID
	f.tagsBufWriter.PutUint32(tagID)
}

// FlushFieldID flushes the relation of fieldName and fieldID to buffer
func (f *metricsMetaFlusher) FlushFieldID(fieldName string, fieldType field.Type, fieldID uint16) {
	if fieldName == "" {
		return
	}
	if len(fieldName) > math.MaxUint8 {
		nameIDIndexFlusherLogger.Error("fieldName too long", zap.Int("length", len(fieldName)))
	}
	// write field-name
	f.fieldBufWriter.PutByte(byte(len(fieldName)))
	f.fieldBufWriter.PutBytes([]byte(fieldName))
	// write fieldType
	f.fieldBufWriter.PutByte(byte(fieldType))
	// write fieldID
	f.fieldBufWriter.PutUInt16(fieldID)
}

// FlushMetricsMeta flushes meta info above to the underlying kv table
func (f *metricsMetaFlusher) FlushMetricMeta(metricID uint32) error {
	defer func() {
		f.valueBufWriter.Reset()
		f.tagsBufWriter.Reset()
		f.fieldBufWriter.Reset()
	}()
	f.buildMetricMeta()
	data, _ := f.valueBufWriter.Bytes()
	return f.flusher.Add(metricID, data)
}

// buildMetricMeta build the meta buffer
func (f *metricsMetaFlusher) buildMetricMeta() {
	// write tags meta length
	f.valueBufWriter.PutUvarint64(uint64(f.tagsBufWriter.Len()))
	// write tags meta
	data, _ := f.tagsBufWriter.Bytes()
	f.valueBufWriter.PutBytes(data)
	// write fields meta length
	f.valueBufWriter.PutUvarint64(uint64(f.fieldBufWriter.Len()))
	data, _ = f.fieldBufWriter.Bytes()
	// write fields meta
	f.valueBufWriter.PutBytes(data)
}

// Commit closes the writer, this will be called after writing all metric meta info.
func (f *metricsMetaFlusher) Commit() error {
	return f.flusher.Commit()
}
