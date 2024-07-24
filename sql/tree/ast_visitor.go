package tree

import (
	"fmt"
	"reflect"

	"github.com/antlr4-go/antlr/v4"

	"github.com/lindb/lindb/pkg/collections"
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
	default:
		return v.VisitChildren(ctx)
	}
}

func (v *AstVisitor) VisitDdlStatement(ctx *grammar.DdlStatementContext) any {
	switch {
	case ctx.CreateDatabase() != nil:
		createDatabaseCtx := ctx.CreateDatabase()
		props := createDatabaseCtx.Properties()
		if props != nil {
			props.Accept(v)
		}
		return &CreateDatabase{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Name: v.getQualifiedName(createDatabaseCtx.GetName()).Name,
		}
	case ctx.CreateBroker() != nil:
		panic("need impl create broker")
	default:
		return v.VisitChildren(ctx)
	}
}

func (v *AstVisitor) VisitProperties(ctx *grammar.PropertiesContext) any {
	for _, prop := range ctx.PropertyAssignments().AllProperty() {
		prop.Accept(v)
	}
	fmt.Println("test.....props")
	return nil
}

func (v *AstVisitor) VisitProperty(ctx *grammar.PropertyContext) any {
	identifer := v.Visit(ctx.GetName()).(*Identifier)
	fmt.Println(reflect.TypeOf(ctx.GetValue()))
	fmt.Println(identifer.Value)
	fmt.Println("test.....")
	return nil
}

func (v *AstVisitor) VisitUseStatement(ctx *grammar.UseStatementContext) any {
	identifer := v.Visit(ctx.GetDatabase()).(*Identifier)
	return &Use{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Database: identifer,
	}
}

func (v *AstVisitor) VisitDmlStatement(ctx *grammar.DmlStatementContext) any {
	if ctx.Query() != nil {
		return v.Visit(ctx.Query())
	}
	return v.VisitChildren(ctx)
}

func (v *AstVisitor) VisitQuery(ctx *grammar.QueryContext) any {
	query := v.Visit(ctx.QueryNoWith()).(*Query)
	return &Query{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		With:      visitIfPresent[*With](ctx.With(), v),
		QueryBody: query.QueryBody,
		OrderBy:   query.OrderBy,
		Limit:     query.Limit,
	}
}

func (v *AstVisitor) VisitWith(ctx *grammar.WithContext) any {
	return &With{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Queries: visit[*WithQuery](ctx.AllNamedQuery(), v),
	}
}

func (v *AstVisitor) VisitNamedQuery(ctx *grammar.NamedQueryContext) any {
	identifer := v.Visit(ctx.GetName()).(*Identifier)
	query := v.Visit(ctx.Query()).(*Query)
	return &WithQuery{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Name:  identifer,
		Query: query,
	}
}

func (v *AstVisitor) VisitQueryNoWith(ctx *grammar.QueryNoWithContext) any {
	term := v.Visit(ctx.QueryTerm()).(QueryBody)
	var orderBy *OrderBy
	if ctx.ORDER() != nil {
		orderBy = &OrderBy{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.ORDER().GetSymbol()),
			},
			SortItems: visit[*SortItem](ctx.OrderBy().AllSortItem(), v),
		}
	}
	var limit *Limit
	if ctx.LIMIT() != nil {
		// TODO: all
		var rowCount Expression
		if ctx.LimitRowCount().INTEGER_VALUE() != nil {
			rowCount = NewLongLiteral(
				getLocation(ctx.LimitRowCount().INTEGER_VALUE().GetSymbol()),
				ctx.LimitRowCount().GetText())
		}
		limit = &Limit{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.LIMIT().GetSymbol()),
			},
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
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			QueryBody: &QuerySpecification{
				BaseNode: BaseNode{
					ID:       v.idAllocator.Next(),
					Location: getLocation(ctx.GetStart()),
				},
				Select:  query.Select,
				From:    query.From,
				Where:   query.Where,
				GroupBy: query.GroupBy,
				Having:  query.Having,
				OrderBy: orderBy,
				Limit:   limit,
			},
		}
	}
	return &Query{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
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
	} else {
		panic("relation cannot be empty")
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
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Target: visitIfPresent[Expression](ctx.PrimaryExpression(), v),
	}
}

func (v *AstVisitor) VisitSelectSingle(ctx *grammar.SelectSingleContext) any {
	expression := v.Visit(ctx.Expression()).(Expression)
	return &SingleColumn{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Expression: expression,
		Aliase:     visitIfPresent[*Identifier](ctx.Identifier(), v),
	}
}

func (v *AstVisitor) VisitJoinRelation(ctx *grammar.JoinRelationContext) any {
	left := v.Visit(ctx.GetLeft()).(Relation)
	if ctx.CROSS() != nil {
		// prase cross join
		right := v.Visit(ctx.GetRight()).(Relation)
		return &Join{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Type:  CROSS,
			Left:  left,
			Right: right,
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
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
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
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Name: v.getQualifiedName(ctx.QualifiedName()),
	}
}

func (v *AstVisitor) VisitSubQueryRelation(ctx *grammar.SubQueryRelationContext) any {
	query := v.Visit(ctx.Query()).(*Query)
	return &TableSubQuery{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Query: query,
	}
}

func (v *AstVisitor) VisitAliasedRelation(ctx *grammar.AliasedRelationContext) any {
	child := v.Visit(ctx.RelationPrimary()).(Relation)
	if ctx.Identifier() == nil {
		return child
	}
	// parese relation aliase
	identifer := v.Visit(ctx.Identifier()).(*Identifier)
	return &AliasedRelation{
		BaseNode: BaseNode{
			Location: getLocation(ctx.GetStart()),
			ID:       v.idAllocator.Next(),
		},
		Relation: child,
		Aliase:   identifer,
	}
}

func (v *AstVisitor) VisitBinaryComparisonPredicate(ctx *grammar.BinaryComparisonPredicateContext) any {
	left := v.Visit(ctx.GetLeft()).(Expression)
	right := v.Visit(ctx.GetRight()).(Expression)
	return &ComparisonExpression{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Operator: ComparisonOperator(ctx.GetOperator().GetText()), // FIXME:
		Left:     left,
		Right:    right,
	}
}

func (v *AstVisitor) VisitRegexpPredicate(ctx *grammar.RegexpPredicateContext) any {
	var result Expression
	value := v.Visit(ctx.GetLeft()).(Expression)
	pattern := v.Visit(ctx.GetPattern()).(Expression)
	result = &RegexPredicate{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value:   value,
		Pattern: pattern,
	}
	if ctx.NEQREGEXP() != nil {
		result = &NotExpression{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Value: result,
		}
	}
	return result
}

func (v *AstVisitor) VisitLikePredicate(ctx *grammar.LikePredicateContext) any {
	var result Expression
	value := v.Visit(ctx.GetLeft()).(Expression)
	pattern := v.Visit(ctx.GetPattern()).(Expression)
	result = &LikePredicate{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value:   value,
		Pattern: pattern,
	}
	if ctx.NOT() != nil {
		result = &NotExpression{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Value: result,
		}
	}
	return result
}

func (v *AstVisitor) VisitInPredicate(ctx *grammar.InPredicateContext) any {
	var result Expression
	value := v.Visit(ctx.GetLeft()).(Expression)
	result = &InPredicate{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value: value,
		ValueList: &InListExpression{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Values: visit[Expression](ctx.AllExpression(), v),
		},
	}

	if ctx.NOT() != nil {
		result = &NotExpression{
			BaseNode: BaseNode{
				ID:       v.idAllocator.Next(),
				Location: getLocation(ctx.GetStart()),
			},
			Value: result,
		}
	}
	return result
}

func (v *AstVisitor) VisitLogicalNot(ctx *grammar.LogicalNotContext) any {
	value := v.Visit(ctx.BooleanExpression()).(Expression)
	return &NotExpression{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value: value,
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
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
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
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Operator: LogicalAND,
		Terms:    visit[Expression](terms, v),
	}
}

func (v *AstVisitor) flatten(root antlr.ParserRuleContext, extractChildren func(ctx antlr.ParserRuleContext) []antlr.ParserRuleContext) (result []antlr.ParserRuleContext) {
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

func (v *AstVisitor) VisitPredicatedExpression(ctx *grammar.PredicatedExpressionContext) any {
	return v.Visit(ctx.Predicate())
}

func (v *AstVisitor) VisitParenExpression(ctx *grammar.ParenExpressionContext) any {
	return v.Visit(ctx.Expression())
}

func (v *AstVisitor) VisitGroupBy(ctx *grammar.GroupByContext) any {
	return &GroupBy{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		GroupingElements: visit[GroupingElement](ctx.AllGroupingElement(), v),
	}
}

func (v *AstVisitor) VisitSingleGroupingSet(ctx *grammar.SingleGroupingSetContext) any {
	return &SimpleGroupBy{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Columns: visit[Expression](ctx.GroupingSet().AllExpression(), v),
	}
}

func (v *AstVisitor) VisitSortItem(ctx *grammar.SortItemContext) any {
	expression := v.Visit(ctx.Expression()).(Expression)
	return &SortItem{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		SortKey:  expression,
		Ordering: getOrderingType(ctx),
	}
}

func (v *AstVisitor) VisitUnquotedIdentifier(ctx *grammar.UnquotedIdentifierContext) any {
	return &Identifier{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value:     ctx.GetText(),
		Delimited: false,
	}
}

func (v *AstVisitor) VisitQuotedIdentifier(ctx *grammar.QuotedIdentifierContext) any {
	token := ctx.GetText()
	identifier, err := strutil.GetStringValue(token)
	if err != nil {
		panic(err)
	}
	return &Identifier{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value:     identifier,
		Delimited: true,
	}
}

func (v *AstVisitor) VisitValueExpressionDefault(ctx *grammar.ValueExpressionDefaultContext) any {
	return v.Visit(ctx.PrimaryExpression())
}

func (v *AstVisitor) VisitValueExpressionPredicate(ctx *grammar.ValueExpressionPredicateContext) any {
	return v.Visit(ctx.ValueExpression())
}

func (v *AstVisitor) VisitDereference(ctx *grammar.DereferenceContext) any {
	base := v.Visit(ctx.GetBase()).(Expression)
	fieldName := v.Visit(ctx.GetFieldName()).(*Identifier)
	return &DereferenceExpression{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Base:  base,
		Field: fieldName,
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
	return &FunctionCall{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
	}
}

// ************** literals **************

func (v *AstVisitor) VisitStringLiteral(ctx *grammar.StringLiteralContext) any {
	return v.Visit(ctx.String_())
}

func (v *AstVisitor) VisitBasicStringLiteral(ctx *grammar.BasicStringLiteralContext) any {
	value, err := strutil.GetStringValue(ctx.STRING().GetText())
	if err != nil {
		panic(err)
	}
	return &StringLiteral{
		BaseNode: BaseNode{
			ID:       v.idAllocator.Next(),
			Location: getLocation(ctx.GetStart()),
		},
		Value: value,
	}
}

func (v *AstVisitor) VisitBooleanLiteral(ctx *grammar.BooleanLiteralContext) any {
	return NewBooleanLiteral(getLocation(ctx.GetStart()), ctx.GetText())
}

func (v *AstVisitor) VisitIntegerLiteral(ctx *grammar.IntegerLiteralContext) any {
	return NewLongLiteral(getLocation(ctx.GetStart()), ctx.GetText())
}

func (v *AstVisitor) VisitDecimalLiteral(ctx *grammar.DecimalLiteralContext) any {
	return NewFloatLiteral(getLocation(ctx.GetStart()), ctx.GetText())
}

func (v *AstVisitor) VisitDoubleLiteral(ctx *grammar.DoubleLiteralContext) any {
	return NewFloatLiteral(getLocation(ctx.GetStart()), ctx.GetText())
}

func (v *AstVisitor) getQualifiedName(ctx grammar.IQualifiedNameContext) *QualifiedName {
	parts := visit[*Identifier](ctx.AllIdentifier(), v)
	return NewQualifiedName(parts)
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

func getLocation(token antlr.Token) *NodeLocation {
	return newNodeLocation(token.GetLine(), token.GetTokenSource().GetCharPositionInLine())
}

func getOrderingType(ctx *grammar.SortItemContext) Ordering {
	if ctx.DESC() != nil {
		return DESCENDING
	}
	return ASCENDING
}
