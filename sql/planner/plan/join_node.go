package plan

import "github.com/lindb/lindb/sql/tree"

type JoinType string
type DistributionType string

var (
	Inner JoinType = "InnerJoin"
	Left  JoinType = "LeftJoin"
	Right JoinType = "RightJoin"
	Full  JoinType = "FullJoin"

	Partitioned DistributionType = "Partitioned"
)

type EqualJoinCriteria struct {
	Left  *Symbol
	Right *Symbol
}

func (n *EqualJoinCriteria) ToExpression() *tree.ComparisonExpression {
	return &tree.ComparisonExpression{
		Left:     n.Left.ToSymbolReference(),
		Operator: tree.ComparisonEqual,
		Right:    n.Right.ToSymbolReference(),
	}
}

type JoinNode struct {
	BaseNode
	DistributionType DistributionType
	Type             JoinType
	Left             PlanNode
	Right            PlanNode

	Criteria []*EqualJoinCriteria
}

func (n *JoinNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *JoinNode) GetSources() []PlanNode {
	return []PlanNode{n.Left, n.Right}
}

func (n *JoinNode) GetOutputSymbols() []*Symbol {
	return nil
}

func (n *JoinNode) IsCrossJoin() bool {
	return n.Criteria != nil && n.Type == Inner //FIXME: check filter?????
}

func (n *JoinNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &JoinNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		DistributionType: n.DistributionType,
		Type:             n.Type,
		Criteria:         n.Criteria,
		Left:             newChildren[0],
		Right:            newChildren[1],
	}
}

func JoinTypeConvert(joinType tree.JoinType) JoinType {
	switch joinType {
	case tree.CROSS, tree.IMPLICIT, tree.INNER:
		return Inner
	case tree.LEFT:
		return Left
	case tree.RIGHT:
		return Right
	case tree.FULL:
		return Full
	default:
		panic("unsupport join type")
	}
}
