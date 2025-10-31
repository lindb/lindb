package stream

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type Expr interface{}

type ComparisonExpr struct {
	column types.ColumnMetadata
	value  string
}

type InExpr struct {
	column types.ColumnMetadata
	values []string
}

type NotExpr struct {
	expr Expr
}

type LogicalExpr struct {
	op    tree.LogicalOperator
	exprs []Expr
}
