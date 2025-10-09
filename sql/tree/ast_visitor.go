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

package tree

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
)

// for testing
var (
	newNodeLocation = NewNodeLocation
)

type Visitor interface {
	Visit(context any, node Node) (r any)
}

type AstVisitor struct {
	grammar.BaseSQLParserVisitor

	idAllocator *NodeIDAllocator
}

func NewAstVisitor(idAllocator *NodeIDAllocator) *AstVisitor {
	return &AstVisitor{idAllocator: idAllocator}
}

func (v *AstVisitor) Visit(ctx antlr.ParseTree) any {
	return ctx.Accept(v)
}

func (v *AstVisitor) VisitStatement(ctx *grammar.StatementContext) any {
	switch {
	case ctx.DmlStatement() != nil:
		return v.Visit(ctx.DmlStatement())
	case ctx.DdlStatement() != nil:
		return v.Visit(ctx.DdlStatement())
	case ctx.UtilityStatement() != nil:
		return v.Visit(ctx.UtilityStatement())
	case ctx.AdminStatement() != nil:
		return v.Visit(ctx.AdminStatement())
	default:
		return v.VisitChildren(ctx)
	}
}

func (v *AstVisitor) VisitAdminStatement(ctx *grammar.AdminStatementContext) any {
	switch {
	case ctx.FlushDatabase() != nil:
		return v.Visit(ctx.FlushDatabase())
	case ctx.CompactDatabase() != nil:
		return v.Visit(ctx.CompactDatabase())
	case ctx.ShowStatement() != nil:
		return v.Visit(ctx.ShowStatement())
	default:
		return v.VisitChildren(ctx)
	}
}

func (v *AstVisitor) VisitFlushDatabase(ctx *grammar.FlushDatabaseContext) any {
	return &FlushDatabase{
		BaseNode: v.createBaseNode(ctx),
		Database: v.getQualifiedName(ctx.GetDatabase()).Name,
	}
}

func (v *AstVisitor) VisitCompactDatabase(ctx *grammar.CompactDatabaseContext) any {
	return &CompactDatabase{
		BaseNode: v.createBaseNode(ctx),
		Database: v.getQualifiedName(ctx.GetDatabase()).Name,
	}
}

func (v *AstVisitor) VisitShowDatabases(ctx *grammar.ShowDatabasesContext) any {
	show := &ShowDatabases{
		BaseNode: v.createBaseNode(ctx),
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitShowNamespaces(ctx *grammar.ShowNamespacesContext) any {
	show := &ShowNamespaces{}
	if ctx.GetNamespace() != nil {
		value, err := strutil.GetStringValue(ctx.GetNamespace().GetText())
		if err != nil {
			panic(err)
		}
		show.LikePattern = value
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitShowTableNames(ctx *grammar.ShowTableNamesContext) any {
	show := &ShowTableNames{}
	if ctx.QualifiedName() != nil {
		show.Namespace = v.getQualifiedName(ctx.QualifiedName())
	}
	if ctx.GetTableName() != nil {
		value, err := strutil.GetStringValue(ctx.GetTableName().GetText())
		if err != nil {
			panic(err)
		}
		show.LikePattern = value
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitShowColumns(ctx *grammar.ShowColumnsContext) any {
	show := &ShowColumns{
		Table: &Table{
			BaseNode: v.createBaseNode(ctx),
			Name:     v.getQualifiedName(ctx.QualifiedName()),
		},
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitShowReplications(ctx *grammar.ShowReplicationsContext) any {
	show := &ShowReplications{
		BaseNode: v.createBaseNode(ctx),
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitShowMemoryDatabases(ctx *grammar.ShowMemoryDatabasesContext) any {
	show := &ShowMemoryDatabases{
		BaseNode: v.createBaseNode(ctx),
	}
	return &Show{
		BaseNode: v.createBaseNode(ctx),
		Body:     show,
	}
}

func (v *AstVisitor) VisitDdlStatement(ctx *grammar.DdlStatementContext) any {
	switch {
	case ctx.CreateDatabase() != nil:
		return v.Visit(ctx.CreateDatabase())
	case ctx.DropDatabase() != nil:
		return v.Visit(ctx.DropDatabase())
	case ctx.CreateBroker() != nil:
		panic("need impl create broker")
	default:
		return v.VisitChildren(ctx)
	}
}

func (v *AstVisitor) VisitDropDatabase(ctx *grammar.DropDatabaseContext) any {
	return &DropDatabase{
		BaseNode: v.createBaseNode(ctx),
		Name:     v.getQualifiedName(ctx.QualifiedName()).Name,
		Exists:   ctx.EXISTS() != nil,
	}
}

func (v *AstVisitor) VisitCreateDatabase(ctx *grammar.CreateDatabaseContext) any {
	createDatabase := &CreateDatabase{
		BaseNode: v.createBaseNode(ctx),
		Name:     v.getQualifiedName(ctx.GetName()).Name,
	}
	options := ctx.AllDatabaseOptions()
	for _, option := range options {
		switch opt := (v.Visit(option)).(type) {
		case []CreateOption:
			createDatabase.CreateOptions = append(createDatabase.CreateOptions, opt...)
		case []*RollupOption:
			createDatabase.Rollup = append(createDatabase.Rollup, opt...)
		case []*Property:
			createDatabase.Props = append(createDatabase.Props, opt...)
		}
	}
	fmt.Printf("props=%v,rollup=%v\n", createDatabase.Props, createDatabase.Rollup)
	return createDatabase
}

func (v *AstVisitor) VisitDatabaseOpts(ctx *grammar.DatabaseOptsContext) any {
	return visit[CreateOption](ctx.AllCreateDatabaseOptions(), v)
}

func (v *AstVisitor) VisitWithProps(ctx *grammar.WithPropsContext) any {
	return visit[*Property](ctx.Properties().PropertyAssignments().AllProperty(), v)
}

func (v *AstVisitor) VisitRollupProps(ctx *grammar.RollupPropsContext) any {
	fmt.Println("rollup props")
	return visit[*RollupOption](ctx.AllRollupOptions(), v)
}

func (v *AstVisitor) VisitEngineOption(ctx *grammar.EngineOptionContext) any {
	engineType := option.Metric

	switch {
	case ctx.METRIC() != nil:
		engineType = option.Metric
	case ctx.LOG() != nil:
		engineType = option.Log
	case ctx.TRACE() != nil:
		engineType = option.Trace

	}
	return &EngineOption{
		Type: engineType,
	}
}

func (v *AstVisitor) VisitRollupOptions(ctx *grammar.RollupOptionsContext) any {
	fmt.Println("rollup props options")
	return &RollupOption{
		BaseNode: v.createBaseNode(ctx),
		Props:    visit[*Property](ctx.Properties().PropertyAssignments().AllProperty(), v),
	}
}

func (v *AstVisitor) VisitProperty(ctx *grammar.PropertyContext) any {
	fmt.Println("visit property.....")
	return &Property{
		BaseNode: v.createBaseNode(ctx),
		Name:     visitIfPresent[*Identifier](ctx.GetName(), v),
		Value:    visitIfPresent[Expression](ctx.GetValue(), v),
	}
}

func (v *AstVisitor) VisitDefaultPropertyValue(ctx *grammar.DefaultPropertyValueContext) any {
	return nil
}

func (v *AstVisitor) VisitNonDefaultPropertyValue(ctx *grammar.NonDefaultPropertyValueContext) any {
	return visitIfPresent[Expression](ctx.Expression(), v)
}

func (v *AstVisitor) VisitUtilityStatement(ctx *grammar.UtilityStatementContext) any {
	switch {
	case ctx.UseStatement() != nil:
		return v.Visit(ctx.UseStatement())
	default:
		panic("unsupported utility statement")
	}
}

func (v *AstVisitor) VisitUseStatement(ctx *grammar.UseStatementContext) any {
	identifier := v.Visit(ctx.GetDatabase()).(*Identifier)
	return &Use{
		BaseNode: v.createBaseNode(ctx),
		Database: identifier,
	}
}

// VisitStatementDefault visits default statement(query statement).
func (v *AstVisitor) VisitStatementDefault(ctx *grammar.StatementDefaultContext) any {
	if ctx.Query() != nil {
		return v.Visit(ctx.Query())
	}
	return v.VisitChildren(ctx)
}

func (v *AstVisitor) VisitExplain(ctx *grammar.ExplainContext) any {
	return &Explain{
		BaseNode:  v.createBaseNode(ctx),
		Options:   visit[ExplainOption](ctx.AllExplainOption(), v),
		Statement: visitIfPresent[Statement](ctx.DmlStatement(), v),
	}
}

func (v *AstVisitor) VisitExplainType(ctx *grammar.ExplainTypeContext) any {
	val := LogicalExplain
	if ctx.DISTRIBUTED() != nil {
		val = DistributedExplain
	}
	return &ExplainType{
		Type: val,
	}
}

func (v *AstVisitor) VisitExplainAnalyze(ctx *grammar.ExplainAnalyzeContext) any {
	return v.VisitChildren(ctx)
}

func (v *AstVisitor) VisitQuery(ctx *grammar.QueryContext) any {
	query := v.Visit(ctx.QueryNoWith()).(*Query)
	return &Query{
		BaseNode:  v.createBaseNode(ctx),
		With:      visitIfPresent[*With](ctx.With(), v),
		QueryBody: query.QueryBody,
		OrderBy:   query.OrderBy,
		Limit:     query.Limit,
	}
}

func (v *AstVisitor) VisitWith(ctx *grammar.WithContext) any {
	return &With{
		BaseNode: v.createBaseNode(ctx),
		Queries:  visit[*WithQuery](ctx.AllNamedQuery(), v),
	}
}

func (v *AstVisitor) VisitNamedQuery(ctx *grammar.NamedQueryContext) any {
	return &WithQuery{
		BaseNode: v.createBaseNode(ctx),
		Name:     visitIfPresent[*Identifier](ctx.GetName(), v),
		Query:    visitIfPresent[*Query](ctx.Query(), v),
	}
}

func (v *AstVisitor) VisitQueryNoWith(ctx *grammar.QueryNoWithContext) any {
	term := visitIfPresent[QueryBody](ctx.QueryTerm(), v)
	var orderBy *OrderBy
	if ctx.ORDER() != nil {
		orderBy = &OrderBy{
			BaseNode:  v.createBaseNode(ctx),
			SortItems: visit[*SortItem](ctx.OrderBy().AllSortItem(), v),
		}
	}
	var limit *Limit
	if ctx.LIMIT() != nil {
		// TODO: all
		var rowCount Expression
		if ctx.LimitRowCount().INTEGER_VALUE() != nil {
			rowCount = NewLongLiteral(
				v.idAllocator.Next(),
				getLocation(ctx.LimitRowCount()),
				ctx.LimitRowCount().GetText())
		}
		limit = &Limit{
			BaseNode: v.createBaseNode(ctx),
			RowCount: rowCount,
		}
	}
	if query, ok := term.(*QuerySpecification); ok {
		// When we have a simple query specification
		// followed by order by, limit,
		// fold the order by, limit clauses
		// into the query specification (analyzer/planner
		// expects this structure to resolve references with respect
		// to columns defined in the query specification)
		return &Query{
			BaseNode: v.createBaseNode(ctx),
			QueryBody: &QuerySpecification{
				BaseNode: v.createBaseNode(ctx),
				Select:   query.Select,
				From:     query.From,
				Where:    query.Where,
				GroupBy:  query.GroupBy,
				Having:   query.Having,
				OrderBy:  orderBy,
				Limit:    limit,
			},
		}
	}
	return &Query{
		BaseNode:  v.createBaseNode(ctx),
		QueryBody: term,
		OrderBy:   orderBy,
		Limit:     limit,
	}
}

func (v *AstVisitor) VisitQueryTermDefault(ctx *grammar.QueryTermDefaultContext) any {
	return v.Visit(ctx.QueryPrimary())
}

func (v *AstVisitor) VisitQueryPrimaryDefault(ctx *grammar.QueryPrimaryDefaultContext) any {
	return v.Visit(ctx.QuerySpecification())
}

func (v *AstVisitor) VisitQuerySpecification(ctx *grammar.QuerySpecificationContext) any {
	// parse select items
	selectItems := visit[SelectItem](ctx.AllSelectItem(), v)
	// parse relations
	relations := visit[Relation](ctx.AllRelation(), v)
	var from Relation
	if len(relations) > 0 {
		// synthesize implicit join nodes
		relation := relations[0]
		i := 1
		for i < len(relations) {
			relation = &Join{
				Type:  IMPLICIT,
				Left:  relation,
				Right: relations[i],
			}
			i++
		}
		from = relation
	}

	return &QuerySpecification{
		Select: &Select{
			SelectItems: selectItems,
		},
		From:    from,
		Where:   visitIfPresent[Expression](ctx.GetWhere(), v),
		GroupBy: visitIfPresent[*GroupBy](ctx.GroupBy(), v),
		Having:  visitIfPresent[Expression](ctx.Having(), v),
	}
}

func (v *AstVisitor) VisitSelectAll(ctx *grammar.SelectAllContext) any {
	return &AllColumns{
		BaseNode: v.createBaseNode(ctx),
		Target:   visitIfPresent[Expression](ctx.PrimaryExpression(), v),
	}
}

func (v *AstVisitor) VisitSelectSingle(ctx *grammar.SelectSingleContext) any {
	return &SingleColumn{
		BaseNode:   v.createBaseNode(ctx),
		Expression: visitIfPresent[Expression](ctx.Expression(), v),
		Aliase:     visitIfPresent[*Identifier](ctx.Identifier(), v),
	}
}

func (v *AstVisitor) VisitJoinRelation(ctx *grammar.JoinRelationContext) any {
	left := v.Visit(ctx.GetLeft()).(Relation)
	if ctx.CROSS() != nil {
		// prase cross join
		right := v.Visit(ctx.GetRight()).(Relation)
		return &Join{
			BaseNode: v.createBaseNode(ctx),
			Type:     CROSS,
			Left:     left,
			Right:    right,
		}
	}
	// parse left/right/inner join
	right := v.Visit(ctx.GetRightRelation()).(Relation)
	var joinCriteria JoinCriteria
	switch {
	case ctx.JoinCriteria().ON() != nil:
		expression := v.Visit(ctx.JoinCriteria().BooleanExpression()).(Expression)
		joinCriteria = &JoinOn{
			Expression: expression,
		}
	case ctx.JoinCriteria().USING() != nil:
		joinCriteria = &JoinUsing{
			Columns: visit[*Identifier](ctx.JoinCriteria().AllIdentifier(), v),
		}
	default:
		panic("unsupported join criteria")
	}
	var joinType JoinType
	switch {
	case ctx.JoinType().LEFT() != nil:
		joinType = LEFT
	case ctx.JoinType().RIGHT() != nil:
		joinType = RIGHT
	default:
		joinType = INNER
	}
	return &Join{
		BaseNode: v.createBaseNode(ctx),
		Type:     joinType,
		Left:     left,
		Right:    right,
		Criteria: joinCriteria,
	}
}

func (v *AstVisitor) VisitRelationDefault(ctx *grammar.RelationDefaultContext) any {
	return v.Visit(ctx.AliasedRelation())
}

func (v *AstVisitor) VisitTableName(ctx *grammar.TableNameContext) any {
	return &Table{
		BaseNode: v.createBaseNode(ctx),
		Name:     v.getQualifiedName(ctx.QualifiedName()),
	}
}

func (v *AstVisitor) VisitSubQueryRelation(ctx *grammar.SubQueryRelationContext) any {
	query := v.Visit(ctx.Query()).(*Query)
	return &TableSubQuery{
		BaseNode: v.createBaseNode(ctx),
		Query:    query,
	}
}

func (v *AstVisitor) VisitAliasedRelation(ctx *grammar.AliasedRelationContext) any {
	child := v.Visit(ctx.RelationPrimary()).(Relation)
	if ctx.Identifier() == nil {
		return child
	}
	// parese relation aliase
	identifier := v.Visit(ctx.Identifier()).(*Identifier)
	return &AliasedRelation{
		BaseNode: v.createBaseNode(ctx),
		Relation: child,
		Aliase:   identifier,
	}
}

func (v *AstVisitor) VisitBinaryComparisonPredicate(ctx *grammar.BinaryComparisonPredicateContext) any {
	return &ComparisonExpression{
		BaseNode: v.createBaseNode(ctx),
		Operator: ComparisonOperator(ctx.GetOperator().GetText()), // FIXME:
		Left:     visitIfPresent[Expression](ctx.GetLeft(), v),
		Right:    visitIfPresent[Expression](ctx.GetRight(), v),
	}
}

func (v *AstVisitor) VisitRegexpPredicate(ctx *grammar.RegexpPredicateContext) any {
	var result Expression
	result = &RegexPredicate{
		BaseNode: v.createBaseNode(ctx),
		Value:    visitIfPresent[Expression](ctx.GetLeft(), v),
		Pattern:  visitIfPresent[Expression](ctx.GetPattern(), v),
	}
	if ctx.NEQREGEXP() != nil {
		result = &NotExpression{
			BaseNode: v.createBaseNode(ctx),
			Value:    result,
		}
	}
	return result
}

func (v *AstVisitor) VisitTimestampPredicate(ctx *grammar.TimestampPredicateContext) any {
	return &TimePredicate{
		BaseNode: v.createBaseNode(ctx),
		Operator: ComparisonOperator(ctx.GetOperator().GetText()), // FIXME:
		Value:    visitIfPresent[Expression](ctx.ValueExpression(), v),
	}
}

func (v *AstVisitor) VisitNullPredicate(ctx *grammar.NullPredicateContext) any {
	predicate := &NullPredicate{
		BaseNode: v.createBaseNode(ctx),
		Value:    visitIfPresent[Expression](ctx.GetLeft(), v),
	}
	if ctx.NOT() != nil {
		predicate.Not = true
	}
	return predicate
}

func (v *AstVisitor) VisitLikePredicate(ctx *grammar.LikePredicateContext) any {
	var result Expression
	result = &LikePredicate{
		BaseNode: v.createBaseNode(ctx),
		Value:    visitIfPresent[Expression](ctx.GetLeft(), v),
		Pattern:  visitIfPresent[Expression](ctx.GetPattern(), v),
	}
	if ctx.NOT() != nil {
		result = &NotExpression{
			BaseNode: v.createBaseNode(ctx),
			Value:    result,
		}
	}
	return result
}

func (v *AstVisitor) VisitInPredicate(ctx *grammar.InPredicateContext) any {
	var result Expression
	result = &InPredicate{
		BaseNode: v.createBaseNode(ctx),
		Value:    visitIfPresent[Expression](ctx.GetLeft(), v),
		ValueList: &InListExpression{
			BaseNode: v.createBaseNode(ctx),
			Values:   visit[Expression](ctx.AllExpression(), v),
		},
	}

	if ctx.NOT() != nil {
		result = &NotExpression{
			BaseNode: v.createBaseNode(ctx),
			Value:    result,
		}
	}
	return result
}

func (v *AstVisitor) VisitLogicalNot(ctx *grammar.LogicalNotContext) any {
	return &NotExpression{
		BaseNode: v.createBaseNode(ctx),
		Value:    visitIfPresent[Expression](ctx.BooleanExpression(), v),
	}
}

func (v *AstVisitor) VisitOr(ctx *grammar.OrContext) any {
	terms := v.flatten(ctx, func(parentCtx antlr.ParserRuleContext) (rs []antlr.ParserRuleContext) {
		if or, ok := parentCtx.(*grammar.OrContext); ok {
			expressions := or.AllBooleanExpression()
			for _, expression := range expressions {
				rs = append(rs, expression)
			}
		}
		return
	})
	return &LogicalExpression{
		BaseNode: v.createBaseNode(ctx),
		Operator: LogicalOR,
		Terms:    visit[Expression](terms, v),
	}
}

func (v *AstVisitor) VisitAnd(ctx *grammar.AndContext) any {
	terms := v.flatten(ctx, func(parentCtx antlr.ParserRuleContext) (rs []antlr.ParserRuleContext) {
		if and, ok := parentCtx.(*grammar.AndContext); ok {
			expressions := and.AllBooleanExpression()
			for _, expression := range expressions {
				rs = append(rs, expression)
			}
		}
		return
	})
	return &LogicalExpression{
		BaseNode: v.createBaseNode(ctx),
		Operator: LogicalAND,
		Terms:    visit[Expression](terms, v),
	}
}

func (v *AstVisitor) flatten(root antlr.ParserRuleContext,
	extractChildren func(ctx antlr.ParserRuleContext) []antlr.ParserRuleContext,
) (result []antlr.ParserRuleContext) {
	pending := collections.NewStack()
	pending.Push(root)
	for pending.Size() > 0 {
		next := pending.Pop().(antlr.ParserRuleContext)
		children := extractChildren(next)
		if len(children) == 0 {
			result = append(result, next)
		} else {
			for i := len(children) - 1; i >= 0; i-- {
				pending.Push(children[i])
			}
		}
	}
	return
}

func (v *AstVisitor) VisitParenExpression(ctx *grammar.ParenExpressionContext) any {
	return v.Visit(ctx.Expression())
}

func (v *AstVisitor) VisitGroupBy(ctx *grammar.GroupByContext) any {
	return &GroupBy{
		BaseNode:         v.createBaseNode(ctx),
		GroupingElements: visit[GroupingElement](ctx.AllGroupingElement(), v),
	}
}

func (v *AstVisitor) VisitSingleGroupingSet(ctx *grammar.SingleGroupingSetContext) any {
	return &SimpleGroupBy{
		BaseNode: v.createBaseNode(ctx),
		Columns:  visit[Expression](ctx.GroupingSet().AllExpression(), v),
	}
}

func (v *AstVisitor) VisitSortItem(ctx *grammar.SortItemContext) any {
	expression := v.Visit(ctx.Expression()).(Expression)
	return &SortItem{
		BaseNode: v.createBaseNode(ctx),
		SortKey:  expression,
		Ordering: getOrderingType(ctx),
	}
}

func (v *AstVisitor) VisitUnquotedIdentifier(ctx *grammar.UnquotedIdentifierContext) any {
	return &Identifier{
		BaseNode:  v.createBaseNode(ctx),
		Value:     ctx.GetText(),
		Delimited: false,
	}
}

func (v *AstVisitor) VisitQuotedIdentifier(ctx *grammar.QuotedIdentifierContext) any {
	identifier, err := strutil.GetStringValue(ctx.GetText())
	if err != nil {
		panic(err)
	}
	return &Identifier{
		BaseNode:  v.createBaseNode(ctx),
		Value:     identifier,
		Delimited: true,
	}
}

func (v *AstVisitor) VisitPredicatedExpression(ctx *grammar.PredicatedExpressionContext) any {
	return v.Visit(ctx.Predicate())
}

func (v *AstVisitor) VisitValueExpressionDefault(ctx *grammar.ValueExpressionDefaultContext) any {
	return v.Visit(ctx.PrimaryExpression())
}

func (v *AstVisitor) VisitValueExpressionPredicate(ctx *grammar.ValueExpressionPredicateContext) any {
	if ctx.ValueExpression() != nil {
		return v.Visit(ctx.ValueExpression())
	}
	fmt.Printf("value path...=%v\n", ctx)
	return v.VisitChildren(ctx)
}

func (v *AstVisitor) VisitDereference(ctx *grammar.DereferenceContext) any {
	return &DereferenceExpression{
		BaseNode: v.createBaseNode(ctx),
		Base:     visitIfPresent[Expression](ctx.GetBase(), v),
		Field:    visitIfPresent[*Identifier](ctx.GetFieldName(), v),
	}
}

func (v *AstVisitor) VisitColumnReference(ctx *grammar.ColumnReferenceContext) any {
	return v.Visit(ctx.Identifier())
}

func (v *AstVisitor) VisitExpression(ctx *grammar.ExpressionContext) any {
	return v.Visit(ctx.BooleanExpression())
}

func (v *AstVisitor) VisitFunctionCall(ctx *grammar.FunctionCallContext) any {
	// FIXME: parse funcion call
	funcName := FuncName(strings.ToLower(v.getQualifiedName(ctx.QualifiedName()).Name)) // TODO: check function name
	return &FunctionCall{
		BaseNode:  v.createBaseNode(ctx),
		Name:      funcName,
		RetType:   GetDefaultFuncReturnType(funcName),
		Arguments: visit[Expression](ctx.AllExpression(), v),
	}
}

func (v *AstVisitor) VisitArithmeticBinary(ctx *grammar.ArithmeticBinaryContext) any {
	return &ArithmeticBinaryExpression{
		BaseNode: v.createBaseNode(ctx),
		Operator: ArithmeticOperator(ctx.GetOperator().GetText()), // TODO: add check
		Left:     v.Visit(ctx.GetLeft()).(Expression),
		Right:    v.Visit(ctx.GetRight()).(Expression),
	}
}

// ************** visit literals **************

func (v *AstVisitor) VisitStringLiteral(ctx *grammar.StringLiteralContext) any {
	return v.Visit(ctx.String_())
}

func (v *AstVisitor) VisitNumericLiteral(ctx *grammar.NumericLiteralContext) any {
	return v.Visit(ctx.Number())
}

func (v *AstVisitor) VisitIntervalLiteral(ctx *grammar.IntervalLiteralContext) any {
	return NewIntervalLiteral(v.idAllocator.Next(), getLocation(ctx),
		ctx.Interval().GetValue().GetText(), IntervalUnit(strings.ToUpper(ctx.Interval().GetUnit().GetText())))
}

func (v *AstVisitor) VisitBasicStringLiteral(ctx *grammar.BasicStringLiteralContext) any {
	value, err := strutil.GetStringValue(ctx.STRING().GetText())
	if err != nil {
		panic(err)
	}
	return &StringLiteral{
		BaseNode: v.createBaseNode(ctx),
		Value:    value,
	}
}

func (v *AstVisitor) VisitBooleanLiteral(ctx *grammar.BooleanLiteralContext) any {
	return NewBooleanLiteral(v.idAllocator.Next(), getLocation(ctx), ctx.GetText())
}

func (v *AstVisitor) VisitIntegerLiteral(ctx *grammar.IntegerLiteralContext) any {
	return NewLongLiteral(v.idAllocator.Next(), getLocation(ctx), ctx.GetText())
}

func (v *AstVisitor) VisitDecimalLiteral(ctx *grammar.DecimalLiteralContext) any {
	return NewFloatLiteral(v.idAllocator.Next(), getLocation(ctx), ctx.GetText())
}

func (v *AstVisitor) VisitDoubleLiteral(ctx *grammar.DoubleLiteralContext) any {
	return NewFloatLiteral(v.idAllocator.Next(), getLocation(ctx), ctx.GetText())
}

func (v *AstVisitor) getQualifiedName(ctx grammar.IQualifiedNameContext) *QualifiedName {
	parts := visit[*Identifier](ctx.AllIdentifier(), v)
	return NewQualifiedName(parts)
}

func (v *AstVisitor) createBaseNode(ctx antlr.ParserRuleContext) BaseNode {
	return BaseNode{
		ID:       v.idAllocator.Next(),
		Location: getLocation(ctx),
	}
}

func visit[R any, C antlr.ParserRuleContext](contexts []C, visitor grammar.SQLParserVisitor) (r []R) {
	for _, ctx := range contexts {
		result := visitor.Visit(ctx)
		if result != nil {
			r = append(r, result.(R))
		}
	}
	return
}

func visitIfPresent[R any, C antlr.ParserRuleContext](ctx C, visitor grammar.SQLParserVisitor) (r R) {
	rv := reflect.ValueOf(ctx)
	if rv.Kind() == reflect.Invalid || (rv.Kind() != reflect.Invalid && rv.IsNil()) {
		return
	}
	result := visitor.Visit(ctx)
	if result != nil {
		if rr, ok := result.(R); ok {
			r = rr
		}
	}
	return
}

func getLocation(ctx antlr.ParserRuleContext) *NodeLocation {
	token := ctx.GetStart()
	return newNodeLocation(token.GetLine(), token.GetTokenSource().GetCharPositionInLine())
}

func getOrderingType(ctx *grammar.SortItemContext) Ordering {
	if ctx.DESC() != nil {
		return DESCENDING
	}
	return ASCENDING
}
