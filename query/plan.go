package query

import "github.com/eleme/lindb/query/aggregation"

type Plan interface {
	Plan() *aggregation.AggregatorStreamSpec
}
