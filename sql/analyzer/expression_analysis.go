package analyzer

import (
	"github.com/lindb/lindb/spi/value"
	"github.com/lindb/lindb/sql/tree"
)

type ExpressionAnalysis struct {
	expressionTypes map[tree.Expression]value.Type
}

func NewExpressionAnalysis() *ExpressionAnalysis {
	return &ExpressionAnalysis{
		expressionTypes: make(map[tree.Expression]value.Type),
	}
}

func (ea *ExpressionAnalysis) GetType(expression tree.Expression) value.Type {
	return nil
}
