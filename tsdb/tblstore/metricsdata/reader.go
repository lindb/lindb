package metricsdata

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package metricsdata

// Reader implements metrics from sstable.
type Reader interface {
}
