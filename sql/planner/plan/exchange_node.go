package plan

type ExchangeScope string
type ExchangeType string

var (
	Local  ExchangeScope = "Local"
	Remote ExchangeScope = "Remote"

	Gather    ExchangeType = "Gather"
	Partition ExchangeType = "Partition"
)

type ExchangeNode struct {
	Type    ExchangeType
	Scope   ExchangeScope
	Sources []PlanNode

	BaseNode
}

func (n *ExchangeNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *ExchangeNode) GetSources() []PlanNode {
	return n.Sources
}

func (n *ExchangeNode) GetOutputSymbols() []*Symbol {
	return n.Sources[0].GetOutputSymbols() //FIXME: fix it
}

func (n *ExchangeNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Type:    n.Type,
		Scope:   n.Scope,
		Sources: newChildren,
	}
}

func GatheringExchange(id PlanNodeID, scope ExchangeScope, child PlanNode) *ExchangeNode {
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Type:    Gather,
		Scope:   scope,
		Sources: []PlanNode{child},
	}
}

func PartitionedExchange(id PlanNodeID, scope ExchangeScope, child PlanNode) *ExchangeNode {
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Type:    Partition,
		Scope:   scope,
		Sources: []PlanNode{child},
	}
}
