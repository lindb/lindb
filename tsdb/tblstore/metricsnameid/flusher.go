package metricsnameid

import (
	"bytes"
	"compress/gzip"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package metricsnameid

// Flusher is a wrapper of kv.Builder, provides the ability to store metricNames and metricIDs to disk.
// The layout is available in `tsdb/doc.go`(Metric NameID Table)
type Flusher interface {
	// FlushNameID flushes a mapping from metricName to metricID
	FlushNameID(metricName string, metricID uint32)
	// FlushMetricsNS flushes a mapping relation of metric-name and metric-ID of a namespace to kv table.
	// NameSpace is a concept for multi-tenancy.
	FlushMetricsNS(nsID uint32, metricIDSeq, tagKeyIDSeq uint32) error
	// Commit closes the writer, this will be called after writing all namespaces.
	Commit() error
}

// flusher implements Flusher
type flusher struct {
	kvFlusher kv.Flusher
	sBuf      bytes.Buffer // stream
	sw        *stream.BufferWriter
	gBuf      bytes.Buffer // gzip
	gw        *gzip.Writer
}

// NewFlusher returns a new MetricsNameIDFlusher
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	f := &flusher{
		kvFlusher: kvFlusher,
		sw:        stream.NewBufferWriter(nil)}
	f.sw = stream.NewBufferWriter(&f.sBuf)
	f.gw, _ = gzip.NewWriterLevel(&f.gBuf, gzip.BestSpeed)
	return f
}

// FlushNameID flushes a mapping from metricName to metricID
func (f *flusher) FlushNameID(metricName string, metricID uint32) {
	// write metricName length
	f.sw.PutUvarint64(uint64(len(metricName)))
	// write metricName
	f.sw.PutBytes([]byte(metricName))
	// write metricID
	f.sw.PutUint32(metricID)
}

// FlushMetricsNS flushes a mapping relation of metric-name and metric-ID to the underlying kv table.
func (f *flusher) FlushMetricsNS(nsID uint32, metricIDSeq, tagKeyIDSeq uint32) error {
	defer f.Reset()
	unCompressed, _ := f.sw.Bytes()
	_, _ = f.gw.Write(unCompressed)
	if err := f.gw.Close(); err != nil {
		return err
	}
	// switch to write to gzipBuffer
	f.sw.SwitchBuffer(&f.gBuf)
	// write metricIDSeq
	f.sw.PutUint32(metricIDSeq)
	// write tagKeyIDSeq
	f.sw.PutUint32(tagKeyIDSeq)
	// write back to stream buffer
	f.sw.SwitchBuffer(&f.sBuf)
	return f.kvFlusher.Add(nsID, f.gBuf.Bytes())
}

// Reset resets the buffer for flushing next name-space.
func (f *flusher) Reset() {
	f.sw.Reset()
	f.gBuf.Reset()
	f.gw.Reset(&f.gBuf)
}

// Commit closes the writer, this will be called after writing all namespaces.
func (f *flusher) Commit() error { return f.kvFlusher.Commit() }
