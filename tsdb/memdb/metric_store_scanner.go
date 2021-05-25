package memdb

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"

	"github.com/lindb/roaring"
)

// metricStoreScanner implements flow.Scanner interface that scans metric data from memory storage.
type metricStoreScanner struct {
	lowContainer     roaring.Container
	timeSeriesStores []tStoreINTF
	fields           field.Metas // sort by field id
}

// newMetricStoreScanner creates a memory storage metric scanner.
func newMetricStoreScanner(lowContainer roaring.Container,
	timeSeriesStores []tStoreINTF,
	fields field.Metas,
) flow.Scanner {
	return &metricStoreScanner{
		lowContainer:     lowContainer,
		timeSeriesStores: timeSeriesStores,
		fields:           fields,
	}
}

// Scan scans the metric data by given series id from memory storage.
func (s *metricStoreScanner) Scan(lowSeriesID uint16) [][]byte {
	// check low series id if exist
	if !s.lowContainer.Contains(lowSeriesID) {
		return nil
	}
	// get the index of low series id in container
	idx := s.lowContainer.Rank(lowSeriesID)
	// scan the data and aggregate the values
	store := s.timeSeriesStores[idx-1]
	return store.scan(s.fields)
}
