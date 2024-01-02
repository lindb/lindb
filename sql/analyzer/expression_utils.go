package analyzer

import (
	"github.com/lindb/lindb/sql/tree"
)

func asQualifiedName(expression tree.Expression) (name *tree.QualifiedName) {
	if identifier, ok := expression.(*tree.Identifier); ok {
		name = tree.NewQualifiedName([]*tree.Identifier{identifier})
	} else if dereference, ok := expression.(*tree.DereferenceExpression); ok {
		//TODO:????
		name = dereference.ToQualifiedName()
	}
	return
}

func ExtractConjuncts(expression tree.Expression) (result []tree.Expression) {
	return ExtractPredicates(tree.LogicalAND, expression, result)
}

func ExtractPredicates(operator tree.LogicalOperator, expression tree.Expression, result []tree.Expression) (r []tree.Expression) {
	if logicalExpression, ok := expression.(*tree.LogicalExpression); ok && logicalExpression.Operator == operator {
		for i := range logicalExpression.Terms {
			term := logicalExpression.Terms[i]
			result = ExtractPredicates(operator, term, result)
		}
	} else {
		result = append(result, expression)
	}
	return result
}
