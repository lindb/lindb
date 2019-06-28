// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/cmd/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSQLListener is a complete listener for a parse tree produced by SQLParser.
type BaseSQLListener struct{}

var _ SQLListener = &BaseSQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseSQLListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseSQLListener) ExitStatement(ctx *StatementContext) {}

// EnterStatement_list is called when production statement_list is entered.
func (s *BaseSQLListener) EnterStatement_list(ctx *Statement_listContext) {}

// ExitStatement_list is called when production statement_list is exited.
func (s *BaseSQLListener) ExitStatement_list(ctx *Statement_listContext) {}

// EnterCreate_database_stmt is called when production create_database_stmt is entered.
func (s *BaseSQLListener) EnterCreate_database_stmt(ctx *Create_database_stmtContext) {}

// ExitCreate_database_stmt is called when production create_database_stmt is exited.
func (s *BaseSQLListener) ExitCreate_database_stmt(ctx *Create_database_stmtContext) {}

// EnterWith_clause_list is called when production with_clause_list is entered.
func (s *BaseSQLListener) EnterWith_clause_list(ctx *With_clause_listContext) {}

// ExitWith_clause_list is called when production with_clause_list is exited.
func (s *BaseSQLListener) ExitWith_clause_list(ctx *With_clause_listContext) {}

// EnterWith_clause is called when production with_clause is entered.
func (s *BaseSQLListener) EnterWith_clause(ctx *With_clauseContext) {}

// ExitWith_clause is called when production with_clause is exited.
func (s *BaseSQLListener) ExitWith_clause(ctx *With_clauseContext) {}

// EnterInterval_define_list is called when production interval_define_list is entered.
func (s *BaseSQLListener) EnterInterval_define_list(ctx *Interval_define_listContext) {}

// ExitInterval_define_list is called when production interval_define_list is exited.
func (s *BaseSQLListener) ExitInterval_define_list(ctx *Interval_define_listContext) {}

// EnterInterval_define is called when production interval_define is entered.
func (s *BaseSQLListener) EnterInterval_define(ctx *Interval_defineContext) {}

// ExitInterval_define is called when production interval_define is exited.
func (s *BaseSQLListener) ExitInterval_define(ctx *Interval_defineContext) {}

// EnterShard_num is called when production shard_num is entered.
func (s *BaseSQLListener) EnterShard_num(ctx *Shard_numContext) {}

// ExitShard_num is called when production shard_num is exited.
func (s *BaseSQLListener) ExitShard_num(ctx *Shard_numContext) {}

// EnterTtl_val is called when production ttl_val is entered.
func (s *BaseSQLListener) EnterTtl_val(ctx *Ttl_valContext) {}

// ExitTtl_val is called when production ttl_val is exited.
func (s *BaseSQLListener) ExitTtl_val(ctx *Ttl_valContext) {}

// EnterMetattl_val is called when production metattl_val is entered.
func (s *BaseSQLListener) EnterMetattl_val(ctx *Metattl_valContext) {}

// ExitMetattl_val is called when production metattl_val is exited.
func (s *BaseSQLListener) ExitMetattl_val(ctx *Metattl_valContext) {}

// EnterPast_val is called when production past_val is entered.
func (s *BaseSQLListener) EnterPast_val(ctx *Past_valContext) {}

// ExitPast_val is called when production past_val is exited.
func (s *BaseSQLListener) ExitPast_val(ctx *Past_valContext) {}

// EnterFuture_val is called when production future_val is entered.
func (s *BaseSQLListener) EnterFuture_val(ctx *Future_valContext) {}

// ExitFuture_val is called when production future_val is exited.
func (s *BaseSQLListener) ExitFuture_val(ctx *Future_valContext) {}

// EnterInterval_name_val is called when production interval_name_val is entered.
func (s *BaseSQLListener) EnterInterval_name_val(ctx *Interval_name_valContext) {}

// ExitInterval_name_val is called when production interval_name_val is exited.
func (s *BaseSQLListener) ExitInterval_name_val(ctx *Interval_name_valContext) {}

// EnterReplica_factor is called when production replica_factor is entered.
func (s *BaseSQLListener) EnterReplica_factor(ctx *Replica_factorContext) {}

// ExitReplica_factor is called when production replica_factor is exited.
func (s *BaseSQLListener) ExitReplica_factor(ctx *Replica_factorContext) {}

// EnterDatabase_name is called when production database_name is entered.
func (s *BaseSQLListener) EnterDatabase_name(ctx *Database_nameContext) {}

// ExitDatabase_name is called when production database_name is exited.
func (s *BaseSQLListener) ExitDatabase_name(ctx *Database_nameContext) {}

// EnterUpdate_database_stmt is called when production update_database_stmt is entered.
func (s *BaseSQLListener) EnterUpdate_database_stmt(ctx *Update_database_stmtContext) {}

// ExitUpdate_database_stmt is called when production update_database_stmt is exited.
func (s *BaseSQLListener) ExitUpdate_database_stmt(ctx *Update_database_stmtContext) {}

// EnterDrop_database_stmt is called when production drop_database_stmt is entered.
func (s *BaseSQLListener) EnterDrop_database_stmt(ctx *Drop_database_stmtContext) {}

// ExitDrop_database_stmt is called when production drop_database_stmt is exited.
func (s *BaseSQLListener) ExitDrop_database_stmt(ctx *Drop_database_stmtContext) {}

// EnterShow_databases_stmt is called when production show_databases_stmt is entered.
func (s *BaseSQLListener) EnterShow_databases_stmt(ctx *Show_databases_stmtContext) {}

// ExitShow_databases_stmt is called when production show_databases_stmt is exited.
func (s *BaseSQLListener) ExitShow_databases_stmt(ctx *Show_databases_stmtContext) {}

// EnterShow_node_stmt is called when production show_node_stmt is entered.
func (s *BaseSQLListener) EnterShow_node_stmt(ctx *Show_node_stmtContext) {}

// ExitShow_node_stmt is called when production show_node_stmt is exited.
func (s *BaseSQLListener) ExitShow_node_stmt(ctx *Show_node_stmtContext) {}

// EnterShow_measurements_stmt is called when production show_measurements_stmt is entered.
func (s *BaseSQLListener) EnterShow_measurements_stmt(ctx *Show_measurements_stmtContext) {}

// ExitShow_measurements_stmt is called when production show_measurements_stmt is exited.
func (s *BaseSQLListener) ExitShow_measurements_stmt(ctx *Show_measurements_stmtContext) {}

// EnterShow_tag_keys_stmt is called when production show_tag_keys_stmt is entered.
func (s *BaseSQLListener) EnterShow_tag_keys_stmt(ctx *Show_tag_keys_stmtContext) {}

// ExitShow_tag_keys_stmt is called when production show_tag_keys_stmt is exited.
func (s *BaseSQLListener) ExitShow_tag_keys_stmt(ctx *Show_tag_keys_stmtContext) {}

// EnterShow_info_stmt is called when production show_info_stmt is entered.
func (s *BaseSQLListener) EnterShow_info_stmt(ctx *Show_info_stmtContext) {}

// ExitShow_info_stmt is called when production show_info_stmt is exited.
func (s *BaseSQLListener) ExitShow_info_stmt(ctx *Show_info_stmtContext) {}

// EnterShow_tag_values_stmt is called when production show_tag_values_stmt is entered.
func (s *BaseSQLListener) EnterShow_tag_values_stmt(ctx *Show_tag_values_stmtContext) {}

// ExitShow_tag_values_stmt is called when production show_tag_values_stmt is exited.
func (s *BaseSQLListener) ExitShow_tag_values_stmt(ctx *Show_tag_values_stmtContext) {}

// EnterShow_tag_values_info_stmt is called when production show_tag_values_info_stmt is entered.
func (s *BaseSQLListener) EnterShow_tag_values_info_stmt(ctx *Show_tag_values_info_stmtContext) {}

// ExitShow_tag_values_info_stmt is called when production show_tag_values_info_stmt is exited.
func (s *BaseSQLListener) ExitShow_tag_values_info_stmt(ctx *Show_tag_values_info_stmtContext) {}

// EnterShow_field_keys_stmt is called when production show_field_keys_stmt is entered.
func (s *BaseSQLListener) EnterShow_field_keys_stmt(ctx *Show_field_keys_stmtContext) {}

// ExitShow_field_keys_stmt is called when production show_field_keys_stmt is exited.
func (s *BaseSQLListener) ExitShow_field_keys_stmt(ctx *Show_field_keys_stmtContext) {}

// EnterShow_queries_stmt is called when production show_queries_stmt is entered.
func (s *BaseSQLListener) EnterShow_queries_stmt(ctx *Show_queries_stmtContext) {}

// ExitShow_queries_stmt is called when production show_queries_stmt is exited.
func (s *BaseSQLListener) ExitShow_queries_stmt(ctx *Show_queries_stmtContext) {}

// EnterShow_stats_stmt is called when production show_stats_stmt is entered.
func (s *BaseSQLListener) EnterShow_stats_stmt(ctx *Show_stats_stmtContext) {}

// ExitShow_stats_stmt is called when production show_stats_stmt is exited.
func (s *BaseSQLListener) ExitShow_stats_stmt(ctx *Show_stats_stmtContext) {}

// EnterWith_measurement_clause is called when production with_measurement_clause is entered.
func (s *BaseSQLListener) EnterWith_measurement_clause(ctx *With_measurement_clauseContext) {}

// ExitWith_measurement_clause is called when production with_measurement_clause is exited.
func (s *BaseSQLListener) ExitWith_measurement_clause(ctx *With_measurement_clauseContext) {}

// EnterWith_tag_clause is called when production with_tag_clause is entered.
func (s *BaseSQLListener) EnterWith_tag_clause(ctx *With_tag_clauseContext) {}

// ExitWith_tag_clause is called when production with_tag_clause is exited.
func (s *BaseSQLListener) ExitWith_tag_clause(ctx *With_tag_clauseContext) {}

// EnterWhere_tag_cascade is called when production where_tag_cascade is entered.
func (s *BaseSQLListener) EnterWhere_tag_cascade(ctx *Where_tag_cascadeContext) {}

// ExitWhere_tag_cascade is called when production where_tag_cascade is exited.
func (s *BaseSQLListener) ExitWhere_tag_cascade(ctx *Where_tag_cascadeContext) {}

// EnterKill_query_stmt is called when production kill_query_stmt is entered.
func (s *BaseSQLListener) EnterKill_query_stmt(ctx *Kill_query_stmtContext) {}

// ExitKill_query_stmt is called when production kill_query_stmt is exited.
func (s *BaseSQLListener) ExitKill_query_stmt(ctx *Kill_query_stmtContext) {}

// EnterQuery_id is called when production query_id is entered.
func (s *BaseSQLListener) EnterQuery_id(ctx *Query_idContext) {}

// ExitQuery_id is called when production query_id is exited.
func (s *BaseSQLListener) ExitQuery_id(ctx *Query_idContext) {}

// EnterServer_id is called when production server_id is entered.
func (s *BaseSQLListener) EnterServer_id(ctx *Server_idContext) {}

// ExitServer_id is called when production server_id is exited.
func (s *BaseSQLListener) ExitServer_id(ctx *Server_idContext) {}

// EnterModule is called when production module is entered.
func (s *BaseSQLListener) EnterModule(ctx *ModuleContext) {}

// ExitModule is called when production module is exited.
func (s *BaseSQLListener) ExitModule(ctx *ModuleContext) {}

// EnterComponent is called when production component is entered.
func (s *BaseSQLListener) EnterComponent(ctx *ComponentContext) {}

// ExitComponent is called when production component is exited.
func (s *BaseSQLListener) ExitComponent(ctx *ComponentContext) {}

// EnterQuery_stmt is called when production query_stmt is entered.
func (s *BaseSQLListener) EnterQuery_stmt(ctx *Query_stmtContext) {}

// ExitQuery_stmt is called when production query_stmt is exited.
func (s *BaseSQLListener) ExitQuery_stmt(ctx *Query_stmtContext) {}

// EnterFields is called when production fields is entered.
func (s *BaseSQLListener) EnterFields(ctx *FieldsContext) {}

// ExitFields is called when production fields is exited.
func (s *BaseSQLListener) ExitFields(ctx *FieldsContext) {}

// EnterField is called when production field is entered.
func (s *BaseSQLListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseSQLListener) ExitField(ctx *FieldContext) {}

// EnterAlias is called when production alias is entered.
func (s *BaseSQLListener) EnterAlias(ctx *AliasContext) {}

// ExitAlias is called when production alias is exited.
func (s *BaseSQLListener) ExitAlias(ctx *AliasContext) {}

// EnterFrom_clause is called when production from_clause is entered.
func (s *BaseSQLListener) EnterFrom_clause(ctx *From_clauseContext) {}

// ExitFrom_clause is called when production from_clause is exited.
func (s *BaseSQLListener) ExitFrom_clause(ctx *From_clauseContext) {}

// EnterWhere_clause is called when production where_clause is entered.
func (s *BaseSQLListener) EnterWhere_clause(ctx *Where_clauseContext) {}

// ExitWhere_clause is called when production where_clause is exited.
func (s *BaseSQLListener) ExitWhere_clause(ctx *Where_clauseContext) {}

// EnterClause_boolean_expr is called when production clause_boolean_expr is entered.
func (s *BaseSQLListener) EnterClause_boolean_expr(ctx *Clause_boolean_exprContext) {}

// ExitClause_boolean_expr is called when production clause_boolean_expr is exited.
func (s *BaseSQLListener) ExitClause_boolean_expr(ctx *Clause_boolean_exprContext) {}

// EnterTag_cascade_expr is called when production tag_cascade_expr is entered.
func (s *BaseSQLListener) EnterTag_cascade_expr(ctx *Tag_cascade_exprContext) {}

// ExitTag_cascade_expr is called when production tag_cascade_expr is exited.
func (s *BaseSQLListener) ExitTag_cascade_expr(ctx *Tag_cascade_exprContext) {}

// EnterTag_equal_expr is called when production tag_equal_expr is entered.
func (s *BaseSQLListener) EnterTag_equal_expr(ctx *Tag_equal_exprContext) {}

// ExitTag_equal_expr is called when production tag_equal_expr is exited.
func (s *BaseSQLListener) ExitTag_equal_expr(ctx *Tag_equal_exprContext) {}

// EnterTag_boolean_expr is called when production tag_boolean_expr is entered.
func (s *BaseSQLListener) EnterTag_boolean_expr(ctx *Tag_boolean_exprContext) {}

// ExitTag_boolean_expr is called when production tag_boolean_expr is exited.
func (s *BaseSQLListener) ExitTag_boolean_expr(ctx *Tag_boolean_exprContext) {}

// EnterTag_value_list is called when production tag_value_list is entered.
func (s *BaseSQLListener) EnterTag_value_list(ctx *Tag_value_listContext) {}

// ExitTag_value_list is called when production tag_value_list is exited.
func (s *BaseSQLListener) ExitTag_value_list(ctx *Tag_value_listContext) {}

// EnterTime_expr is called when production time_expr is entered.
func (s *BaseSQLListener) EnterTime_expr(ctx *Time_exprContext) {}

// ExitTime_expr is called when production time_expr is exited.
func (s *BaseSQLListener) ExitTime_expr(ctx *Time_exprContext) {}

// EnterTime_boolean_expr is called when production time_boolean_expr is entered.
func (s *BaseSQLListener) EnterTime_boolean_expr(ctx *Time_boolean_exprContext) {}

// ExitTime_boolean_expr is called when production time_boolean_expr is exited.
func (s *BaseSQLListener) ExitTime_boolean_expr(ctx *Time_boolean_exprContext) {}

// EnterNow_expr is called when production now_expr is entered.
func (s *BaseSQLListener) EnterNow_expr(ctx *Now_exprContext) {}

// ExitNow_expr is called when production now_expr is exited.
func (s *BaseSQLListener) ExitNow_expr(ctx *Now_exprContext) {}

// EnterNow_func is called when production now_func is entered.
func (s *BaseSQLListener) EnterNow_func(ctx *Now_funcContext) {}

// ExitNow_func is called when production now_func is exited.
func (s *BaseSQLListener) ExitNow_func(ctx *Now_funcContext) {}

// EnterGroup_by_clause is called when production group_by_clause is entered.
func (s *BaseSQLListener) EnterGroup_by_clause(ctx *Group_by_clauseContext) {}

// ExitGroup_by_clause is called when production group_by_clause is exited.
func (s *BaseSQLListener) ExitGroup_by_clause(ctx *Group_by_clauseContext) {}

// EnterDimensions is called when production dimensions is entered.
func (s *BaseSQLListener) EnterDimensions(ctx *DimensionsContext) {}

// ExitDimensions is called when production dimensions is exited.
func (s *BaseSQLListener) ExitDimensions(ctx *DimensionsContext) {}

// EnterDimension is called when production dimension is entered.
func (s *BaseSQLListener) EnterDimension(ctx *DimensionContext) {}

// ExitDimension is called when production dimension is exited.
func (s *BaseSQLListener) ExitDimension(ctx *DimensionContext) {}

// EnterFill_option is called when production fill_option is entered.
func (s *BaseSQLListener) EnterFill_option(ctx *Fill_optionContext) {}

// ExitFill_option is called when production fill_option is exited.
func (s *BaseSQLListener) ExitFill_option(ctx *Fill_optionContext) {}

// EnterOrder_by_clause is called when production order_by_clause is entered.
func (s *BaseSQLListener) EnterOrder_by_clause(ctx *Order_by_clauseContext) {}

// ExitOrder_by_clause is called when production order_by_clause is exited.
func (s *BaseSQLListener) ExitOrder_by_clause(ctx *Order_by_clauseContext) {}

// EnterInterval_by_clause is called when production interval_by_clause is entered.
func (s *BaseSQLListener) EnterInterval_by_clause(ctx *Interval_by_clauseContext) {}

// ExitInterval_by_clause is called when production interval_by_clause is exited.
func (s *BaseSQLListener) ExitInterval_by_clause(ctx *Interval_by_clauseContext) {}

// EnterSort_field is called when production sort_field is entered.
func (s *BaseSQLListener) EnterSort_field(ctx *Sort_fieldContext) {}

// ExitSort_field is called when production sort_field is exited.
func (s *BaseSQLListener) ExitSort_field(ctx *Sort_fieldContext) {}

// EnterSort_fields is called when production sort_fields is entered.
func (s *BaseSQLListener) EnterSort_fields(ctx *Sort_fieldsContext) {}

// ExitSort_fields is called when production sort_fields is exited.
func (s *BaseSQLListener) ExitSort_fields(ctx *Sort_fieldsContext) {}

// EnterHaving_clause is called when production having_clause is entered.
func (s *BaseSQLListener) EnterHaving_clause(ctx *Having_clauseContext) {}

// ExitHaving_clause is called when production having_clause is exited.
func (s *BaseSQLListener) ExitHaving_clause(ctx *Having_clauseContext) {}

// EnterBool_expr is called when production bool_expr is entered.
func (s *BaseSQLListener) EnterBool_expr(ctx *Bool_exprContext) {}

// ExitBool_expr is called when production bool_expr is exited.
func (s *BaseSQLListener) ExitBool_expr(ctx *Bool_exprContext) {}

// EnterBool_expr_logical_op is called when production bool_expr_logical_op is entered.
func (s *BaseSQLListener) EnterBool_expr_logical_op(ctx *Bool_expr_logical_opContext) {}

// ExitBool_expr_logical_op is called when production bool_expr_logical_op is exited.
func (s *BaseSQLListener) ExitBool_expr_logical_op(ctx *Bool_expr_logical_opContext) {}

// EnterBool_expr_atom is called when production bool_expr_atom is entered.
func (s *BaseSQLListener) EnterBool_expr_atom(ctx *Bool_expr_atomContext) {}

// ExitBool_expr_atom is called when production bool_expr_atom is exited.
func (s *BaseSQLListener) ExitBool_expr_atom(ctx *Bool_expr_atomContext) {}

// EnterBool_expr_binary is called when production bool_expr_binary is entered.
func (s *BaseSQLListener) EnterBool_expr_binary(ctx *Bool_expr_binaryContext) {}

// ExitBool_expr_binary is called when production bool_expr_binary is exited.
func (s *BaseSQLListener) ExitBool_expr_binary(ctx *Bool_expr_binaryContext) {}

// EnterBool_expr_binary_operator is called when production bool_expr_binary_operator is entered.
func (s *BaseSQLListener) EnterBool_expr_binary_operator(ctx *Bool_expr_binary_operatorContext) {}

// ExitBool_expr_binary_operator is called when production bool_expr_binary_operator is exited.
func (s *BaseSQLListener) ExitBool_expr_binary_operator(ctx *Bool_expr_binary_operatorContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSQLListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSQLListener) ExitExpr(ctx *ExprContext) {}

// EnterDuration_lit is called when production duration_lit is entered.
func (s *BaseSQLListener) EnterDuration_lit(ctx *Duration_litContext) {}

// ExitDuration_lit is called when production duration_lit is exited.
func (s *BaseSQLListener) ExitDuration_lit(ctx *Duration_litContext) {}

// EnterInterval_item is called when production interval_item is entered.
func (s *BaseSQLListener) EnterInterval_item(ctx *Interval_itemContext) {}

// ExitInterval_item is called when production interval_item is exited.
func (s *BaseSQLListener) ExitInterval_item(ctx *Interval_itemContext) {}

// EnterExpr_func is called when production expr_func is entered.
func (s *BaseSQLListener) EnterExpr_func(ctx *Expr_funcContext) {}

// ExitExpr_func is called when production expr_func is exited.
func (s *BaseSQLListener) ExitExpr_func(ctx *Expr_funcContext) {}

// EnterExpr_func_params is called when production expr_func_params is entered.
func (s *BaseSQLListener) EnterExpr_func_params(ctx *Expr_func_paramsContext) {}

// ExitExpr_func_params is called when production expr_func_params is exited.
func (s *BaseSQLListener) ExitExpr_func_params(ctx *Expr_func_paramsContext) {}

// EnterFunc_param is called when production func_param is entered.
func (s *BaseSQLListener) EnterFunc_param(ctx *Func_paramContext) {}

// ExitFunc_param is called when production func_param is exited.
func (s *BaseSQLListener) ExitFunc_param(ctx *Func_paramContext) {}

// EnterExpr_atom is called when production expr_atom is entered.
func (s *BaseSQLListener) EnterExpr_atom(ctx *Expr_atomContext) {}

// ExitExpr_atom is called when production expr_atom is exited.
func (s *BaseSQLListener) ExitExpr_atom(ctx *Expr_atomContext) {}

// EnterIdent_filter is called when production ident_filter is entered.
func (s *BaseSQLListener) EnterIdent_filter(ctx *Ident_filterContext) {}

// ExitIdent_filter is called when production ident_filter is exited.
func (s *BaseSQLListener) ExitIdent_filter(ctx *Ident_filterContext) {}

// EnterInt_number is called when production int_number is entered.
func (s *BaseSQLListener) EnterInt_number(ctx *Int_numberContext) {}

// ExitInt_number is called when production int_number is exited.
func (s *BaseSQLListener) ExitInt_number(ctx *Int_numberContext) {}

// EnterDec_number is called when production dec_number is entered.
func (s *BaseSQLListener) EnterDec_number(ctx *Dec_numberContext) {}

// ExitDec_number is called when production dec_number is exited.
func (s *BaseSQLListener) ExitDec_number(ctx *Dec_numberContext) {}

// EnterLimit_clause is called when production limit_clause is entered.
func (s *BaseSQLListener) EnterLimit_clause(ctx *Limit_clauseContext) {}

// ExitLimit_clause is called when production limit_clause is exited.
func (s *BaseSQLListener) ExitLimit_clause(ctx *Limit_clauseContext) {}

// EnterMetric_name is called when production metric_name is entered.
func (s *BaseSQLListener) EnterMetric_name(ctx *Metric_nameContext) {}

// ExitMetric_name is called when production metric_name is exited.
func (s *BaseSQLListener) ExitMetric_name(ctx *Metric_nameContext) {}

// EnterTag_key is called when production tag_key is entered.
func (s *BaseSQLListener) EnterTag_key(ctx *Tag_keyContext) {}

// ExitTag_key is called when production tag_key is exited.
func (s *BaseSQLListener) ExitTag_key(ctx *Tag_keyContext) {}

// EnterTag_value is called when production tag_value is entered.
func (s *BaseSQLListener) EnterTag_value(ctx *Tag_valueContext) {}

// ExitTag_value is called when production tag_value is exited.
func (s *BaseSQLListener) ExitTag_value(ctx *Tag_valueContext) {}

// EnterTag_value_pattern is called when production tag_value_pattern is entered.
func (s *BaseSQLListener) EnterTag_value_pattern(ctx *Tag_value_patternContext) {}

// ExitTag_value_pattern is called when production tag_value_pattern is exited.
func (s *BaseSQLListener) ExitTag_value_pattern(ctx *Tag_value_patternContext) {}

// EnterIdent is called when production ident is entered.
func (s *BaseSQLListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production ident is exited.
func (s *BaseSQLListener) ExitIdent(ctx *IdentContext) {}

// EnterNon_reserved_words is called when production non_reserved_words is entered.
func (s *BaseSQLListener) EnterNon_reserved_words(ctx *Non_reserved_wordsContext) {}

// ExitNon_reserved_words is called when production non_reserved_words is exited.
func (s *BaseSQLListener) ExitNon_reserved_words(ctx *Non_reserved_wordsContext) {}
