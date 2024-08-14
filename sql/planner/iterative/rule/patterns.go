package rule

import (
	"reflect"

	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/plan"
)

func project() *matching.Pattern {
	return matching.TypeOf(reflect.TypeFor[*plan.ProjectionNode]())
}

func output() *matching.Pattern {
	return matching.TypeOf(reflect.TypeFor[*plan.OutputNode]())
}

func aggregation() *matching.Pattern {
	return matching.TypeOf(reflect.TypeFor[*plan.AggregationNode]())
}

func exchange() *matching.Pattern {
	return matching.TypeOf(reflect.TypeFor[*plan.ExchangeNode]())
}
