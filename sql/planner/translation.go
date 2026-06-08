// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package planner

import (
	"fmt"
	"maps"

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/samber/lo"

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
	astToSymbols map[string]*plan.Symbol // expression format string => symbol
	outerContext *TranslationMap

	fieldSymbols []*plan.Symbol
}

func (t *TranslationMap) Rewrite(expr tree.Expression) tree.Expression {
	// TODO: check symbol referencea are not allowed/expr if analyzed
	return t.translate(expr, true)
}

func (t *TranslationMap) withNewMappings(mappings map[string]*plan.Symbol, fields []*plan.Symbol) *TranslationMap {
	return &TranslationMap{
		context:      t.context,
		scope:        t.scope,
		outerContext: t.outerContext,
		astToSymbols: mappings,
		fieldSymbols: fields,
	}
}

func (t *TranslationMap) withAdditionalMapping(mappings map[string]*plan.Symbol) *TranslationMap {
	newMappings := make(map[string]*plan.Symbol)
	maps.Copy(newMappings, t.astToSymbols)
	maps.Copy(newMappings, mappings)
	return &TranslationMap{
		scope:        t.scope,
		context:      t.context,
		outerContext: t.outerContext,
		astToSymbols: newMappings, // TODO: verify ast expression
		fieldSymbols: t.fieldSymbols,
	}
}

func (t *TranslationMap) tryGetMapping(node tree.Expression) *tree.SymbolReference {
	if len(t.astToSymbols) == 0 {
		return nil
	}
	symbol, ok := t.astToSymbols[node.String()]
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
	if t.scope.Dynamic || arrow.TypeEqual(field.Field.DataType, larrow.ExtensionTypes.Dynamic) {
		// Preserve the Hidden flag so fixed-schema columns (e.g. timestamp) that are
		// marked hidden in the table schema remain hidden after translation.
		// Without this, every field lookup through the Dynamic branch would produce
		// Hidden=false, causing hidden columns to appear in the query output.
		return &plan.Symbol{
			Name:     field.Field.Name,
			DataType: field.Field.DataType,
			Hidden:   field.Field.Hidden,
		}
	}
	isLocalScope := t.scope.IsLocalScope(field.Scope)
	if isLocalScope {
		return t.fieldSymbols[field.HierarchyFieldIndex]
	}

	if t.outerContext != nil {
		return plan.SymbolFrom(t.outerContext.Rewrite(node))
	}

	return nil
}

func (t *TranslationMap) CanTranslate(node tree.Expression) bool {
	// TODO: check symbol referencea are not allowed
	if _, ok := t.astToSymbols[node.String()]; ok {
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
					ID:   node.GetID(),
					Text: expr.Text,
				},
				Type:  arrow.PrimitiveTypes.Int64,
				Value: expr.Value,
			}
		case *tree.IntervalLiteral:
			result = &tree.Constant{
				BaseNode: tree.BaseNode{
					ID:   node.GetID(),
					Text: expr.Text,
				},
				Type:  arrow.FixedWidthTypes.Duration_ns,
				Value: expr.Value,
			}
		case *tree.StringLiteral:
			result = &tree.Constant{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID:   node.GetID(),
					Text: expr.Text,
				},
				Type:  arrow.BinaryTypes.String,
				Value: expr.Value,
			}
		case *tree.ArithmeticBinaryExpression:
			exceptedType := t.context.AnalyzerContext.Analysis.GetType(expr)
			result = &tree.FunctionCall{
				// TODO: replace
				BaseNode: tree.BaseNode{
					ID:   node.GetID(),
					Text: expr.Text,
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
		case *tree.LikePredicate:
			expr.Value = t.translate(expr.Value, false)
			expr.Pattern = t.translate(expr.Pattern, false)
			result = expr
		case *tree.RegexPredicate:
			expr.Value = t.translate(expr.Value, false)
			expr.Pattern = t.translate(expr.Pattern, false)
			result = expr
		case *tree.NullPredicate:
			expr.Value = t.translate(expr.Value, false)
			result = expr
		case *tree.ComparisonExpression:
			expr.Left = t.translate(expr.Left, false)
			expr.Right = t.translate(expr.Right, false)
			// TODO:
			result = expr
		case *tree.NotExpression:
			result = expr
		case *tree.Constant:
			result = expr
		case *tree.InPredicate:
			expr.Value = t.translate(expr.Value, false)
			if inListExpression, ok := expr.ValueList.(*tree.InListExpression); ok {
				lo.ForEach(inListExpression.Values, func(item tree.Expression, index int) {
					inListExpression.Values[index] = t.translate(item, false)
				})
			}
			result = expr
		case *tree.FunctionCall:
			if !tree.IsFuncSupported(expr.Name) {
				panic(fmt.Sprintf("function %s is not supported", expr.Name))
			}
			expr.RetType = t.context.AnalyzerContext.Analysis.GetType(expr)
			// TODO: expr.Name = t.context.AnalyzerContext.Analysis.GetResolvedFunction(expr)
			expr.Arguments = lo.Map(expr.Arguments, func(arg tree.Expression, index int) tree.Expression {
				return t.translate(arg, false)
			})
			result = expr
		case *tree.SymbolReference:
			result = expr
		case *tree.SubscriptExpression:
			expr.Base = t.translate(expr.Base, false)
			expr.Key = t.translate(expr.Key, false)
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
