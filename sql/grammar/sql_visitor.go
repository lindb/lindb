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

// Code generated from ./sql/grammar/SQL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by SQLParser.
type SQLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SQLParser#useStmt.
	VisitUseStmt(ctx *UseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#setLimitStmt.
	VisitSetLimitStmt(ctx *SetLimitStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStmt.
	VisitShowStmt(ctx *ShowStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMasterStmt.
	VisitShowMasterStmt(ctx *ShowMasterStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showRequestsStmt.
	VisitShowRequestsStmt(ctx *ShowRequestsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showRequestStmt.
	VisitShowRequestStmt(ctx *ShowRequestStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showBrokersStmt.
	VisitShowBrokersStmt(ctx *ShowBrokersStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showLimitStmt.
	VisitShowLimitStmt(ctx *ShowLimitStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetadataTypesStmt.
	VisitShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showRootMetaStmt.
	VisitShowRootMetaStmt(ctx *ShowRootMetaStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showBrokerMetaStmt.
	VisitShowBrokerMetaStmt(ctx *ShowBrokerMetaStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMasterMetaStmt.
	VisitShowMasterMetaStmt(ctx *ShowMasterMetaStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStorageMetaStmt.
	VisitShowStorageMetaStmt(ctx *ShowStorageMetaStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showAliveStmt.
	VisitShowAliveStmt(ctx *ShowAliveStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showReplicationStmt.
	VisitShowReplicationStmt(ctx *ShowReplicationStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMemoryDatabaseStmt.
	VisitShowMemoryDatabaseStmt(ctx *ShowMemoryDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showRootMetricStmt.
	VisitShowRootMetricStmt(ctx *ShowRootMetricStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showBrokerMetricStmt.
	VisitShowBrokerMetricStmt(ctx *ShowBrokerMetricStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showStorageMetricStmt.
	VisitShowStorageMetricStmt(ctx *ShowStorageMetricStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#createStorageStmt.
	VisitCreateStorageStmt(ctx *CreateStorageStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#createBrokerStmt.
	VisitCreateBrokerStmt(ctx *CreateBrokerStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#recoverStorageStmt.
	VisitRecoverStorageStmt(ctx *RecoverStorageStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showSchemasStmt.
	VisitShowSchemasStmt(ctx *ShowSchemasStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#createDatabaseStmt.
	VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#dropDatabaseStmt.
	VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showDatabaseStmt.
	VisitShowDatabaseStmt(ctx *ShowDatabaseStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showNameSpacesStmt.
	VisitShowNameSpacesStmt(ctx *ShowNameSpacesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetricsStmt.
	VisitShowMetricsStmt(ctx *ShowMetricsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showFieldsStmt.
	VisitShowFieldsStmt(ctx *ShowFieldsStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagKeysStmt.
	VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#showTagValuesStmt.
	VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#prefix.
	VisitPrefix(ctx *PrefixContext) interface{}

	// Visit a parse tree produced by SQLParser#withTagKey.
	VisitWithTagKey(ctx *WithTagKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#namespace.
	VisitNamespace(ctx *NamespaceContext) interface{}

	// Visit a parse tree produced by SQLParser#databaseName.
	VisitDatabaseName(ctx *DatabaseNameContext) interface{}

	// Visit a parse tree produced by SQLParser#storageName.
	VisitStorageName(ctx *StorageNameContext) interface{}

	// Visit a parse tree produced by SQLParser#requestID.
	VisitRequestID(ctx *RequestIDContext) interface{}

	// Visit a parse tree produced by SQLParser#source.
	VisitSource(ctx *SourceContext) interface{}

	// Visit a parse tree produced by SQLParser#optionClause.
	VisitOptionClause(ctx *OptionClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#optionPairs.
	VisitOptionPairs(ctx *OptionPairsContext) interface{}

	// Visit a parse tree produced by SQLParser#closedOptionPairs.
	VisitClosedOptionPairs(ctx *ClosedOptionPairsContext) interface{}

	// Visit a parse tree produced by SQLParser#optionPair.
	VisitOptionPair(ctx *OptionPairContext) interface{}

	// Visit a parse tree produced by SQLParser#optionKey.
	VisitOptionKey(ctx *OptionKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#optionValue.
	VisitOptionValue(ctx *OptionValueContext) interface{}

	// Visit a parse tree produced by SQLParser#queryStmt.
	VisitQueryStmt(ctx *QueryStmtContext) interface{}

	// Visit a parse tree produced by SQLParser#sourceAndSelect.
	VisitSourceAndSelect(ctx *SourceAndSelectContext) interface{}

	// Visit a parse tree produced by SQLParser#selectExpr.
	VisitSelectExpr(ctx *SelectExprContext) interface{}

	// Visit a parse tree produced by SQLParser#fields.
	VisitFields(ctx *FieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by SQLParser#alias.
	VisitAlias(ctx *AliasContext) interface{}

	// Visit a parse tree produced by SQLParser#brokerFilter.
	VisitBrokerFilter(ctx *BrokerFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#databaseFilter.
	VisitDatabaseFilter(ctx *DatabaseFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#typeFilter.
	VisitTypeFilter(ctx *TypeFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#fromClause.
	VisitFromClause(ctx *FromClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#whereClause.
	VisitWhereClause(ctx *WhereClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#conditionExpr.
	VisitConditionExpr(ctx *ConditionExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagFilterExpr.
	VisitTagFilterExpr(ctx *TagFilterExprContext) interface{}

	// Visit a parse tree produced by SQLParser#tagValueList.
	VisitTagValueList(ctx *TagValueListContext) interface{}

	// Visit a parse tree produced by SQLParser#metricListFilter.
	VisitMetricListFilter(ctx *MetricListFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#metricList.
	VisitMetricList(ctx *MetricListContext) interface{}

	// Visit a parse tree produced by SQLParser#timeRangeExpr.
	VisitTimeRangeExpr(ctx *TimeRangeExprContext) interface{}

	// Visit a parse tree produced by SQLParser#timeExpr.
	VisitTimeExpr(ctx *TimeExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowExpr.
	VisitNowExpr(ctx *NowExprContext) interface{}

	// Visit a parse tree produced by SQLParser#nowFunc.
	VisitNowFunc(ctx *NowFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByClause.
	VisitGroupByClause(ctx *GroupByClauseContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByKeys.
	VisitGroupByKeys(ctx *GroupByKeysContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByKey.
	VisitGroupByKey(ctx *GroupByKeyContext) interface{}

	// Visit a parse tree produced by SQLParser#fillOption.
	VisitFillOption(ctx *FillOptionContext) interface{}

	// Visit a parse tree produced by SQLParser#orderByClause.
	VisitOrderByClause(ctx *OrderByClauseContext) interface{}

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

	// Visit a parse tree produced by SQLParser#binaryExpr.
	VisitBinaryExpr(ctx *BinaryExprContext) interface{}

	// Visit a parse tree produced by SQLParser#binaryOperator.
	VisitBinaryOperator(ctx *BinaryOperatorContext) interface{}

	// Visit a parse tree produced by SQLParser#fieldExpr.
	VisitFieldExpr(ctx *FieldExprContext) interface{}

	// Visit a parse tree produced by SQLParser#star.
	VisitStar(ctx *StarContext) interface{}

	// Visit a parse tree produced by SQLParser#durationLit.
	VisitDurationLit(ctx *DurationLitContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalItem.
	VisitIntervalItem(ctx *IntervalItemContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFunc.
	VisitExprFunc(ctx *ExprFuncContext) interface{}

	// Visit a parse tree produced by SQLParser#funcName.
	VisitFuncName(ctx *FuncNameContext) interface{}

	// Visit a parse tree produced by SQLParser#exprFuncParams.
	VisitExprFuncParams(ctx *ExprFuncParamsContext) interface{}

	// Visit a parse tree produced by SQLParser#funcParam.
	VisitFuncParam(ctx *FuncParamContext) interface{}

	// Visit a parse tree produced by SQLParser#exprAtom.
	VisitExprAtom(ctx *ExprAtomContext) interface{}

	// Visit a parse tree produced by SQLParser#identFilter.
	VisitIdentFilter(ctx *IdentFilterContext) interface{}

	// Visit a parse tree produced by SQLParser#json.
	VisitJson(ctx *JsonContext) interface{}

	// Visit a parse tree produced by SQLParser#toml.
	VisitToml(ctx *TomlContext) interface{}

	// Visit a parse tree produced by SQLParser#obj.
	VisitObj(ctx *ObjContext) interface{}

	// Visit a parse tree produced by SQLParser#pair.
	VisitPair(ctx *PairContext) interface{}

	// Visit a parse tree produced by SQLParser#arr.
	VisitArr(ctx *ArrContext) interface{}

	// Visit a parse tree produced by SQLParser#value.
	VisitValue(ctx *ValueContext) interface{}

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

	// Visit a parse tree produced by SQLParser#ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by SQLParser#nonReservedWords.
	VisitNonReservedWords(ctx *NonReservedWordsContext) interface{}
}
