package parallel

import "github.com/lindb/lindb/pkg/field"

//go:generate mockgen -source=./executor.go -destination=./executor_mock.go -package=parallel

// Executor represents a query executor both storage/broker side.
// When returning query results the following is the order in which processing takes place:
// 1) filtering
// 2) Scanning
// 3) Grouping if need
// 4) Down sampling
// 5) Aggregation
// 6) Functions
// 7) Expressions
type Executor interface {
	// Execute execute query
	// 1) plan query language
	// 2) aggregator data from time series(memory/file/network)
	Execute() <-chan field.GroupedTimeSeries

	// Error returns the execution error
	Error() error
}
