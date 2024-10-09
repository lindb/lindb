// Code generated from ./sql/grammar/SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // SQLParser
import "github.com/antlr4-go/antlr/v4"

// BaseSQLParserListener is a complete listener for a parse tree produced by SQLParser.
type BaseSQLParserListener struct{}

var _ SQLParserListener = &BaseSQLParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseSQLParserListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseSQLParserListener) ExitStatement(ctx *StatementContext) {}

// EnterDdlStatement is called when production ddlStatement is entered.
func (s *BaseSQLParserListener) EnterDdlStatement(ctx *DdlStatementContext) {}

// ExitDdlStatement is called when production ddlStatement is exited.
func (s *BaseSQLParserListener) ExitDdlStatement(ctx *DdlStatementContext) {}

// EnterStatementDefault is called when production statementDefault is entered.
func (s *BaseSQLParserListener) EnterStatementDefault(ctx *StatementDefaultContext) {}

// ExitStatementDefault is called when production statementDefault is exited.
func (s *BaseSQLParserListener) ExitStatementDefault(ctx *StatementDefaultContext) {}

// EnterExplain is called when production explain is entered.
func (s *BaseSQLParserListener) EnterExplain(ctx *ExplainContext) {}

// ExitExplain is called when production explain is exited.
func (s *BaseSQLParserListener) ExitExplain(ctx *ExplainContext) {}

// EnterExplainAnalyze is called when production explainAnalyze is entered.
func (s *BaseSQLParserListener) EnterExplainAnalyze(ctx *ExplainAnalyzeContext) {}

// ExitExplainAnalyze is called when production explainAnalyze is exited.
func (s *BaseSQLParserListener) ExitExplainAnalyze(ctx *ExplainAnalyzeContext) {}

// EnterAdminStatement is called when production adminStatement is entered.
func (s *BaseSQLParserListener) EnterAdminStatement(ctx *AdminStatementContext) {}

// ExitAdminStatement is called when production adminStatement is exited.
func (s *BaseSQLParserListener) ExitAdminStatement(ctx *AdminStatementContext) {}

// EnterUtilityStatement is called when production utilityStatement is entered.
func (s *BaseSQLParserListener) EnterUtilityStatement(ctx *UtilityStatementContext) {}

// ExitUtilityStatement is called when production utilityStatement is exited.
func (s *BaseSQLParserListener) ExitUtilityStatement(ctx *UtilityStatementContext) {}

// EnterCreateDatabase is called when production createDatabase is entered.
func (s *BaseSQLParserListener) EnterCreateDatabase(ctx *CreateDatabaseContext) {}

// ExitCreateDatabase is called when production createDatabase is exited.
func (s *BaseSQLParserListener) ExitCreateDatabase(ctx *CreateDatabaseContext) {}

// EnterRollupOptions is called when production rollupOptions is entered.
func (s *BaseSQLParserListener) EnterRollupOptions(ctx *RollupOptionsContext) {}

// ExitRollupOptions is called when production rollupOptions is exited.
func (s *BaseSQLParserListener) ExitRollupOptions(ctx *RollupOptionsContext) {}

// EnterDropDatabase is called when production dropDatabase is entered.
func (s *BaseSQLParserListener) EnterDropDatabase(ctx *DropDatabaseContext) {}

// ExitDropDatabase is called when production dropDatabase is exited.
func (s *BaseSQLParserListener) ExitDropDatabase(ctx *DropDatabaseContext) {}

// EnterCreateBroker is called when production createBroker is entered.
func (s *BaseSQLParserListener) EnterCreateBroker(ctx *CreateBrokerContext) {}

// ExitCreateBroker is called when production createBroker is exited.
func (s *BaseSQLParserListener) ExitCreateBroker(ctx *CreateBrokerContext) {}

// EnterFlushDatabase is called when production flushDatabase is entered.
func (s *BaseSQLParserListener) EnterFlushDatabase(ctx *FlushDatabaseContext) {}

// ExitFlushDatabase is called when production flushDatabase is exited.
func (s *BaseSQLParserListener) ExitFlushDatabase(ctx *FlushDatabaseContext) {}

// EnterCompactDatabase is called when production compactDatabase is entered.
func (s *BaseSQLParserListener) EnterCompactDatabase(ctx *CompactDatabaseContext) {}

// ExitCompactDatabase is called when production compactDatabase is exited.
func (s *BaseSQLParserListener) ExitCompactDatabase(ctx *CompactDatabaseContext) {}

// EnterShowMaster is called when production showMaster is entered.
func (s *BaseSQLParserListener) EnterShowMaster(ctx *ShowMasterContext) {}

// ExitShowMaster is called when production showMaster is exited.
func (s *BaseSQLParserListener) ExitShowMaster(ctx *ShowMasterContext) {}

// EnterShowBrokers is called when production showBrokers is entered.
func (s *BaseSQLParserListener) EnterShowBrokers(ctx *ShowBrokersContext) {}

// ExitShowBrokers is called when production showBrokers is exited.
func (s *BaseSQLParserListener) ExitShowBrokers(ctx *ShowBrokersContext) {}

// EnterShowRequests is called when production showRequests is entered.
func (s *BaseSQLParserListener) EnterShowRequests(ctx *ShowRequestsContext) {}

// ExitShowRequests is called when production showRequests is exited.
func (s *BaseSQLParserListener) ExitShowRequests(ctx *ShowRequestsContext) {}

// EnterShowLimit is called when production showLimit is entered.
func (s *BaseSQLParserListener) EnterShowLimit(ctx *ShowLimitContext) {}

// ExitShowLimit is called when production showLimit is exited.
func (s *BaseSQLParserListener) ExitShowLimit(ctx *ShowLimitContext) {}

// EnterShowMetadataTypes is called when production showMetadataTypes is entered.
func (s *BaseSQLParserListener) EnterShowMetadataTypes(ctx *ShowMetadataTypesContext) {}

// ExitShowMetadataTypes is called when production showMetadataTypes is exited.
func (s *BaseSQLParserListener) ExitShowMetadataTypes(ctx *ShowMetadataTypesContext) {}

// EnterShowMetadatas is called when production showMetadatas is entered.
func (s *BaseSQLParserListener) EnterShowMetadatas(ctx *ShowMetadatasContext) {}

// ExitShowMetadatas is called when production showMetadatas is exited.
func (s *BaseSQLParserListener) ExitShowMetadatas(ctx *ShowMetadatasContext) {}

// EnterShowAlive is called when production showAlive is entered.
func (s *BaseSQLParserListener) EnterShowAlive(ctx *ShowAliveContext) {}

// ExitShowAlive is called when production showAlive is exited.
func (s *BaseSQLParserListener) ExitShowAlive(ctx *ShowAliveContext) {}

// EnterShowReplications is called when production showReplications is entered.
func (s *BaseSQLParserListener) EnterShowReplications(ctx *ShowReplicationsContext) {}

// ExitShowReplications is called when production showReplications is exited.
func (s *BaseSQLParserListener) ExitShowReplications(ctx *ShowReplicationsContext) {}

// EnterShowState is called when production showState is entered.
func (s *BaseSQLParserListener) EnterShowState(ctx *ShowStateContext) {}

// ExitShowState is called when production showState is exited.
func (s *BaseSQLParserListener) ExitShowState(ctx *ShowStateContext) {}

// EnterShowDatabases is called when production showDatabases is entered.
func (s *BaseSQLParserListener) EnterShowDatabases(ctx *ShowDatabasesContext) {}

// ExitShowDatabases is called when production showDatabases is exited.
func (s *BaseSQLParserListener) ExitShowDatabases(ctx *ShowDatabasesContext) {}

// EnterUseStatement is called when production useStatement is entered.
func (s *BaseSQLParserListener) EnterUseStatement(ctx *UseStatementContext) {}

// ExitUseStatement is called when production useStatement is exited.
func (s *BaseSQLParserListener) ExitUseStatement(ctx *UseStatementContext) {}

// EnterShowNamespaces is called when production showNamespaces is entered.
func (s *BaseSQLParserListener) EnterShowNamespaces(ctx *ShowNamespacesContext) {}

// ExitShowNamespaces is called when production showNamespaces is exited.
func (s *BaseSQLParserListener) ExitShowNamespaces(ctx *ShowNamespacesContext) {}

// EnterShowMetrics is called when production showMetrics is entered.
func (s *BaseSQLParserListener) EnterShowMetrics(ctx *ShowMetricsContext) {}

// ExitShowMetrics is called when production showMetrics is exited.
func (s *BaseSQLParserListener) ExitShowMetrics(ctx *ShowMetricsContext) {}

// EnterShowFields is called when production showFields is entered.
func (s *BaseSQLParserListener) EnterShowFields(ctx *ShowFieldsContext) {}

// ExitShowFields is called when production showFields is exited.
func (s *BaseSQLParserListener) ExitShowFields(ctx *ShowFieldsContext) {}

// EnterShowTagKeys is called when production showTagKeys is entered.
func (s *BaseSQLParserListener) EnterShowTagKeys(ctx *ShowTagKeysContext) {}

// ExitShowTagKeys is called when production showTagKeys is exited.
func (s *BaseSQLParserListener) ExitShowTagKeys(ctx *ShowTagKeysContext) {}

// EnterShowTagValues is called when production showTagValues is entered.
func (s *BaseSQLParserListener) EnterShowTagValues(ctx *ShowTagValuesContext) {}

// ExitShowTagValues is called when production showTagValues is exited.
func (s *BaseSQLParserListener) ExitShowTagValues(ctx *ShowTagValuesContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseSQLParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseSQLParserListener) ExitQuery(ctx *QueryContext) {}

// EnterWith is called when production with is entered.
func (s *BaseSQLParserListener) EnterWith(ctx *WithContext) {}

// ExitWith is called when production with is exited.
func (s *BaseSQLParserListener) ExitWith(ctx *WithContext) {}

// EnterNamedQuery is called when production namedQuery is entered.
func (s *BaseSQLParserListener) EnterNamedQuery(ctx *NamedQueryContext) {}

// ExitNamedQuery is called when production namedQuery is exited.
func (s *BaseSQLParserListener) ExitNamedQuery(ctx *NamedQueryContext) {}

// EnterQueryNoWith is called when production queryNoWith is entered.
func (s *BaseSQLParserListener) EnterQueryNoWith(ctx *QueryNoWithContext) {}

// ExitQueryNoWith is called when production queryNoWith is exited.
func (s *BaseSQLParserListener) ExitQueryNoWith(ctx *QueryNoWithContext) {}

// EnterQueryTermDefault is called when production queryTermDefault is entered.
func (s *BaseSQLParserListener) EnterQueryTermDefault(ctx *QueryTermDefaultContext) {}

// ExitQueryTermDefault is called when production queryTermDefault is exited.
func (s *BaseSQLParserListener) ExitQueryTermDefault(ctx *QueryTermDefaultContext) {}

// EnterQueryPrimaryDefault is called when production queryPrimaryDefault is entered.
func (s *BaseSQLParserListener) EnterQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) {}

// ExitQueryPrimaryDefault is called when production queryPrimaryDefault is exited.
func (s *BaseSQLParserListener) ExitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) {}

// EnterSubquery is called when production subquery is entered.
func (s *BaseSQLParserListener) EnterSubquery(ctx *SubqueryContext) {}

// ExitSubquery is called when production subquery is exited.
func (s *BaseSQLParserListener) ExitSubquery(ctx *SubqueryContext) {}

// EnterQuerySpecification is called when production querySpecification is entered.
func (s *BaseSQLParserListener) EnterQuerySpecification(ctx *QuerySpecificationContext) {}

// ExitQuerySpecification is called when production querySpecification is exited.
func (s *BaseSQLParserListener) ExitQuerySpecification(ctx *QuerySpecificationContext) {}

// EnterSelectSingle is called when production selectSingle is entered.
func (s *BaseSQLParserListener) EnterSelectSingle(ctx *SelectSingleContext) {}

// ExitSelectSingle is called when production selectSingle is exited.
func (s *BaseSQLParserListener) ExitSelectSingle(ctx *SelectSingleContext) {}

// EnterSelectAll is called when production selectAll is entered.
func (s *BaseSQLParserListener) EnterSelectAll(ctx *SelectAllContext) {}

// ExitSelectAll is called when production selectAll is exited.
func (s *BaseSQLParserListener) ExitSelectAll(ctx *SelectAllContext) {}

// EnterRelationDefault is called when production relationDefault is entered.
func (s *BaseSQLParserListener) EnterRelationDefault(ctx *RelationDefaultContext) {}

// ExitRelationDefault is called when production relationDefault is exited.
func (s *BaseSQLParserListener) ExitRelationDefault(ctx *RelationDefaultContext) {}

// EnterJoinRelation is called when production joinRelation is entered.
func (s *BaseSQLParserListener) EnterJoinRelation(ctx *JoinRelationContext) {}

// ExitJoinRelation is called when production joinRelation is exited.
func (s *BaseSQLParserListener) ExitJoinRelation(ctx *JoinRelationContext) {}

// EnterJoinType is called when production joinType is entered.
func (s *BaseSQLParserListener) EnterJoinType(ctx *JoinTypeContext) {}

// ExitJoinType is called when production joinType is exited.
func (s *BaseSQLParserListener) ExitJoinType(ctx *JoinTypeContext) {}

// EnterJoinCriteria is called when production joinCriteria is entered.
func (s *BaseSQLParserListener) EnterJoinCriteria(ctx *JoinCriteriaContext) {}

// ExitJoinCriteria is called when production joinCriteria is exited.
func (s *BaseSQLParserListener) ExitJoinCriteria(ctx *JoinCriteriaContext) {}

// EnterAliasedRelation is called when production aliasedRelation is entered.
func (s *BaseSQLParserListener) EnterAliasedRelation(ctx *AliasedRelationContext) {}

// ExitAliasedRelation is called when production aliasedRelation is exited.
func (s *BaseSQLParserListener) ExitAliasedRelation(ctx *AliasedRelationContext) {}

// EnterTableName is called when production tableName is entered.
func (s *BaseSQLParserListener) EnterTableName(ctx *TableNameContext) {}

// ExitTableName is called when production tableName is exited.
func (s *BaseSQLParserListener) ExitTableName(ctx *TableNameContext) {}

// EnterSubQueryRelation is called when production subQueryRelation is entered.
func (s *BaseSQLParserListener) EnterSubQueryRelation(ctx *SubQueryRelationContext) {}

// ExitSubQueryRelation is called when production subQueryRelation is exited.
func (s *BaseSQLParserListener) ExitSubQueryRelation(ctx *SubQueryRelationContext) {}

// EnterGroupBy is called when production groupBy is entered.
func (s *BaseSQLParserListener) EnterGroupBy(ctx *GroupByContext) {}

// ExitGroupBy is called when production groupBy is exited.
func (s *BaseSQLParserListener) ExitGroupBy(ctx *GroupByContext) {}

// EnterSingleGroupingSet is called when production singleGroupingSet is entered.
func (s *BaseSQLParserListener) EnterSingleGroupingSet(ctx *SingleGroupingSetContext) {}

// ExitSingleGroupingSet is called when production singleGroupingSet is exited.
func (s *BaseSQLParserListener) ExitSingleGroupingSet(ctx *SingleGroupingSetContext) {}

// EnterGroupByAllColumns is called when production groupByAllColumns is entered.
func (s *BaseSQLParserListener) EnterGroupByAllColumns(ctx *GroupByAllColumnsContext) {}

// ExitGroupByAllColumns is called when production groupByAllColumns is exited.
func (s *BaseSQLParserListener) ExitGroupByAllColumns(ctx *GroupByAllColumnsContext) {}

// EnterGroupingSet is called when production groupingSet is entered.
func (s *BaseSQLParserListener) EnterGroupingSet(ctx *GroupingSetContext) {}

// ExitGroupingSet is called when production groupingSet is exited.
func (s *BaseSQLParserListener) ExitGroupingSet(ctx *GroupingSetContext) {}

// EnterHaving is called when production having is entered.
func (s *BaseSQLParserListener) EnterHaving(ctx *HavingContext) {}

// ExitHaving is called when production having is exited.
func (s *BaseSQLParserListener) ExitHaving(ctx *HavingContext) {}

// EnterOrderBy is called when production orderBy is entered.
func (s *BaseSQLParserListener) EnterOrderBy(ctx *OrderByContext) {}

// ExitOrderBy is called when production orderBy is exited.
func (s *BaseSQLParserListener) ExitOrderBy(ctx *OrderByContext) {}

// EnterSortItem is called when production sortItem is entered.
func (s *BaseSQLParserListener) EnterSortItem(ctx *SortItemContext) {}

// ExitSortItem is called when production sortItem is exited.
func (s *BaseSQLParserListener) ExitSortItem(ctx *SortItemContext) {}

// EnterLimitRowCount is called when production limitRowCount is entered.
func (s *BaseSQLParserListener) EnterLimitRowCount(ctx *LimitRowCountContext) {}

// ExitLimitRowCount is called when production limitRowCount is exited.
func (s *BaseSQLParserListener) ExitLimitRowCount(ctx *LimitRowCountContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseSQLParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseSQLParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterLogicalNot is called when production logicalNot is entered.
func (s *BaseSQLParserListener) EnterLogicalNot(ctx *LogicalNotContext) {}

// ExitLogicalNot is called when production logicalNot is exited.
func (s *BaseSQLParserListener) ExitLogicalNot(ctx *LogicalNotContext) {}

// EnterPredicatedExpression is called when production predicatedExpression is entered.
func (s *BaseSQLParserListener) EnterPredicatedExpression(ctx *PredicatedExpressionContext) {}

// ExitPredicatedExpression is called when production predicatedExpression is exited.
func (s *BaseSQLParserListener) ExitPredicatedExpression(ctx *PredicatedExpressionContext) {}

// EnterOr is called when production or is entered.
func (s *BaseSQLParserListener) EnterOr(ctx *OrContext) {}

// ExitOr is called when production or is exited.
func (s *BaseSQLParserListener) ExitOr(ctx *OrContext) {}

// EnterAnd is called when production and is entered.
func (s *BaseSQLParserListener) EnterAnd(ctx *AndContext) {}

// ExitAnd is called when production and is exited.
func (s *BaseSQLParserListener) ExitAnd(ctx *AndContext) {}

// EnterValueExpressionDefault is called when production valueExpressionDefault is entered.
func (s *BaseSQLParserListener) EnterValueExpressionDefault(ctx *ValueExpressionDefaultContext) {}

// ExitValueExpressionDefault is called when production valueExpressionDefault is exited.
func (s *BaseSQLParserListener) ExitValueExpressionDefault(ctx *ValueExpressionDefaultContext) {}

// EnterArithmeticBinary is called when production arithmeticBinary is entered.
func (s *BaseSQLParserListener) EnterArithmeticBinary(ctx *ArithmeticBinaryContext) {}

// ExitArithmeticBinary is called when production arithmeticBinary is exited.
func (s *BaseSQLParserListener) ExitArithmeticBinary(ctx *ArithmeticBinaryContext) {}

// EnterDereference is called when production dereference is entered.
func (s *BaseSQLParserListener) EnterDereference(ctx *DereferenceContext) {}

// ExitDereference is called when production dereference is exited.
func (s *BaseSQLParserListener) ExitDereference(ctx *DereferenceContext) {}

// EnterColumnReference is called when production columnReference is entered.
func (s *BaseSQLParserListener) EnterColumnReference(ctx *ColumnReferenceContext) {}

// ExitColumnReference is called when production columnReference is exited.
func (s *BaseSQLParserListener) ExitColumnReference(ctx *ColumnReferenceContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseSQLParserListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseSQLParserListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseSQLParserListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseSQLParserListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterParenExpression is called when production parenExpression is entered.
func (s *BaseSQLParserListener) EnterParenExpression(ctx *ParenExpressionContext) {}

// ExitParenExpression is called when production parenExpression is exited.
func (s *BaseSQLParserListener) ExitParenExpression(ctx *ParenExpressionContext) {}

// EnterNumericLiteral is called when production numericLiteral is entered.
func (s *BaseSQLParserListener) EnterNumericLiteral(ctx *NumericLiteralContext) {}

// ExitNumericLiteral is called when production numericLiteral is exited.
func (s *BaseSQLParserListener) ExitNumericLiteral(ctx *NumericLiteralContext) {}

// EnterBooleanLiteral is called when production booleanLiteral is entered.
func (s *BaseSQLParserListener) EnterBooleanLiteral(ctx *BooleanLiteralContext) {}

// ExitBooleanLiteral is called when production booleanLiteral is exited.
func (s *BaseSQLParserListener) ExitBooleanLiteral(ctx *BooleanLiteralContext) {}

// EnterBinaryComparisonPredicate is called when production binaryComparisonPredicate is entered.
func (s *BaseSQLParserListener) EnterBinaryComparisonPredicate(ctx *BinaryComparisonPredicateContext) {
}

// ExitBinaryComparisonPredicate is called when production binaryComparisonPredicate is exited.
func (s *BaseSQLParserListener) ExitBinaryComparisonPredicate(ctx *BinaryComparisonPredicateContext) {
}

// EnterInPredicate is called when production inPredicate is entered.
func (s *BaseSQLParserListener) EnterInPredicate(ctx *InPredicateContext) {}

// ExitInPredicate is called when production inPredicate is exited.
func (s *BaseSQLParserListener) ExitInPredicate(ctx *InPredicateContext) {}

// EnterLikePredicate is called when production likePredicate is entered.
func (s *BaseSQLParserListener) EnterLikePredicate(ctx *LikePredicateContext) {}

// ExitLikePredicate is called when production likePredicate is exited.
func (s *BaseSQLParserListener) ExitLikePredicate(ctx *LikePredicateContext) {}

// EnterRegexpPredicate is called when production regexpPredicate is entered.
func (s *BaseSQLParserListener) EnterRegexpPredicate(ctx *RegexpPredicateContext) {}

// ExitRegexpPredicate is called when production regexpPredicate is exited.
func (s *BaseSQLParserListener) ExitRegexpPredicate(ctx *RegexpPredicateContext) {}

// EnterValueExpressionPredicate is called when production valueExpressionPredicate is entered.
func (s *BaseSQLParserListener) EnterValueExpressionPredicate(ctx *ValueExpressionPredicateContext) {}

// ExitValueExpressionPredicate is called when production valueExpressionPredicate is exited.
func (s *BaseSQLParserListener) ExitValueExpressionPredicate(ctx *ValueExpressionPredicateContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BaseSQLParserListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BaseSQLParserListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseSQLParserListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseSQLParserListener) ExitFilter(ctx *FilterContext) {}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BaseSQLParserListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BaseSQLParserListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterProperties is called when production properties is entered.
func (s *BaseSQLParserListener) EnterProperties(ctx *PropertiesContext) {}

// ExitProperties is called when production properties is exited.
func (s *BaseSQLParserListener) ExitProperties(ctx *PropertiesContext) {}

// EnterPropertyAssignments is called when production propertyAssignments is entered.
func (s *BaseSQLParserListener) EnterPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// ExitPropertyAssignments is called when production propertyAssignments is exited.
func (s *BaseSQLParserListener) ExitPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// EnterProperty is called when production property is entered.
func (s *BaseSQLParserListener) EnterProperty(ctx *PropertyContext) {}

// ExitProperty is called when production property is exited.
func (s *BaseSQLParserListener) ExitProperty(ctx *PropertyContext) {}

// EnterDefaultPropertyValue is called when production defaultPropertyValue is entered.
func (s *BaseSQLParserListener) EnterDefaultPropertyValue(ctx *DefaultPropertyValueContext) {}

// ExitDefaultPropertyValue is called when production defaultPropertyValue is exited.
func (s *BaseSQLParserListener) ExitDefaultPropertyValue(ctx *DefaultPropertyValueContext) {}

// EnterNonDefaultPropertyValue is called when production nonDefaultPropertyValue is entered.
func (s *BaseSQLParserListener) EnterNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) {}

// ExitNonDefaultPropertyValue is called when production nonDefaultPropertyValue is exited.
func (s *BaseSQLParserListener) ExitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) {}

// EnterBooleanValue is called when production booleanValue is entered.
func (s *BaseSQLParserListener) EnterBooleanValue(ctx *BooleanValueContext) {}

// ExitBooleanValue is called when production booleanValue is exited.
func (s *BaseSQLParserListener) ExitBooleanValue(ctx *BooleanValueContext) {}

// EnterBasicStringLiteral is called when production basicStringLiteral is entered.
func (s *BaseSQLParserListener) EnterBasicStringLiteral(ctx *BasicStringLiteralContext) {}

// ExitBasicStringLiteral is called when production basicStringLiteral is exited.
func (s *BaseSQLParserListener) ExitBasicStringLiteral(ctx *BasicStringLiteralContext) {}

// EnterUnquotedIdentifier is called when production unquotedIdentifier is entered.
func (s *BaseSQLParserListener) EnterUnquotedIdentifier(ctx *UnquotedIdentifierContext) {}

// ExitUnquotedIdentifier is called when production unquotedIdentifier is exited.
func (s *BaseSQLParserListener) ExitUnquotedIdentifier(ctx *UnquotedIdentifierContext) {}

// EnterQuotedIdentifier is called when production quotedIdentifier is entered.
func (s *BaseSQLParserListener) EnterQuotedIdentifier(ctx *QuotedIdentifierContext) {}

// ExitQuotedIdentifier is called when production quotedIdentifier is exited.
func (s *BaseSQLParserListener) ExitQuotedIdentifier(ctx *QuotedIdentifierContext) {}

// EnterBackQuotedIdentifier is called when production backQuotedIdentifier is entered.
func (s *BaseSQLParserListener) EnterBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) {}

// ExitBackQuotedIdentifier is called when production backQuotedIdentifier is exited.
func (s *BaseSQLParserListener) ExitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) {}

// EnterDigitIdentifier is called when production digitIdentifier is entered.
func (s *BaseSQLParserListener) EnterDigitIdentifier(ctx *DigitIdentifierContext) {}

// ExitDigitIdentifier is called when production digitIdentifier is exited.
func (s *BaseSQLParserListener) ExitDigitIdentifier(ctx *DigitIdentifierContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseSQLParserListener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseSQLParserListener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterDoubleLiteral is called when production doubleLiteral is entered.
func (s *BaseSQLParserListener) EnterDoubleLiteral(ctx *DoubleLiteralContext) {}

// ExitDoubleLiteral is called when production doubleLiteral is exited.
func (s *BaseSQLParserListener) ExitDoubleLiteral(ctx *DoubleLiteralContext) {}

// EnterIntegerLiteral is called when production integerLiteral is entered.
func (s *BaseSQLParserListener) EnterIntegerLiteral(ctx *IntegerLiteralContext) {}

// ExitIntegerLiteral is called when production integerLiteral is exited.
func (s *BaseSQLParserListener) ExitIntegerLiteral(ctx *IntegerLiteralContext) {}

// EnterNonReserved is called when production nonReserved is entered.
func (s *BaseSQLParserListener) EnterNonReserved(ctx *NonReservedContext) {}

// ExitNonReserved is called when production nonReserved is exited.
func (s *BaseSQLParserListener) ExitNonReserved(ctx *NonReservedContext) {}
