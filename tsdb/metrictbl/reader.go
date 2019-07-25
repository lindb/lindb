package metrictbl

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package metrictbl

// TableReader reads metrics from sstable.
type TableReader interface {
	ContainsMetric(metricID uint32, startTime, endTime int) bool
	ContainsTSEntry(uint32) bool
}

// todo: @codingcrush, implement this
