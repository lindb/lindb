// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SQLListener is a complete listener for a parse tree produced by SQLParser.
type SQLListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterStatementList is called when entering the statementList production.
	EnterStatementList(c *StatementListContext)

	// EnterCreateDatabaseStmt is called when entering the createDatabaseStmt production.
	EnterCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// EnterWithClauseList is called when entering the withClauseList production.
	EnterWithClauseList(c *WithClauseListContext)

	// EnterWithClause is called when entering the withClause production.
	EnterWithClause(c *WithClauseContext)

	// EnterIntervalDefineList is called when entering the intervalDefineList production.
	EnterIntervalDefineList(c *IntervalDefineListContext)

	// EnterIntervalDefine is called when entering the intervalDefine production.
	EnterIntervalDefine(c *IntervalDefineContext)

	// EnterShardNum is called when entering the shardNum production.
	EnterShardNum(c *ShardNumContext)

	// EnterTtlVal is called when entering the ttlVal production.
	EnterTtlVal(c *TtlValContext)

	// EnterMetattlVal is called when entering the metattlVal production.
	EnterMetattlVal(c *MetattlValContext)

	// EnterPastVal is called when entering the pastVal production.
	EnterPastVal(c *PastValContext)

	// EnterFutureVal is called when entering the futureVal production.
	EnterFutureVal(c *FutureValContext)

	// EnterIntervalNameVal is called when entering the intervalNameVal production.
	EnterIntervalNameVal(c *IntervalNameValContext)

	// EnterReplicaFactor is called when entering the replicaFactor production.
	EnterReplicaFactor(c *ReplicaFactorContext)

	// EnterDatabaseName is called when entering the databaseName production.
	EnterDatabaseName(c *DatabaseNameContext)

	// EnterUpdateDatabaseStmt is called when entering the updateDatabaseStmt production.
	EnterUpdateDatabaseStmt(c *UpdateDatabaseStmtContext)

	// EnterDropDatabaseStmt is called when entering the dropDatabaseStmt production.
	EnterDropDatabaseStmt(c *DropDatabaseStmtContext)

	// EnterShowDatabasesStmt is called when entering the showDatabasesStmt production.
	EnterShowDatabasesStmt(c *ShowDatabasesStmtContext)

	// EnterShowNodeStmt is called when entering the showNodeStmt production.
	EnterShowNodeStmt(c *ShowNodeStmtContext)

	// EnterShowMeasurementsStmt is called when entering the showMeasurementsStmt production.
	EnterShowMeasurementsStmt(c *ShowMeasurementsStmtContext)

	// EnterShowTagKeysStmt is called when entering the showTagKeysStmt production.
	EnterShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// EnterShowInfoStmt is called when entering the showInfoStmt production.
	EnterShowInfoStmt(c *ShowInfoStmtContext)

	// EnterShowTagValuesStmt is called when entering the showTagValuesStmt production.
	EnterShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// EnterShowTagValuesInfoStmt is called when entering the showTagValuesInfoStmt production.
	EnterShowTagValuesInfoStmt(c *ShowTagValuesInfoStmtContext)

	// EnterShowFieldKeysStmt is called when entering the showFieldKeysStmt production.
	EnterShowFieldKeysStmt(c *ShowFieldKeysStmtContext)

	// EnterShowQueriesStmt is called when entering the showQueriesStmt production.
	EnterShowQueriesStmt(c *ShowQueriesStmtContext)

	// EnterShowStatsStmt is called when entering the showStatsStmt production.
	EnterShowStatsStmt(c *ShowStatsStmtContext)

	// EnterWithMeasurementClause is called when entering the withMeasurementClause production.
	EnterWithMeasurementClause(c *WithMeasurementClauseContext)

	// EnterWithTagClause is called when entering the withTagClause production.
	EnterWithTagClause(c *WithTagClauseContext)

	// EnterWhereTagCascade is called when entering the whereTagCascade production.
	EnterWhereTagCascade(c *WhereTagCascadeContext)

	// EnterKillQueryStmt is called when entering the killQueryStmt production.
	EnterKillQueryStmt(c *KillQueryStmtContext)

	// EnterQueryId is called when entering the queryId production.
	EnterQueryId(c *QueryIdContext)

	// EnterServerId is called when entering the serverId production.
	EnterServerId(c *ServerIdContext)

	// EnterModule is called when entering the module production.
	EnterModule(c *ModuleContext)

	// EnterComponent is called when entering the component production.
	EnterComponent(c *ComponentContext)

	// EnterQueryStmt is called when entering the queryStmt production.
	EnterQueryStmt(c *QueryStmtContext)

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

	// EnterClauseBooleanExpr is called when entering the clauseBooleanExpr production.
	EnterClauseBooleanExpr(c *ClauseBooleanExprContext)

	// EnterTagCascadeExpr is called when entering the tagCascadeExpr production.
	EnterTagCascadeExpr(c *TagCascadeExprContext)

	// EnterTagEqualExpr is called when entering the tagEqualExpr production.
	EnterTagEqualExpr(c *TagEqualExprContext)

	// EnterTagBooleanExpr is called when entering the tagBooleanExpr production.
	EnterTagBooleanExpr(c *TagBooleanExprContext)

	// EnterTagValueList is called when entering the tagValueList production.
	EnterTagValueList(c *TagValueListContext)

	// EnterTimeExpr is called when entering the timeExpr production.
	EnterTimeExpr(c *TimeExprContext)

	// EnterTimeBooleanExpr is called when entering the timeBooleanExpr production.
	EnterTimeBooleanExpr(c *TimeBooleanExprContext)

	// EnterNowExpr is called when entering the nowExpr production.
	EnterNowExpr(c *NowExprContext)

	// EnterNowFunc is called when entering the nowFunc production.
	EnterNowFunc(c *NowFuncContext)

	// EnterGroupByClause is called when entering the groupByClause production.
	EnterGroupByClause(c *GroupByClauseContext)

	// EnterDimensions is called when entering the dimensions production.
	EnterDimensions(c *DimensionsContext)

	// EnterDimension is called when entering the dimension production.
	EnterDimension(c *DimensionContext)

	// EnterFillOption is called when entering the fillOption production.
	EnterFillOption(c *FillOptionContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterIntervalByClause is called when entering the intervalByClause production.
	EnterIntervalByClause(c *IntervalByClauseContext)

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

	// EnterBoolExprBinary is called when entering the boolExprBinary production.
	EnterBoolExprBinary(c *BoolExprBinaryContext)

	// EnterBoolExprBinaryOperator is called when entering the boolExprBinaryOperator production.
	EnterBoolExprBinaryOperator(c *BoolExprBinaryOperatorContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

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

	// EnterTagValuePattern is called when entering the tagValuePattern production.
	EnterTagValuePattern(c *TagValuePatternContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// EnterNonReservedWords is called when entering the nonReservedWords production.
	EnterNonReservedWords(c *NonReservedWordsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitStatementList is called when exiting the statementList production.
	ExitStatementList(c *StatementListContext)

	// ExitCreateDatabaseStmt is called when exiting the createDatabaseStmt production.
	ExitCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// ExitWithClauseList is called when exiting the withClauseList production.
	ExitWithClauseList(c *WithClauseListContext)

	// ExitWithClause is called when exiting the withClause production.
	ExitWithClause(c *WithClauseContext)

	// ExitIntervalDefineList is called when exiting the intervalDefineList production.
	ExitIntervalDefineList(c *IntervalDefineListContext)

	// ExitIntervalDefine is called when exiting the intervalDefine production.
	ExitIntervalDefine(c *IntervalDefineContext)

	// ExitShardNum is called when exiting the shardNum production.
	ExitShardNum(c *ShardNumContext)

	// ExitTtlVal is called when exiting the ttlVal production.
	ExitTtlVal(c *TtlValContext)

	// ExitMetattlVal is called when exiting the metattlVal production.
	ExitMetattlVal(c *MetattlValContext)

	// ExitPastVal is called when exiting the pastVal production.
	ExitPastVal(c *PastValContext)

	// ExitFutureVal is called when exiting the futureVal production.
	ExitFutureVal(c *FutureValContext)

	// ExitIntervalNameVal is called when exiting the intervalNameVal production.
	ExitIntervalNameVal(c *IntervalNameValContext)

	// ExitReplicaFactor is called when exiting the replicaFactor production.
	ExitReplicaFactor(c *ReplicaFactorContext)

	// ExitDatabaseName is called when exiting the databaseName production.
	ExitDatabaseName(c *DatabaseNameContext)

	// ExitUpdateDatabaseStmt is called when exiting the updateDatabaseStmt production.
	ExitUpdateDatabaseStmt(c *UpdateDatabaseStmtContext)

	// ExitDropDatabaseStmt is called when exiting the dropDatabaseStmt production.
	ExitDropDatabaseStmt(c *DropDatabaseStmtContext)

	// ExitShowDatabasesStmt is called when exiting the showDatabasesStmt production.
	ExitShowDatabasesStmt(c *ShowDatabasesStmtContext)

	// ExitShowNodeStmt is called when exiting the showNodeStmt production.
	ExitShowNodeStmt(c *ShowNodeStmtContext)

	// ExitShowMeasurementsStmt is called when exiting the showMeasurementsStmt production.
	ExitShowMeasurementsStmt(c *ShowMeasurementsStmtContext)

	// ExitShowTagKeysStmt is called when exiting the showTagKeysStmt production.
	ExitShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// ExitShowInfoStmt is called when exiting the showInfoStmt production.
	ExitShowInfoStmt(c *ShowInfoStmtContext)

	// ExitShowTagValuesStmt is called when exiting the showTagValuesStmt production.
	ExitShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// ExitShowTagValuesInfoStmt is called when exiting the showTagValuesInfoStmt production.
	ExitShowTagValuesInfoStmt(c *ShowTagValuesInfoStmtContext)

	// ExitShowFieldKeysStmt is called when exiting the showFieldKeysStmt production.
	ExitShowFieldKeysStmt(c *ShowFieldKeysStmtContext)

	// ExitShowQueriesStmt is called when exiting the showQueriesStmt production.
	ExitShowQueriesStmt(c *ShowQueriesStmtContext)

	// ExitShowStatsStmt is called when exiting the showStatsStmt production.
	ExitShowStatsStmt(c *ShowStatsStmtContext)

	// ExitWithMeasurementClause is called when exiting the withMeasurementClause production.
	ExitWithMeasurementClause(c *WithMeasurementClauseContext)

	// ExitWithTagClause is called when exiting the withTagClause production.
	ExitWithTagClause(c *WithTagClauseContext)

	// ExitWhereTagCascade is called when exiting the whereTagCascade production.
	ExitWhereTagCascade(c *WhereTagCascadeContext)

	// ExitKillQueryStmt is called when exiting the killQueryStmt production.
	ExitKillQueryStmt(c *KillQueryStmtContext)

	// ExitQueryId is called when exiting the queryId production.
	ExitQueryId(c *QueryIdContext)

	// ExitServerId is called when exiting the serverId production.
	ExitServerId(c *ServerIdContext)

	// ExitModule is called when exiting the module production.
	ExitModule(c *ModuleContext)

	// ExitComponent is called when exiting the component production.
	ExitComponent(c *ComponentContext)

	// ExitQueryStmt is called when exiting the queryStmt production.
	ExitQueryStmt(c *QueryStmtContext)

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

	// ExitClauseBooleanExpr is called when exiting the clauseBooleanExpr production.
	ExitClauseBooleanExpr(c *ClauseBooleanExprContext)

	// ExitTagCascadeExpr is called when exiting the tagCascadeExpr production.
	ExitTagCascadeExpr(c *TagCascadeExprContext)

	// ExitTagEqualExpr is called when exiting the tagEqualExpr production.
	ExitTagEqualExpr(c *TagEqualExprContext)

	// ExitTagBooleanExpr is called when exiting the tagBooleanExpr production.
	ExitTagBooleanExpr(c *TagBooleanExprContext)

	// ExitTagValueList is called when exiting the tagValueList production.
	ExitTagValueList(c *TagValueListContext)

	// ExitTimeExpr is called when exiting the timeExpr production.
	ExitTimeExpr(c *TimeExprContext)

	// ExitTimeBooleanExpr is called when exiting the timeBooleanExpr production.
	ExitTimeBooleanExpr(c *TimeBooleanExprContext)

	// ExitNowExpr is called when exiting the nowExpr production.
	ExitNowExpr(c *NowExprContext)

	// ExitNowFunc is called when exiting the nowFunc production.
	ExitNowFunc(c *NowFuncContext)

	// ExitGroupByClause is called when exiting the groupByClause production.
	ExitGroupByClause(c *GroupByClauseContext)

	// ExitDimensions is called when exiting the dimensions production.
	ExitDimensions(c *DimensionsContext)

	// ExitDimension is called when exiting the dimension production.
	ExitDimension(c *DimensionContext)

	// ExitFillOption is called when exiting the fillOption production.
	ExitFillOption(c *FillOptionContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitIntervalByClause is called when exiting the intervalByClause production.
	ExitIntervalByClause(c *IntervalByClauseContext)

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

	// ExitBoolExprBinary is called when exiting the boolExprBinary production.
	ExitBoolExprBinary(c *BoolExprBinaryContext)

	// ExitBoolExprBinaryOperator is called when exiting the boolExprBinaryOperator production.
	ExitBoolExprBinaryOperator(c *BoolExprBinaryOperatorContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

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

	// ExitTagValuePattern is called when exiting the tagValuePattern production.
	ExitTagValuePattern(c *TagValuePatternContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)

	// ExitNonReservedWords is called when exiting the nonReservedWords production.
	ExitNonReservedWords(c *NonReservedWordsContext)
}
