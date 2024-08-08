package analyzer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/function"
	"github.com/lindb/lindb/sql/tree"
)

// Plan Fragemnt 0
//
//	Output[columnNames = [idle]]
//	│ Layout: [xxx]
//	│ idle := xxx
//	└─ Remote[sourceFragmentIds = [1]]
//	     Layout: []
//
// Plan Fragemnt 1
//
//	Projection[]
//	│ Layout: []
//	│ detail 1
//	│ detail 2
//	└─ ScanFilterProjection[filterPredicate = (role = 'Broker')]
//	   │ Layout: []
//	   │ detail 1
//	   │ detail 2
//	   └─ TableScan[database = _internal]
//	        Layout: [idle, nice, system, user, irq, steal, softirq, iowait]
//	        Partitions: [10.73.59.79:2891=[0]]

var log = logger.GetLogger("Analyzer", "Statement")

type StatementAnalyzer struct {
	ctx             *AnalyzerContext
	metadataMgr     spi.MetadataManager
	funcionResolver *function.FunctionResolver // FIXME:???
}

func NewStatementAnalyzer(ctx *AnalyzerContext, metadataMgr spi.MetadataManager) *StatementAnalyzer {
	return &StatementAnalyzer{
		ctx:             ctx,
		metadataMgr:     metadataMgr,
		funcionResolver: function.NewFunctionResolver(),
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
	isTopLevel      bool
	analyzer        *StatementAnalyzer
}

func NewStatementVisitor(outerQueryScope *Scope, analyzer *StatementAnalyzer, isTopLevel bool) *StatementVisitor {
	return &StatementVisitor{
		outerQueryScope: outerQueryScope,
		analyzer:        analyzer,
		isTopLevel:      isTopLevel,
	}
}

// TODO: check state
func (v *StatementVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.Query:
		return v.visitQuery(context, node)
	case *tree.QuerySpecification:
		return v.visitQuerySpecification(context, node)
	case *tree.Table:
		return v.visitTable(context, node)
	case *tree.AliasedRelation:
		return v.visitAliasedRelation(context, node)
	case *tree.Join:
		return v.visitJoin(context, node)
	case *tree.FunctionCall:
		return v.visitFunctionCall(context, node)
	default:
		panic(fmt.Sprintf("statement analyzer unsupport node: %T", n))
	}
}

func (v *StatementVisitor) visitQuery(context any, node *tree.Query) (r any) {
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

func (v *StatementVisitor) visitQuerySpecification(context any, node *tree.QuerySpecification) (r any) {
	scope := context.(*Scope)
	// analyze from(relation)
	sourceScope := v.analyzeFrom(node, scope)
	if node.Where != nil {
		// analyze where condition
		v.analyzeWhere(node, sourceScope, node.Where)
	}

	outputExpressions := v.analyzeSelect(node, sourceScope)
	groupByAnalysis := v.analyzeGroupBy(node, sourceScope, outputExpressions)
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

	fmt.Printf("select express.....%v\n", selectExpressions)
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

func (v *StatementVisitor) visitJoin(context any, node *tree.Join) (r any) {
	fmt.Println("join table...")
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
		v.analyzeExpression(expression, output)
		// FIXME:
		// panic("impl it")

		v.analyzer.ctx.Analysis.SetJoinCriteria(node, expression)
	}
	fmt.Println("jjjjjjj")
	// fmt.Println(output.RelationType.Fields[0])

	return output
}

func (v *StatementVisitor) analyzeJoinUsing(node *tree.Join, columns []*tree.Identifier,
	scope, left, right *Scope,
) *Scope {
	fmt.Println("fdd..........")
	// return &Scope{
	// 	RelationType: NewRelationType(nil),
	// }
	panic("using")
}

func (v *StatementVisitor) visitAliasedRelation(context any, relation *tree.AliasedRelation) (r any) {
	scope := context.(*Scope)
	aliased := tree.NewQualifiedName([]*tree.Identifier{relation.Aliase})
	v.analyzer.ctx.Analysis.SetRelationName(relation, aliased)
	// v.analyzer.analysis.AddAliased(relation, aliased)

	relationScope := relation.Relation.Accept(scope, v).(*Scope)
	relationType := relationScope.RelationType
	descriptor := relationType.withAlias(relation.Aliase.Value)
	// FIXME: add column

	return v.createAndAssignScope(relation, scope, descriptor)
}

func (v *StatementVisitor) visitTable(ctx any, table *tree.Table) (r any) {
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
	for _, col := range tableMetadata.Schema.Columns {
		// TODO: check agg????
		outputFields = append(outputFields, &tree.Field{
			Name:          col.Name, // TODO: dup tag name/field name
			DataType:      col.DataType,
			RelationAlias: table.Name,
		})
	}

	// TODO: check table type
	v.analyzer.ctx.Analysis.SetRelationName(table, table.Name)
	v.analyzer.ctx.Analysis.RegisterTableMetadata(table, tableMetadata)
	// FIXME: table fields??

	return v.createAndAssignScope(table, scope, NewRelation(TableRelation, outputFields))
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
			// relation := v.analyzer.analysis.GetRelationByAliased(prefix)
			if identifierChain == nil {
				panic(fmt.Sprintf("unable to resolve reference %s", prefix.Name))
			}
			if identifierChain.Type == TABLE {
				relationType := identifierChain.RelationType
				// relationScope := v.analyzer.analysis.GetScope(relation)
				fmt.Println("table========" + prefix.Name)
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
	for i := range relationType.Fields {
		fieldRef := &tree.FieldReference{
			BaseNode: tree.BaseNode{
				ID: v.analyzer.ctx.IDAllocator.Next(),
			},
			FieldIndex: i,
		}
		v.analyzeExpression(fieldRef, scope)
		outputExpressions = append(outputExpressions, fieldRef)
		selectExpressions = append(selectExpressions, &SelectExpression{
			Expression: fieldRef,
		})
	}
	// FIXME: ???
	v.analyzer.ctx.Analysis.SetSelectAllResultFields(allColumns, relationType.Fields)
	return outputExpressions, selectExpressions
}

// ----- relation ------
func (v *StatementVisitor) analyzeFrom(node *tree.QuerySpecification, scope *Scope) *Scope {
	return node.From.Accept(scope, v).(*Scope)
}

func (v *StatementVisitor) analyzeWhere(node tree.Node, scope *Scope, predicate tree.Expression) {
	// FIXME: verify no aggregate and group by function
	fmt.Println("analyze where")
	v.analyzeExpression(predicate, scope)
	// v.analyzer.ctx.Analysis.RecordSubQueries(node, expressionAnalysis)

	// FIXME: check predicate type
	// predicateType := expressionAnalysis.GetType(predicate)

	v.analyzer.ctx.Analysis.SetWhere(node, predicate)
}

func (v *StatementVisitor) analyzeGroupBy(node *tree.QuerySpecification, scope *Scope,
	outputExpressions []tree.Expression,
) *GroupingSetAnalysis {
	if node.GroupBy != nil {
		var (
			groupingExpressions []tree.Expression
			complexExpressions  []tree.Expression
			sets                [][]*FieldID
		)

		for _, groupingElement := range node.GroupBy.GroupingElements {
			switch groupByEle := groupingElement.(type) {
			// TODO: gropu by *
			case *tree.SimpleGroupBy:
				for _, column := range groupByEle.Columns {
					switch column.(type) {
					case *tree.LongLiteral:
						// TODO: fixme index field
						panic("impl long group key")
					default:
						v.analyzeExpression(column, scope)
					}

					field := v.analyzer.ctx.Analysis.GetColumnReferenceField(column)
					if field != nil {
						sets = append(sets, []*FieldID{field.FieldID()})
					} else {
						// TODO: field sets
						complexExpressions = append(complexExpressions, column)
					}

					groupingExpressions = append(groupingExpressions, column)
				}
			}
		}

		groupingSets := NewGroupingSetAnalysis(groupingExpressions, sets, complexExpressions)
		v.analyzer.ctx.Analysis.SetGroupingSets(node, groupingSets)

		return groupingSets
	}
	// TODO: has aggs
	return nil
}

func (v *StatementVisitor) analyzeGroupingOperations(node *tree.QuerySpecification,
	outputExpressions, orderByExpressions []tree.Expression) {
}

func (v *StatementVisitor) analyzeAggregations(node *tree.QuerySpecification, sourceScope, orderByScope *Scope,
	groupByAnalysis *GroupingSetAnalysis, outputExpressions, orderByExpressions []tree.Expression,
) {
	var expr []tree.Expression
	expr = append(expr, outputExpressions...)
	expr = append(expr, orderByExpressions...)
	// TODO:
	ExtractAggregationFunctions(expr, v.analyzer.funcionResolver)
	for _, selectExpr := range outputExpressions {
		if ident, ok := selectExpr.(*tree.Identifier); ok {
			// transfer filed builtin aggregation
			resolvedField := sourceScope.resolveField(node, tree.NewQualifiedName([]*tree.Identifier{ident}), true)
			if resolvedField.Field.DataType.CanAggregatin() {
				fn := &tree.FunctionCall{
					Name:      tree.QualifiedName{Suffix: resolvedField.Field.DataType.String()},
					Arguments: []tree.Expression{selectExpr},
					RefField:  resolvedField.Field,
				}
				v.analyzer.ctx.Analysis.SetAggregates(node, []*tree.FunctionCall{fn})
				resolvedFn := v.analyzer.funcionResolver.ResolveFunction(&fn.Name)
				v.analyzer.ctx.Analysis.AddResolvedFunction(fn, resolvedFn)
				v.analyzer.ctx.Analysis.AddType(fn, resolvedFn.Signature.ReturnType)
			}
		}
	}
	// TODO: extract agg func
	if v.analyzer.ctx.Analysis.IsGroupingSets(node) {
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
}

func (v *StatementVisitor) analyzeOrderBy(node tree.Node,
	sortItems []*tree.SortItem, orderByScope *Scope,
) (orderByExpressions []tree.Expression) {
	return
}

func (v *StatementVisitor) analyzeLimit(node *tree.Limit, scope *Scope) {
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

	return v.createAndAssignScope(table, scope, NewRelation(UnknownRelation, fields))
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
	scope, sourceScope *Scope,
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
			field := item.Aliase
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

			outputFields = append(outputFields, &tree.Field{
				Name: fieldName,
			})
		default:
			panic(fmt.Sprintf("unsupported selec type type: %s", reflect.TypeOf(item)))
		}
	}
	return v.createAndAssignScope(node, scope, NewRelation(UnknownRelation, outputFields))
}

func (v *StatementVisitor) computeAndAssignOrderByScope(node *tree.OrderBy,
	sourceScope, outputSource *Scope, fields []*tree.Field,
) *Scope {
	return &Scope{}
}

func (v *StatementVisitor) descriptorToFields(scope *Scope) (selectExpressions []*SelectExpression) {
	for i := range scope.RelationType.Fields {
		expression := &tree.FieldReference{
			BaseNode: tree.BaseNode{
				ID: v.analyzer.ctx.IDAllocator.Next(),
			},
			FieldIndex: i,
		}
		selectExpressions = append(selectExpressions, &SelectExpression{
			Expression: expression,
		})
		v.analyzeExpression(expression, scope)
	}
	return
}
