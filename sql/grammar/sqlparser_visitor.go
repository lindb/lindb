// Code generated from ./sql/grammar/SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // SQLParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by SQLParser.
type SQLParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SQLParser#ddlStatement.
	VisitDdlStatement(ctx *DdlStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#statementDefault.
	VisitStatementDefault(ctx *StatementDefaultContext) interface{}

	// Visit a parse tree produced by SQLParser#explain.
	VisitExplain(ctx *ExplainContext) interface{}

	// Visit a parse tree produced by SQLParser#explainAnalyze.
	VisitExplainAnalyze(ctx *ExplainAnalyzeContext) interface{}

	// Visit a parse tree produced by SQLParser#adminStatement.
	VisitAdminStatement(ctx *AdminStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#utilityStatement.
	VisitUtilityStatement(ctx *UtilityStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#explainType.
	VisitExplainType(ctx *ExplainTypeContext) interface{}

	// Visit a parse tree produced by SQLParser#createDatabase.
	VisitCreateDatabase(ctx *CreateDatabaseContext) interface{}

	// Visit a parse tree produced by SQLParser#engineOption.
	VisitEngineOption(ctx *EngineOptionContext) interface{}

	// Visit a parse tree produced by SQLParser#rollupOptions.
	VisitRollupOptions(ctx *RollupOptionsContext) interface{}

	// Visit a parse tree produced by SQLParser#dropDatabase.
	VisitDropDatabase(ctx *DropDatabaseContext) interface{}

	// Visit a parse tree produced by SQLParser#createBroker.
	VisitCreateBroker(ctx *CreateBrokerContext) interface{}

	// Visit a parse tree produced by SQLParser#flushDatabase.
	VisitFlushDatabase(ctx *FlushDatabaseContext) interface{}

	// Visit a parse tree produced by SQLParser#compactDatabase.
	VisitCompactDatabase(ctx *CompactDatabaseContext) interface{}

	// Visit a parse tree produced by SQLParser#showMaster.
	VisitShowMaster(ctx *ShowMasterContext) interface{}

	// Visit a parse tree produced by SQLParser#showBrokers.
	VisitShowBrokers(ctx *ShowBrokersContext) interface{}

	// Visit a parse tree produced by SQLParser#showRequests.
	VisitShowRequests(ctx *ShowRequestsContext) interface{}

	// Visit a parse tree produced by SQLParser#showLimit.
	VisitShowLimit(ctx *ShowLimitContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetadataTypes.
	VisitShowMetadataTypes(ctx *ShowMetadataTypesContext) interface{}

	// Visit a parse tree produced by SQLParser#showMetadatas.
	VisitShowMetadatas(ctx *ShowMetadatasContext) interface{}

	// Visit a parse tree produced by SQLParser#showReplications.
	VisitShowReplications(ctx *ShowReplicationsContext) interface{}

	// Visit a parse tree produced by SQLParser#showState.
	VisitShowState(ctx *ShowStateContext) interface{}

	// Visit a parse tree produced by SQLParser#showDatabases.
	VisitShowDatabases(ctx *ShowDatabasesContext) interface{}

	// Visit a parse tree produced by SQLParser#showNamespaces.
	VisitShowNamespaces(ctx *ShowNamespacesContext) interface{}

	// Visit a parse tree produced by SQLParser#showTableNames.
	VisitShowTableNames(ctx *ShowTableNamesContext) interface{}

	// Visit a parse tree produced by SQLParser#showColumns.
	VisitShowColumns(ctx *ShowColumnsContext) interface{}

	// Visit a parse tree produced by SQLParser#useStatement.
	VisitUseStatement(ctx *UseStatementContext) interface{}

	// Visit a parse tree produced by SQLParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by SQLParser#with.
	VisitWith(ctx *WithContext) interface{}

	// Visit a parse tree produced by SQLParser#namedQuery.
	VisitNamedQuery(ctx *NamedQueryContext) interface{}

	// Visit a parse tree produced by SQLParser#queryNoWith.
	VisitQueryNoWith(ctx *QueryNoWithContext) interface{}

	// Visit a parse tree produced by SQLParser#queryTermDefault.
	VisitQueryTermDefault(ctx *QueryTermDefaultContext) interface{}

	// Visit a parse tree produced by SQLParser#queryPrimaryDefault.
	VisitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) interface{}

	// Visit a parse tree produced by SQLParser#subquery.
	VisitSubquery(ctx *SubqueryContext) interface{}

	// Visit a parse tree produced by SQLParser#querySpecification.
	VisitQuerySpecification(ctx *QuerySpecificationContext) interface{}

	// Visit a parse tree produced by SQLParser#selectSingle.
	VisitSelectSingle(ctx *SelectSingleContext) interface{}

	// Visit a parse tree produced by SQLParser#selectAll.
	VisitSelectAll(ctx *SelectAllContext) interface{}

	// Visit a parse tree produced by SQLParser#relationDefault.
	VisitRelationDefault(ctx *RelationDefaultContext) interface{}

	// Visit a parse tree produced by SQLParser#joinRelation.
	VisitJoinRelation(ctx *JoinRelationContext) interface{}

	// Visit a parse tree produced by SQLParser#joinType.
	VisitJoinType(ctx *JoinTypeContext) interface{}

	// Visit a parse tree produced by SQLParser#joinCriteria.
	VisitJoinCriteria(ctx *JoinCriteriaContext) interface{}

	// Visit a parse tree produced by SQLParser#aliasedRelation.
	VisitAliasedRelation(ctx *AliasedRelationContext) interface{}

	// Visit a parse tree produced by SQLParser#tableName.
	VisitTableName(ctx *TableNameContext) interface{}

	// Visit a parse tree produced by SQLParser#subQueryRelation.
	VisitSubQueryRelation(ctx *SubQueryRelationContext) interface{}

	// Visit a parse tree produced by SQLParser#groupBy.
	VisitGroupBy(ctx *GroupByContext) interface{}

	// Visit a parse tree produced by SQLParser#singleGroupingSet.
	VisitSingleGroupingSet(ctx *SingleGroupingSetContext) interface{}

	// Visit a parse tree produced by SQLParser#groupByAllColumns.
	VisitGroupByAllColumns(ctx *GroupByAllColumnsContext) interface{}

	// Visit a parse tree produced by SQLParser#groupingSet.
	VisitGroupingSet(ctx *GroupingSetContext) interface{}

	// Visit a parse tree produced by SQLParser#having.
	VisitHaving(ctx *HavingContext) interface{}

	// Visit a parse tree produced by SQLParser#orderBy.
	VisitOrderBy(ctx *OrderByContext) interface{}

	// Visit a parse tree produced by SQLParser#sortItem.
	VisitSortItem(ctx *SortItemContext) interface{}

	// Visit a parse tree produced by SQLParser#limitRowCount.
	VisitLimitRowCount(ctx *LimitRowCountContext) interface{}

	// Visit a parse tree produced by SQLParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by SQLParser#logicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by SQLParser#predicatedExpression.
	VisitPredicatedExpression(ctx *PredicatedExpressionContext) interface{}

	// Visit a parse tree produced by SQLParser#or.
	VisitOr(ctx *OrContext) interface{}

	// Visit a parse tree produced by SQLParser#and.
	VisitAnd(ctx *AndContext) interface{}

	// Visit a parse tree produced by SQLParser#valueExpressionDefault.
	VisitValueExpressionDefault(ctx *ValueExpressionDefaultContext) interface{}

	// Visit a parse tree produced by SQLParser#arithmeticBinary.
	VisitArithmeticBinary(ctx *ArithmeticBinaryContext) interface{}

	// Visit a parse tree produced by SQLParser#dereference.
	VisitDereference(ctx *DereferenceContext) interface{}

	// Visit a parse tree produced by SQLParser#columnReference.
	VisitColumnReference(ctx *ColumnReferenceContext) interface{}

	// Visit a parse tree produced by SQLParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by SQLParser#parenExpression.
	VisitParenExpression(ctx *ParenExpressionContext) interface{}

	// Visit a parse tree produced by SQLParser#numericLiteral.
	VisitNumericLiteral(ctx *NumericLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalLiteral.
	VisitIntervalLiteral(ctx *IntervalLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#timestampPredicate.
	VisitTimestampPredicate(ctx *TimestampPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#binaryComparisonPredicate.
	VisitBinaryComparisonPredicate(ctx *BinaryComparisonPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#betweenPredicate.
	VisitBetweenPredicate(ctx *BetweenPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#inPredicate.
	VisitInPredicate(ctx *InPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#likePredicate.
	VisitLikePredicate(ctx *LikePredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#regexpPredicate.
	VisitRegexpPredicate(ctx *RegexpPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#valueExpressionPredicate.
	VisitValueExpressionPredicate(ctx *ValueExpressionPredicateContext) interface{}

	// Visit a parse tree produced by SQLParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by SQLParser#qualifiedName.
	VisitQualifiedName(ctx *QualifiedNameContext) interface{}

	// Visit a parse tree produced by SQLParser#properties.
	VisitProperties(ctx *PropertiesContext) interface{}

	// Visit a parse tree produced by SQLParser#propertyAssignments.
	VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{}

	// Visit a parse tree produced by SQLParser#property.
	VisitProperty(ctx *PropertyContext) interface{}

	// Visit a parse tree produced by SQLParser#defaultPropertyValue.
	VisitDefaultPropertyValue(ctx *DefaultPropertyValueContext) interface{}

	// Visit a parse tree produced by SQLParser#nonDefaultPropertyValue.
	VisitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) interface{}

	// Visit a parse tree produced by SQLParser#booleanValue.
	VisitBooleanValue(ctx *BooleanValueContext) interface{}

	// Visit a parse tree produced by SQLParser#basicStringLiteral.
	VisitBasicStringLiteral(ctx *BasicStringLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#unquotedIdentifier.
	VisitUnquotedIdentifier(ctx *UnquotedIdentifierContext) interface{}

	// Visit a parse tree produced by SQLParser#quotedIdentifier.
	VisitQuotedIdentifier(ctx *QuotedIdentifierContext) interface{}

	// Visit a parse tree produced by SQLParser#backQuotedIdentifier.
	VisitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) interface{}

	// Visit a parse tree produced by SQLParser#digitIdentifier.
	VisitDigitIdentifier(ctx *DigitIdentifierContext) interface{}

	// Visit a parse tree produced by SQLParser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#doubleLiteral.
	VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by SQLParser#interval.
	VisitInterval(ctx *IntervalContext) interface{}

	// Visit a parse tree produced by SQLParser#intervalUnit.
	VisitIntervalUnit(ctx *IntervalUnitContext) interface{}

	// Visit a parse tree produced by SQLParser#nonReserved.
	VisitNonReserved(ctx *NonReservedContext) interface{}
}
