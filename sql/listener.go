package sql

import (
	"github.com/eleme/lindb/sql/grammar"
	"github.com/eleme/lindb/sql/stmt"
)

type listener struct {
	*grammar.BaseSQLListener
	stmt *queryStmtParse
}

// EnterQueryStmt is called when production queryStmt is entered.
func (l *listener) EnterQueryStmt(ctx *grammar.QueryStmtContext) {
	l.stmt = newQueryStmtParse()
}

// EnterMetricName is called when production metricName is entered.
func (l *listener) EnterMetricName(ctx *grammar.MetricNameContext) {
	if l.stmt != nil {
		l.stmt.visitMetricName(ctx)
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

// EnterExprFunc is called when production exprFunc is entered.
func (l *listener) EnterExprFunc(ctx *grammar.ExprFuncContext) {
	if l.stmt != nil {
		l.stmt.visitExprFunc(ctx)
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
	if l.stmt != nil {
		l.stmt.visitLimit(ctx)
	}
}

// EnterTagFilterExpr is called when production tagFilterExpr is entered.
func (l *listener) EnterTagFilterExpr(ctx *grammar.TagFilterExprContext) {
	if l.stmt != nil {
		l.stmt.visitTagFilterExpr(ctx)
	}
}

// ExitTagFilterExpr is called when production tagValueList is exited.
func (l *listener) ExitTagFilterExpr(ctx *grammar.TagFilterExprContext) {
	if l.stmt != nil {
		l.stmt.completeTagFilterExpr()
	}
}

// EnterTagValue is called when production tagValue is entered.
func (l *listener) EnterTagValue(ctx *grammar.TagValueContext) {
	if l.stmt != nil {
		l.stmt.visitTagValue(ctx)
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
func (l *listener) statement() (*stmt.Query, error) {
	if l.stmt != nil {
		return l.stmt.build()

	}
	return nil, nil
}
