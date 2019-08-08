package indextbl

import (
	art "github.com/plar/go-adaptive-radix-tree"
)

//go:generate mockgen -source ./metrics_flusher.go -destination=./metrics_flusher_mock.go -package indextbl

// MetricsTreeFlusher is a wrapper of kv.Builder, provides the ability to store radix-tree of metrics to disk.
// The layout is available in `tsdb/doc.go`(Metric Tree Table)
type MetricsTreeFlusher interface {
	FlushMetricsTree(radixTreeID uint32, artTree art.Tree)
	// Commit closes the writer, this will be called after writing all metricKeys.
	Commit() error
}

// MetricsMetaFlusher is a wrapper of kv.Builder, provides the ability to store meta info of a metricID.
// The layout is available in `tsdb/doc.go`(Metric Meta Table)
type MetricsMetaFlusher interface {
	FlushMetricsMeta(metricID uint32, metricInfo []byte)
	// Commit closes the writer, this will be called after writing all metric meta info.
	Commit() error
}
