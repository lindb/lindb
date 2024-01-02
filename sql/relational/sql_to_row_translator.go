package relational

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type SQLToRowExpressionTranslator struct{}

func (t *SQLToRowExpressionTranslator) Translate(expression tree.Expression) any {
	visitor := &SQLToRowExpressionVisitor{}
	result := expression.Accept(nil, visitor)
	fmt.Println(result)
	return result
}

type SQLToRowExpressionVisitor struct {
	tree.Visitor
}

func (v *SQLToRowExpressionVisitor) Visit(context any, node tree.Node) any {
	return node.Accept(context, v)
}

func (v *SQLToRowExpressionVisitor) VisitStringLiteral(context any, node *tree.StringLiteral) (r any) {
	return nil
}
