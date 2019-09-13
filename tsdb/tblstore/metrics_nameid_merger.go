package tblstore

import (
	"fmt"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/stream"
)

type metricsNameIDMerger struct {
	sr           *stream.Reader
	reader       *metricsNameIDReader
	flusher      *metricsNameIDFlusher
	nopKVFlusher *kv.NopFlusher
}

// NewMetricsNameIDMerger returns a merger to compact MetricNameIDTable
func NewMetricsNameIDMerger() kv.Merger {
	m := &metricsNameIDMerger{
		sr:           stream.NewReader(nil),
		nopKVFlusher: kv.NewNopFlusher(),
		reader:       NewMetricsNameIDReader(nil).(*metricsNameIDReader)}
	m.flusher = NewMetricsNameIDFlusher(m.nopKVFlusher).(*metricsNameIDFlusher)
	return m
}

func (m *metricsNameIDMerger) Merge(key uint32, value [][]byte) ([]byte, error) {
	defer m.flusher.Reset()

	var (
		contents       [][]byte
		maxMetricIDSeq uint32
		maxTagKeyIDSeq uint32
	)
	for _, block := range value {
		content, metricIDSeq, tagKeyIDSeq, thisOK := m.reader.ReadBlock(block)
		if !thisOK {
			return nil, fmt.Errorf("failed parsing block")
		}
		contents = append(contents, content)
		if metricIDSeq > maxMetricIDSeq {
			maxMetricIDSeq = metricIDSeq
		}
		if tagKeyIDSeq > maxTagKeyIDSeq {
			maxTagKeyIDSeq = tagKeyIDSeq
		}
	}
	if len(contents) == 0 {
		return nil, fmt.Errorf("no available blocks for compacting")
	}

	for _, content := range contents {
		decoded, err := m.reader.DeCompress(content)
		if err != nil {
			return nil, err
		}
		m.sr.Reset(decoded)
		for !m.sr.Empty() {
			// read length of metricName
			size := m.sr.ReadUvarint64()
			metricName := m.sr.ReadSlice(int(size))
			metricID := m.sr.ReadUint32()
			if m.sr.Error() != nil {
				return nil, m.sr.Error()
			}
			m.flusher.FlushNameID(string(metricName), metricID)
		}
	}
	_ = m.flusher.FlushMetricsNS(key, maxMetricIDSeq, maxTagKeyIDSeq)
	return m.nopKVFlusher.Bytes(), nil
}
