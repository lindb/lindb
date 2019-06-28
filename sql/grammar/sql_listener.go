// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/cmd/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SQLListener is a complete listener for a parse tree produced by SQLParser.
type SQLListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterStatement_list is called when entering the statement_list production.
	EnterStatement_list(c *Statement_listContext)

	// EnterCreate_database_stmt is called when entering the create_database_stmt production.
	EnterCreate_database_stmt(c *Create_database_stmtContext)

	// EnterWith_clause_list is called when entering the with_clause_list production.
	EnterWith_clause_list(c *With_clause_listContext)

	// EnterWith_clause is called when entering the with_clause production.
	EnterWith_clause(c *With_clauseContext)

	// EnterInterval_define_list is called when entering the interval_define_list production.
	EnterInterval_define_list(c *Interval_define_listContext)

	// EnterInterval_define is called when entering the interval_define production.
	EnterInterval_define(c *Interval_defineContext)

	// EnterShard_num is called when entering the shard_num production.
	EnterShard_num(c *Shard_numContext)

	// EnterTtl_val is called when entering the ttl_val production.
	EnterTtl_val(c *Ttl_valContext)

	// EnterMetattl_val is called when entering the metattl_val production.
	EnterMetattl_val(c *Metattl_valContext)

	// EnterPast_val is called when entering the past_val production.
	EnterPast_val(c *Past_valContext)

	// EnterFuture_val is called when entering the future_val production.
	EnterFuture_val(c *Future_valContext)

	// EnterInterval_name_val is called when entering the interval_name_val production.
	EnterInterval_name_val(c *Interval_name_valContext)

	// EnterReplica_factor is called when entering the replica_factor production.
	EnterReplica_factor(c *Replica_factorContext)

	// EnterDatabase_name is called when entering the database_name production.
	EnterDatabase_name(c *Database_nameContext)

	// EnterUpdate_database_stmt is called when entering the update_database_stmt production.
	EnterUpdate_database_stmt(c *Update_database_stmtContext)

	// EnterDrop_database_stmt is called when entering the drop_database_stmt production.
	EnterDrop_database_stmt(c *Drop_database_stmtContext)

	// EnterShow_databases_stmt is called when entering the show_databases_stmt production.
	EnterShow_databases_stmt(c *Show_databases_stmtContext)

	// EnterShow_node_stmt is called when entering the show_node_stmt production.
	EnterShow_node_stmt(c *Show_node_stmtContext)

	// EnterShow_measurements_stmt is called when entering the show_measurements_stmt production.
	EnterShow_measurements_stmt(c *Show_measurements_stmtContext)

	// EnterShow_tag_keys_stmt is called when entering the show_tag_keys_stmt production.
	EnterShow_tag_keys_stmt(c *Show_tag_keys_stmtContext)

	// EnterShow_info_stmt is called when entering the show_info_stmt production.
	EnterShow_info_stmt(c *Show_info_stmtContext)

	// EnterShow_tag_values_stmt is called when entering the show_tag_values_stmt production.
	EnterShow_tag_values_stmt(c *Show_tag_values_stmtContext)

	// EnterShow_tag_values_info_stmt is called when entering the show_tag_values_info_stmt production.
	EnterShow_tag_values_info_stmt(c *Show_tag_values_info_stmtContext)

	// EnterShow_field_keys_stmt is called when entering the show_field_keys_stmt production.
	EnterShow_field_keys_stmt(c *Show_field_keys_stmtContext)

	// EnterShow_queries_stmt is called when entering the show_queries_stmt production.
	EnterShow_queries_stmt(c *Show_queries_stmtContext)

	// EnterShow_stats_stmt is called when entering the show_stats_stmt production.
	EnterShow_stats_stmt(c *Show_stats_stmtContext)

	// EnterWith_measurement_clause is called when entering the with_measurement_clause production.
	EnterWith_measurement_clause(c *With_measurement_clauseContext)

	// EnterWith_tag_clause is called when entering the with_tag_clause production.
	EnterWith_tag_clause(c *With_tag_clauseContext)

	// EnterWhere_tag_cascade is called when entering the where_tag_cascade production.
	EnterWhere_tag_cascade(c *Where_tag_cascadeContext)

	// EnterKill_query_stmt is called when entering the kill_query_stmt production.
	EnterKill_query_stmt(c *Kill_query_stmtContext)

	// EnterQuery_id is called when entering the query_id production.
	EnterQuery_id(c *Query_idContext)

	// EnterServer_id is called when entering the server_id production.
	EnterServer_id(c *Server_idContext)

	// EnterModule is called when entering the module production.
	EnterModule(c *ModuleContext)

	// EnterComponent is called when entering the component production.
	EnterComponent(c *ComponentContext)

	// EnterQuery_stmt is called when entering the query_stmt production.
	EnterQuery_stmt(c *Query_stmtContext)

	// EnterFields is called when entering the fields production.
	EnterFields(c *FieldsContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterAlias is called when entering the alias production.
	EnterAlias(c *AliasContext)

	// EnterFrom_clause is called when entering the from_clause production.
	EnterFrom_clause(c *From_clauseContext)

	// EnterWhere_clause is called when entering the where_clause production.
	EnterWhere_clause(c *Where_clauseContext)

	// EnterClause_boolean_expr is called when entering the clause_boolean_expr production.
	EnterClause_boolean_expr(c *Clause_boolean_exprContext)

	// EnterTag_cascade_expr is called when entering the tag_cascade_expr production.
	EnterTag_cascade_expr(c *Tag_cascade_exprContext)

	// EnterTag_equal_expr is called when entering the tag_equal_expr production.
	EnterTag_equal_expr(c *Tag_equal_exprContext)

	// EnterTag_boolean_expr is called when entering the tag_boolean_expr production.
	EnterTag_boolean_expr(c *Tag_boolean_exprContext)

	// EnterTag_value_list is called when entering the tag_value_list production.
	EnterTag_value_list(c *Tag_value_listContext)

	// EnterTime_expr is called when entering the time_expr production.
	EnterTime_expr(c *Time_exprContext)

	// EnterTime_boolean_expr is called when entering the time_boolean_expr production.
	EnterTime_boolean_expr(c *Time_boolean_exprContext)

	// EnterNow_expr is called when entering the now_expr production.
	EnterNow_expr(c *Now_exprContext)

	// EnterNow_func is called when entering the now_func production.
	EnterNow_func(c *Now_funcContext)

	// EnterGroup_by_clause is called when entering the group_by_clause production.
	EnterGroup_by_clause(c *Group_by_clauseContext)

	// EnterDimensions is called when entering the dimensions production.
	EnterDimensions(c *DimensionsContext)

	// EnterDimension is called when entering the dimension production.
	EnterDimension(c *DimensionContext)

	// EnterFill_option is called when entering the fill_option production.
	EnterFill_option(c *Fill_optionContext)

	// EnterOrder_by_clause is called when entering the order_by_clause production.
	EnterOrder_by_clause(c *Order_by_clauseContext)

	// EnterInterval_by_clause is called when entering the interval_by_clause production.
	EnterInterval_by_clause(c *Interval_by_clauseContext)

	// EnterSort_field is called when entering the sort_field production.
	EnterSort_field(c *Sort_fieldContext)

	// EnterSort_fields is called when entering the sort_fields production.
	EnterSort_fields(c *Sort_fieldsContext)

	// EnterHaving_clause is called when entering the having_clause production.
	EnterHaving_clause(c *Having_clauseContext)

	// EnterBool_expr is called when entering the bool_expr production.
	EnterBool_expr(c *Bool_exprContext)

	// EnterBool_expr_logical_op is called when entering the bool_expr_logical_op production.
	EnterBool_expr_logical_op(c *Bool_expr_logical_opContext)

	// EnterBool_expr_atom is called when entering the bool_expr_atom production.
	EnterBool_expr_atom(c *Bool_expr_atomContext)

	// EnterBool_expr_binary is called when entering the bool_expr_binary production.
	EnterBool_expr_binary(c *Bool_expr_binaryContext)

	// EnterBool_expr_binary_operator is called when entering the bool_expr_binary_operator production.
	EnterBool_expr_binary_operator(c *Bool_expr_binary_operatorContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterDuration_lit is called when entering the duration_lit production.
	EnterDuration_lit(c *Duration_litContext)

	// EnterInterval_item is called when entering the interval_item production.
	EnterInterval_item(c *Interval_itemContext)

	// EnterExpr_func is called when entering the expr_func production.
	EnterExpr_func(c *Expr_funcContext)

	// EnterExpr_func_params is called when entering the expr_func_params production.
	EnterExpr_func_params(c *Expr_func_paramsContext)

	// EnterFunc_param is called when entering the func_param production.
	EnterFunc_param(c *Func_paramContext)

	// EnterExpr_atom is called when entering the expr_atom production.
	EnterExpr_atom(c *Expr_atomContext)

	// EnterIdent_filter is called when entering the ident_filter production.
	EnterIdent_filter(c *Ident_filterContext)

	// EnterInt_number is called when entering the int_number production.
	EnterInt_number(c *Int_numberContext)

	// EnterDec_number is called when entering the dec_number production.
	EnterDec_number(c *Dec_numberContext)

	// EnterLimit_clause is called when entering the limit_clause production.
	EnterLimit_clause(c *Limit_clauseContext)

	// EnterMetric_name is called when entering the metric_name production.
	EnterMetric_name(c *Metric_nameContext)

	// EnterTag_key is called when entering the tag_key production.
	EnterTag_key(c *Tag_keyContext)

	// EnterTag_value is called when entering the tag_value production.
	EnterTag_value(c *Tag_valueContext)

	// EnterTag_value_pattern is called when entering the tag_value_pattern production.
	EnterTag_value_pattern(c *Tag_value_patternContext)

	// EnterIdent is called when entering the ident production.
	EnterIdent(c *IdentContext)

	// EnterNon_reserved_words is called when entering the non_reserved_words production.
	EnterNon_reserved_words(c *Non_reserved_wordsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitStatement_list is called when exiting the statement_list production.
	ExitStatement_list(c *Statement_listContext)

	// ExitCreate_database_stmt is called when exiting the create_database_stmt production.
	ExitCreate_database_stmt(c *Create_database_stmtContext)

	// ExitWith_clause_list is called when exiting the with_clause_list production.
	ExitWith_clause_list(c *With_clause_listContext)

	// ExitWith_clause is called when exiting the with_clause production.
	ExitWith_clause(c *With_clauseContext)

	// ExitInterval_define_list is called when exiting the interval_define_list production.
	ExitInterval_define_list(c *Interval_define_listContext)

	// ExitInterval_define is called when exiting the interval_define production.
	ExitInterval_define(c *Interval_defineContext)

	// ExitShard_num is called when exiting the shard_num production.
	ExitShard_num(c *Shard_numContext)

	// ExitTtl_val is called when exiting the ttl_val production.
	ExitTtl_val(c *Ttl_valContext)

	// ExitMetattl_val is called when exiting the metattl_val production.
	ExitMetattl_val(c *Metattl_valContext)

	// ExitPast_val is called when exiting the past_val production.
	ExitPast_val(c *Past_valContext)

	// ExitFuture_val is called when exiting the future_val production.
	ExitFuture_val(c *Future_valContext)

	// ExitInterval_name_val is called when exiting the interval_name_val production.
	ExitInterval_name_val(c *Interval_name_valContext)

	// ExitReplica_factor is called when exiting the replica_factor production.
	ExitReplica_factor(c *Replica_factorContext)

	// ExitDatabase_name is called when exiting the database_name production.
	ExitDatabase_name(c *Database_nameContext)

	// ExitUpdate_database_stmt is called when exiting the update_database_stmt production.
	ExitUpdate_database_stmt(c *Update_database_stmtContext)

	// ExitDrop_database_stmt is called when exiting the drop_database_stmt production.
	ExitDrop_database_stmt(c *Drop_database_stmtContext)

	// ExitShow_databases_stmt is called when exiting the show_databases_stmt production.
	ExitShow_databases_stmt(c *Show_databases_stmtContext)

	// ExitShow_node_stmt is called when exiting the show_node_stmt production.
	ExitShow_node_stmt(c *Show_node_stmtContext)

	// ExitShow_measurements_stmt is called when exiting the show_measurements_stmt production.
	ExitShow_measurements_stmt(c *Show_measurements_stmtContext)

	// ExitShow_tag_keys_stmt is called when exiting the show_tag_keys_stmt production.
	ExitShow_tag_keys_stmt(c *Show_tag_keys_stmtContext)

	// ExitShow_info_stmt is called when exiting the show_info_stmt production.
	ExitShow_info_stmt(c *Show_info_stmtContext)

	// ExitShow_tag_values_stmt is called when exiting the show_tag_values_stmt production.
	ExitShow_tag_values_stmt(c *Show_tag_values_stmtContext)

	// ExitShow_tag_values_info_stmt is called when exiting the show_tag_values_info_stmt production.
	ExitShow_tag_values_info_stmt(c *Show_tag_values_info_stmtContext)

	// ExitShow_field_keys_stmt is called when exiting the show_field_keys_stmt production.
	ExitShow_field_keys_stmt(c *Show_field_keys_stmtContext)

	// ExitShow_queries_stmt is called when exiting the show_queries_stmt production.
	ExitShow_queries_stmt(c *Show_queries_stmtContext)

	// ExitShow_stats_stmt is called when exiting the show_stats_stmt production.
	ExitShow_stats_stmt(c *Show_stats_stmtContext)

	// ExitWith_measurement_clause is called when exiting the with_measurement_clause production.
	ExitWith_measurement_clause(c *With_measurement_clauseContext)

	// ExitWith_tag_clause is called when exiting the with_tag_clause production.
	ExitWith_tag_clause(c *With_tag_clauseContext)

	// ExitWhere_tag_cascade is called when exiting the where_tag_cascade production.
	ExitWhere_tag_cascade(c *Where_tag_cascadeContext)

	// ExitKill_query_stmt is called when exiting the kill_query_stmt production.
	ExitKill_query_stmt(c *Kill_query_stmtContext)

	// ExitQuery_id is called when exiting the query_id production.
	ExitQuery_id(c *Query_idContext)

	// ExitServer_id is called when exiting the server_id production.
	ExitServer_id(c *Server_idContext)

	// ExitModule is called when exiting the module production.
	ExitModule(c *ModuleContext)

	// ExitComponent is called when exiting the component production.
	ExitComponent(c *ComponentContext)

	// ExitQuery_stmt is called when exiting the query_stmt production.
	ExitQuery_stmt(c *Query_stmtContext)

	// ExitFields is called when exiting the fields production.
	ExitFields(c *FieldsContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitAlias is called when exiting the alias production.
	ExitAlias(c *AliasContext)

	// ExitFrom_clause is called when exiting the from_clause production.
	ExitFrom_clause(c *From_clauseContext)

	// ExitWhere_clause is called when exiting the where_clause production.
	ExitWhere_clause(c *Where_clauseContext)

	// ExitClause_boolean_expr is called when exiting the clause_boolean_expr production.
	ExitClause_boolean_expr(c *Clause_boolean_exprContext)

	// ExitTag_cascade_expr is called when exiting the tag_cascade_expr production.
	ExitTag_cascade_expr(c *Tag_cascade_exprContext)

	// ExitTag_equal_expr is called when exiting the tag_equal_expr production.
	ExitTag_equal_expr(c *Tag_equal_exprContext)

	// ExitTag_boolean_expr is called when exiting the tag_boolean_expr production.
	ExitTag_boolean_expr(c *Tag_boolean_exprContext)

	// ExitTag_value_list is called when exiting the tag_value_list production.
	ExitTag_value_list(c *Tag_value_listContext)

	// ExitTime_expr is called when exiting the time_expr production.
	ExitTime_expr(c *Time_exprContext)

	// ExitTime_boolean_expr is called when exiting the time_boolean_expr production.
	ExitTime_boolean_expr(c *Time_boolean_exprContext)

	// ExitNow_expr is called when exiting the now_expr production.
	ExitNow_expr(c *Now_exprContext)

	// ExitNow_func is called when exiting the now_func production.
	ExitNow_func(c *Now_funcContext)

	// ExitGroup_by_clause is called when exiting the group_by_clause production.
	ExitGroup_by_clause(c *Group_by_clauseContext)

	// ExitDimensions is called when exiting the dimensions production.
	ExitDimensions(c *DimensionsContext)

	// ExitDimension is called when exiting the dimension production.
	ExitDimension(c *DimensionContext)

	// ExitFill_option is called when exiting the fill_option production.
	ExitFill_option(c *Fill_optionContext)

	// ExitOrder_by_clause is called when exiting the order_by_clause production.
	ExitOrder_by_clause(c *Order_by_clauseContext)

	// ExitInterval_by_clause is called when exiting the interval_by_clause production.
	ExitInterval_by_clause(c *Interval_by_clauseContext)

	// ExitSort_field is called when exiting the sort_field production.
	ExitSort_field(c *Sort_fieldContext)

	// ExitSort_fields is called when exiting the sort_fields production.
	ExitSort_fields(c *Sort_fieldsContext)

	// ExitHaving_clause is called when exiting the having_clause production.
	ExitHaving_clause(c *Having_clauseContext)

	// ExitBool_expr is called when exiting the bool_expr production.
	ExitBool_expr(c *Bool_exprContext)

	// ExitBool_expr_logical_op is called when exiting the bool_expr_logical_op production.
	ExitBool_expr_logical_op(c *Bool_expr_logical_opContext)

	// ExitBool_expr_atom is called when exiting the bool_expr_atom production.
	ExitBool_expr_atom(c *Bool_expr_atomContext)

	// ExitBool_expr_binary is called when exiting the bool_expr_binary production.
	ExitBool_expr_binary(c *Bool_expr_binaryContext)

	// ExitBool_expr_binary_operator is called when exiting the bool_expr_binary_operator production.
	ExitBool_expr_binary_operator(c *Bool_expr_binary_operatorContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitDuration_lit is called when exiting the duration_lit production.
	ExitDuration_lit(c *Duration_litContext)

	// ExitInterval_item is called when exiting the interval_item production.
	ExitInterval_item(c *Interval_itemContext)

	// ExitExpr_func is called when exiting the expr_func production.
	ExitExpr_func(c *Expr_funcContext)

	// ExitExpr_func_params is called when exiting the expr_func_params production.
	ExitExpr_func_params(c *Expr_func_paramsContext)

	// ExitFunc_param is called when exiting the func_param production.
	ExitFunc_param(c *Func_paramContext)

	// ExitExpr_atom is called when exiting the expr_atom production.
	ExitExpr_atom(c *Expr_atomContext)

	// ExitIdent_filter is called when exiting the ident_filter production.
	ExitIdent_filter(c *Ident_filterContext)

	// ExitInt_number is called when exiting the int_number production.
	ExitInt_number(c *Int_numberContext)

	// ExitDec_number is called when exiting the dec_number production.
	ExitDec_number(c *Dec_numberContext)

	// ExitLimit_clause is called when exiting the limit_clause production.
	ExitLimit_clause(c *Limit_clauseContext)

	// ExitMetric_name is called when exiting the metric_name production.
	ExitMetric_name(c *Metric_nameContext)

	// ExitTag_key is called when exiting the tag_key production.
	ExitTag_key(c *Tag_keyContext)

	// ExitTag_value is called when exiting the tag_value production.
	ExitTag_value(c *Tag_valueContext)

	// ExitTag_value_pattern is called when exiting the tag_value_pattern production.
	ExitTag_value_pattern(c *Tag_value_patternContext)

	// ExitIdent is called when exiting the ident production.
	ExitIdent(c *IdentContext)

	// ExitNon_reserved_words is called when exiting the non_reserved_words production.
	ExitNon_reserved_words(c *Non_reserved_wordsContext)
}
