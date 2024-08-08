package plan

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("plan.PlanNode", &encoding.JSONEncoder[PlanNode]{})
	jsoniter.RegisterTypeDecoder("plan.PlanNode", &encoding.JSONDecoder[PlanNode]{})

	encoding.RegisterNodeType(AggregationNode{})
	encoding.RegisterNodeType(ExchangeNode{})
	encoding.RegisterNodeType(FilterNode{})
	encoding.RegisterNodeType(JoinNode{})
	encoding.RegisterNodeType(ProjectionNode{})
	encoding.RegisterNodeType(TableScanNode{})
	encoding.RegisterNodeType(OutputNode{})
	encoding.RegisterNodeType(Symbol{})
}
