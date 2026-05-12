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

package analyzer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

var log = logger.GetLogger("Analyzer", "Statement")

type StatementAnalyzer struct {
	ctx         *AnalyzerContext
	metadataMgr spi.MetadataManager
}

func NewStatementAnalyzer(ctx *AnalyzerContext, metadataMgr spi.MetadataManager) *StatementAnalyzer {
	return &StatementAnalyzer{
		ctx:         ctx,
		metadataMgr: metadataMgr,
	}
}

func (sa *StatementAnalyzer) Analyze(node tree.Node) *Scope {
	return sa.analyze(node, nil, true)
}

func (sa *StatementAnalyzer) analyze(node tree.Node, outerQueryScope *Scope, isTopLevel bool) *Scope {
	visitor := NewStatementVisitor(outerQueryScope, sa, isTopLevel)
	scope := node.Accept(nil, visitor)
	if scope == nil {
		return nil
	}
	return scope.(*Scope)
}

type StatementVisitor struct {
	outerQueryScope *Scope
	analyzer        *StatementAnalyzer

	isTopLevel bool
}

func NewStatementVisitor(outerQueryScope *Scope, analyzer *StatementAnalyzer, isTopLevel bool) *StatementVisitor {
	return &StatementVisitor{
		outerQueryScope: outerQueryScope,
		analyzer:        analyzer,
		isTopLevel:      isTopLevel,
	}
}

func (v *StatementVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.Insert:
		return v.visitInsert(context, node)
	case *tree.Query:
		return v.visitQuery(context, node)
	case *tree.QuerySpecification:
		return v.visitQuerySpecification(context, node)
	case *tree.Table:
		return v.visitTable(context, node)
	case *tree.AliasedRelation:
		return v.visitAliasedRelation(context, node)
	case *tree.Values:
		return v.visitValues(context, node)
	case *tree.Join:
		return v.visitJoin(context, node)
	case *tree.FunctionCall:
		return v.visitFunctionCall(context, node)
	default:
		panic(fmt.Sprintf("statement analyzer unsupport node: %T", n))
	}
}

func (v *StatementVisitor) visitInsert(context any, node *tree.Insert) *Scope {
	// analyze query that creates data
	queryScope := node.Query.Accept(context, v).(*Scope)

	v.analyzer.ctx.Analysis.SetInsert(&Insert{
		Table: node.Table,
	})

	// TODO: need add output scope
	return v.createAndAssignScope(node, queryScope, NewRelation([]*tree.Field{}))
}

func (v *StatementVisitor) visitQuery(context any, node *tree.Query) *Scope {
	var scope *Scope
	if context != nil {
		scope = context.(*Scope)
	}

	// analyze named queries
	withScope := v.analyzeWith(node, scope)

	// analyze query body
	queryBodyScope := node.QueryBody.Accept(withScope, v).(*Scope)

	// analyze order by
	var orderByExpressions []tree.Expression
	if node.OrderBy != nil {
		orderByExpressions = v.analyzeOrderBy(node, node.OrderBy.SortItems, queryBodyScope)
		// FIXME: not the root scope and ORDER BY is ineffective
	}
	v.analyzer.ctx.Analysis.SetOrderByExpressions(node, orderByExpressions)

	// analyze limit
	if node.Limit != nil {
		v.analyzeLimit(node.Limit, queryBodyScope)
	}

	// input fields == Output fields
	v.analyzer.ctx.Analysis.SetSelectExpressions(node, v.descriptorToFields(queryBodyScope))

	queryScope := NewScopeBuilder(withScope).
		withRelation(NewRelationID(node), queryBodyScope.RelationType).
		build()

	v.analyzer.ctx.Analysis.SetScope(node, queryScope)
	return queryScope
}

func (v *StatementVisitor) visitQuerySpecification(context any, node *tree.QuerySpecification) *Scope {
	scope := context.(*Scope)
	// analyze from(relation)
	sourceScope := v.analyzeFrom(node, scope)
	if node.Where != nil {
		// analyze where condition
		v.analyzeWhere(node, sourceScope, node.Where)
	}

	v.analyzeSelect(node, sourceScope)
	groupByAnalysis := v.analyzeGroupBy(node, sourceScope)
	v.analyzeHaving(node, sourceScope)

	outputScope := v.computeAndAssignOutputScope(node, scope, sourceScope)

	var orderByExpressions []tree.Expression
	var orderByScope *Scope
	if node.OrderBy != nil {
		// FIXME: create order by scope
		orderByScope = v.computeAndAssignOrderByScope(node.OrderBy, sourceScope, outputScope, nil)
		orderByExpressions = v.analyzeOrderBy(node, node.OrderBy.SortItems, orderByScope)
		// FIXME: not the root scope and ORDER BY is ineffective
	}
	v.analyzer.ctx.Analysis.SetOrderByExpressions(node, orderByExpressions)

	// analyze limit
	if node.Limit != nil {
		v.analyzeLimit(node.Limit, outputScope)
	}

	var sourceExpressions []tree.Expression
	selectExpressions := v.analyzer.ctx.Analysis.GetSelectExpressions(node)
	for _, selectExpr := range selectExpressions {
		sourceExpressions = append(sourceExpressions, selectExpr.Expression)
	}

	// FIXME: select

	if node.Having != nil {
		// if has having expression, add to source expressions
		sourceExpressions = append(sourceExpressions, node.Having)
	}
	v.analyzeGroupingOperations(node, sourceExpressions, orderByExpressions)
	v.analyzeAggregations(node, sourceScope, orderByScope, groupByAnalysis, sourceExpressions, orderByExpressions)

	// FIXME: order agg

	return outputScope
}

func (v *StatementVisitor) visitJoin(context any, node *tree.Join) *Scope {
	scope := context.(*Scope)
	left := node.Left.Accept(scope, v).(*Scope)
	right := node.Right.Accept(scope, v).(*Scope)
	criteria := node.Criteria
	if joinUsing, ok := criteria.(*tree.JoinUsing); ok {
		return v.analyzeJoinUsing(node, joinUsing.Columns, scope, left, right)
	}
	output := v.createAndAssignScope(node, scope, left.RelationType.joinWith(right.RelationType))
	if node.Type == tree.CROSS || node.Type == tree.IMPLICIT {
		return output
	}

	if joinOn, ok := criteria.(*tree.JoinOn); ok {
		expression := joinOn.Expression
		// FIXME: impl it
		v.analyzeExpression(expression, output)

		v.analyzer.ctx.Analysis.SetJoinCriteria(node, expression)
	}

	return output
}

func (v *StatementVisitor) analyzeJoinUsing(node *tree.Join, columns []*tree.Identifier,
	scope, left, right *Scope,
) *Scope {
	panic("using")
}

func (v *StatementVisitor) visitAliasedRelation(context any, relation *tree.AliasedRelation) *Scope {
	scope := context.(*Scope)
	aliased := tree.NewQualifiedName([]*tree.Identifier{relation.Aliase})
	v.analyzer.ctx.Analysis.SetRelationName(relation, aliased)
	// TODO: v.analyzer.analysis.AddAliased(relation, aliased)

	relationScope := relation.Relation.Accept(scope, v).(*Scope)
	columnAliases := lo.Map(relation.ColumnNames, func(item *tree.Identifier, index int) string {
		return item.Value
	})
	relationType := relationScope.RelationType
	descriptor := relationType.withAlias(relation.Aliase.Value, columnAliases)

	return v.createAndAssignScope(relation, scope, descriptor)
}

func (v *StatementVisitor) visitValues(context any, values *tree.Values) *Scope {
	scope := context.(*Scope)
	schema := values.Rows.Schema()
	var fields []*tree.Field
	for i, column := range schema.Fields() {
		fields = append(fields, &tree.Field{
			Name:     column.Name,
			DataType: column.Type,
			Index:    tree.FieldIndex(i),
			// FIXME: Hidden:   column.Hidden,???
		})
	}

	return v.createAndAssignScope(values, scope, NewRelation(fields))
}

func (v *StatementVisitor) visitTable(ctx any, table *tree.Table) *Scope {
	scope := ctx.(*Scope)
	if table.Name.Prefix == nil {
		name := strings.ToLower(table.Name.Suffix)
		// if reference to a WITH query
		withQuery := createScope(scope).getNameQuery(name)
		if withQuery != nil {
			// analyze named query
			v.analyzer.ctx.Analysis.SetRelationName(table, table.Name)
			return v.createScopeForCommonTableExpression(table, scope, withQuery)
		}
	}
	database := table.GetDatabase(v.analyzer.ctx.Database)
	namespace := table.GetNamespace()
	tableMetadata, err := v.analyzer.metadataMgr.GetTableMetadata(database,
		namespace, table.GetTableName())
	if err != nil {
		log.Warn("get table metadata fail", logger.String("database", database), logger.String("ns", namespace),
			logger.String("table", table.GetTableName()), logger.Error(err))
		// TODO: remove
		panic(err)
	}

	// analyze table
	var outputFields []*tree.Field
	for i, col := range tableMetadata.Schema.Fields() {
		// Check the Arrow field metadata for the "hidden" marker (set by the storage
		// layer, e.g. the implicit timestamp column in metric queries).
		hidden := col.Metadata.FindKey("hidden") >= 0 && col.Metadata.Values()[col.Metadata.FindKey("hidden")] == "true"
		// TODO: check agg????
		outputFields = append(outputFields, &tree.Field{
			Index:         tree.FieldIndex(i),
			Name:          col.Name, // TODO: dup tag name/field name
			DataType:      col.Type,
			Hidden:        hidden,
			RelationAlias: table.Name.Name, // TODO: relation alias
		})
	}

	// TODO: check table type
	v.analyzer.ctx.Analysis.SetRelationName(table, table.Name)
	tableHandle := v.analyzer.metadataMgr.GetTableHandle(database, namespace, table.GetTableName())
	v.analyzer.ctx.Analysis.RegisterTableHandle(table, tableHandle)
	v.analyzer.ctx.Analysis.RegisterTableMetadata(tableHandle.String(), tableMetadata)
	// FIXME: table fields??

	subScope := v.createAndAssignScope(table, scope, NewRelation(outputFields))
	subScope.Dynamic = tableMetadata.SupportDynamicField
	return subScope
}

func (v *StatementVisitor) visitFunctionCall(context any, node *tree.FunctionCall) (r any) {
	panic("impl func.....")
}

func (v *StatementVisitor) analyzeWith(node *tree.Query, scope *Scope) *Scope {
	if !node.HasWith() {
		return createScope(scope)
	}
	// analyze with clause
	with := node.With
	withScopeBuilder := NewScopeBuilder(scope)
	for i := range with.Queries {
		withQuery := with.Queries[i]
		name := strings.ToLower(withQuery.Name.Value)
		// check name if duplicate
		if withScopeBuilder.containsNamedQuery(name) {
			panic(fmt.Sprintf("with query name '%s' specified more than once", name))
		}
		// analyze query statement
		v.analyzer.analyze(withQuery.Query, withScopeBuilder.build(), false)
		// store name query under scope
		withScopeBuilder.withNameQuery(name, withQuery)
	}
	withScope := withScopeBuilder.build()
	v.analyzer.ctx.Analysis.SetScope(with, withScope)
	return withScope
}

func (v *StatementVisitor) analyzeSelect(node *tree.QuerySpecification,
	scope *Scope,
) (outputExpressions []tree.Expression) {
	var selectExpressions []*SelectExpression
	for i := range node.Select.SelectItems {
		selectItem := node.Select.SelectItems[i]
		switch item := selectItem.(type) {
		case *tree.AllColumns:
			outputExpressions, selectExpressions = v.analyzeSelectAllColumns(item, node, scope,
				outputExpressions, selectExpressions)
		case *tree.SingleColumn:
			outputExpressions, selectExpressions = v.analyzeSelectSingleColumn(item, node, scope,
				outputExpressions, selectExpressions)
		default:
			panic(fmt.Sprintf("unsupported select type type: %s", reflect.TypeOf(item)))
		}
	}

	v.analyzer.ctx.Analysis.SetSelectExpressions(node, selectExpressions)
	return
}

func (v *StatementVisitor) analyzeSelectSingleColumn(singleColumn *tree.SingleColumn, node *tree.QuerySpecification,
	scope *Scope, outputExpressions []tree.Expression, selectExpressions []*SelectExpression,
) (outputs []tree.Expression, selects []*SelectExpression) {
	expression := singleColumn.Expression
	v.analyzeExpression(expression, scope)
	outputExpressions = append(outputExpressions, expression)
	selectExpressions = append(selectExpressions, &SelectExpression{
		Alias:      singleColumn.Alias,
		Expression: expression,
	})
	// TODO: check distinct
	return outputExpressions, selectExpressions
}

func (v *StatementVisitor) analyzeSelectAllColumns(allColumns *tree.AllColumns, node *tree.QuerySpecification,
	scope *Scope, outputExpressions []tree.Expression, selectExpressions []*SelectExpression,
) (outputs []tree.Expression, selects []*SelectExpression) {
	// expand * and expression.*
	if allColumns.Target != nil {
		// analyze all columns with target expression(expression.*)
		expression := allColumns.Target
		prefix := asQualifiedName(expression)
		if prefix != nil {
			// analyze prefix as an 'asterisked identifier chain'
			// ref table
			identifierChain := scope.resolveAsteriskedIdentifierChain(prefix, allColumns)
			// TODO: relation := v.analyzer.analysis.GetRelationByAliased(prefix)
			if identifierChain == nil {
				panic(fmt.Sprintf("unable to resolve reference %s", prefix.Name))
			}
			if identifierChain.Type == TABLE {
				relationType := identifierChain.RelationType
				// TODO: relationScope := v.analyzer.analysis.GetScope(relation)
				// FIXME:????? scope from
				outputExpressions, selectExpressions = v.analyzeAllColumnsFromTable(allColumns, node, scope,
					outputExpressions, selectExpressions, relationType, prefix)
				return outputExpressions, selectExpressions
			}
		}
	} else {
		// analyze all columns without target expression('*')
		// TODO: add check
		outputExpressions, selectExpressions = v.analyzeAllColumnsFromTable(allColumns, node, scope,
			outputExpressions, selectExpressions, scope.RelationType, nil)
	}
	return outputExpressions, selectExpressions
}

func (v *StatementVisitor) analyzeAllColumnsFromTable(allColumns *tree.AllColumns, node *tree.QuerySpecification,
	scope *Scope, outputExpressions []tree.Expression, selectExpressions []*SelectExpression,
	relationType *Relation, relationAlias *tree.QualifiedName,
) (outputs []tree.Expression, selects []*SelectExpression) {
	fields := relationType.Fields
	for _, field := range fields {
		fieldRef := &tree.FieldReference{
			BaseNode: tree.BaseNode{
				ID: v.analyzer.ctx.IDAllocator.Next(),
			},
			FieldIndex: field.Index,
		}
		v.analyzeExpression(fieldRef, scope)
		outputExpressions = append(outputExpressions, fieldRef)
		selectExpressions = append(selectExpressions, &SelectExpression{
			Expression: fieldRef,
		})
	}
	// FIXME: ???
	v.analyzer.ctx.Analysis.SetSelectAllResultFields(allColumns, fields)
	return outputExpressions, selectExpressions
}

// ----- relation ------
func (v *StatementVisitor) analyzeFrom(node *tree.QuerySpecification, scope *Scope) *Scope {
	if node.From != nil {
		return node.From.Accept(scope, v).(*Scope)
	}
	result := createScope(scope)
	v.analyzer.ctx.Analysis.SetImplicitFromScope(node, result)
	return result
}

func (v *StatementVisitor) analyzeWhere(node *tree.QuerySpecification, scope *Scope, predicate tree.Expression) {
	if timePredicate, ok := predicate.(*tree.TimePredicate); ok {
		v.analyzer.ctx.Analysis.SetTimePredicates(node, []*tree.TimePredicate{timePredicate})
		return
	}
	// extract time predicates from where clause expressions
	timePredicates, newPredicate := tree.ExtractTimePredicates(predicate)
	if len(timePredicates) > 0 {
		v.analyzer.ctx.Analysis.SetTimePredicates(node, timePredicates)
	}

	if newPredicate == nil {
		// new predicate is nil, means no where clause after extract time predicates
		return
	}

	// FIXME: verify no aggregate and group by function
	v.analyzeExpression(newPredicate, scope)
	// TODO: v.analyzer.ctx.Analysis.RecordSubQueries(node, expressionAnalysis)

	// FIXME: check predicate type
	// predicateType := expressionAnalysis.GetType(predicate)

	v.analyzer.ctx.Analysis.SetWhere(node, newPredicate)
}

func (v *StatementVisitor) analyzeGroupBy(node *tree.QuerySpecification, scope *Scope) *GroupingSetAnalysis {
	var (
		groupingExpressions []tree.Expression
		complexExpressions  []tree.Expression
		sets                [][]*FieldID
	)
	selectExpressions := v.analyzer.ctx.Analysis.GetSelectExpressions(node)
	if node.GroupBy != nil {
		for _, groupingElement := range node.GroupBy.GroupingElements {
			switch groupByEle := groupingElement.(type) {
			case *tree.GroupByAllColumns:
				panic("impl group by all columns")
			// TODO: group by *
			case *tree.SimpleGroupBy:
				var field *ResolvedField
				for _, column := range groupByEle.Columns {
					switch item := column.(type) {
					case *tree.LongLiteral:
						// fixme: index field
						panic("impl long group key index ref to select item")
					case *tree.IntervalLiteral:
						// set grouping interval
						v.analyzer.ctx.Analysis.SetGroupingInterval(node, item)
						// ignore grouping interval
						goto Next
					case *tree.Identifier:
						selectExpression, ok := lo.Find(selectExpressions, func(selectExpr *SelectExpression) bool {
							return selectExpr.Alias != nil && selectExpr.Alias.Value == item.Value
						})
						if ok {
							// rewrite group by column with select expression if alias matched
							column = selectExpression.Expression
						} else {
							v.analyzeExpression(column, scope)
						}
					default:
						v.analyzeExpression(column, scope)
					}

					field = v.analyzer.ctx.Analysis.GetColumnReferenceField(column)
					if field != nil {
						// TODO: check field if aggregate
						if arrow.TypeEqual(field.Field.DataType, arrow.FixedWidthTypes.Timestamp_ns) {
							panic(fmt.Sprintf("aggregate/timestamp field[%v] cannot appear in group by", field.Field.Name))
						}
						sets = append(sets, []*FieldID{field.FieldID()})
					} else {
						// TODO: field sets
						complexExpressions = append(complexExpressions, column)
					}
					groupingExpressions = append(groupingExpressions, column)

				Next:
				}
			}
		}
		if len(groupingExpressions) == 0 {
			// no grouping column
			return nil
		}
	} else if v.hasAggregates(node) {
		for _, item := range node.Select.SelectItems {
			if single, ok := item.(*tree.SingleColumn); ok {
				field := v.analyzer.ctx.Analysis.GetColumnReferenceField(single.Expression)
				if field != nil && arrow.TypeEqual(field.Field.DataType, arrow.FixedWidthTypes.Timestamp_ns) {
					sets = append(sets, []*FieldID{field.FieldID()})
				}
			}
		}
	} else {
		return nil
	}

	groupingSets := NewGroupingSetAnalysis(groupingExpressions, sets, complexExpressions)
	v.analyzer.ctx.Analysis.SetGroupingSets(node, groupingSets)
	return groupingSets
}

func (v *StatementVisitor) analyzeGroupingOperations(node *tree.QuerySpecification,
	outputExpressions, orderByExpressions []tree.Expression) {
}

func (v *StatementVisitor) analyzeAggregations(query *tree.QuerySpecification, sourceScope, orderByScope *Scope,
	groupByAnalysis *GroupingSetAnalysis, outputExpressions, orderByExpressions []tree.Expression,
) {
	var expr []tree.Node
	for _, output := range outputExpressions {
		expr = append(expr, output)
	}
	for _, orderBy := range orderByExpressions {
		expr = append(expr, orderBy)
	}
	var functions []*tree.FunctionCall
	var stack []tree.Expression
	isFuncArg := func() bool {
		if len(stack) == 0 {
			return false
		}
		_, ok := stack[len(stack)-1].(*tree.FunctionCall)
		return ok
	}
	// TODO:
	tree.ExtractAggregationFunctions(expr, func(n tree.Node) {
		switch node := n.(type) {
		case *tree.Identifier:
			// transfer filed builtin aggregation
			resolvedField := sourceScope.resolveField(n, tree.NewQualifiedName([]*tree.Identifier{node}), true)
			if aggType, ok := resolvedField.Field.DataType.(*larray.AggregationType); ok && !isFuncArg() {
				// agg field and field is not function arg, add builtin agg func for this field
				fn := &tree.FunctionCall{
					Name: tree.FuncName(string(aggType.Kind())),
					Arguments: []tree.Expression{&tree.SymbolReference{
						Name:     resolvedField.Field.Name,
						DataType: resolvedField.Field.DataType,
						Hidden:   resolvedField.Field.Hidden,
					}},
					RefField: resolvedField.Field,
				}
				functions = append(functions, fn)
				v.analyzer.ctx.Analysis.AddResolvedFunction(fn, fn.Name)
				v.analyzer.ctx.Analysis.AddType(fn, resolvedField.Field.DataType) // TODO: remove it
			}
		case *tree.FunctionCall:
			if tree.IsAggFunc(node.Name) {
				functions = append(functions, node)
				// TODO: need do other func
				v.analyzer.ctx.Analysis.AddResolvedFunction(node, node.Name)
			}
		}

		// add node into expression stack
		stack = append(stack, n)
	})
	v.analyzer.ctx.Analysis.SetAggregates(query, functions) // TODO: remove
	// TODO: extract agg func
	if v.analyzer.ctx.Analysis.IsGroupingSets(query) {
		// ensure SELECT, ORDER BY and HAVING are constant with respect to group
		// e.g, these are all valid expressions:
		//     SELECT f(a) GROUP BY a
		//     SELECT f(a + 1) GROUP BY a + 1
		//     SELECT a + sum(b) GROUP BY a
		distinctGroupingColumns := groupByAnalysis.GetOriginalExpression()
		verifySourceAggregations(v.analyzer.ctx.Analysis, distinctGroupingColumns, outputExpressions)

		if len(orderByExpressions) > 0 {
			verifyOrderByAggregations(v.analyzer.ctx.Analysis, distinctGroupingColumns, orderByExpressions)
		}
	}
}

func (v *StatementVisitor) analyzeHaving(node *tree.QuerySpecification, scope *Scope) {
	// FIXME:impl it
}

func (v *StatementVisitor) analyzeOrderBy(node tree.Node,
	sortItems []*tree.SortItem, orderByScope *Scope,
) (orderByExpressions []tree.Expression) {
	panic("implement it analyzeOrderBy")
}

func (v *StatementVisitor) analyzeLimit(node *tree.Limit, _ *Scope) {
	var rowCount int64

	if long, ok := node.RowCount.(*tree.LongLiteral); ok {
		rowCount = long.Value
	}
	if rowCount < 0 {
		panic(fmt.Sprintf("limit row count must be greater or equal to 0 (actual value: %d)", rowCount))
	}

	v.analyzer.ctx.Analysis.SetLimit(node, rowCount)
}

func (v *StatementVisitor) analyzeExpression(expression tree.Expression, scope *Scope) {
	analyzer := NewExpressionAnalyzer(v.analyzer.ctx)
	analyzer.Analyze(expression, scope)
}

func (v *StatementVisitor) createScopeForCommonTableExpression(table *tree.Table, scope *Scope,
	withQuery *tree.WithQuery,
) *Scope {
	query := withQuery.Query
	v.analyzer.ctx.Analysis.RegisterNamedQuery(table, query)
	// FIXME: analyze field
	var fields []*tree.Field

	return v.createAndAssignScope(table, scope, NewRelation(fields))
}

func (v *StatementVisitor) createAndAssignScope(node tree.Node, parent *Scope, relationType *Relation) *Scope {
	if relationType == nil {
		panic("nil....")
	}
	scope := NewScopeBuilder(parent).
		withRelation(NewRelationID(node), relationType).
		build()
	v.analyzer.ctx.Analysis.SetScope(node, scope)
	return scope
}

func (v *StatementVisitor) computeAndAssignOutputScope(node *tree.QuerySpecification,
	scope, _ *Scope,
) *Scope {
	var outputFields []*tree.Field
	selectItems := node.Select.SelectItems
	for i := range selectItems {
		selectItem := selectItems[i]
		switch item := selectItem.(type) {
		case *tree.AllColumns:
			fields := v.analyzer.ctx.Analysis.GetSelectAllResultFields(item)
			outputFields = append(outputFields, fields...)
		case *tree.SingleColumn:
			expression := item.Expression
			field := item.Alias
			var name *tree.QualifiedName
			switch expr := expression.(type) {
			case *tree.Identifier:
				name = tree.NewQualifiedName([]*tree.Identifier{{Value: expr.Value}})
			case *tree.DereferenceExpression:
				name = expr.ToQualifiedName()
			}

			if field == nil {
				if name != nil {
					field = name.OriginalParts[len(name.OriginalParts)-1] // get last value
				}
			}

			var fieldName string
			if field != nil {
				fieldName = field.Value
			}

			// NOTE: field name is empty when expression/function call
			outputFields = append(outputFields, &tree.Field{
				Name:  fieldName,
				Index: tree.FieldIndex(i),
			})
		default:
			panic(fmt.Sprintf("unsupported selec type type: %s", reflect.TypeOf(item)))
		}
	}
	return v.createAndAssignScope(node, scope, NewRelation(outputFields))
}

func (v *StatementVisitor) computeAndAssignOrderByScope(_ *tree.OrderBy,
	sourceScope, outputSource *Scope, fields []*tree.Field,
) *Scope {
	panic("impl ordery by")
}

func (v *StatementVisitor) descriptorToFields(scope *Scope) (selectExpressions []*SelectExpression) {
	for _, field := range scope.RelationType.Fields {
		expression := &tree.FieldReference{
			BaseNode: tree.BaseNode{
				ID: v.analyzer.ctx.IDAllocator.Next(),
			},
			FieldIndex: field.Index,
		}
		selectExpressions = append(selectExpressions, &SelectExpression{
			Expression: expression,
		})
		v.analyzeExpression(expression, scope)
	}
	return
}

func (v *StatementVisitor) hasAggregates(node *tree.QuerySpecification) bool {
	var toExtract []tree.Node
	for _, selectItem := range node.Select.SelectItems {
		toExtract = append(toExtract, selectItem)
	}
	var aggregates []tree.Expression
	tree.ExtractAggregationFunctions(toExtract, func(n tree.Node) {
		if _, ok := n.(*tree.FunctionCall); ok {
			aggregates = append(aggregates, n)
		}
	})
	return len(aggregates) != 0
}
