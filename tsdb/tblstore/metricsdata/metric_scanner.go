package metricsdata

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
)

// metricScanner implements flow.Scanner interface that scans metric data from file storage.
type metricScanner struct {
	reader        Reader
	fieldAggs     []*fieldAggregator
	lowContainer  roaring.Container
	seriesOffsets *encoding.FixedOffsetDecoder

	tsd *encoding.TSDDecoder
}

// newMetricScanner creates a file storage metric scanner.
func newMetricScanner(reader Reader,
	fieldAggs []*fieldAggregator,
	lowContainer roaring.Container,
	seriesOffsets *encoding.FixedOffsetDecoder,
) flow.Scanner {
	return &metricScanner{
		reader:        reader,
		fieldAggs:     fieldAggs,
		lowContainer:  lowContainer,
		seriesOffsets: seriesOffsets,
		tsd:           encoding.GetTSDDecoder(),
	}
}

// Scan scans the metric data by given series id from file storage.
func (s *metricScanner) Scan(lowSeriesID uint16) {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values
	seriesPos, _ := s.seriesOffsets.Get(idx - 1)
	// read series data and agg it
	s.reader.readSeriesData(seriesPos, s.tsd, s.fieldAggs)
}

// Close closes the resource of file storage metric scanner.
func (s *metricScanner) Close() error {
	encoding.ReleaseTSDDecoder(s.tsd)
	return nil
}
