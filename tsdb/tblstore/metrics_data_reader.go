package tblstore

//go:generate mockgen -source ./metrics_data_reader.go -destination=./metrics_data_reader_mock.go -package tblstore

// MetricsDataReader reads metrics from sstable.
type MetricsDataReader interface {
	ContainsMetric(metricID uint32, startTime, endTime int) bool
	ContainsTSEntry(uint32) bool
}

// todo: @codingcrush, implement this
