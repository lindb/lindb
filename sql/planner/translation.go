package planner

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// TranslationMap keeps mapping of fields and AST expressions to symbols
// in the current plan within query boundary.
type TranslationMap struct {
	scope        *analyzer.Scope
	context      *context.PlannerContext
	astToSymbols map[tree.NodeID]*plan.Symbol
	outerContext *TranslationMap

	fieldSymbols []*plan.Symbol
}

func (t *TranslationMap) Rewrite(expr tree.Expression) tree.Expression {
	// TODO: check symbol referencea are not allowed/expr if analyzed
	fmt.Sprintf("rewrite expression=%v,%T\n", expr)
	return t.translate(expr, true)
}

func (t *TranslationMap) withNewMappings(mappings map[tree.NodeID]*plan.Symbol, fields []*plan.Symbol) *TranslationMap {
	return &TranslationMap{
		context:      t.context,
		scope:        t.scope,
		outerContext: t.outerContext,
		astToSymbols: mappings,
		fieldSymbols: fields,
	}
}

func (t *TranslationMap) withAdditionalMapping(mappings map[tree.NodeID]*plan.Symbol) *TranslationMap {
	newMappings := make(map[tree.NodeID]*plan.Symbol)
	for k, v := range t.astToSymbols {
		newMappings[k] = v
	}
	for k, v := range mappings {
		newMappings[k] = v
	}
	fmt.Printf("addition mapping=%v,%v\n", newMappings, t)
	return &TranslationMap{
		scope:        t.scope,
		context:      t.context,
		outerContext: t.outerContext,
		astToSymbols: newMappings, // TODO: verify ast expression
		fieldSymbols: t.fieldSymbols,
	}
}

func (t *TranslationMap) tryGetMapping(node tree.Expression) *tree.SymbolReference {
	fmt.Printf("try get maping=%v,%T\n", t.astToSymbols, node)
	if t.astToSymbols == nil {
		return nil
	}
	symbol, ok := t.astToSymbols[node.GetID()]
	if ok {
		return symbol.ToSymbolReference()
	}
	return nil
}

func (t *TranslationMap) getSymbolForColumn(node tree.Expression) *plan.Symbol {
	field := t.context.AnalyzerContext.Analysis.GetColumnReferenceField(node)
	if field == nil {
		return nil
	}
	if t.scope.IsLocalScope(field.Scope) {
		fmt.Printf("look........%v\n", field.HierarchyFieldIndex)
		return t.fieldSymbols[field.HierarchyFieldIndex]
	}

	if t.outerContext != nil {
		return plan.SymbolFrom(t.outerContext.Rewrite(node))
	}

	return nil
}

func (t *TranslationMap) CanTranslate(node tree.Expression) bool {
	// TODO: check symbol referencea are not allowed
	if _, ok := t.astToSymbols[node.GetID()]; ok {
		return true
	}
	if _, ok := node.(*tree.FieldReference); ok {
		return true
	}

	if field := t.context.AnalyzerContext.Analysis.GetColumnReferenceField(node); field != nil {
		return t.scope.IsLocalScope(field.Scope)
	}
	return false
}

func (t *TranslationMap) translate(node tree.Expression, isRoot bool) (result tree.Expression) {
	fmt.Printf("translate(%T)\n", node)
	mapped := t.tryGetMapping(node)
	if mapped != nil {
		result = mapped
	} else {
		switch expr := node.(type) {
		case *tree.FieldReference:
			result = t.getSymbolForColumn(expr).ToSymbolReference()
		case *tree.DereferenceExpression:
			if t.context.AnalyzerContext.Analysis.IsColumnReference(node) {
				symbol := t.getSymbolForColumn(node)
				if symbol == nil {
					panic(fmt.Sprintf("no mapping for %T", node))
				}
				result = symbol.ToSymbolReference()
			}
			// TODO: add
		case *tree.Identifier:
			result = t.getSymbolForColumn(node).ToSymbolReference()
		case *tree.LongLiteral:
			result = &tree.Constant{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID: node.GetID(),
				},
				Type:  types.DTInt,
				Value: expr.Value,
			}
		case *tree.IntervalLiteral:
			result = &tree.Constant{
				BaseNode: tree.BaseNode{
					ID: node.GetID(),
				},
				Type:  types.DTDuration,
				Value: expr.Value,
			}
		case *tree.StringLiteral:
			result = &tree.Constant{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID: node.GetID(),
				},
				Type:  types.DTString,
				Value: expr.Value,
			}
		case *tree.ArithmeticBinaryExpression:
			exceptedType := t.context.AnalyzerContext.Analysis.GetType(expr)
			result = &tree.FunctionCall{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID: node.GetID(),
				},
				Name:      expr.Operator.FunctionName(),
				RetType:   exceptedType,
				Arguments: []tree.Expression{t.translate(expr.Left, false), t.translate(expr.Right, false)},
			}
		case *tree.LogicalExpression:
			for _, term := range expr.Terms {
				t.translate(term, false)
			}
			result = expr
		case *tree.TimePredicate:
			t.translate(expr.Value, false)
			result = expr
		case *tree.ComparisonExpression:
			t.translate(expr.Left, false)
			t.translate(expr.Right, false)
			// TODO:
			result = expr
		case *tree.FunctionCall:
			if !expression.IsFuncSupported(expr.Name) {
				panic(fmt.Sprintf("function %s is not supported", expr.Name))
			}
			expr.RetType = t.context.AnalyzerContext.Analysis.GetType(expr)
			// expr.Name = t.context.AnalyzerContext.Analysis.GetResolvedFunction(expr)
			expr.Arguments = lo.Map(expr.Arguments, func(arg tree.Expression, index int) tree.Expression {
				return t.translate(arg, false)
			})
			result = expr
		default:
			panic(fmt.Sprintf("translate not supported: %T", node))
		}
	}
	if isRoot {
		return result
	}
	// TODO: need refact
	return coerceIfNecessary(t.context.AnalyzerContext.Analysis, node, result)
}
