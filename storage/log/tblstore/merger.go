package tblstore

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var LogIndexMerger kv.MergerType = "LogIndexMerger"

// init registers metric data merger create function
func init() {
	kv.RegisterMerger(LogIndexMerger, NewMerger)
}

// merger implements kv.Merger for merging series data for each metric
type merger struct {
	flusher kv.Flusher

	result, temp *roaring.Bitmap
}

// NewMerger creates a metric data merger
func NewMerger(flusher kv.Flusher) (kv.Merger, error) {
	return &merger{
		flusher: flusher,
		result:  roaring.New(),
		temp:    roaring.New(),
	}, nil
}

// Init implements kv.Merger.
func (m *merger) Init(params map[string]any) {
}

// Merge implements kv.Merger.
func (m *merger) Merge(key uint32, values [][]byte) error {
	m.result.Clear()

	for _, op := range values {
		m.temp.Clear()
		encoding.BitmapUnmarshal(m.temp, op)

		// merge
		m.result.Or(m.temp)
	}

	d, err := m.result.ToBytes()
	if err != nil {
		return err
	}
	return m.flusher.Add(key, d)
}
