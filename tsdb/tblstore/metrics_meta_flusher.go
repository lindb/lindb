package tblstore

import (
	"math"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./metrics_meta_flusher.go -destination=./metrics_meta_flusher_mock.go -package tblstore

var (
	metaFlusherLogger = logger.GetLogger("tsdb", "MetricsMetaFlusher")
)

// MetricsMetaFlusher is a wrapper of kv.Builder, provides the ability to store meta info of a metricID.
// The layout is available in `tsdb/doc.go`(Metric Meta Table)
type MetricsMetaFlusher interface {
	// FlushTagMeta flushes the relation of tagKey and tagID to buffer
	FlushTagMeta(tagMeta tag.Meta)
	// FlushFieldMeta flushes the relation of fieldName and fieldID to buffer
	// make sure tagKey are flushed before
	FlushFieldMeta(fieldMeta field.Meta)
	// FlushMetricsMeta flushes meta info above to the underlying kv table
	FlushMetricMeta(metricID uint32) error
	// Commit closes the writer, this will be called after writing all metric meta info.
	Commit() error
}

// metricsMetaFlusher implements MetricsMetaFlusher
type metricsMetaFlusher struct {
	flusher      kv.Flusher
	writer       *stream.BufferWriter
	fieldMetaPos int
}

// NewMetricsMetaFlusher returns a new MetricsMetaFlusher
func NewMetricsMetaFlusher(flusher kv.Flusher) MetricsMetaFlusher {
	return &metricsMetaFlusher{
		flusher: flusher,
		writer:  stream.NewBufferWriter(nil)}
}

// FlushTagMeta flushes the relation of tagKey and tagID to buffer
func (f *metricsMetaFlusher) FlushTagMeta(tagMeta tag.Meta) {
	if tagMeta.Key == "" {
		return
	}
	if len(tagMeta.Key) > math.MaxUint8 {
		metaFlusherLogger.Error("tagKey too long", zap.Int("length", len(tagMeta.Key)))
	}
	// write tagKey
	f.writer.PutByte(byte(len(tagMeta.Key)))
	f.writer.PutBytes([]byte(tagMeta.Key))
	// write tagKeyID
	f.writer.PutUint32(tagMeta.ID)

	f.fieldMetaPos = f.writer.Len()
}

// FlushFieldMeta flushes the relation of fieldName and fieldID to buffer
func (f *metricsMetaFlusher) FlushFieldMeta(fieldMeta field.Meta) {
	if fieldMeta.Name == "" {
		return
	}
	if len(fieldMeta.Name) > math.MaxUint8 {
		metaFlusherLogger.Error("fieldName too long", zap.Int("length", len(fieldMeta.Name)))
	}
	// write fieldID
	f.writer.PutUInt16(fieldMeta.ID)
	// write fieldType
	f.writer.PutByte(byte(fieldMeta.Type))
	// write field-name
	f.writer.PutUvarint64(uint64(len(fieldMeta.Name)))
	f.writer.PutBytes([]byte(fieldMeta.Name))
}

// FlushMetricsMeta flushes meta info above to the underlying kv table
func (f *metricsMetaFlusher) FlushMetricMeta(metricID uint32) error {
	defer f.Reset()
	// write pos of field-meta
	f.writer.PutUint32(uint32(f.fieldMetaPos))
	data, _ := f.writer.Bytes()
	return f.flusher.Add(metricID, data)
}

// Commit closes the writer, this will be called after writing all metric meta info.
func (f *metricsMetaFlusher) Commit() error {
	return f.flusher.Commit()
}

// Reset resets the writers
func (f *metricsMetaFlusher) Reset() {
	f.writer.Reset()
	f.fieldMetaPos = 0
}
