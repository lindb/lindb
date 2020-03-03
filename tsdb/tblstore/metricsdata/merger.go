package metricsdata

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

// merger implements kv.Merger for merging series data for each metric
type merger struct {
	dataFlusher  Flusher
	flusher      *kv.NopFlusher
	seriesMerger SeriesMerger
}

// NewMerger creates a metric data merger
func NewMerger() kv.Merger {
	flusher := kv.NewNopFlusher()
	dataFlusher := NewFlusher(flusher)
	return &merger{
		flusher:      flusher,
		dataFlusher:  dataFlusher,
		seriesMerger: newSeriesMerger(dataFlusher),
	}
}

// Merge merges the multi metric data into one target metric data for same metric id
func (m *merger) Merge(key uint32, values [][]byte) ([]byte, error) {
	blockCount := len(values)
	scanners := make([]*dataScanner, blockCount)
	seriesIDs := roaring.New()    // target series ids
	targetFields := field.Metas{} // target fields
	targetStart := uint16(0)
	targetEnd := uint16(0)
	// 1. prepare readers and metric level data(field/time slot/series ids)
	for idx, value := range values {
		reader, err := NewReader(value)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(reader.GetSeriesIDs())
		// get target slot range(start/end)
		start, end := reader.GetTimeRange()
		if len(targetFields) == 0 {
			targetStart = start
			targetEnd = end
		} else {
			if targetStart > start {
				targetStart = start
			}
			if targetEnd < end {
				targetEnd = end
			}
		}
		// merge target fields under metric level
		for _, f := range reader.GetFields() {
			_, ok := targetFields.GetFromID(f.ID)
			if !ok {
				targetFields = targetFields.Insert(f)
			}
		}
		// create data scanner
		scanners[idx] = newDataScanner(reader)
	}
	// 2. sort by field id
	sort.Slice(targetFields, func(i, j int) bool { return targetFields[i].ID < targetFields[j].ID })
	// 3. flush fields
	m.dataFlusher.FlushFieldMetas(targetFields)
	// 3. merge series data by roaring container
	highKeys := seriesIDs.GetHighKeys()
	decodeStreams := make([]*encoding.TSDDecoder, blockCount) // make decodeStreams for reuse
	defer func() {
		for _, stream := range decodeStreams {
			encoding.ReleaseTSDDecoder(stream)
		}
	}()
	encodeStream := encoding.TSDEncodeFunc(targetStart)
	fieldReaders := make([]FieldReader, blockCount)
	for idx, highKey := range highKeys {
		container := seriesIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowSeriesID := it.Next()
			for blockIdx, scanner := range scanners {
				seriesPos := scanner.scan(highKey, lowSeriesID)
				if seriesPos >= 0 {
					start, end := scanner.slotRange()
					if fieldReaders[blockIdx] == nil {
						fieldReaders[blockIdx] = newFieldReader(values[blockIdx], seriesPos, start, end)
					} else {
						fieldReaders[blockIdx].reset(values[blockIdx], seriesPos, start, end)
					}
				}
			}

			if err := m.seriesMerger.merge(targetFields, decodeStreams, encodeStream, fieldReaders, targetStart, targetEnd); err != nil {
				return nil, err
			}
			// flush series id
			hk := uint32(highKey) << 16
			m.dataFlusher.FlushSeries(encoding.ValueWithHighLowBits(hk, lowSeriesID))
		}
	}
	// flush metric data
	if err := m.dataFlusher.FlushMetric(key, targetStart, targetEnd); err != nil {
		return nil, err
	}
	return m.flusher.Bytes(), nil
}
