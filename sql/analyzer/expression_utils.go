package analyzer

import (
	"github.com/lindb/lindb/sql/tree"
)

func asQualifiedName(expression tree.Expression) (name *tree.QualifiedName) {
	if identifier, ok := expression.(*tree.Identifier); ok {
		name = tree.NewQualifiedName([]*tree.Identifier{identifier})
	} else if dereference, ok := expression.(*tree.DereferenceExpression); ok {
		// TODO:????
		name = dereference.ToQualifiedName()
	}
	return
}

func ExtractConjuncts(expression tree.Expression) (result []tree.Expression) {
	return ExtractPredicates(tree.LogicalAND, expression, result)
}

func ExtractTimePredicates(expression tree.Expression) (result []*tree.TimePredicate, newExpr tree.Expression) {
	if logicalExpression, ok := expression.(*tree.LogicalExpression); ok {
		var newTerms []tree.Expression
		for _, term := range logicalExpression.Terms {
			if timePredicate, ok := term.(*tree.TimePredicate); ok {
				result = append(result, timePredicate)
			} else {
				subResult, newTerm := ExtractTimePredicates(term)
				if len(subResult) > 0 && len(result) > 0 {
					panic("time predicate not support nested")
				}
				result = subResult
				if newTerm != nil {
					newTerms = append(newTerms, term)
				}
			}
		}
		if len(newTerms) == 0 {
			newExpr = nil
			return
		}
		if len(newTerms) != len(logicalExpression.Terms) {
			logicalExpression.Terms = newTerms
			newExpr = logicalExpression
		} else {
			newExpr = logicalExpression
		}
	} else {
		newExpr = expression
	}
	return
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

func ExtractAggregationFunctions(nodes []tree.Expression, handle func(node tree.Node)) {
	visitor := &tree.DefaultTraversalVisitor{
		PreProcess: handle,
	}
	for _, node := range nodes {
		visitor.Visit(nil, node)
	}
}
