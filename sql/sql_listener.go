package sql

import (
	"fmt"

	parser "github.com/eleme/lindb/sql/grammar"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/eleme/lindb/pkg/proto"

	"github.com/eleme/lindb/sql/util"

	"strconv"
	"strings"
)

type Listener struct {
	*parser.BaseSQLListener
	statement *proto.Stmt
	stmt      *QueryStatement
}

func (ql *Listener) InitSQLListener() {
	ql.statement = &proto.Stmt{}
}

func (ql *Listener) VisitTerminal(node antlr.TerminalNode) {
	if node.GetSymbol().GetTokenType() == antlr.TokenEOF {
		return
	}
}

// EnterDropDatabaseStmt override lindb sql drop database stmt
func (ql *Listener) EnterDropDatabaseStmt(ctx *parser.DropDatabaseStmtContext) {
	fmt.Println("EnterDropDatabaseStmt")
	dropDatabase := new(proto.DropDatabase)
	databaseName := ctx.DatabaseName().(*parser.DatabaseNameContext)
	dropDatabase.Database = util.GetStringValue(databaseName.Ident().GetText())
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_DropDatabase)
	if ok {
		x.DropDatabase = dropDatabase
	}
}

// EnterShowStatsStmt override lindb sql show stats stmt
func (ql *Listener) EnterShowStatsStmt(ctx *parser.ShowStatsStmtContext) {
	fmt.Println("EnterShowStatsStmt")
	showStats := new(proto.ShowStats)
	iModuleContext := ctx.Module()
	if iModuleContext != nil {
		moduleContext := iModuleContext.(*parser.ModuleContext)
		s := moduleContext.Ident().GetText()
		showStats.Module = util.GetStringValue(s)
	}
	iComponentCtx := ctx.Component()
	if iComponentCtx != nil {
		componentCtx := iComponentCtx.(*parser.ComponentContext)
		showStats.Component = util.GetStringValue(componentCtx.Ident().GetText())
	}
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowStats)
	if ok {
		x.ShowStats = showStats
	}
}

// EnterShowDatabasesStmt override lindb sql show databases stmt
func (ql *Listener) EnterShowDatabasesStmt(ctx *parser.ShowDatabasesStmtContext) {
	//todo
	fmt.Println("EnterShowDatabasesStmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowDatabases)
	if ok {
		x.ShowDatabases = new(proto.ShowDatabases)
	}
}

// EnterShowNodeStmt override lindb sql show node stmt
func (ql *Listener) EnterShowNodeStmt(ctx *parser.ShowNodeStmtContext) {
	//todo
	fmt.Println("EnterShowNodeStmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowNode)
	if ok {
		x.ShowNode = new(proto.ShowNode)
	}
}

// EnterShowQueriesStmt override lindb sql show queries stmt
func (ql *Listener) EnterShowQueriesStmt(ctx *parser.ShowQueriesStmtContext) {
	//todo
	fmt.Println("EnterShowQueriesStmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowQueries)
	if ok {
		x.ShowQueries = new(proto.ShowQueries)
	}
}

// EnterKillQueryStmt override lindb sql kill query stmt
func (ql *Listener) EnterKillQueryStmt(ctx *parser.KillQueryStmtContext) {
	fmt.Println("EnterKillQueryStmt")
	killQuery := new(proto.KillQuery)
	queryID := ctx.QueryId().GetText()
	id, _ := strconv.ParseInt(queryID, 10, 64)
	killQuery.QueryId = id
	iServerIDContext := ctx.ServerId()
	if iServerIDContext != nil {
		serverIDContext := iServerIDContext.(*parser.ServerIdContext)
		serverID, _ := strconv.ParseInt(serverIDContext.L_INT().GetText(), 10, 32)
		killQuery.ServerId = int32(serverID)
	}
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_KillQuery)
	if ok {
		x.KillQuery = killQuery
	}
}

// EnterShowMeasurementsStmt override lindb sql show measurements stmt
func (ql *Listener) EnterShowMeasurementsStmt(ctx *parser.ShowMeasurementsStmtContext) {
	fmt.Println("EnterShowMeasurementsStmt")
	showMetric := new(proto.ShowMetric)
	iWithMeasurementClauseContext := ctx.WithMeasurementClause()
	if iWithMeasurementClauseContext != nil {
		withMeasurementClauseContext := iWithMeasurementClauseContext.(*parser.WithMeasurementClauseContext)
		name := withMeasurementClauseContext.MetricName().GetText()
		showMetric.Name = name
	}
	limit := parseLimitDefault(ctx.LimitClause())
	showMetric.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowMetric)
	if ok {
		x.ShowMetric = showMetric
	}
}

// EnterShowFieldKeysStmt override lindb sql show field keys stmt
func (ql *Listener) EnterShowFieldKeysStmt(ctx *parser.ShowFieldKeysStmtContext) {
	fmt.Println("EnterShowFieldKeysStmt")
	metric := ctx.MetricName().GetText()
	showFieldKeys := new(proto.ShowFieldKeys)
	showFieldKeys.Measurement = util.GetStringValue(metric)
	limit := parseLimitDefault(ctx.LimitClause())
	showFieldKeys.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowFieldKeys)
	if ok {
		x.ShowFieldKeys = showFieldKeys
	}
}

// EnterShowTagKeysStmt override lindb sql show tag keys stmt
func (ql *Listener) EnterShowTagKeysStmt(ctx *parser.ShowTagKeysStmtContext) {
	fmt.Println("EnterShowTagKeysStmt")
	metric := ctx.MetricName().GetText()
	showTagKeys := new(proto.ShowTagKeys)
	showTagKeys.Measurement = util.GetStringValue(metric)

	limit := parseLimitDefault(ctx.LimitClause())
	showTagKeys.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowTagKeys)
	if ok {
		x.ShowTagKeys = showTagKeys
	}
}

// EnterShowInfoStmt override lindb sql show info stmt
func (ql *Listener) EnterShowInfoStmt(ctx *parser.ShowInfoStmtContext) {
	fmt.Println("EnterShowInfoStmt")
	metric := ctx.MetricName().GetText()
	showInfo := new(proto.ShowInfo)
	showInfo.Measurement = util.GetStringValue(metric)
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowInfo)
	if ok {
		x.ShowInfo = showInfo
	}
}

// EnterShowTagValuesStmt override lindb sql show tag values stmt
func (ql *Listener) EnterShowTagValuesStmt(ctx *parser.ShowTagValuesStmtContext) {
	fmt.Println("EnterShowTagValuesStmt")
	metric := ctx.MetricName().GetText()
	withTagClause := ctx.WithTagClause().(*parser.WithTagClauseContext)
	tagKey := withTagClause.TagKey().GetText()
	showTagValues := new(proto.ShowTagValues)
	showTagValues.Measurement = util.GetStringValue(metric)
	showTagValues.TagKey = util.GetStringValue(tagKey)

	limit := parseLimitDefault(ctx.LimitClause())
	showTagValues.Limit = limit
	iWhereTagValueClauseContext := ctx.WhereTagCascade()
	if iWhereTagValueClauseContext != nil {
		whereTagValueClauseContext := iWhereTagValueClauseContext.(*parser.WhereTagCascadeContext)
		iTagCascadeExprContext := whereTagValueClauseContext.TagCascadeExpr()
		if iTagCascadeExprContext != nil {
			tagCascadeExprContext := iTagCascadeExprContext.(*parser.TagCascadeExprContext)
			tagEqualExpr := tagCascadeExprContext.TagEqualExpr().(*parser.TagEqualExprContext)
			var tagValuePattern = tagEqualExpr.TagValuePattern().GetText()
			if len(tagValuePattern) > 0 {
				tagValuePattern = util.GetStringValue(tagValuePattern)
				index := strings.Index(tagValuePattern, "*")
				if index >= 0 {
					showTagValues.TagValue = tagValuePattern[0:index]
				}
			}
			tagBooleanExprContext := tagCascadeExprContext.TagBooleanExpr()
			if tagBooleanExprContext != nil {
				condition := new(proto.Condition)
				showTagValues.Condition = condition
			}
		}
	}
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowTagValues)
	if ok {
		x.ShowTagValues = showTagValues
	}
}

// EnterShowTagValuesInfoStmt override lindb sql show tag values info stmt
func (ql *Listener) EnterShowTagValuesInfoStmt(ctx *parser.ShowTagValuesInfoStmtContext) {
	metric := ctx.MetricName().GetText()
	withTagClause := ctx.WithTagClause().(*parser.WithTagClauseContext)
	tagKey := withTagClause.TagKey().GetText()
	showTagValuesInfo := new(proto.ShowTagValuesInfo)
	showTagValuesInfo.Measurement = util.GetStringValue(metric)
	showTagValuesInfo.TagKey = util.GetStringValue(tagKey)

	iWhereTagValueClauseContext := ctx.WhereTagCascade()
	if iWhereTagValueClauseContext != nil {
		whereTagValueClauseContext := iWhereTagValueClauseContext.(*parser.WhereTagCascadeContext)
		iTagCascadeExprContext := whereTagValueClauseContext.TagCascadeExpr()
		if iTagCascadeExprContext != nil {
			tagCascadeExprContext := iTagCascadeExprContext.(*parser.TagCascadeExprContext)
			tagEqualExpr := tagCascadeExprContext.TagEqualExpr().(*parser.TagEqualExprContext)
			var tagValuePattern = tagEqualExpr.TagValuePattern().GetText()
			if len(tagValuePattern) > 0 {
				tagValuePattern = util.GetStringValue(tagValuePattern)
				showTagValuesInfo.TagValue = tagValuePattern
				//todo
				x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowTagValuesInfo)
				if ok {
					x.ShowTagValuesInfo = showTagValuesInfo
				}
				return
			}
		}
	}
	panic(fmt.Sprintln("show tag values info is not valid"))
}

// EnterQueryStmt override lindb sql query stmt
func (ql *Listener) EnterQueryStmt(ctx *parser.QueryStmtContext) {
	fmt.Println("EnterQueryStmt")
	ql.stmt = NewDefaultQueryStatement()
	ql.stmt.Parse(ctx)
}

func parseLimitDefault(ctx parser.ILimitClauseContext) int32 {
	return parseLimit(ctx, 50)
}

func parseLimit(ctx parser.ILimitClauseContext, defaultValue int32) int32 {
	if ctx == nil {
		return defaultValue
	}
	limitClauseContext := ctx.(*parser.LimitClauseContext)
	limit := limitClauseContext.L_INT().GetText()
	l, _ := strconv.ParseInt(limit, 10, 32)
	return int32(l)
}

func (ql *Listener) GetStatement() *proto.Stmt {
	if ql.stmt != nil {
		x, ok := ql.statement.GetStmt().(*proto.Stmt_Query)
		if ok {
			x.Query = ql.stmt.build()
		}

	}
	return ql.statement
}
