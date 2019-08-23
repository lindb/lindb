package indextbl

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/logger"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./metrics_flusher.go -destination=./metrics_flusher_mock.go -package indextbl

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
	var variableBuf [4]byte
	// write metricIDSeq
	binary.BigEndian.PutUint32(variableBuf[:], metricIDSeq)
	data = append(data, variableBuf[:]...)
	// write tagIDSeq
	binary.BigEndian.PutUint32(variableBuf[:], tagIDSeq)
	data = append(data, variableBuf[:]...)
	return f.flusher.Add(nsID, data)
}

// Commit closes the writer, this will be called after writing all namespaces.
func (f *metricsNameIDFlusher) Commit() error { return f.flusher.Commit() }

// metricsMetaFlusher implements MetricsMetaFlusher
type metricsMetaFlusher struct {
	flusher     kv.Flusher
	valueBuf    bytes.Buffer
	tagsBuf     bytes.Buffer
	fieldBuf    bytes.Buffer
	variableBuf [8]byte
}

// NewMetricsMetaFlusher returns a new MetricsMetaFlusher
func NewMetricsMetaFlusher(flusher kv.Flusher) MetricsMetaFlusher {
	return &metricsMetaFlusher{flusher: flusher}
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
	f.tagsBuf.WriteByte(byte(len(tagKey)))
	f.tagsBuf.WriteString(tagKey)
	// write tagID
	binary.BigEndian.PutUint32(f.variableBuf[:], tagID)
	f.tagsBuf.Write(f.variableBuf[:tagIDSize])
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
	f.fieldBuf.WriteByte(byte(len(fieldName)))
	f.fieldBuf.WriteString(fieldName)
	// write fieldType
	f.fieldBuf.WriteByte(byte(fieldType))
	// write fieldID
	binary.BigEndian.PutUint16(f.variableBuf[:], fieldID)
	f.fieldBuf.Write(f.variableBuf[:fieldIDSize])
}

// FlushMetricsMeta flushes meta info above to the underlying kv table
func (f *metricsMetaFlusher) FlushMetricMeta(metricID uint32) error {
	defer func() {
		f.valueBuf.Reset()
		f.tagsBuf.Reset()
		f.fieldBuf.Reset()
	}()
	f.buildMetricMeta()
	return f.flusher.Add(metricID, f.valueBuf.Bytes())
}

// buildMetricMeta build the meta buffer
func (f *metricsMetaFlusher) buildMetricMeta() {
	// write tags meta length
	size1 := binary.PutUvarint(f.variableBuf[:], uint64(f.tagsBuf.Len()))
	f.valueBuf.Write(f.variableBuf[:size1])
	// write tags meta
	f.valueBuf.Write(f.tagsBuf.Bytes())
	// write fields meta length
	size2 := binary.PutUvarint(f.variableBuf[:], uint64(f.fieldBuf.Len()))
	f.valueBuf.Write(f.variableBuf[:size2])
	// write fields meta
	f.valueBuf.Write(f.fieldBuf.Bytes())
}

// Commit closes the writer, this will be called after writing all metric meta info.
func (f *metricsMetaFlusher) Commit() error {
	return f.flusher.Commit()
}
