package analyzer

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ExpressionAnalysis struct {
	expressionTypes map[tree.Expression]types.Type
}

func NewExpressionAnalysis() *ExpressionAnalysis {
	return &ExpressionAnalysis{
		expressionTypes: make(map[tree.Expression]types.Type),
	}
}

func (ea *ExpressionAnalysis) GetType(expression tree.Expression) types.Type {
	return nil
}
