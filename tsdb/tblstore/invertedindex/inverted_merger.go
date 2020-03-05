package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var SeriesInvertedMerger kv.MergerType = "SeriesInvertedMerger"

// init registers series inverted merger create function
func init() {
	kv.RegisterMerger(SeriesInvertedMerger, NewInvertedMerger)
}

// invertedMerger implements kv.Merger for merging inverted index data for each tag key
type invertedMerger struct {
	invertedFlusher InvertedFlusher
	flusher         *kv.NopFlusher
}

// NewInvertedMerger creates a inverted merger
func NewInvertedMerger() kv.Merger {
	flusher := kv.NewNopFlusher()
	return &invertedMerger{
		flusher:         flusher,
		invertedFlusher: NewInvertedFlusher(flusher),
	}
}

func (m *invertedMerger) Init(params map[string]interface{}) {
	// do nothing
}

// Merge merges the multi inverted index data into a inverted index for same tag key id
func (m *invertedMerger) Merge(key uint32, values [][]byte) ([]byte, error) {
	var scanners []*tagInvertedScanner
	targetTagValueIDs := roaring.New() // target merged tag value ids
	// 1. prepare tag inverted scanner
	for _, value := range values {
		reader, err := newTagInvertedReader(value)
		if err != nil {
			return nil, err
		}
		targetTagValueIDs.Or(reader.keys)
		scanners = append(scanners, newTagInvertedScanner(reader))
	}

	// 2. merge inverted index by roaring container
	highKeys := targetTagValueIDs.GetHighKeys()
	seriesIDs := roaring.New()
	for idx, highKey := range highKeys {
		container := targetTagValueIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// scan index data then merge series ids
			for _, scanner := range scanners {
				if err := scanner.scan(highKey, lowTagValueID, seriesIDs); err != nil {
					return nil, err
				}
			}

			hk := uint32(highKey) << 16
			// flush tag value id=>series ids mapping
			if err := m.invertedFlusher.
				FlushInvertedIndex(encoding.ValueWithHighLowBits(hk, lowTagValueID), seriesIDs); err != nil {
				return nil, err
			}
			seriesIDs.Clear() // clear target series ids
		}
	}
	if err := m.invertedFlusher.FlushTagKeyID(key); err != nil {
		return nil, err
	}
	return m.flusher.Bytes(), nil
}
