package analyzer

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

func verifySourceAggregations(analysis *Analysis, groupByExpressions, expressions []tree.Expression) {
	analyzer := NewAggregationAnalyzer(analysis)
	for _, expression := range expressions {
		analyzer.analyze(expression)
	}
}

func verifyOrderByAggregations(analysis *Analysis, groupByExpressions, expressions []tree.Expression) {
	analyzer := NewAggregationAnalyzer(analysis)
	for _, expression := range expressions {
		analyzer.analyze(expression)
	}
}

type AggregationAnanlyzer struct {
	analysis *Analysis
}

func NewAggregationAnalyzer(analysis *Analysis) *AggregationAnanlyzer {
	return &AggregationAnanlyzer{
		analysis: analysis,
	}
}

func (aa *AggregationAnanlyzer) analyze(expression tree.Expression) {
	visitor := NewAggregationAnalyzeVisitor(aa.analysis)
	if r, ok := expression.Accept(nil, visitor).(bool); ok {
		if !r {
			panic(fmt.Sprintf("'%s' must be an aggregate expression or appear in GROUP BY clause",
				tree.FormatExpression(expression)))
		}
	}
}

type aggregationAnalyzeVisitor struct {
	analysis *Analysis
}

func NewAggregationAnalyzeVisitor(analysis *Analysis) tree.Visitor {
	return &aggregationAnalyzeVisitor{
		analysis: analysis,
	}
}

func (v *aggregationAnalyzeVisitor) Visit(context any, n tree.Node) (r any) {
	switch node := n.(type) {
	case *tree.DereferenceExpression:
		return v.visitDereferenceExpression(node)
	default:
		panic(fmt.Sprintf("unsupported node<%T> when aggregation ananlyzer", n))
	}
}

func (v *aggregationAnalyzeVisitor) visitDereferenceExpression(node *tree.DereferenceExpression) (r any) {
	// 		            ExpressionAnalyzer.LabelPrefixedReference labelDereference = analysis.getLabelDereference(node);
	// if (labelDereference != null) {
	//     return labelDereference.getColumn().map(this::process).orElse(true);
	// }
	//
	// if (!hasReferencesToScope(node, analysis, sourceScope)) {
	//     // reference to outer scope is group-invariant
	//     return true;
	// }
	//
	// if (columnReferences.containsKey(NodeRef.<Expression>of(node))) {
	//     return isGroupingKey(node);
	// }
	//

	// Allow SELECT col1.f1 FROM table1 GROUP BY col1
	return node.Base.Accept(nil, v)
}

func (v *aggregationAnalyzeVisitor) isGroupingKey(node tree.Expression) bool {
	return false
}

// private boolean isGroupingKey(Expression node)
//       {
//           FieldId fieldId = requireNonNull(columnReferences.get(NodeRef.of(node)), () -> "No field for " + node).getFieldId();
//
//           if (orderByScope.isPresent() && isFieldFromScope(fieldId, orderByScope.get())) {
//               return true;
//           }
//
//           return groupingFields.contains(fieldId);
//       }
