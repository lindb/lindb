package parallel

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

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
	Execute() <-chan *series.TimeSeriesEvent
	// Statement returns the query statement
	Statement() *stmt.Query

	// Error returns the execution error
	Error() error
}
