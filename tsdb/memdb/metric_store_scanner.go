package memdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
)

// metricStoreScanner implements flow.Scanner interface that scans metric data from memory storage.
type metricStoreScanner struct {
	lowContainer     roaring.Container
	timeSeriesStores []tStoreINTF

	scanCtx *memScanContext
}

// newMetricStoreScanner creates a memory storage metric scanner.
func newMetricStoreScanner(lowContainer roaring.Container,
	timeSeriesStores []tStoreINTF,
	aggs []*fieldAggregator,
) flow.Scanner {
	return &metricStoreScanner{
		lowContainer:     lowContainer,
		timeSeriesStores: timeSeriesStores,
		scanCtx: &memScanContext{
			tsd:       encoding.GetTSDDecoder(),
			fieldAggs: aggs,
		},
	}
}

// Scan scans the metric data by given series id from memory storage.
func (s *metricStoreScanner) Scan(lowSeriesID uint16) {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values
	store := s.timeSeriesStores[idx-1]
	store.scan(s.scanCtx)
}

// Close closes the resource of memory storage metric scanner.
func (s *metricStoreScanner) Close() error {
	encoding.ReleaseTSDDecoder(s.scanCtx.tsd)
	return nil
}
