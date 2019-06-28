package sql

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/eleme/lindb/pkg/proto"
	"github.com/eleme/lindb/sql/grammar"
	"github.com/eleme/lindb/sql/util"
	"strconv"
	"strings"
)

type Listener struct {
	*parser.BaseSQLListener
	statement *proto.Stmt
	stmt      *queryStatement
}

func (ql *Listener) InitSQLListener() {
	ql.statement = &proto.Stmt{}
}

func (ql *Listener) VisitTerminal(node antlr.TerminalNode) {
	if node.GetSymbol().GetTokenType() == antlr.TokenEOF {
		return
	}
}

// EnterDrop_database_stmt override lindb sql drop database stmt
func (ql *Listener) EnterDrop_database_stmt(ctx *parser.Drop_database_stmtContext) {
	fmt.Println("EnterDrop_database_stmt")
	dropDatabase := new(proto.DropDatabase)
	databaseName := ctx.Database_name().(*parser.Database_nameContext)
	dropDatabase.Database = util.GetStringValue(databaseName.Ident().GetText())
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_DropDatabase)
	if ok {
		x.DropDatabase = dropDatabase
	}
}

// EnterShow_stats_stmt override lindb sql show stats stmt
func (ql *Listener) EnterShow_stats_stmt(ctx *parser.Show_stats_stmtContext) {
	fmt.Println("EnterShow_stats_stmt")
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

// EnterShow_databases_stmt override lindb sql show databases stmt
func (ql *Listener) EnterShow_databases_stmt(ctx *parser.Show_databases_stmtContext) {
	//todo
	fmt.Println("EnterShow_databases_stmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowDatabases)
	if ok {
		x.ShowDatabases = new(proto.ShowDatabases)
	}
}

// EnterShow_node_stmt override lindb sql show node stmt
func (ql *Listener) EnterShow_node_stmt(ctx *parser.Show_node_stmtContext) {
	//todo
	fmt.Println("EnterShow_node_stmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowNode)
	if ok {
		x.ShowNode = new(proto.ShowNode)
	}
}

// EnterShow_queries_stmt override lindb sql show queries stmt
func (ql *Listener) EnterShow_queries_stmt(ctx *parser.Show_queries_stmtContext) {
	//todo
	fmt.Println("EnterShow_queries_stmt")
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowQueries)
	if ok {
		x.ShowQueries = new(proto.ShowQueries)
	}
}

// EnterKill_query_stmt override lindb sql kill query stmt
func (ql *Listener) EnterKill_query_stmt(ctx *parser.Kill_query_stmtContext) {
	fmt.Println("EnterKill_query_stmt")
	killQuery := new(proto.KillQuery)
	queryId := ctx.Query_id().GetText()
	id, _ := strconv.ParseInt(queryId, 10, 64)
	killQuery.QueryId = id
	iServerIdContext := ctx.Server_id()
	if iServerIdContext != nil {
		serverIdContext := iServerIdContext.(*parser.Server_idContext)
		serverId, _ := strconv.ParseInt(serverIdContext.L_INT().GetText(), 10, 32)
		killQuery.ServerId = int32(serverId)
	}
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_KillQuery)
	if ok {
		x.KillQuery = killQuery
	}
}

// EnterShow_measurements_stmt override lindb sql show measurements stmt
func (ql *Listener) EnterShow_measurements_stmt(ctx *parser.Show_measurements_stmtContext) {
	fmt.Println("EnterShow_measurements_stmt")
	showMetric := new(proto.ShowMetric)
	iWithMeasurementClauseContext := ctx.With_measurement_clause()
	if iWithMeasurementClauseContext != nil {
		withMeasurementClauseContext := iWithMeasurementClauseContext.(*parser.With_measurement_clauseContext)
		name := withMeasurementClauseContext.Metric_name().GetText()
		showMetric.Name = name
	}
	limit := parseLimitDefault(ctx.Limit_clause())
	showMetric.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowMetric)
	if ok {
		x.ShowMetric = showMetric
	}
}

// EnterShow_field_keys_stmt override lindb sql show field keys stmt
func (ql *Listener) EnterShow_field_keys_stmt(ctx *parser.Show_field_keys_stmtContext) {
	fmt.Println("EnterShow_field_keys_stmt")
	metric := ctx.Metric_name().GetText()
	showFieldKeys := new(proto.ShowFieldKeys)
	showFieldKeys.Measurement = util.GetStringValue(metric)
	limit := parseLimitDefault(ctx.Limit_clause())
	showFieldKeys.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowFieldKeys)
	if ok {
		x.ShowFieldKeys = showFieldKeys
	}
}

// EnterShow_tag_keys_stmt override lindb sql show tag keys stmt
func (ql *Listener) EnterShow_tag_keys_stmt(ctx *parser.Show_tag_keys_stmtContext) {
	fmt.Println("EnterShow_tag_keys_stmt")
	metric := ctx.Metric_name().GetText()
	showTagKeys := new(proto.ShowTagKeys)
	showTagKeys.Measurement = util.GetStringValue(metric)

	limit := parseLimitDefault(ctx.Limit_clause())
	showTagKeys.Limit = limit
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowTagKeys)
	if ok {
		x.ShowTagKeys = showTagKeys
	}
}

// EnterShow_info_stmt override lindb sql show info stmt
func (ql *Listener) EnterShow_info_stmt(ctx *parser.Show_info_stmtContext) {
	fmt.Println("EnterShow_info_stmt")
	metric := ctx.Metric_name().GetText()
	showInfo := new(proto.ShowInfo)
	showInfo.Measurement = util.GetStringValue(metric)
	//todo
	x, ok := ql.statement.GetStmt().(*proto.Stmt_ShowInfo)
	if ok {
		x.ShowInfo = showInfo
	}
}

// EnterShow_tag_values_stmt override lindb sql show tag values stmt
func (ql *Listener) EnterShow_tag_values_stmt(ctx *parser.Show_tag_values_stmtContext) {
	fmt.Println("EnterShow_tag_values_stmt")
	metric := ctx.Metric_name().GetText()
	withTagClause := ctx.With_tag_clause().(*parser.With_tag_clauseContext)
	tagKey := withTagClause.Tag_key().GetText()
	showTagValues := new(proto.ShowTagValues)
	showTagValues.Measurement = util.GetStringValue(metric)
	showTagValues.TagKey = util.GetStringValue(tagKey)

	limit := parseLimitDefault(ctx.Limit_clause())
	showTagValues.Limit = limit
	iWhereTagValueClauseContext := ctx.Where_tag_cascade()
	if iWhereTagValueClauseContext != nil {
		whereTagValueClauseContext := iWhereTagValueClauseContext.(*parser.Where_tag_cascadeContext)
		iTagCascadeExprContext := whereTagValueClauseContext.Tag_cascade_expr()
		if iTagCascadeExprContext != nil {
			tagCascadeExprContext := iTagCascadeExprContext.(*parser.Tag_cascade_exprContext)
			tagEqualExpr := tagCascadeExprContext.Tag_equal_expr().(*parser.Tag_equal_exprContext)
			var tagValuePattern = tagEqualExpr.Tag_value_pattern().GetText()
			if len(tagValuePattern) > 0 {
				tagValuePattern = util.GetStringValue(tagValuePattern)
				index := strings.Index(tagValuePattern, "*")
				if index >= 0 {
					showTagValues.TagValue = tagValuePattern[0:index]
				}
			}
			tagBooleanExprContext := tagCascadeExprContext.Tag_boolean_expr()
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

// EnterShow_tag_values_info_stmt override lindb sql show tag values info stmt
func (ql *Listener) EnterShow_tag_values_info_stmt(ctx *parser.Show_tag_values_info_stmtContext) {
	metric := ctx.Metric_name().GetText()
	withTagClause := ctx.With_tag_clause().(*parser.With_tag_clauseContext)
	tagKey := withTagClause.Tag_key().GetText()
	showTagValuesInfo := new(proto.ShowTagValuesInfo)
	showTagValuesInfo.Measurement = util.GetStringValue(metric)
	showTagValuesInfo.TagKey = util.GetStringValue(tagKey)

	iWhereTagValueClauseContext := ctx.Where_tag_cascade()
	if iWhereTagValueClauseContext != nil {
		whereTagValueClauseContext := iWhereTagValueClauseContext.(*parser.Where_tag_cascadeContext)
		iTagCascadeExprContext := whereTagValueClauseContext.Tag_cascade_expr()
		if iTagCascadeExprContext != nil {
			tagCascadeExprContext := iTagCascadeExprContext.(*parser.Tag_cascade_exprContext)
			tagEqualExpr := tagCascadeExprContext.Tag_equal_expr().(*parser.Tag_equal_exprContext)
			var tagValuePattern = tagEqualExpr.Tag_value_pattern().GetText()
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

// EnterQuery_stmt override lindb sql query stmt
func (ql *Listener) EnterQuery_stmt(ctx *parser.Query_stmtContext) {
	fmt.Println("EnterQuery_stmt")
	ql.stmt = NewDefaultQueryStatement()
	ql.stmt.Parse(ctx)
}

func parseLimitDefault(ctx parser.ILimit_clauseContext) int32 {
	return parseLimit(ctx, 50)
}

func parseLimit(ctx parser.ILimit_clauseContext, defaultValue int32) int32 {
	if ctx == nil {
		return defaultValue
	}
	limitClauseContext := ctx.(*parser.Limit_clauseContext)
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
