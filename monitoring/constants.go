package monitoring

import "github.com/prometheus/client_golang/prometheus"

// DefaultHistogramBuckets represents default prometheus histogram buckets in LinDB
var DefaultHistogramBuckets = []float64{
	0.0, 10.0, 25.0, 50.0, 75.0,
	100.0, 200.0, 300.0, 400.0, 500.0, 600.0, 800.0,
	1000.0, 2000.0, 5000.0,
}

var (
	// StorageRegistry/StorageGatherer represents prometheus metric registerer/ and gatherer in storage side
	storageRegistry                       = prometheus.NewRegistry()
	StorageRegistry prometheus.Registerer = storageRegistry
	StorageGatherer prometheus.Gatherer   = storageRegistry

	// BrokerRegistry/BrokerGatherer represents prometheus metric registerer/ and gatherer in broker side
	brokerRegistry                       = prometheus.NewRegistry()
	BrokerRegistry prometheus.Registerer = brokerRegistry
	BrokerGatherer prometheus.Gatherer   = brokerRegistry
)
