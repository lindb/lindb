package tree

import (
	"github.com/lindb/common/constants"

	"github.com/lindb/lindb/spi/types"
)

type Relation interface {
	Node
}

type Values struct {
	BaseNode

	Rows *types.Page
}

func (n *Values) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type AliasedRelation struct {
	BaseNode

	Relation    Relation
	Aliase      *Identifier
	ColumnNames []*Identifier
}

func (n *AliasedRelation) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

type Table struct {
	BaseNode
	Name *QualifiedName
}

func (n *Table) GetDatabase(defaultDB string) string {
	if len(n.Name.Parts) == 3 {
		return n.Name.Parts[0]
	}
	return defaultDB
}

func (n *Table) GetNamespace() string {
	switch len(n.Name.Parts) {
	case 3:
		return n.Name.Parts[1]
	case 2:
		return n.Name.Parts[0]
	default:
		return constants.DefaultNamespace
	}
}

func (n *Table) GetTableName() string {
	return n.Name.Parts[len(n.Name.Parts)-1]
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
