package tree

type JoinType string

var (
	CROSS    JoinType = "CROSS"
	IMPLICIT JoinType = "IMPLICIT"
	INNER    JoinType = "INNER"
	LEFT     JoinType = "LEFT"
	RIGHT    JoinType = "RIGHT"
	FULL     JoinType = "FULL"
)

type Join struct {
	BaseNode
	Type     JoinType
	Left     Relation
	Right    Relation
	Criteria JoinCriteria
}

type JoinCriteria interface{}

type JoinUsing struct {
	Columns []*Identifier
}

type JoinOn struct {
	Expression Expression
}

func (n *Join) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}
