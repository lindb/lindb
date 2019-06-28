// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/cmd/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"
// A complete Visitor for a parse tree produced by SQLParser.
type SQLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SQLParser#statement_list.
	VisitStatement_list(ctx *Statement_listContext) interface{}

	// Visit a parse tree produced by SQLParser#create_database_stmt.
	VisitCreate_database_stmt(ctx *Create_database_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#with_clause_list.
	VisitWith_clause_list(ctx *With_clause_listContext) interface{}

	// Visit a parse tree produced by SQLParser#with_clause.
	VisitWith_clause(ctx *With_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#interval_define_list.
	VisitInterval_define_list(ctx *Interval_define_listContext) interface{}

	// Visit a parse tree produced by SQLParser#interval_define.
	VisitInterval_define(ctx *Interval_defineContext) interface{}

	// Visit a parse tree produced by SQLParser#shard_num.
	VisitShard_num(ctx *Shard_numContext) interface{}

	// Visit a parse tree produced by SQLParser#ttl_val.
	VisitTtl_val(ctx *Ttl_valContext) interface{}

	// Visit a parse tree produced by SQLParser#metattl_val.
	VisitMetattl_val(ctx *Metattl_valContext) interface{}

	// Visit a parse tree produced by SQLParser#past_val.
	VisitPast_val(ctx *Past_valContext) interface{}

	// Visit a parse tree produced by SQLParser#future_val.
	VisitFuture_val(ctx *Future_valContext) interface{}

	// Visit a parse tree produced by SQLParser#interval_name_val.
	VisitInterval_name_val(ctx *Interval_name_valContext) interface{}

	// Visit a parse tree produced by SQLParser#replica_factor.
	VisitReplica_factor(ctx *Replica_factorContext) interface{}

	// Visit a parse tree produced by SQLParser#database_name.
	VisitDatabase_name(ctx *Database_nameContext) interface{}

	// Visit a parse tree produced by SQLParser#update_database_stmt.
	VisitUpdate_database_stmt(ctx *Update_database_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#drop_database_stmt.
	VisitDrop_database_stmt(ctx *Drop_database_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_databases_stmt.
	VisitShow_databases_stmt(ctx *Show_databases_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_node_stmt.
	VisitShow_node_stmt(ctx *Show_node_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_measurements_stmt.
	VisitShow_measurements_stmt(ctx *Show_measurements_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_tag_keys_stmt.
	VisitShow_tag_keys_stmt(ctx *Show_tag_keys_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_info_stmt.
	VisitShow_info_stmt(ctx *Show_info_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_tag_values_stmt.
	VisitShow_tag_values_stmt(ctx *Show_tag_values_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_tag_values_info_stmt.
	VisitShow_tag_values_info_stmt(ctx *Show_tag_values_info_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_field_keys_stmt.
	VisitShow_field_keys_stmt(ctx *Show_field_keys_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_queries_stmt.
	VisitShow_queries_stmt(ctx *Show_queries_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#show_stats_stmt.
	VisitShow_stats_stmt(ctx *Show_stats_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#with_measurement_clause.
	VisitWith_measurement_clause(ctx *With_measurement_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#with_tag_clause.
	VisitWith_tag_clause(ctx *With_tag_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#where_tag_cascade.
	VisitWhere_tag_cascade(ctx *Where_tag_cascadeContext) interface{}

	// Visit a parse tree produced by SQLParser#kill_query_stmt.
	VisitKill_query_stmt(ctx *Kill_query_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#query_id.
	VisitQuery_id(ctx *Query_idContext) interface{}

	// Visit a parse tree produced by SQLParser#server_id.
	VisitServer_id(ctx *Server_idContext) interface{}

	// Visit a parse tree produced by SQLParser#module.
	VisitModule(ctx *ModuleContext) interface{}

	// Visit a parse tree produced by SQLParser#component.
	VisitComponent(ctx *ComponentContext) interface{}

	// Visit a parse tree produced by SQLParser#query_stmt.
	VisitQuery_stmt(ctx *Query_stmtContext) interface{}

	// Visit a parse tree produced by SQLParser#fields.
	VisitFields(ctx *FieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by SQLParser#alias.
	VisitAlias(ctx *AliasContext) interface{}

	// Visit a parse tree produced by SQLParser#from_clause.
	VisitFrom_clause(ctx *From_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#where_clause.
	VisitWhere_clause(ctx *Where_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#clause_boolean_expr.
	VisitClause_boolean_expr(ctx *Clause_boolean_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_cascade_expr.
	VisitTag_cascade_expr(ctx *Tag_cascade_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_equal_expr.
	VisitTag_equal_expr(ctx *Tag_equal_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_boolean_expr.
	VisitTag_boolean_expr(ctx *Tag_boolean_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_value_list.
	VisitTag_value_list(ctx *Tag_value_listContext) interface{}

	// Visit a parse tree produced by SQLParser#time_expr.
	VisitTime_expr(ctx *Time_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#time_boolean_expr.
	VisitTime_boolean_expr(ctx *Time_boolean_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#now_expr.
	VisitNow_expr(ctx *Now_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#now_func.
	VisitNow_func(ctx *Now_funcContext) interface{}

	// Visit a parse tree produced by SQLParser#group_by_clause.
	VisitGroup_by_clause(ctx *Group_by_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#dimensions.
	VisitDimensions(ctx *DimensionsContext) interface{}

	// Visit a parse tree produced by SQLParser#dimension.
	VisitDimension(ctx *DimensionContext) interface{}

	// Visit a parse tree produced by SQLParser#fill_option.
	VisitFill_option(ctx *Fill_optionContext) interface{}

	// Visit a parse tree produced by SQLParser#order_by_clause.
	VisitOrder_by_clause(ctx *Order_by_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#interval_by_clause.
	VisitInterval_by_clause(ctx *Interval_by_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#sort_field.
	VisitSort_field(ctx *Sort_fieldContext) interface{}

	// Visit a parse tree produced by SQLParser#sort_fields.
	VisitSort_fields(ctx *Sort_fieldsContext) interface{}

	// Visit a parse tree produced by SQLParser#having_clause.
	VisitHaving_clause(ctx *Having_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#bool_expr.
	VisitBool_expr(ctx *Bool_exprContext) interface{}

	// Visit a parse tree produced by SQLParser#bool_expr_logical_op.
	VisitBool_expr_logical_op(ctx *Bool_expr_logical_opContext) interface{}

	// Visit a parse tree produced by SQLParser#bool_expr_atom.
	VisitBool_expr_atom(ctx *Bool_expr_atomContext) interface{}

	// Visit a parse tree produced by SQLParser#bool_expr_binary.
	VisitBool_expr_binary(ctx *Bool_expr_binaryContext) interface{}

	// Visit a parse tree produced by SQLParser#bool_expr_binary_operator.
	VisitBool_expr_binary_operator(ctx *Bool_expr_binary_operatorContext) interface{}

	// Visit a parse tree produced by SQLParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by SQLParser#duration_lit.
	VisitDuration_lit(ctx *Duration_litContext) interface{}

	// Visit a parse tree produced by SQLParser#interval_item.
	VisitInterval_item(ctx *Interval_itemContext) interface{}

	// Visit a parse tree produced by SQLParser#expr_func.
	VisitExpr_func(ctx *Expr_funcContext) interface{}

	// Visit a parse tree produced by SQLParser#expr_func_params.
	VisitExpr_func_params(ctx *Expr_func_paramsContext) interface{}

	// Visit a parse tree produced by SQLParser#func_param.
	VisitFunc_param(ctx *Func_paramContext) interface{}

	// Visit a parse tree produced by SQLParser#expr_atom.
	VisitExpr_atom(ctx *Expr_atomContext) interface{}

	// Visit a parse tree produced by SQLParser#ident_filter.
	VisitIdent_filter(ctx *Ident_filterContext) interface{}

	// Visit a parse tree produced by SQLParser#int_number.
	VisitInt_number(ctx *Int_numberContext) interface{}

	// Visit a parse tree produced by SQLParser#dec_number.
	VisitDec_number(ctx *Dec_numberContext) interface{}

	// Visit a parse tree produced by SQLParser#limit_clause.
	VisitLimit_clause(ctx *Limit_clauseContext) interface{}

	// Visit a parse tree produced by SQLParser#metric_name.
	VisitMetric_name(ctx *Metric_nameContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_key.
	VisitTag_key(ctx *Tag_keyContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_value.
	VisitTag_value(ctx *Tag_valueContext) interface{}

	// Visit a parse tree produced by SQLParser#tag_value_pattern.
	VisitTag_value_pattern(ctx *Tag_value_patternContext) interface{}

	// Visit a parse tree produced by SQLParser#ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by SQLParser#non_reserved_words.
	VisitNon_reserved_words(ctx *Non_reserved_wordsContext) interface{}

}