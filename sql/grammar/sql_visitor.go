// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by SQLParser.
type SQLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SQLParser#statementList.
	VisitStatementList(ctx *StatementListContext) interface{}

	// Visit a parse tree produced by SQLParser#createDatabaseStmt.
	VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#withClauseList.
	VisitWithClauseList(ctx *WithClauseListContext) interface{}

	// Visit a parse tree produced by SQLParser#withClause.
	VisitWithClause(ctx *WithClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalDefineList.
	VisitIntervalDefineList(ctx *IntervalDefineListContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalDefine.
	VisitIntervalDefine(ctx *IntervalDefineContext) interface{}

	// Visit a parse tree produced by SQLParser#shardNum.
	VisitShardNum(ctx *ShardNumContext) interface{}

	// Visit a parse tree produced by SQLParser#ttlVal.
	VisitTtlVal(ctx *TtlValContext) interface{}

	// Visit a parse tree produced by SQLParser#metattlVal.
	VisitMetattlVal(ctx *MetattlValContext) interface{}

	// Visit a parse tree produced by SQLParser#pastVal.
	VisitPastVal(ctx *PastValContext) interface{}

	// Visit a parse tree produced by SQLParser#futureVal.
	VisitFutureVal(ctx *FutureValContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalNameVal.
	VisitIntervalNameVal(ctx *IntervalNameValContext) interface{}

	// Visit a parse tree produced by SQLParser#replicaFactor.
	VisitReplicaFactor(ctx *ReplicaFactorContext) interface{}

	// Visit a parse tree produced by SQLParser#databaseName.
	VisitDatabaseName(ctx *DatabaseNameContext) interface{}

	// Visit a parse tree produced by SQLParser#updateDatabaseStmt.
	VisitUpdateDatabaseStmt(ctx *UpdateDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#dropDatabaseStmt.
	VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showDatabasesStmt.
	VisitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showNodeStmt.
	VisitShowNodeStmt(ctx *ShowNodeStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMeasurementsStmt.
	VisitShowMeasurementsStmt(ctx *ShowMeasurementsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagKeysStmt.
	VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showInfoStmt.
	VisitShowInfoStmt(ctx *ShowInfoStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagValuesStmt.
	VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagValuesInfoStmt.
	VisitShowTagValuesInfoStmt(ctx *ShowTagValuesInfoStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showFieldKeysStmt.
	VisitShowFieldKeysStmt(ctx *ShowFieldKeysStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showQueriesStmt.
	VisitShowQueriesStmt(ctx *ShowQueriesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStatsStmt.
	VisitShowStatsStmt(ctx *ShowStatsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#withMeasurementClause.
	VisitWithMeasurementClause(ctx *WithMeasurementClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#withTagClause.
	VisitWithTagClause(ctx *WithTagClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#whereTagCascade.
	VisitWhereTagCascade(ctx *WhereTagCascadeContext) interface{}

	// Visit a parse tree produced by SQLParser#killQueryStmt.
	VisitKillQueryStmt(ctx *KillQueryStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#queryId.
	VisitQueryId(ctx *QueryIdContext) interface{}

	// Visit a parse tree produced by SQLParser#serverId.
	VisitServerId(ctx *ServerIdContext) interface{}

	// Visit a parse tree produced by SQLParser#module.
	VisitModule(ctx *ModuleContext) interface{}

	// Visit a parse tree produced by SQLParser#component.
	VisitComponent(ctx *ComponentContext) interface{}

	// Visit a parse tree produced by SQLParser#queryStmt.
	VisitQueryStmt(ctx *QueryStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#fields.
	VisitFields(ctx *FieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by SQLParser#alias.
	VisitAlias(ctx *AliasContext) interface{}

	// Visit a parse tree produced by SQLParser#fromClause.
	VisitFromClause(ctx *FromClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#whereClause.
	VisitWhereClause(ctx *WhereClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#clauseBooleanExpr.
	VisitClauseBooleanExpr(ctx *ClauseBooleanExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagCascadeExpr.
	VisitTagCascadeExpr(ctx *TagCascadeExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagEqualExpr.
	VisitTagEqualExpr(ctx *TagEqualExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagBooleanExpr.
	VisitTagBooleanExpr(ctx *TagBooleanExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagValueList.
	VisitTagValueList(ctx *TagValueListContext) interface{}

	// Visit a parse tree produced by SQLParser#timeExpr.
	VisitTimeExpr(ctx *TimeExprContext) interface{}

	// Visit a parse tree produced by SQLParser#timeBooleanExpr.
	VisitTimeBooleanExpr(ctx *TimeBooleanExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowExpr.
	VisitNowExpr(ctx *NowExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowFunc.
	VisitNowFunc(ctx *NowFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByClause.
	VisitGroupByClause(ctx *GroupByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#dimensions.
	VisitDimensions(ctx *DimensionsContext) interface{}

	// Visit a parse tree produced by SQLParser#dimension.
	VisitDimension(ctx *DimensionContext) interface{}

	// Visit a parse tree produced by SQLParser#fillOption.
	VisitFillOption(ctx *FillOptionContext) interface{}

	// Visit a parse tree produced by SQLParser#orderByClause.
	VisitOrderByClause(ctx *OrderByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalByClause.
	VisitIntervalByClause(ctx *IntervalByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#sortField.
	VisitSortField(ctx *SortFieldContext) interface{}

	// Visit a parse tree produced by SQLParser#sortFields.
	VisitSortFields(ctx *SortFieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#havingClause.
	VisitHavingClause(ctx *HavingClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExpr.
	VisitBoolExpr(ctx *BoolExprContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprLogicalOp.
	VisitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprAtom.
	VisitBoolExprAtom(ctx *BoolExprAtomContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprBinary.
	VisitBoolExprBinary(ctx *BoolExprBinaryContext) interface{}

	// Visit a parse tree produced by SQLParser#boolExprBinaryOperator.
	VisitBoolExprBinaryOperator(ctx *BoolExprBinaryOperatorContext) interface{}

	// Visit a parse tree produced by SQLParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by SQLParser#durationLit.
	VisitDurationLit(ctx *DurationLitContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalItem.
	VisitIntervalItem(ctx *IntervalItemContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFunc.
	VisitExprFunc(ctx *ExprFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFuncParams.
	VisitExprFuncParams(ctx *ExprFuncParamsContext) interface{}

	// Visit a parse tree produced by SQLParser#funcParam.
	VisitFuncParam(ctx *FuncParamContext) interface{}

	// Visit a parse tree produced by SQLParser#exprAtom.
	VisitExprAtom(ctx *ExprAtomContext) interface{}

	// Visit a parse tree produced by SQLParser#identFilter.
	VisitIdentFilter(ctx *IdentFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#intNumber.
	VisitIntNumber(ctx *IntNumberContext) interface{}

	// Visit a parse tree produced by SQLParser#decNumber.
	VisitDecNumber(ctx *DecNumberContext) interface{}

	// Visit a parse tree produced by SQLParser#limitClause.
	VisitLimitClause(ctx *LimitClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#metricName.
	VisitMetricName(ctx *MetricNameContext) interface{}

	// Visit a parse tree produced by SQLParser#tagKey.
	VisitTagKey(ctx *TagKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#tagValue.
	VisitTagValue(ctx *TagValueContext) interface{}

	// Visit a parse tree produced by SQLParser#tagValuePattern.
	VisitTagValuePattern(ctx *TagValuePatternContext) interface{}

	// Visit a parse tree produced by SQLParser#ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by SQLParser#nonReservedWords.
	VisitNonReservedWords(ctx *NonReservedWordsContext) interface{}
}
