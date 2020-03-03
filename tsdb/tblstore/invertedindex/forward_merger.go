package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
)

// forwardMerger implements kv.Merger for merging forward index data for each tag key
type forwardMerger struct {
	forwardFlusher ForwardFlusher
	flusher        *kv.NopFlusher
}

// NewForwardMerger creates a forward merger
func NewForwardMerger() kv.Merger {
	flusher := kv.NewNopFlusher()
	return &forwardMerger{
		flusher:        flusher,
		forwardFlusher: NewForwardFlusher(flusher),
	}
}

// Merge merges the multi forward index data into a forward index for same tag key id
func (m *forwardMerger) Merge(key uint32, values [][]byte) ([]byte, error) {
	var scanners []*tagForwardScanner
	seriesIDs := roaring.New() // target merged series ids
	// 1. prepare tag forward scanner
	for _, value := range values {
		reader, err := NewTagForwardReader(value)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(reader.getSeriesIDs())
		scanners = append(scanners, newTagForwardScanner(reader))
	}

	// 2. merge forward index by roaring container
	highKeys := seriesIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		container := seriesIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		var tagValueIDs []uint32
		for it.HasNext() {
			lowSeriesID := it.Next()
			// scan index data then merge tag value ids, sort by series id
			for _, scanner := range scanners {
				tagValueIDs = scanner.scan(highKey, lowSeriesID, tagValueIDs)
			}
		}
		// flush tag value ids by one container
		m.forwardFlusher.FlushForwardIndex(tagValueIDs)
	}
	// flush all series ids under this tag key
	if err := m.forwardFlusher.FlushTagKeyID(key, seriesIDs); err != nil {
		return nil, err
	}
	return m.flusher.Bytes(), nil
}
