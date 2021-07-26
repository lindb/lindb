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

package sql

import (
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

type listener struct {
	*grammar.BaseSQLListener
	stmt *queryStmtParse

	metaStmt *metaStmtParser
}

// EnterQueryStmt is called when production queryStmt is entered.
func (l *listener) EnterQueryStmt(ctx *grammar.QueryStmtContext) {
	l.stmt = newQueryStmtParse(ctx.T_EXPLAIN() != nil)
}

// EnterShowDatabaseStmt is called when production showDatabaseStmt is entered.
func (l *listener) EnterShowDatabaseStmt(ctx *grammar.ShowDatabaseStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.Database)
}

// EnterShowNameSpacesStmt is called when production showNameSpacesStmt is entered.
func (l *listener) EnterShowNameSpacesStmt(ctx *grammar.ShowNameSpacesStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.Namespace)
}

// EnterShowMetricsStmt is called when production showMetricsStmt is entered.
func (l *listener) EnterShowMetricsStmt(ctx *grammar.ShowMetricsStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.Metric)
}

// EnterShowFieldsStmt is called when production showFieldsStmt is entered.
func (l *listener) EnterShowFieldsStmt(ctx *grammar.ShowFieldsStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.Field)
}

// EnterShowTagKeysStmt is called when production showTagKeysStmt is entered.
func (l *listener) EnterShowTagKeysStmt(ctx *grammar.ShowTagKeysStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.TagKey)
}

// EnterShowTagValuesStmt is called when production showTagValuesStmt is entered.
func (l *listener) EnterShowTagValuesStmt(ctx *grammar.ShowTagValuesStmtContext) {
	l.metaStmt = newMetaStmtParser(stmt.TagValue)
}

// EnterNamespace is called when production namespace is entered.
func (l *listener) EnterNamespace(ctx *grammar.NamespaceContext) {
	switch {
	case l.stmt != nil:
		l.stmt.visitNamespace(ctx)
	case l.metaStmt != nil:
		l.metaStmt.visitNamespace(ctx)
	}
}

// EnterWithTagKey is called when production withTagKey is entered.
func (l *listener) EnterWithTagKey(ctx *grammar.WithTagKeyContext) {
	if l.metaStmt != nil {
		l.metaStmt.visitWithTagKey(ctx)
	}
}

// EnterPrefix is called when production prefix is entered.
func (l *listener) EnterPrefix(ctx *grammar.PrefixContext) {
	if l.metaStmt != nil {
		l.metaStmt.visitPrefix(ctx)
	}
}

// EnterMetricName is called when production metricName is entered.
func (l *listener) EnterMetricName(ctx *grammar.MetricNameContext) {
	switch {
	case l.stmt != nil:
		l.stmt.visitMetricName(ctx)
	case l.metaStmt != nil:
		l.metaStmt.visitMetricName(ctx)
	}
}

// EnterSelectExpr is called when production selectExpr is entered.
func (l *listener) EnterSelectExpr(ctx *grammar.SelectExprContext) {
	if l.stmt != nil {
		l.stmt.resetExprStack()
	}
}

// EnterWhereClause is called when production whereClause is entered.
func (l *listener) EnterWhereClause(ctx *grammar.WhereClauseContext) {
	if l.stmt != nil {
		l.stmt.resetExprStack()
	}
}

// EnterFieldExpr is called when production fieldExpr is entered.
func (l *listener) EnterFieldExpr(ctx *grammar.FieldExprContext) {
	if l.stmt != nil {
		l.stmt.visitFieldExpr(ctx)
	}
}

// ExitFieldExpr is called when production fieldExpr is exited.
func (l *listener) ExitFieldExpr(ctx *grammar.FieldExprContext) {
	if l.stmt != nil {
		l.stmt.completeFieldExpr(ctx)
	}
}

// EnterFuncName is called when production exprFunc is entered.
func (l *listener) EnterFuncName(ctx *grammar.FuncNameContext) {
	if l.stmt != nil {
		l.stmt.visitFuncName(ctx)
	}
}

// ExitExprFunc is called when production exprFunc is exited.
func (l *listener) ExitExprFunc(ctx *grammar.ExprFuncContext) {
	if l.stmt != nil {
		l.stmt.completeFuncExpr()
	}
}

// EnterExprAtom is called when production exprAtom is entered.
func (l *listener) EnterExprAtom(ctx *grammar.ExprAtomContext) {
	if l.stmt != nil {
		l.stmt.visitExprAtom(ctx)
	}
}

// EnterAlias is called when production alias is entered.
func (l *listener) EnterAlias(ctx *grammar.AliasContext) {
	if l.stmt != nil {
		l.stmt.visitAlias(ctx)
	}
}

// EnterLimitClause is called when production limitClause is entered.
func (l *listener) EnterLimitClause(ctx *grammar.LimitClauseContext) {
	switch {
	case l.stmt != nil:
		l.stmt.visitLimit(ctx)
	case l.metaStmt != nil:
		l.metaStmt.visitLimit(ctx)
	}
}

// EnterTagFilterExpr is called when production tagFilterExpr is entered.
func (l *listener) EnterTagFilterExpr(ctx *grammar.TagFilterExprContext) {
	switch {
	case l.stmt != nil:
		l.stmt.visitTagFilterExpr(ctx)
	case l.metaStmt != nil:
		l.metaStmt.visitTagFilterExpr(ctx)
	}
}

// ExitTagFilterExpr is called when production tagValueList is exited.
func (l *listener) ExitTagFilterExpr(ctx *grammar.TagFilterExprContext) {
	switch {
	case l.stmt != nil:
		l.stmt.completeTagFilterExpr()
	case l.metaStmt != nil:
		l.metaStmt.completeTagFilterExpr()
	}
}

// EnterTagValue is called when production tagValue is entered.
func (l *listener) EnterTagValue(ctx *grammar.TagValueContext) {
	switch {
	case l.stmt != nil:
		l.stmt.visitTagValue(ctx)
	case l.metaStmt != nil:
		l.metaStmt.visitTagValue(ctx)
	}
}

// EnterTimeRangeExpr is called when production timeRangeExpr is entered.
func (l *listener) EnterTimeRangeExpr(ctx *grammar.TimeRangeExprContext) {
	if l.stmt != nil {
		l.stmt.visitTimeRangeExpr(ctx)
	}
}

// EnterGroupByClause is called when production groupByClause is entered.
func (l *listener) EnterGroupByKey(ctx *grammar.GroupByKeyContext) {
	if l.stmt != nil {
		l.stmt.visitGroupByKey(ctx)
	}
}

// statement returns query statement, if failure return error
func (l *listener) statement() (stmt.Statement, error) {
	if l.stmt != nil {
		return l.stmt.build()
	} else if l.metaStmt != nil {
		return l.metaStmt.build()
	}
	return nil, nil
}
