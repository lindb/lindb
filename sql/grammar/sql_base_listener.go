// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSQLListener is a complete listener for a parse tree produced by SQLParser.
type BaseSQLListener struct{}

var _ SQLListener = &BaseSQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseSQLListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseSQLListener) ExitStatement(ctx *StatementContext) {}

// EnterStatementList is called when production statementList is entered.
func (s *BaseSQLListener) EnterStatementList(ctx *StatementListContext) {}

// ExitStatementList is called when production statementList is exited.
func (s *BaseSQLListener) ExitStatementList(ctx *StatementListContext) {}

// EnterCreateDatabaseStmt is called when production createDatabaseStmt is entered.
func (s *BaseSQLListener) EnterCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) {}

// ExitCreateDatabaseStmt is called when production createDatabaseStmt is exited.
func (s *BaseSQLListener) ExitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) {}

// EnterWithClauseList is called when production withClauseList is entered.
func (s *BaseSQLListener) EnterWithClauseList(ctx *WithClauseListContext) {}

// ExitWithClauseList is called when production withClauseList is exited.
func (s *BaseSQLListener) ExitWithClauseList(ctx *WithClauseListContext) {}

// EnterWithClause is called when production withClause is entered.
func (s *BaseSQLListener) EnterWithClause(ctx *WithClauseContext) {}

// ExitWithClause is called when production withClause is exited.
func (s *BaseSQLListener) ExitWithClause(ctx *WithClauseContext) {}

// EnterIntervalDefineList is called when production intervalDefineList is entered.
func (s *BaseSQLListener) EnterIntervalDefineList(ctx *IntervalDefineListContext) {}

// ExitIntervalDefineList is called when production intervalDefineList is exited.
func (s *BaseSQLListener) ExitIntervalDefineList(ctx *IntervalDefineListContext) {}

// EnterIntervalDefine is called when production intervalDefine is entered.
func (s *BaseSQLListener) EnterIntervalDefine(ctx *IntervalDefineContext) {}

// ExitIntervalDefine is called when production intervalDefine is exited.
func (s *BaseSQLListener) ExitIntervalDefine(ctx *IntervalDefineContext) {}

// EnterShardNum is called when production shardNum is entered.
func (s *BaseSQLListener) EnterShardNum(ctx *ShardNumContext) {}

// ExitShardNum is called when production shardNum is exited.
func (s *BaseSQLListener) ExitShardNum(ctx *ShardNumContext) {}

// EnterTtlVal is called when production ttlVal is entered.
func (s *BaseSQLListener) EnterTtlVal(ctx *TtlValContext) {}

// ExitTtlVal is called when production ttlVal is exited.
func (s *BaseSQLListener) ExitTtlVal(ctx *TtlValContext) {}

// EnterMetattlVal is called when production metattlVal is entered.
func (s *BaseSQLListener) EnterMetattlVal(ctx *MetattlValContext) {}

// ExitMetattlVal is called when production metattlVal is exited.
func (s *BaseSQLListener) ExitMetattlVal(ctx *MetattlValContext) {}

// EnterPastVal is called when production pastVal is entered.
func (s *BaseSQLListener) EnterPastVal(ctx *PastValContext) {}

// ExitPastVal is called when production pastVal is exited.
func (s *BaseSQLListener) ExitPastVal(ctx *PastValContext) {}

// EnterFutureVal is called when production futureVal is entered.
func (s *BaseSQLListener) EnterFutureVal(ctx *FutureValContext) {}

// ExitFutureVal is called when production futureVal is exited.
func (s *BaseSQLListener) ExitFutureVal(ctx *FutureValContext) {}

// EnterIntervalNameVal is called when production intervalNameVal is entered.
func (s *BaseSQLListener) EnterIntervalNameVal(ctx *IntervalNameValContext) {}

// ExitIntervalNameVal is called when production intervalNameVal is exited.
func (s *BaseSQLListener) ExitIntervalNameVal(ctx *IntervalNameValContext) {}

// EnterReplicaFactor is called when production replicaFactor is entered.
func (s *BaseSQLListener) EnterReplicaFactor(ctx *ReplicaFactorContext) {}

// ExitReplicaFactor is called when production replicaFactor is exited.
func (s *BaseSQLListener) ExitReplicaFactor(ctx *ReplicaFactorContext) {}

// EnterDatabaseName is called when production databaseName is entered.
func (s *BaseSQLListener) EnterDatabaseName(ctx *DatabaseNameContext) {}

// ExitDatabaseName is called when production databaseName is exited.
func (s *BaseSQLListener) ExitDatabaseName(ctx *DatabaseNameContext) {}

// EnterUpdateDatabaseStmt is called when production updateDatabaseStmt is entered.
func (s *BaseSQLListener) EnterUpdateDatabaseStmt(ctx *UpdateDatabaseStmtContext) {}

// ExitUpdateDatabaseStmt is called when production updateDatabaseStmt is exited.
func (s *BaseSQLListener) ExitUpdateDatabaseStmt(ctx *UpdateDatabaseStmtContext) {}

// EnterDropDatabaseStmt is called when production dropDatabaseStmt is entered.
func (s *BaseSQLListener) EnterDropDatabaseStmt(ctx *DropDatabaseStmtContext) {}

// ExitDropDatabaseStmt is called when production dropDatabaseStmt is exited.
func (s *BaseSQLListener) ExitDropDatabaseStmt(ctx *DropDatabaseStmtContext) {}

// EnterShowDatabasesStmt is called when production showDatabasesStmt is entered.
func (s *BaseSQLListener) EnterShowDatabasesStmt(ctx *ShowDatabasesStmtContext) {}

// ExitShowDatabasesStmt is called when production showDatabasesStmt is exited.
func (s *BaseSQLListener) ExitShowDatabasesStmt(ctx *ShowDatabasesStmtContext) {}

// EnterShowNodeStmt is called when production showNodeStmt is entered.
func (s *BaseSQLListener) EnterShowNodeStmt(ctx *ShowNodeStmtContext) {}

// ExitShowNodeStmt is called when production showNodeStmt is exited.
func (s *BaseSQLListener) ExitShowNodeStmt(ctx *ShowNodeStmtContext) {}

// EnterShowMeasurementsStmt is called when production showMeasurementsStmt is entered.
func (s *BaseSQLListener) EnterShowMeasurementsStmt(ctx *ShowMeasurementsStmtContext) {}

// ExitShowMeasurementsStmt is called when production showMeasurementsStmt is exited.
func (s *BaseSQLListener) ExitShowMeasurementsStmt(ctx *ShowMeasurementsStmtContext) {}

// EnterShowTagKeysStmt is called when production showTagKeysStmt is entered.
func (s *BaseSQLListener) EnterShowTagKeysStmt(ctx *ShowTagKeysStmtContext) {}

// ExitShowTagKeysStmt is called when production showTagKeysStmt is exited.
func (s *BaseSQLListener) ExitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) {}

// EnterShowInfoStmt is called when production showInfoStmt is entered.
func (s *BaseSQLListener) EnterShowInfoStmt(ctx *ShowInfoStmtContext) {}

// ExitShowInfoStmt is called when production showInfoStmt is exited.
func (s *BaseSQLListener) ExitShowInfoStmt(ctx *ShowInfoStmtContext) {}

// EnterShowTagValuesStmt is called when production showTagValuesStmt is entered.
func (s *BaseSQLListener) EnterShowTagValuesStmt(ctx *ShowTagValuesStmtContext) {}

// ExitShowTagValuesStmt is called when production showTagValuesStmt is exited.
func (s *BaseSQLListener) ExitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) {}

// EnterShowTagValuesInfoStmt is called when production showTagValuesInfoStmt is entered.
func (s *BaseSQLListener) EnterShowTagValuesInfoStmt(ctx *ShowTagValuesInfoStmtContext) {}

// ExitShowTagValuesInfoStmt is called when production showTagValuesInfoStmt is exited.
func (s *BaseSQLListener) ExitShowTagValuesInfoStmt(ctx *ShowTagValuesInfoStmtContext) {}

// EnterShowFieldKeysStmt is called when production showFieldKeysStmt is entered.
func (s *BaseSQLListener) EnterShowFieldKeysStmt(ctx *ShowFieldKeysStmtContext) {}

// ExitShowFieldKeysStmt is called when production showFieldKeysStmt is exited.
func (s *BaseSQLListener) ExitShowFieldKeysStmt(ctx *ShowFieldKeysStmtContext) {}

// EnterShowQueriesStmt is called when production showQueriesStmt is entered.
func (s *BaseSQLListener) EnterShowQueriesStmt(ctx *ShowQueriesStmtContext) {}

// ExitShowQueriesStmt is called when production showQueriesStmt is exited.
func (s *BaseSQLListener) ExitShowQueriesStmt(ctx *ShowQueriesStmtContext) {}

// EnterShowStatsStmt is called when production showStatsStmt is entered.
func (s *BaseSQLListener) EnterShowStatsStmt(ctx *ShowStatsStmtContext) {}

// ExitShowStatsStmt is called when production showStatsStmt is exited.
func (s *BaseSQLListener) ExitShowStatsStmt(ctx *ShowStatsStmtContext) {}

// EnterWithMeasurementClause is called when production withMeasurementClause is entered.
func (s *BaseSQLListener) EnterWithMeasurementClause(ctx *WithMeasurementClauseContext) {}

// ExitWithMeasurementClause is called when production withMeasurementClause is exited.
func (s *BaseSQLListener) ExitWithMeasurementClause(ctx *WithMeasurementClauseContext) {}

// EnterWithTagClause is called when production withTagClause is entered.
func (s *BaseSQLListener) EnterWithTagClause(ctx *WithTagClauseContext) {}

// ExitWithTagClause is called when production withTagClause is exited.
func (s *BaseSQLListener) ExitWithTagClause(ctx *WithTagClauseContext) {}

// EnterWhereTagCascade is called when production whereTagCascade is entered.
func (s *BaseSQLListener) EnterWhereTagCascade(ctx *WhereTagCascadeContext) {}

// ExitWhereTagCascade is called when production whereTagCascade is exited.
func (s *BaseSQLListener) ExitWhereTagCascade(ctx *WhereTagCascadeContext) {}

// EnterKillQueryStmt is called when production killQueryStmt is entered.
func (s *BaseSQLListener) EnterKillQueryStmt(ctx *KillQueryStmtContext) {}

// ExitKillQueryStmt is called when production killQueryStmt is exited.
func (s *BaseSQLListener) ExitKillQueryStmt(ctx *KillQueryStmtContext) {}

// EnterQueryId is called when production queryId is entered.
func (s *BaseSQLListener) EnterQueryId(ctx *QueryIdContext) {}

// ExitQueryId is called when production queryId is exited.
func (s *BaseSQLListener) ExitQueryId(ctx *QueryIdContext) {}

// EnterServerId is called when production serverId is entered.
func (s *BaseSQLListener) EnterServerId(ctx *ServerIdContext) {}

// ExitServerId is called when production serverId is exited.
func (s *BaseSQLListener) ExitServerId(ctx *ServerIdContext) {}

// EnterModule is called when production module is entered.
func (s *BaseSQLListener) EnterModule(ctx *ModuleContext) {}

// ExitModule is called when production module is exited.
func (s *BaseSQLListener) ExitModule(ctx *ModuleContext) {}

// EnterComponent is called when production component is entered.
func (s *BaseSQLListener) EnterComponent(ctx *ComponentContext) {}

// ExitComponent is called when production component is exited.
func (s *BaseSQLListener) ExitComponent(ctx *ComponentContext) {}

// EnterQueryStmt is called when production queryStmt is entered.
func (s *BaseSQLListener) EnterQueryStmt(ctx *QueryStmtContext) {}

// ExitQueryStmt is called when production queryStmt is exited.
func (s *BaseSQLListener) ExitQueryStmt(ctx *QueryStmtContext) {}

// EnterFields is called when production fields is entered.
func (s *BaseSQLListener) EnterFields(ctx *FieldsContext) {}

// ExitFields is called when production fields is exited.
func (s *BaseSQLListener) ExitFields(ctx *FieldsContext) {}

// EnterField is called when production field is entered.
func (s *BaseSQLListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseSQLListener) ExitField(ctx *FieldContext) {}

// EnterAlias is called when production alias is entered.
func (s *BaseSQLListener) EnterAlias(ctx *AliasContext) {}

// ExitAlias is called when production alias is exited.
func (s *BaseSQLListener) ExitAlias(ctx *AliasContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BaseSQLListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BaseSQLListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseSQLListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseSQLListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterClauseBooleanExpr is called when production clauseBooleanExpr is entered.
func (s *BaseSQLListener) EnterClauseBooleanExpr(ctx *ClauseBooleanExprContext) {}

// ExitClauseBooleanExpr is called when production clauseBooleanExpr is exited.
func (s *BaseSQLListener) ExitClauseBooleanExpr(ctx *ClauseBooleanExprContext) {}

// EnterTagCascadeExpr is called when production tagCascadeExpr is entered.
func (s *BaseSQLListener) EnterTagCascadeExpr(ctx *TagCascadeExprContext) {}

// ExitTagCascadeExpr is called when production tagCascadeExpr is exited.
func (s *BaseSQLListener) ExitTagCascadeExpr(ctx *TagCascadeExprContext) {}

// EnterTagEqualExpr is called when production tagEqualExpr is entered.
func (s *BaseSQLListener) EnterTagEqualExpr(ctx *TagEqualExprContext) {}

// ExitTagEqualExpr is called when production tagEqualExpr is exited.
func (s *BaseSQLListener) ExitTagEqualExpr(ctx *TagEqualExprContext) {}

// EnterTagBooleanExpr is called when production tagBooleanExpr is entered.
func (s *BaseSQLListener) EnterTagBooleanExpr(ctx *TagBooleanExprContext) {}

// ExitTagBooleanExpr is called when production tagBooleanExpr is exited.
func (s *BaseSQLListener) ExitTagBooleanExpr(ctx *TagBooleanExprContext) {}

// EnterTagValueList is called when production tagValueList is entered.
func (s *BaseSQLListener) EnterTagValueList(ctx *TagValueListContext) {}

// ExitTagValueList is called when production tagValueList is exited.
func (s *BaseSQLListener) ExitTagValueList(ctx *TagValueListContext) {}

// EnterTimeExpr is called when production timeExpr is entered.
func (s *BaseSQLListener) EnterTimeExpr(ctx *TimeExprContext) {}

// ExitTimeExpr is called when production timeExpr is exited.
func (s *BaseSQLListener) ExitTimeExpr(ctx *TimeExprContext) {}

// EnterTimeBooleanExpr is called when production timeBooleanExpr is entered.
func (s *BaseSQLListener) EnterTimeBooleanExpr(ctx *TimeBooleanExprContext) {}

// ExitTimeBooleanExpr is called when production timeBooleanExpr is exited.
func (s *BaseSQLListener) ExitTimeBooleanExpr(ctx *TimeBooleanExprContext) {}

// EnterNowExpr is called when production nowExpr is entered.
func (s *BaseSQLListener) EnterNowExpr(ctx *NowExprContext) {}

// ExitNowExpr is called when production nowExpr is exited.
func (s *BaseSQLListener) ExitNowExpr(ctx *NowExprContext) {}

// EnterNowFunc is called when production nowFunc is entered.
func (s *BaseSQLListener) EnterNowFunc(ctx *NowFuncContext) {}

// ExitNowFunc is called when production nowFunc is exited.
func (s *BaseSQLListener) ExitNowFunc(ctx *NowFuncContext) {}

// EnterGroupByClause is called when production groupByClause is entered.
func (s *BaseSQLListener) EnterGroupByClause(ctx *GroupByClauseContext) {}

// ExitGroupByClause is called when production groupByClause is exited.
func (s *BaseSQLListener) ExitGroupByClause(ctx *GroupByClauseContext) {}

// EnterDimensions is called when production dimensions is entered.
func (s *BaseSQLListener) EnterDimensions(ctx *DimensionsContext) {}

// ExitDimensions is called when production dimensions is exited.
func (s *BaseSQLListener) ExitDimensions(ctx *DimensionsContext) {}

// EnterDimension is called when production dimension is entered.
func (s *BaseSQLListener) EnterDimension(ctx *DimensionContext) {}

// ExitDimension is called when production dimension is exited.
func (s *BaseSQLListener) ExitDimension(ctx *DimensionContext) {}

// EnterFillOption is called when production fillOption is entered.
func (s *BaseSQLListener) EnterFillOption(ctx *FillOptionContext) {}

// ExitFillOption is called when production fillOption is exited.
func (s *BaseSQLListener) ExitFillOption(ctx *FillOptionContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseSQLListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseSQLListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterIntervalByClause is called when production intervalByClause is entered.
func (s *BaseSQLListener) EnterIntervalByClause(ctx *IntervalByClauseContext) {}

// ExitIntervalByClause is called when production intervalByClause is exited.
func (s *BaseSQLListener) ExitIntervalByClause(ctx *IntervalByClauseContext) {}

// EnterSortField is called when production sortField is entered.
func (s *BaseSQLListener) EnterSortField(ctx *SortFieldContext) {}

// ExitSortField is called when production sortField is exited.
func (s *BaseSQLListener) ExitSortField(ctx *SortFieldContext) {}

// EnterSortFields is called when production sortFields is entered.
func (s *BaseSQLListener) EnterSortFields(ctx *SortFieldsContext) {}

// ExitSortFields is called when production sortFields is exited.
func (s *BaseSQLListener) ExitSortFields(ctx *SortFieldsContext) {}

// EnterHavingClause is called when production havingClause is entered.
func (s *BaseSQLListener) EnterHavingClause(ctx *HavingClauseContext) {}

// ExitHavingClause is called when production havingClause is exited.
func (s *BaseSQLListener) ExitHavingClause(ctx *HavingClauseContext) {}

// EnterBoolExpr is called when production boolExpr is entered.
func (s *BaseSQLListener) EnterBoolExpr(ctx *BoolExprContext) {}

// ExitBoolExpr is called when production boolExpr is exited.
func (s *BaseSQLListener) ExitBoolExpr(ctx *BoolExprContext) {}

// EnterBoolExprLogicalOp is called when production boolExprLogicalOp is entered.
func (s *BaseSQLListener) EnterBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) {}

// ExitBoolExprLogicalOp is called when production boolExprLogicalOp is exited.
func (s *BaseSQLListener) ExitBoolExprLogicalOp(ctx *BoolExprLogicalOpContext) {}

// EnterBoolExprAtom is called when production boolExprAtom is entered.
func (s *BaseSQLListener) EnterBoolExprAtom(ctx *BoolExprAtomContext) {}

// ExitBoolExprAtom is called when production boolExprAtom is exited.
func (s *BaseSQLListener) ExitBoolExprAtom(ctx *BoolExprAtomContext) {}

// EnterBoolExprBinary is called when production boolExprBinary is entered.
func (s *BaseSQLListener) EnterBoolExprBinary(ctx *BoolExprBinaryContext) {}

// ExitBoolExprBinary is called when production boolExprBinary is exited.
func (s *BaseSQLListener) ExitBoolExprBinary(ctx *BoolExprBinaryContext) {}

// EnterBoolExprBinaryOperator is called when production boolExprBinaryOperator is entered.
func (s *BaseSQLListener) EnterBoolExprBinaryOperator(ctx *BoolExprBinaryOperatorContext) {}

// ExitBoolExprBinaryOperator is called when production boolExprBinaryOperator is exited.
func (s *BaseSQLListener) ExitBoolExprBinaryOperator(ctx *BoolExprBinaryOperatorContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSQLListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSQLListener) ExitExpr(ctx *ExprContext) {}

// EnterDurationLit is called when production durationLit is entered.
func (s *BaseSQLListener) EnterDurationLit(ctx *DurationLitContext) {}

// ExitDurationLit is called when production durationLit is exited.
func (s *BaseSQLListener) ExitDurationLit(ctx *DurationLitContext) {}

// EnterIntervalItem is called when production intervalItem is entered.
func (s *BaseSQLListener) EnterIntervalItem(ctx *IntervalItemContext) {}

// ExitIntervalItem is called when production intervalItem is exited.
func (s *BaseSQLListener) ExitIntervalItem(ctx *IntervalItemContext) {}

// EnterExprFunc is called when production exprFunc is entered.
func (s *BaseSQLListener) EnterExprFunc(ctx *ExprFuncContext) {}

// ExitExprFunc is called when production exprFunc is exited.
func (s *BaseSQLListener) ExitExprFunc(ctx *ExprFuncContext) {}

// EnterExprFuncParams is called when production exprFuncParams is entered.
func (s *BaseSQLListener) EnterExprFuncParams(ctx *ExprFuncParamsContext) {}

// ExitExprFuncParams is called when production exprFuncParams is exited.
func (s *BaseSQLListener) ExitExprFuncParams(ctx *ExprFuncParamsContext) {}

// EnterFuncParam is called when production funcParam is entered.
func (s *BaseSQLListener) EnterFuncParam(ctx *FuncParamContext) {}

// ExitFuncParam is called when production funcParam is exited.
func (s *BaseSQLListener) ExitFuncParam(ctx *FuncParamContext) {}

// EnterExprAtom is called when production exprAtom is entered.
func (s *BaseSQLListener) EnterExprAtom(ctx *ExprAtomContext) {}

// ExitExprAtom is called when production exprAtom is exited.
func (s *BaseSQLListener) ExitExprAtom(ctx *ExprAtomContext) {}

// EnterIdentFilter is called when production identFilter is entered.
func (s *BaseSQLListener) EnterIdentFilter(ctx *IdentFilterContext) {}

// ExitIdentFilter is called when production identFilter is exited.
func (s *BaseSQLListener) ExitIdentFilter(ctx *IdentFilterContext) {}

// EnterIntNumber is called when production intNumber is entered.
func (s *BaseSQLListener) EnterIntNumber(ctx *IntNumberContext) {}

// ExitIntNumber is called when production intNumber is exited.
func (s *BaseSQLListener) ExitIntNumber(ctx *IntNumberContext) {}

// EnterDecNumber is called when production decNumber is entered.
func (s *BaseSQLListener) EnterDecNumber(ctx *DecNumberContext) {}

// ExitDecNumber is called when production decNumber is exited.
func (s *BaseSQLListener) ExitDecNumber(ctx *DecNumberContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseSQLListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseSQLListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterMetricName is called when production metricName is entered.
func (s *BaseSQLListener) EnterMetricName(ctx *MetricNameContext) {}

// ExitMetricName is called when production metricName is exited.
func (s *BaseSQLListener) ExitMetricName(ctx *MetricNameContext) {}

// EnterTagKey is called when production tagKey is entered.
func (s *BaseSQLListener) EnterTagKey(ctx *TagKeyContext) {}

// ExitTagKey is called when production tagKey is exited.
func (s *BaseSQLListener) ExitTagKey(ctx *TagKeyContext) {}

// EnterTagValue is called when production tagValue is entered.
func (s *BaseSQLListener) EnterTagValue(ctx *TagValueContext) {}

// ExitTagValue is called when production tagValue is exited.
func (s *BaseSQLListener) ExitTagValue(ctx *TagValueContext) {}

// EnterTagValuePattern is called when production tagValuePattern is entered.
func (s *BaseSQLListener) EnterTagValuePattern(ctx *TagValuePatternContext) {}

// ExitTagValuePattern is called when production tagValuePattern is exited.
func (s *BaseSQLListener) ExitTagValuePattern(ctx *TagValuePatternContext) {}

// EnterIdent is called when production ident is entered.
func (s *BaseSQLListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BaseSQLListener) ExitIdent(ctx *IdentContext) {}

// EnterNonReservedWords is called when production nonReservedWords is entered.
func (s *BaseSQLListener) EnterNonReservedWords(ctx *NonReservedWordsContext) {}

// ExitNonReservedWords is called when production nonReservedWords is exited.
func (s *BaseSQLListener) ExitNonReservedWords(ctx *NonReservedWordsContext) {}
