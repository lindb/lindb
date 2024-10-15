package utils

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type QueryBuilder struct {
	idAllocator *tree.NodeIDAllocator
}

func NewQueryBuilder(idAllocator *tree.NodeIDAllocator) *QueryBuilder {
	return &QueryBuilder{
		idAllocator: idAllocator,
	}
}

func (b *QueryBuilder) SingleValueQuery(columnName, value string) *tree.Query {
	page := types.NewPage()
	column := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{Name: columnName, DataType: types.DTString}, column)
	column.AppendString(value)
	body := &tree.QuerySpecification{
		Select: b.SelectAll(),
		From:   b.Values(page),
	}
	b.setNodeID(body)
	query := b.Query(body)
	return query
}

func (b *QueryBuilder) SimpleQuery(selectItem *tree.Select, from tree.Relation, where tree.Expression) *tree.Query {
	body := &tree.QuerySpecification{
		Select: selectItem,
		From:   from,
		Where:  where,
	}
	b.setNodeID(body)
	query := b.Query(body)
	return query
}

func (b *QueryBuilder) Query(body tree.QueryBody) *tree.Query {
	query := &tree.Query{
		QueryBody: body,
	}
	b.setNodeID(query)
	return query
}

func (b *QueryBuilder) SelectItems(columns ...string) *tree.Select {
	return &tree.Select{
		SelectItems: lo.Map(columns, func(item string, index int) tree.SelectItem {
			return b.SingleSelect(item)
		}),
	}
}

func (b *QueryBuilder) SingleSelect(columnName string) *tree.SingleColumn {
	singleColumn := &tree.SingleColumn{
		Expression: b.Identifier(columnName),
	}
	b.setNodeID(singleColumn)
	return singleColumn
}

func (b *QueryBuilder) Table(database, namespace, tableName string) *tree.Table {
	table := &tree.Table{
		Name: tree.NewQualifiedName([]*tree.Identifier{
			b.Identifier(database),
			b.Identifier(namespace),
			b.Identifier(tableName),
		}),
	}
	b.setNodeID(table)
	return table
}

func (b *QueryBuilder) SelectAll() *tree.Select {
	all := &tree.AllColumns{}
	b.setNodeID(all)
	return &tree.Select{
		SelectItems: []tree.SelectItem{all},
	}
}

func (b *QueryBuilder) Values(rows *types.Page) *tree.Values {
	values := &tree.Values{
		Rows: rows,
	}
	b.setNodeID(values)
	return values
}

func (b *QueryBuilder) LogicalAnd(term ...tree.Expression) *tree.LogicalExpression {
	logical := &tree.LogicalExpression{Terms: term, Operator: tree.LogicalAND}
	b.setNodeID(logical)
	return logical
}

func (b *QueryBuilder) StringEqual(column, value string) tree.Expression {
	return b.Equal(b.Identifier(column), b.StringLiteral(value))
}

func (b *QueryBuilder) Like(column, pattern string) *tree.LikePredicate {
	like := &tree.LikePredicate{
		Value:   b.Identifier(column),
		Pattern: b.StringLiteral(pattern),
	}
	b.setNodeID(like)
	return like
}

func (b *QueryBuilder) Equal(left, right tree.Expression) tree.Expression {
	expr := &tree.ComparisonExpression{
		Left:     left,
		Right:    right,
		Operator: tree.ComparisonEQ,
	}
	b.setNodeID(expr)
	return expr
}

func (b *QueryBuilder) Identifier(value string) *tree.Identifier {
	ident := &tree.Identifier{
		Value: value,
	}
	b.setNodeID(ident)
	return ident
}

func (b *QueryBuilder) StringLiteral(value string) *tree.StringLiteral {
	ident := &tree.StringLiteral{
		Value: value,
	}
	b.setNodeID(ident)
	return ident
}

func (b *QueryBuilder) setNodeID(node tree.Node) {
	node.SetID(b.idAllocator.Next())
}
