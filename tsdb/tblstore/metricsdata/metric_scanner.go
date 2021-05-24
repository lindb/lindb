package metricsdata

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
)

// metricScanner implements flow.Scanner interface that scans metric data from file storage.
type metricScanner struct {
	reader        Reader
	lowContainer  roaring.Container
	seriesOffsets *encoding.FixedOffsetDecoder
}

// newMetricScanner creates a file storage metric scanner.
func newMetricScanner(reader Reader,
	lowContainer roaring.Container,
	seriesOffsets *encoding.FixedOffsetDecoder,
) flow.Scanner {
	return &metricScanner{
		reader:        reader,
		lowContainer:  lowContainer,
		seriesOffsets: seriesOffsets,
	}
}

// Scan scans the metric data by given series id from file storage.
func (s *metricScanner) Scan(lowSeriesID uint16) [][]byte {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return nil
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values
	seriesPos, _ := s.seriesOffsets.Get(idx - 1)
	// read series data of fields
	return s.reader.readSeriesData(seriesPos)
}
