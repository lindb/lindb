// Code generated from /Users/huangjie/Documents/github/lindb/sql/grammar/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// SQLListener is a complete listener for a parse tree produced by SQLParser.
type SQLListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterStatementList is called when entering the statementList production.
	EnterStatementList(c *StatementListContext)

	// EnterQueryStmt is called when entering the queryStmt production.
	EnterQueryStmt(c *QueryStmtContext)

	// EnterSelectExpr is called when entering the selectExpr production.
	EnterSelectExpr(c *SelectExprContext)

	// EnterFields is called when entering the fields production.
	EnterFields(c *FieldsContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterAlias is called when entering the alias production.
	EnterAlias(c *AliasContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterConditionExpr is called when entering the conditionExpr production.
	EnterConditionExpr(c *ConditionExprContext)

	// EnterTagFilterExpr is called when entering the tagFilterExpr production.
	EnterTagFilterExpr(c *TagFilterExprContext)

	// EnterTagValueList is called when entering the tagValueList production.
	EnterTagValueList(c *TagValueListContext)

	// EnterTimeRangeExpr is called when entering the timeRangeExpr production.
	EnterTimeRangeExpr(c *TimeRangeExprContext)

	// EnterTimeExpr is called when entering the timeExpr production.
	EnterTimeExpr(c *TimeExprContext)

	// EnterNowExpr is called when entering the nowExpr production.
	EnterNowExpr(c *NowExprContext)

	// EnterNowFunc is called when entering the nowFunc production.
	EnterNowFunc(c *NowFuncContext)

	// EnterGroupByClause is called when entering the groupByClause production.
	EnterGroupByClause(c *GroupByClauseContext)

	// EnterGroupByKeys is called when entering the groupByKeys production.
	EnterGroupByKeys(c *GroupByKeysContext)

	// EnterGroupByKey is called when entering the groupByKey production.
	EnterGroupByKey(c *GroupByKeyContext)

	// EnterFillOption is called when entering the fillOption production.
	EnterFillOption(c *FillOptionContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterSortField is called when entering the sortField production.
	EnterSortField(c *SortFieldContext)

	// EnterSortFields is called when entering the sortFields production.
	EnterSortFields(c *SortFieldsContext)

	// EnterHavingClause is called when entering the havingClause production.
	EnterHavingClause(c *HavingClauseContext)

	// EnterBoolExpr is called when entering the boolExpr production.
	EnterBoolExpr(c *BoolExprContext)

	// EnterBoolExprLogicalOp is called when entering the boolExprLogicalOp production.
	EnterBoolExprLogicalOp(c *BoolExprLogicalOpContext)

	// EnterBoolExprAtom is called when entering the boolExprAtom production.
	EnterBoolExprAtom(c *BoolExprAtomContext)

	// EnterBinaryExpr is called when entering the binaryExpr production.
	EnterBinaryExpr(c *BinaryExprContext)

	// EnterBinaryOperator is called when entering the binaryOperator production.
	EnterBinaryOperator(c *BinaryOperatorContext)

	// EnterFieldExpr is called when entering the fieldExpr production.
	EnterFieldExpr(c *FieldExprContext)

	// EnterDurationLit is called when entering the durationLit production.
	EnterDurationLit(c *DurationLitContext)

	// EnterIntervalItem is called when entering the intervalItem production.
	EnterIntervalItem(c *IntervalItemContext)

	// EnterExprFunc is called when entering the exprFunc production.
	EnterExprFunc(c *ExprFuncContext)

	// EnterExprFuncParams is called when entering the exprFuncParams production.
	EnterExprFuncParams(c *ExprFuncParamsContext)

	// EnterFuncParam is called when entering the funcParam production.
	EnterFuncParam(c *FuncParamContext)

	// EnterExprAtom is called when entering the exprAtom production.
	EnterExprAtom(c *ExprAtomContext)

	// EnterIdentFilter is called when entering the identFilter production.
	EnterIdentFilter(c *IdentFilterContext)

	// EnterIntNumber is called when entering the intNumber production.
	EnterIntNumber(c *IntNumberContext)

	// EnterDecNumber is called when entering the decNumber production.
	EnterDecNumber(c *DecNumberContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterMetricName is called when entering the metricName production.
	EnterMetricName(c *MetricNameContext)

	// EnterTagKey is called when entering the tagKey production.
	EnterTagKey(c *TagKeyContext)

	// EnterTagValue is called when entering the tagValue production.
	EnterTagValue(c *TagValueContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// EnterNonReservedWords is called when entering the nonReservedWords production.
	EnterNonReservedWords(c *NonReservedWordsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitStatementList is called when exiting the statementList production.
	ExitStatementList(c *StatementListContext)

	// ExitQueryStmt is called when exiting the queryStmt production.
	ExitQueryStmt(c *QueryStmtContext)

	// ExitSelectExpr is called when exiting the selectExpr production.
	ExitSelectExpr(c *SelectExprContext)

	// ExitFields is called when exiting the fields production.
	ExitFields(c *FieldsContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitAlias is called when exiting the alias production.
	ExitAlias(c *AliasContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitConditionExpr is called when exiting the conditionExpr production.
	ExitConditionExpr(c *ConditionExprContext)

	// ExitTagFilterExpr is called when exiting the tagFilterExpr production.
	ExitTagFilterExpr(c *TagFilterExprContext)

	// ExitTagValueList is called when exiting the tagValueList production.
	ExitTagValueList(c *TagValueListContext)

	// ExitTimeRangeExpr is called when exiting the timeRangeExpr production.
	ExitTimeRangeExpr(c *TimeRangeExprContext)

	// ExitTimeExpr is called when exiting the timeExpr production.
	ExitTimeExpr(c *TimeExprContext)

	// ExitNowExpr is called when exiting the nowExpr production.
	ExitNowExpr(c *NowExprContext)

	// ExitNowFunc is called when exiting the nowFunc production.
	ExitNowFunc(c *NowFuncContext)

	// ExitGroupByClause is called when exiting the groupByClause production.
	ExitGroupByClause(c *GroupByClauseContext)

	// ExitGroupByKeys is called when exiting the groupByKeys production.
	ExitGroupByKeys(c *GroupByKeysContext)

	// ExitGroupByKey is called when exiting the groupByKey production.
	ExitGroupByKey(c *GroupByKeyContext)

	// ExitFillOption is called when exiting the fillOption production.
	ExitFillOption(c *FillOptionContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitSortField is called when exiting the sortField production.
	ExitSortField(c *SortFieldContext)

	// ExitSortFields is called when exiting the sortFields production.
	ExitSortFields(c *SortFieldsContext)

	// ExitHavingClause is called when exiting the havingClause production.
	ExitHavingClause(c *HavingClauseContext)

	// ExitBoolExpr is called when exiting the boolExpr production.
	ExitBoolExpr(c *BoolExprContext)

	// ExitBoolExprLogicalOp is called when exiting the boolExprLogicalOp production.
	ExitBoolExprLogicalOp(c *BoolExprLogicalOpContext)

	// ExitBoolExprAtom is called when exiting the boolExprAtom production.
	ExitBoolExprAtom(c *BoolExprAtomContext)

	// ExitBinaryExpr is called when exiting the binaryExpr production.
	ExitBinaryExpr(c *BinaryExprContext)

	// ExitBinaryOperator is called when exiting the binaryOperator production.
	ExitBinaryOperator(c *BinaryOperatorContext)

	// ExitFieldExpr is called when exiting the fieldExpr production.
	ExitFieldExpr(c *FieldExprContext)

	// ExitDurationLit is called when exiting the durationLit production.
	ExitDurationLit(c *DurationLitContext)

	// ExitIntervalItem is called when exiting the intervalItem production.
	ExitIntervalItem(c *IntervalItemContext)

	// ExitExprFunc is called when exiting the exprFunc production.
	ExitExprFunc(c *ExprFuncContext)

	// ExitExprFuncParams is called when exiting the exprFuncParams production.
	ExitExprFuncParams(c *ExprFuncParamsContext)

	// ExitFuncParam is called when exiting the funcParam production.
	ExitFuncParam(c *FuncParamContext)

	// ExitExprAtom is called when exiting the exprAtom production.
	ExitExprAtom(c *ExprAtomContext)

	// ExitIdentFilter is called when exiting the identFilter production.
	ExitIdentFilter(c *IdentFilterContext)

	// ExitIntNumber is called when exiting the intNumber production.
	ExitIntNumber(c *IntNumberContext)

	// ExitDecNumber is called when exiting the decNumber production.
	ExitDecNumber(c *DecNumberContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitMetricName is called when exiting the metricName production.
	ExitMetricName(c *MetricNameContext)

	// ExitTagKey is called when exiting the tagKey production.
	ExitTagKey(c *TagKeyContext)

	// ExitTagValue is called when exiting the tagValue production.
	ExitTagValue(c *TagValueContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)

	// ExitNonReservedWords is called when exiting the nonReservedWords production.
	ExitNonReservedWords(c *NonReservedWordsContext)
}
