package tree

type Query struct {
	BaseNode
	With      *With
	QueryBody QueryBody
	OrderBy   *OrderBy
	Limit     *Limit
}

func (n *Query) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *Query) HasWith() bool {
	return n.With != nil && len(n.With.Queries) > 0
}

type QuerySpecification struct {
	BaseNode
	Select  *Select
	From    Relation
	Where   Expression
	GroupBy *GroupBy
	Having  Expression
	OrderBy *OrderBy
	Limit   *Limit
}

func (n *QuerySpecification) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}
