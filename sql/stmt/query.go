package stmt

import (
	"github.com/eleme/lindb/pkg/timeutil"
)

// Query represents search statement
type Query struct {
	MetricName  string             // like table name
	SelectItems []Expr             // select list, such as field, function call, math expression etc.
	Condition   Expr               // tag filter condition expression
	TimeRange   timeutil.TimeRange // query time range
	Interval    int64              // down sampling interval
	GroupBy     []string           // group by
	Limit       int                // num. of time series list for result
}
