package planner

import (
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// TranslationMap keeps mapping of fields and AST expressions to symbols
// in the current plan within query boundary.
type TranslationMap struct {
	scope        *analyzer.Scope
	context      *context.PlannerContext
	astToSymbols map[tree.NodeID]*plan.Symbol

	fieldSymbols []*plan.Symbol
}

func (t *TranslationMap) Rewrite(root tree.Expression) tree.Expression {
	// TODO: check symbol referencea are not allowed/expr if analyzed
	// return tree.RewriteExpression(nil, &expressionRewriter{translation: t}, node)
	return t.translate(root, true)
}

func (t *TranslationMap) withNewMappings(mappings map[tree.NodeID]*plan.Symbol, fields []*plan.Symbol) *TranslationMap {
	return &TranslationMap{
		context:      t.context,
		scope:        t.scope,
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
		astToSymbols: newMappings, // TODO: verify ast expression
		fieldSymbols: t.fieldSymbols,
	}
}

func (t *TranslationMap) tryGetMapping(node tree.Expression) *tree.SymbolReference {
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

	// TODO: out context

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
				Type:  types.DataTypeInt,
				Value: expr.Value,
			}
		case *tree.ArithmeticBinaryExpression:
			exceptedType := t.context.AnalyzerContext.Analysis.GetType(expr)
			// coercion, ok := t.context.AnalyzerContext.Analysis.GetCoercion(expr)
			// if ok {
			// 	exceptedType = coercion
			// }

			result = &tree.Call{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID: node.GetID(),
				},
				Function: expr.Operator.FunctionName(),
				RetType:  exceptedType,
				Args:     []tree.Expression{t.translate(expr.Left, false), t.translate(expr.Right, false)},
			}
		default:
			panic(fmt.Sprintf("expression rewrite unimplemented: %T", node))
		}
	}
	if isRoot {
		return result
	}
	// TODO: need refact
	return t.coerceIfNecessary(node, result)
}

func (t *TranslationMap) coerceIfNecessary(origianl, rewritten tree.Expression) tree.Expression {
	coercion, ok := t.context.AnalyzerContext.Analysis.GetCoercion(origianl)
	fmt.Printf("check coercion%T=%v,=%v\n", origianl, coercion, ok)
	if !ok {
		return rewritten
	}
	fmt.Println("cast ....... rewrite")
	return &tree.Cast{
		BaseNode: tree.BaseNode{
			ID: origianl.GetID(),
		},
		Type:       coercion,
		Expression: rewritten,
	}
}

type expressionRewriter struct {
	translation *TranslationMap
}

func (e *expressionRewriter) RewriteExpression(context any, node tree.Expression) tree.Expression {
	switch expr := node.(type) {
	case *tree.FieldReference:
		return e.rewriteFieldReference(expr)
	case *tree.DereferenceExpression:
		return e.rewriteDereferenceExpression(expr)
	case *tree.Identifier:
		return e.rewriteIndentifier(expr)
	case *tree.Cast:
		return e.RewriteExpression(context, expr.Expression)
	case *tree.ArithmeticBinaryExpression:
		return e.rewriteArithmeticBinaryExpression(expr)
	default:
		panic(fmt.Sprintf("expression rewrite unimplemented: %T", node))
	}
}

func (e *expressionRewriter) rewriteArithmeticBinaryExpression(node *tree.ArithmeticBinaryExpression) tree.Expression {
	mapped := e.translation.tryGetMapping(node)
	fmt.Printf("mapped %v\n", mapped)
	if mapped != nil {
		return e.coerceIfNecessary(node, mapped)
	}
	// symbol := e.translation.getSymbolForColumn(node)
	// fmt.Printf("symbol %v\n", symbol)
	// if symbol == nil {
	// 	return e.coerceIfNecessary(node, node)
	// }
	// TODO: add default rewrite
	rewrittenExpression := node
	return e.coerceIfNecessary(node, rewrittenExpression)
}

func (e *expressionRewriter) rewriteIndentifier(node *tree.Identifier) tree.Expression {
	fmt.Printf("rewrite====%v\n", node.Value)
	mapped := e.translation.tryGetMapping(node)
	fmt.Printf("mapped %v\n", mapped)
	if mapped != nil {
		return e.coerceIfNecessary(node, mapped)
	}
	symbol := e.translation.getSymbolForColumn(node)
	fmt.Printf("symbol %v\n", symbol)
	if symbol == nil {
		return e.coerceIfNecessary(node, node)
	}
	return e.coerceIfNecessary(node, symbol.ToSymbolReference())
}

func (e *expressionRewriter) rewriteDereferenceExpression(node *tree.DereferenceExpression) tree.Expression {
	mapped := e.translation.tryGetMapping(node)
	if mapped != nil {
		return e.coerceIfNecessary(node, mapped)
	}
	if e.translation.context.AnalyzerContext.Analysis.IsColumnReference(node) {
		symbol := e.translation.getSymbolForColumn(node)
		if symbol == nil {
			panic(fmt.Sprintf("no mapping for %T", node))
		}
		fmt.Println("hahahahahh..")
		return e.coerceIfNecessary(node, symbol.ToSymbolReference())
	}

	return nil
}

func (e *expressionRewriter) rewriteFieldReference(node *tree.FieldReference) tree.Expression {
	mapped := e.translation.tryGetMapping(node)
	if mapped != nil {
		return e.coerceIfNecessary(node, mapped)
	}
	symbol := e.translation.getSymbolForColumn(node)
	if symbol != nil {
		return symbol.ToSymbolReference()
	}
	panic(fmt.Sprintf("no symbol mpapping for node '%T' (%d)", node, node.FieldIndex))
}

func (e *expressionRewriter) coerceIfNecessary(origianl, rewritten tree.Expression) tree.Expression {
	coercion, ok := e.translation.context.AnalyzerContext.Analysis.GetCoercion(origianl)
	fmt.Println("check coercion")
	if !ok {
		return rewritten
	}
	fmt.Println("cast ....... rewrite")
	return &tree.Cast{
		BaseNode: tree.BaseNode{
			ID: e.translation.context.AnalyzerContext.IDAllocator.Next(),
		},
		Type:       coercion,
		Expression: rewritten,
	}
}
