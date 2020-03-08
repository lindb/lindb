package constants

import "errors"

var (
	ErrDatabaseNotFound = errors.New("database not found")
	ErrShardNotFound    = errors.New("shard not found")

	// ErrNotFound represents the data not found
	ErrNotFound = errors.New("not found")

	// ErrNilMetric represents write nil metric error
	ErrNilMetric = errors.New("metric is nil")
	// ErrEmptyMetricName represents metric name is empty when write data
	ErrEmptyMetricName = errors.New("metric name is empty")
	// ErrEmptyField represents field is empty when write data
	ErrEmptyField = errors.New("field is empty")

	// ErrDataFileCorruption represents data in tsdb's file is corrupted
	ErrDataFileCorruption = errors.New("data corruption")
)
