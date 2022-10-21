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

// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

type BaseSQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitUseStmt(ctx *UseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStmt(ctx *ShowStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMasterStmt(ctx *ShowMasterStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowRequestsStmt(ctx *ShowRequestsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowRequestStmt(ctx *ShowRequestStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStoragesStmt(ctx *ShowStoragesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMetadataTypesStmt(ctx *ShowMetadataTypesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowBrokerMetaStmt(ctx *ShowBrokerMetaStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMasterMetaStmt(ctx *ShowMasterMetaStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStorageMetaStmt(ctx *ShowStorageMetaStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowAliveStmt(ctx *ShowAliveStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowReplicationStmt(ctx *ShowReplicationStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMemoryDatabaseStmt(ctx *ShowMemoryDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowBrokerMetricStmt(ctx *ShowBrokerMetricStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowStorageMetricStmt(ctx *ShowStorageMetricStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreateStorageStmt(ctx *CreateStorageStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowSchemasStmt(ctx *ShowSchemasStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreateDatabaseStmt(ctx *CreateDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDropDatabaseStmt(ctx *DropDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowDatabaseStmt(ctx *ShowDatabaseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowNameSpacesStmt(ctx *ShowNameSpacesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowMetricsStmt(ctx *ShowMetricsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowFieldsStmt(ctx *ShowFieldsStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagKeysStmt(ctx *ShowTagKeysStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShowTagValuesStmt(ctx *ShowTagValuesStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPrefix(ctx *PrefixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWithTagKey(ctx *WithTagKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDatabaseName(ctx *DatabaseNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitRequestID(ctx *RequestIDContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSource(ctx *SourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQueryStmt(ctx *QueryStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSourceAndSelect(ctx *SourceAndSelectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSelectExpr(ctx *SelectExprContext) interface{} {
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

func (v *BaseSQLVisitor) VisitStorageFilter(ctx *StorageFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDatabaseFilter(ctx *DatabaseFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTypeFilter(ctx *TypeFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFromClause(ctx *FromClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitConditionExpr(ctx *ConditionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagFilterExpr(ctx *TagFilterExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTagValueList(ctx *TagValueListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetricListFilter(ctx *MetricListFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetricList(ctx *MetricListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTimeRangeExpr(ctx *TimeRangeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTimeExpr(ctx *TimeExprContext) interface{} {
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

func (v *BaseSQLVisitor) VisitGroupByKeys(ctx *GroupByKeysContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroupByKey(ctx *GroupByKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFillOption(ctx *FillOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitOrderByClause(ctx *OrderByClauseContext) interface{} {
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

func (v *BaseSQLVisitor) VisitBinaryExpr(ctx *BinaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBinaryOperator(ctx *BinaryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFieldExpr(ctx *FieldExprContext) interface{} {
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

func (v *BaseSQLVisitor) VisitFuncName(ctx *FuncNameContext) interface{} {
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

func (v *BaseSQLVisitor) VisitJson(ctx *JsonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitObj(ctx *ObjContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPair(ctx *PairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitArr(ctx *ArrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitValue(ctx *ValueContext) interface{} {
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

func (v *BaseSQLVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNonReservedWords(ctx *NonReservedWordsContext) interface{} {
	return v.VisitChildren(ctx)
}
