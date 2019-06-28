// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/cmd/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseSQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitStatement_list(ctx *Statement_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitCreate_database_stmt(ctx *Create_database_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWith_clause_list(ctx *With_clause_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWith_clause(ctx *With_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInterval_define_list(ctx *Interval_define_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInterval_define(ctx *Interval_defineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShard_num(ctx *Shard_numContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTtl_val(ctx *Ttl_valContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetattl_val(ctx *Metattl_valContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitPast_val(ctx *Past_valContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFuture_val(ctx *Future_valContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInterval_name_val(ctx *Interval_name_valContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitReplica_factor(ctx *Replica_factorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDatabase_name(ctx *Database_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitUpdate_database_stmt(ctx *Update_database_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDrop_database_stmt(ctx *Drop_database_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_databases_stmt(ctx *Show_databases_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_node_stmt(ctx *Show_node_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_measurements_stmt(ctx *Show_measurements_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_tag_keys_stmt(ctx *Show_tag_keys_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_info_stmt(ctx *Show_info_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_tag_values_stmt(ctx *Show_tag_values_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_tag_values_info_stmt(ctx *Show_tag_values_info_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_field_keys_stmt(ctx *Show_field_keys_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_queries_stmt(ctx *Show_queries_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitShow_stats_stmt(ctx *Show_stats_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWith_measurement_clause(ctx *With_measurement_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWith_tag_clause(ctx *With_tag_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhere_tag_cascade(ctx *Where_tag_cascadeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitKill_query_stmt(ctx *Kill_query_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQuery_id(ctx *Query_idContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitServer_id(ctx *Server_idContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitModule(ctx *ModuleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitComponent(ctx *ComponentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitQuery_stmt(ctx *Query_stmtContext) interface{} {
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

func (v *BaseSQLVisitor) VisitFrom_clause(ctx *From_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitWhere_clause(ctx *Where_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitClause_boolean_expr(ctx *Clause_boolean_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_cascade_expr(ctx *Tag_cascade_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_equal_expr(ctx *Tag_equal_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_boolean_expr(ctx *Tag_boolean_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_value_list(ctx *Tag_value_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTime_expr(ctx *Time_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTime_boolean_expr(ctx *Time_boolean_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNow_expr(ctx *Now_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNow_func(ctx *Now_funcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitGroup_by_clause(ctx *Group_by_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDimensions(ctx *DimensionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDimension(ctx *DimensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFill_option(ctx *Fill_optionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitOrder_by_clause(ctx *Order_by_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInterval_by_clause(ctx *Interval_by_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSort_field(ctx *Sort_fieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitSort_fields(ctx *Sort_fieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitHaving_clause(ctx *Having_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBool_expr(ctx *Bool_exprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBool_expr_logical_op(ctx *Bool_expr_logical_opContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBool_expr_atom(ctx *Bool_expr_atomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBool_expr_binary(ctx *Bool_expr_binaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitBool_expr_binary_operator(ctx *Bool_expr_binary_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDuration_lit(ctx *Duration_litContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInterval_item(ctx *Interval_itemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpr_func(ctx *Expr_funcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpr_func_params(ctx *Expr_func_paramsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitFunc_param(ctx *Func_paramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitExpr_atom(ctx *Expr_atomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIdent_filter(ctx *Ident_filterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitInt_number(ctx *Int_numberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitDec_number(ctx *Dec_numberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitLimit_clause(ctx *Limit_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitMetric_name(ctx *Metric_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_key(ctx *Tag_keyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_value(ctx *Tag_valueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitTag_value_pattern(ctx *Tag_value_patternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLVisitor) VisitNon_reserved_words(ctx *Non_reserved_wordsContext) interface{} {
	return v.VisitChildren(ctx)
}
