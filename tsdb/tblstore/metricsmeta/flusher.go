package metricsmeta

import (
	"math"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsmeta

var (
	metaFlusherLogger = logger.GetLogger("tsdb", "MetricsMetaFlusher")
)

// Flusher is a wrapper of kv.Builder, provides the ability to store meta info of a metricID.
// The layout is available in `tsdb/doc.go`(Metric Meta Table)
type Flusher interface {
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

// flusher implements Flusher
type flusher struct {
	kvFlusher    kv.Flusher
	writer       *stream.BufferWriter
	fieldMetaPos int
}

// NewFlusher returns a new MetricsMetaFlusher
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	return &flusher{
		kvFlusher: kvFlusher,
		writer:    stream.NewBufferWriter(nil)}
}

// FlushTagMeta flushes the relation of tagKey and tagID to buffer
func (f *flusher) FlushTagMeta(tagMeta tag.Meta) {
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
func (f *flusher) FlushFieldMeta(fieldMeta field.Meta) {
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
func (f *flusher) FlushMetricMeta(metricID uint32) error {
	defer f.Reset()
	// write pos of field-meta
	f.writer.PutUint32(uint32(f.fieldMetaPos))
	data, _ := f.writer.Bytes()
	return f.kvFlusher.Add(metricID, data)
}

// Commit closes the writer, this will be called after writing all metric meta info.
func (f *flusher) Commit() error {
	return f.kvFlusher.Commit()
}

// Reset resets the writers
func (f *flusher) Reset() {
	f.writer.Reset()
	f.fieldMetaPos = 0
}
