package tree

type Relation interface {
	Node
}

type AliasedRelation struct {
	BaseNode

	Relation Relation
	Aliase   *Identifier
}

func (n *AliasedRelation) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type Table struct {
	BaseNode
	Name *QualifiedName
}

func (n *Table) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type TableSubQuery struct {
	BaseNode
	Query *Query
}

func (n *TableSubQuery) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}
