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

// Code generated from /Users/jacklhuang/Documents/code/gopath/src/github.com/lindb/lindb/sql/grammar/SQL.g4 by ANTLR 4.9.2. DO NOT EDIT.

package grammar // SQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// SQLListener is a complete listener for a parse tree produced by SQLParser.
type SQLListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterStatementList is called when entering the statementList production.
	EnterStatementList(c *StatementListContext)

	// EnterUseStmt is called when entering the useStmt production.
	EnterUseStmt(c *UseStmtContext)

	// EnterShowMasterStmt is called when entering the showMasterStmt production.
	EnterShowMasterStmt(c *ShowMasterStmtContext)

	// EnterCreateDatabaseStmt is called when entering the createDatabaseStmt production.
	EnterCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// EnterShowDatabaseStmt is called when entering the showDatabaseStmt production.
	EnterShowDatabaseStmt(c *ShowDatabaseStmtContext)

	// EnterShowNameSpacesStmt is called when entering the showNameSpacesStmt production.
	EnterShowNameSpacesStmt(c *ShowNameSpacesStmtContext)

	// EnterShowMetricsStmt is called when entering the showMetricsStmt production.
	EnterShowMetricsStmt(c *ShowMetricsStmtContext)

	// EnterShowFieldsStmt is called when entering the showFieldsStmt production.
	EnterShowFieldsStmt(c *ShowFieldsStmtContext)

	// EnterShowTagKeysStmt is called when entering the showTagKeysStmt production.
	EnterShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// EnterShowTagValuesStmt is called when entering the showTagValuesStmt production.
	EnterShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// EnterPrefix is called when entering the prefix production.
	EnterPrefix(c *PrefixContext)

	// EnterWithTagKey is called when entering the withTagKey production.
	EnterWithTagKey(c *WithTagKeyContext)

	// EnterNamespace is called when entering the namespace production.
	EnterNamespace(c *NamespaceContext)

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

	// EnterFuncName is called when entering the funcName production.
	EnterFuncName(c *FuncNameContext)

	// EnterExprFuncParams is called when entering the exprFuncParams production.
	EnterExprFuncParams(c *ExprFuncParamsContext)

	// EnterFuncParam is called when entering the funcParam production.
	EnterFuncParam(c *FuncParamContext)

	// EnterExprAtom is called when entering the exprAtom production.
	EnterExprAtom(c *ExprAtomContext)

	// EnterIdentFilter is called when entering the identFilter production.
	EnterIdentFilter(c *IdentFilterContext)

	// EnterJson is called when entering the json production.
	EnterJson(c *JsonContext)

	// EnterObj is called when entering the obj production.
	EnterObj(c *ObjContext)

	// EnterPair is called when entering the pair production.
	EnterPair(c *PairContext)

	// EnterArr is called when entering the arr production.
	EnterArr(c *ArrContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

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

	// ExitUseStmt is called when exiting the useStmt production.
	ExitUseStmt(c *UseStmtContext)

	// ExitShowMasterStmt is called when exiting the showMasterStmt production.
	ExitShowMasterStmt(c *ShowMasterStmtContext)

	// ExitCreateDatabaseStmt is called when exiting the createDatabaseStmt production.
	ExitCreateDatabaseStmt(c *CreateDatabaseStmtContext)

	// ExitShowDatabaseStmt is called when exiting the showDatabaseStmt production.
	ExitShowDatabaseStmt(c *ShowDatabaseStmtContext)

	// ExitShowNameSpacesStmt is called when exiting the showNameSpacesStmt production.
	ExitShowNameSpacesStmt(c *ShowNameSpacesStmtContext)

	// ExitShowMetricsStmt is called when exiting the showMetricsStmt production.
	ExitShowMetricsStmt(c *ShowMetricsStmtContext)

	// ExitShowFieldsStmt is called when exiting the showFieldsStmt production.
	ExitShowFieldsStmt(c *ShowFieldsStmtContext)

	// ExitShowTagKeysStmt is called when exiting the showTagKeysStmt production.
	ExitShowTagKeysStmt(c *ShowTagKeysStmtContext)

	// ExitShowTagValuesStmt is called when exiting the showTagValuesStmt production.
	ExitShowTagValuesStmt(c *ShowTagValuesStmtContext)

	// ExitPrefix is called when exiting the prefix production.
	ExitPrefix(c *PrefixContext)

	// ExitWithTagKey is called when exiting the withTagKey production.
	ExitWithTagKey(c *WithTagKeyContext)

	// ExitNamespace is called when exiting the namespace production.
	ExitNamespace(c *NamespaceContext)

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

	// ExitFuncName is called when exiting the funcName production.
	ExitFuncName(c *FuncNameContext)

	// ExitExprFuncParams is called when exiting the exprFuncParams production.
	ExitExprFuncParams(c *ExprFuncParamsContext)

	// ExitFuncParam is called when exiting the funcParam production.
	ExitFuncParam(c *FuncParamContext)

	// ExitExprAtom is called when exiting the exprAtom production.
	ExitExprAtom(c *ExprAtomContext)

	// ExitIdentFilter is called when exiting the identFilter production.
	ExitIdentFilter(c *IdentFilterContext)

	// ExitJson is called when exiting the json production.
	ExitJson(c *JsonContext)

	// ExitObj is called when exiting the obj production.
	ExitObj(c *ObjContext)

	// ExitPair is called when exiting the pair production.
	ExitPair(c *PairContext)

	// ExitArr is called when exiting the arr production.
	ExitArr(c *ArrContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

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
