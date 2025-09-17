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
}

// NewMerger creates a metric data merger
func NewMerger(flusher kv.Flusher) (kv.Merger, error) {
	return &merger{
		flusher: flusher,
	}, nil
}

// Init implements kv.Merger.
func (m *merger) Init(params map[string]any) {
}

// Merge implements kv.Merger.
func (m *merger) Merge(key uint32, values [][]byte) error {
	value := roaring.New()
	value2 := roaring.New()
	for _, op := range values {
		encoding.BitmapUnmarshal(value2, op)
		value.Or(value2)
		value2.Clear()
	}
	d, err := value.ToBytes()
	if err != nil {
		return err
	}
	return m.flusher.Add(key, d)
}
