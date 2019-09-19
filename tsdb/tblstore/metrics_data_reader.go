package tblstore

//go:generate mockgen -source ./metrics_data_reader.go -destination=./metrics_data_reader_mock.go -package tblstore

// MetricsDataReader implements metrics from sstable.
type MetricsDataReader interface {
}
