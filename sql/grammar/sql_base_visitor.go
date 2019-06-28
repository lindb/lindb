// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseSQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitStatementList(ctx *StatementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithClauseList(ctx *WithClauseListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithClause(ctx *WithClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalDefineList(ctx *IntervalDefineListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalDefine(ctx *IntervalDefineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShardNum(ctx *ShardNumContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTtlVal(ctx *TtlValContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetattlVal(ctx *MetattlValContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPastVal(ctx *PastValContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFutureVal(ctx *FutureValContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalNameVal(ctx *IntervalNameValContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitReplicaFactor(ctx *ReplicaFactorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDatabaseName(ctx *DatabaseNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitUpdateDatabaseStmt(ctx *UpdateDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowNodeStmt(ctx *ShowNodeStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMeasurementsStmt(ctx *ShowMeasurementsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowInfoStmt(ctx *ShowInfoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagValuesInfoStmt(ctx *ShowTagValuesInfoStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowFieldKeysStmt(ctx *ShowFieldKeysStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowQueriesStmt(ctx *ShowQueriesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStatsStmt(ctx *ShowStatsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithMeasurementClause(ctx *WithMeasurementClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithTagClause(ctx *WithTagClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhereTagCascade(ctx *WhereTagCascadeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitKillQueryStmt(ctx *KillQueryStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQueryId(ctx *QueryIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitServerId(ctx *ServerIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitModule(ctx *ModuleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitComponent(ctx *ComponentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQueryStmt(ctx *QueryStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFields(ctx *FieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitField(ctx *FieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitAlias(ctx *AliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFromClause(ctx *FromClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitClauseBooleanExpr(ctx *ClauseBooleanExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagCascadeExpr(ctx *TagCascadeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagEqualExpr(ctx *TagEqualExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagBooleanExpr(ctx *TagBooleanExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagValueList(ctx *TagValueListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTimeExpr(ctx *TimeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTimeBooleanExpr(ctx *TimeBooleanExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNowExpr(ctx *NowExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNowFunc(ctx *NowFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroupByClause(ctx *GroupByClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDimensions(ctx *DimensionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDimension(ctx *DimensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFillOption(ctx *FillOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitOrderByClause(ctx *OrderByClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalByClause(ctx *IntervalByClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSortField(ctx *SortFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSortFields(ctx *SortFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitHavingClause(ctx *HavingClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExpr(ctx *BoolExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprAtom(ctx *BoolExprAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprBinary(ctx *BoolExprBinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBoolExprBinaryOperator(ctx *BoolExprBinaryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDurationLit(ctx *DurationLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntervalItem(ctx *IntervalItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprFunc(ctx *ExprFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprFuncParams(ctx *ExprFuncParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFuncParam(ctx *FuncParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExprAtom(ctx *ExprAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIdentFilter(ctx *IdentFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIntNumber(ctx *IntNumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDecNumber(ctx *DecNumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetricName(ctx *MetricNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagKey(ctx *TagKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagValue(ctx *TagValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagValuePattern(ctx *TagValuePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNonReservedWords(ctx *NonReservedWordsContext) interface{} {
	return v.VisitChildren(ctx)
}
