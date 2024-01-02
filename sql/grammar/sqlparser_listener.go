// Code generated from ./sql/grammar/SQLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package grammar // SQLParser
import "github.com/antlr4-go/antlr/v4"

// SQLParserListener is a complete listener for a parse tree produced by SQLParser.
type SQLParserListener interface {
	antlr.ParseTreeListener

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterDdlStatement is called when entering the ddlStatement production.
	EnterDdlStatement(c *DdlStatementContext)

	// EnterDmlStatement is called when entering the dmlStatement production.
	EnterDmlStatement(c *DmlStatementContext)

	// EnterAdminStatement is called when entering the adminStatement production.
	EnterAdminStatement(c *AdminStatementContext)

	// EnterUtilityStatement is called when entering the utilityStatement production.
	EnterUtilityStatement(c *UtilityStatementContext)

	// EnterCreateDatabase is called when entering the createDatabase production.
	EnterCreateDatabase(c *CreateDatabaseContext)

	// EnterRollupOptions is called when entering the rollupOptions production.
	EnterRollupOptions(c *RollupOptionsContext)

	// EnterDropDatabase is called when entering the dropDatabase production.
	EnterDropDatabase(c *DropDatabaseContext)

	// EnterCreateBroker is called when entering the createBroker production.
	EnterCreateBroker(c *CreateBrokerContext)

	// EnterFlushDatabase is called when entering the flushDatabase production.
	EnterFlushDatabase(c *FlushDatabaseContext)

	// EnterCompactDatabase is called when entering the compactDatabase production.
	EnterCompactDatabase(c *CompactDatabaseContext)

	// EnterShowMaster is called when entering the showMaster production.
	EnterShowMaster(c *ShowMasterContext)

	// EnterShowBrokers is called when entering the showBrokers production.
	EnterShowBrokers(c *ShowBrokersContext)

	// EnterShowRequests is called when entering the showRequests production.
	EnterShowRequests(c *ShowRequestsContext)

	// EnterShowLimit is called when entering the showLimit production.
	EnterShowLimit(c *ShowLimitContext)

	// EnterShowMetadataTypes is called when entering the showMetadataTypes production.
	EnterShowMetadataTypes(c *ShowMetadataTypesContext)

	// EnterShowMetadatas is called when entering the showMetadatas production.
	EnterShowMetadatas(c *ShowMetadatasContext)

	// EnterShowAlive is called when entering the showAlive production.
	EnterShowAlive(c *ShowAliveContext)

	// EnterShowReplications is called when entering the showReplications production.
	EnterShowReplications(c *ShowReplicationsContext)

	// EnterShowState is called when entering the showState production.
	EnterShowState(c *ShowStateContext)

	// EnterShowDatabases is called when entering the showDatabases production.
	EnterShowDatabases(c *ShowDatabasesContext)

	// EnterUseStatement is called when entering the useStatement production.
	EnterUseStatement(c *UseStatementContext)

	// EnterShowNamespaces is called when entering the showNamespaces production.
	EnterShowNamespaces(c *ShowNamespacesContext)

	// EnterShowMetrics is called when entering the showMetrics production.
	EnterShowMetrics(c *ShowMetricsContext)

	// EnterShowFields is called when entering the showFields production.
	EnterShowFields(c *ShowFieldsContext)

	// EnterShowTagKeys is called when entering the showTagKeys production.
	EnterShowTagKeys(c *ShowTagKeysContext)

	// EnterShowTagValues is called when entering the showTagValues production.
	EnterShowTagValues(c *ShowTagValuesContext)

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterWith is called when entering the with production.
	EnterWith(c *WithContext)

	// EnterNamedQuery is called when entering the namedQuery production.
	EnterNamedQuery(c *NamedQueryContext)

	// EnterQueryNoWith is called when entering the queryNoWith production.
	EnterQueryNoWith(c *QueryNoWithContext)

	// EnterQueryTermDefault is called when entering the queryTermDefault production.
	EnterQueryTermDefault(c *QueryTermDefaultContext)

	// EnterQueryPrimaryDefault is called when entering the queryPrimaryDefault production.
	EnterQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// EnterSubquery is called when entering the subquery production.
	EnterSubquery(c *SubqueryContext)

	// EnterQuerySpecification is called when entering the querySpecification production.
	EnterQuerySpecification(c *QuerySpecificationContext)

	// EnterSelectSingle is called when entering the selectSingle production.
	EnterSelectSingle(c *SelectSingleContext)

	// EnterSelectAll is called when entering the selectAll production.
	EnterSelectAll(c *SelectAllContext)

	// EnterRelationDefault is called when entering the relationDefault production.
	EnterRelationDefault(c *RelationDefaultContext)

	// EnterJoinRelation is called when entering the joinRelation production.
	EnterJoinRelation(c *JoinRelationContext)

	// EnterJoinType is called when entering the joinType production.
	EnterJoinType(c *JoinTypeContext)

	// EnterJoinCriteria is called when entering the joinCriteria production.
	EnterJoinCriteria(c *JoinCriteriaContext)

	// EnterAliasedRelation is called when entering the aliasedRelation production.
	EnterAliasedRelation(c *AliasedRelationContext)

	// EnterTableName is called when entering the tableName production.
	EnterTableName(c *TableNameContext)

	// EnterSubQueryRelation is called when entering the subQueryRelation production.
	EnterSubQueryRelation(c *SubQueryRelationContext)

	// EnterGroupBy is called when entering the groupBy production.
	EnterGroupBy(c *GroupByContext)

	// EnterSingleGroupingSet is called when entering the singleGroupingSet production.
	EnterSingleGroupingSet(c *SingleGroupingSetContext)

	// EnterGroupByAllColumns is called when entering the groupByAllColumns production.
	EnterGroupByAllColumns(c *GroupByAllColumnsContext)

	// EnterGroupingSet is called when entering the groupingSet production.
	EnterGroupingSet(c *GroupingSetContext)

	// EnterHaving is called when entering the having production.
	EnterHaving(c *HavingContext)

	// EnterOrderBy is called when entering the orderBy production.
	EnterOrderBy(c *OrderByContext)

	// EnterSortItem is called when entering the sortItem production.
	EnterSortItem(c *SortItemContext)

	// EnterLimitRowCount is called when entering the limitRowCount production.
	EnterLimitRowCount(c *LimitRowCountContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterLogicalNot is called when entering the logicalNot production.
	EnterLogicalNot(c *LogicalNotContext)

	// EnterPredicatedExpression is called when entering the predicatedExpression production.
	EnterPredicatedExpression(c *PredicatedExpressionContext)

	// EnterOr is called when entering the or production.
	EnterOr(c *OrContext)

	// EnterAnd is called when entering the and production.
	EnterAnd(c *AndContext)

	// EnterValueExpressionDefault is called when entering the valueExpressionDefault production.
	EnterValueExpressionDefault(c *ValueExpressionDefaultContext)

	// EnterArithmeticBinary is called when entering the arithmeticBinary production.
	EnterArithmeticBinary(c *ArithmeticBinaryContext)

	// EnterDereference is called when entering the dereference production.
	EnterDereference(c *DereferenceContext)

	// EnterColumnReference is called when entering the columnReference production.
	EnterColumnReference(c *ColumnReferenceContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterParenExpression is called when entering the parenExpression production.
	EnterParenExpression(c *ParenExpressionContext)

	// EnterNumericLiteral is called when entering the numericLiteral production.
	EnterNumericLiteral(c *NumericLiteralContext)

	// EnterBooleanLiteral is called when entering the booleanLiteral production.
	EnterBooleanLiteral(c *BooleanLiteralContext)

	// EnterBinaryComparisonPredicate is called when entering the binaryComparisonPredicate production.
	EnterBinaryComparisonPredicate(c *BinaryComparisonPredicateContext)

	// EnterInPredicate is called when entering the inPredicate production.
	EnterInPredicate(c *InPredicateContext)

	// EnterLikePredicate is called when entering the likePredicate production.
	EnterLikePredicate(c *LikePredicateContext)

	// EnterRegexpPredicate is called when entering the regexpPredicate production.
	EnterRegexpPredicate(c *RegexpPredicateContext)

	// EnterValueExpressionPredicate is called when entering the valueExpressionPredicate production.
	EnterValueExpressionPredicate(c *ValueExpressionPredicateContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterProperties is called when entering the properties production.
	EnterProperties(c *PropertiesContext)

	// EnterPropertyAssignments is called when entering the propertyAssignments production.
	EnterPropertyAssignments(c *PropertyAssignmentsContext)

	// EnterProperty is called when entering the property production.
	EnterProperty(c *PropertyContext)

	// EnterDefaultPropertyValue is called when entering the defaultPropertyValue production.
	EnterDefaultPropertyValue(c *DefaultPropertyValueContext)

	// EnterNonDefaultPropertyValue is called when entering the nonDefaultPropertyValue production.
	EnterNonDefaultPropertyValue(c *NonDefaultPropertyValueContext)

	// EnterBooleanValue is called when entering the booleanValue production.
	EnterBooleanValue(c *BooleanValueContext)

	// EnterBasicStringLiteral is called when entering the basicStringLiteral production.
	EnterBasicStringLiteral(c *BasicStringLiteralContext)

	// EnterUnquotedIdentifier is called when entering the unquotedIdentifier production.
	EnterUnquotedIdentifier(c *UnquotedIdentifierContext)

	// EnterQuotedIdentifier is called when entering the quotedIdentifier production.
	EnterQuotedIdentifier(c *QuotedIdentifierContext)

	// EnterBackQuotedIdentifier is called when entering the backQuotedIdentifier production.
	EnterBackQuotedIdentifier(c *BackQuotedIdentifierContext)

	// EnterDigitIdentifier is called when entering the digitIdentifier production.
	EnterDigitIdentifier(c *DigitIdentifierContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterDoubleLiteral is called when entering the doubleLiteral production.
	EnterDoubleLiteral(c *DoubleLiteralContext)

	// EnterIntegerLiteral is called when entering the integerLiteral production.
	EnterIntegerLiteral(c *IntegerLiteralContext)

	// EnterNonReserved is called when entering the nonReserved production.
	EnterNonReserved(c *NonReservedContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitDdlStatement is called when exiting the ddlStatement production.
	ExitDdlStatement(c *DdlStatementContext)

	// ExitDmlStatement is called when exiting the dmlStatement production.
	ExitDmlStatement(c *DmlStatementContext)

	// ExitAdminStatement is called when exiting the adminStatement production.
	ExitAdminStatement(c *AdminStatementContext)

	// ExitUtilityStatement is called when exiting the utilityStatement production.
	ExitUtilityStatement(c *UtilityStatementContext)

	// ExitCreateDatabase is called when exiting the createDatabase production.
	ExitCreateDatabase(c *CreateDatabaseContext)

	// ExitRollupOptions is called when exiting the rollupOptions production.
	ExitRollupOptions(c *RollupOptionsContext)

	// ExitDropDatabase is called when exiting the dropDatabase production.
	ExitDropDatabase(c *DropDatabaseContext)

	// ExitCreateBroker is called when exiting the createBroker production.
	ExitCreateBroker(c *CreateBrokerContext)

	// ExitFlushDatabase is called when exiting the flushDatabase production.
	ExitFlushDatabase(c *FlushDatabaseContext)

	// ExitCompactDatabase is called when exiting the compactDatabase production.
	ExitCompactDatabase(c *CompactDatabaseContext)

	// ExitShowMaster is called when exiting the showMaster production.
	ExitShowMaster(c *ShowMasterContext)

	// ExitShowBrokers is called when exiting the showBrokers production.
	ExitShowBrokers(c *ShowBrokersContext)

	// ExitShowRequests is called when exiting the showRequests production.
	ExitShowRequests(c *ShowRequestsContext)

	// ExitShowLimit is called when exiting the showLimit production.
	ExitShowLimit(c *ShowLimitContext)

	// ExitShowMetadataTypes is called when exiting the showMetadataTypes production.
	ExitShowMetadataTypes(c *ShowMetadataTypesContext)

	// ExitShowMetadatas is called when exiting the showMetadatas production.
	ExitShowMetadatas(c *ShowMetadatasContext)

	// ExitShowAlive is called when exiting the showAlive production.
	ExitShowAlive(c *ShowAliveContext)

	// ExitShowReplications is called when exiting the showReplications production.
	ExitShowReplications(c *ShowReplicationsContext)

	// ExitShowState is called when exiting the showState production.
	ExitShowState(c *ShowStateContext)

	// ExitShowDatabases is called when exiting the showDatabases production.
	ExitShowDatabases(c *ShowDatabasesContext)

	// ExitUseStatement is called when exiting the useStatement production.
	ExitUseStatement(c *UseStatementContext)

	// ExitShowNamespaces is called when exiting the showNamespaces production.
	ExitShowNamespaces(c *ShowNamespacesContext)

	// ExitShowMetrics is called when exiting the showMetrics production.
	ExitShowMetrics(c *ShowMetricsContext)

	// ExitShowFields is called when exiting the showFields production.
	ExitShowFields(c *ShowFieldsContext)

	// ExitShowTagKeys is called when exiting the showTagKeys production.
	ExitShowTagKeys(c *ShowTagKeysContext)

	// ExitShowTagValues is called when exiting the showTagValues production.
	ExitShowTagValues(c *ShowTagValuesContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitWith is called when exiting the with production.
	ExitWith(c *WithContext)

	// ExitNamedQuery is called when exiting the namedQuery production.
	ExitNamedQuery(c *NamedQueryContext)

	// ExitQueryNoWith is called when exiting the queryNoWith production.
	ExitQueryNoWith(c *QueryNoWithContext)

	// ExitQueryTermDefault is called when exiting the queryTermDefault production.
	ExitQueryTermDefault(c *QueryTermDefaultContext)

	// ExitQueryPrimaryDefault is called when exiting the queryPrimaryDefault production.
	ExitQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// ExitSubquery is called when exiting the subquery production.
	ExitSubquery(c *SubqueryContext)

	// ExitQuerySpecification is called when exiting the querySpecification production.
	ExitQuerySpecification(c *QuerySpecificationContext)

	// ExitSelectSingle is called when exiting the selectSingle production.
	ExitSelectSingle(c *SelectSingleContext)

	// ExitSelectAll is called when exiting the selectAll production.
	ExitSelectAll(c *SelectAllContext)

	// ExitRelationDefault is called when exiting the relationDefault production.
	ExitRelationDefault(c *RelationDefaultContext)

	// ExitJoinRelation is called when exiting the joinRelation production.
	ExitJoinRelation(c *JoinRelationContext)

	// ExitJoinType is called when exiting the joinType production.
	ExitJoinType(c *JoinTypeContext)

	// ExitJoinCriteria is called when exiting the joinCriteria production.
	ExitJoinCriteria(c *JoinCriteriaContext)

	// ExitAliasedRelation is called when exiting the aliasedRelation production.
	ExitAliasedRelation(c *AliasedRelationContext)

	// ExitTableName is called when exiting the tableName production.
	ExitTableName(c *TableNameContext)

	// ExitSubQueryRelation is called when exiting the subQueryRelation production.
	ExitSubQueryRelation(c *SubQueryRelationContext)

	// ExitGroupBy is called when exiting the groupBy production.
	ExitGroupBy(c *GroupByContext)

	// ExitSingleGroupingSet is called when exiting the singleGroupingSet production.
	ExitSingleGroupingSet(c *SingleGroupingSetContext)

	// ExitGroupByAllColumns is called when exiting the groupByAllColumns production.
	ExitGroupByAllColumns(c *GroupByAllColumnsContext)

	// ExitGroupingSet is called when exiting the groupingSet production.
	ExitGroupingSet(c *GroupingSetContext)

	// ExitHaving is called when exiting the having production.
	ExitHaving(c *HavingContext)

	// ExitOrderBy is called when exiting the orderBy production.
	ExitOrderBy(c *OrderByContext)

	// ExitSortItem is called when exiting the sortItem production.
	ExitSortItem(c *SortItemContext)

	// ExitLimitRowCount is called when exiting the limitRowCount production.
	ExitLimitRowCount(c *LimitRowCountContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitLogicalNot is called when exiting the logicalNot production.
	ExitLogicalNot(c *LogicalNotContext)

	// ExitPredicatedExpression is called when exiting the predicatedExpression production.
	ExitPredicatedExpression(c *PredicatedExpressionContext)

	// ExitOr is called when exiting the or production.
	ExitOr(c *OrContext)

	// ExitAnd is called when exiting the and production.
	ExitAnd(c *AndContext)

	// ExitValueExpressionDefault is called when exiting the valueExpressionDefault production.
	ExitValueExpressionDefault(c *ValueExpressionDefaultContext)

	// ExitArithmeticBinary is called when exiting the arithmeticBinary production.
	ExitArithmeticBinary(c *ArithmeticBinaryContext)

	// ExitDereference is called when exiting the dereference production.
	ExitDereference(c *DereferenceContext)

	// ExitColumnReference is called when exiting the columnReference production.
	ExitColumnReference(c *ColumnReferenceContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitParenExpression is called when exiting the parenExpression production.
	ExitParenExpression(c *ParenExpressionContext)

	// ExitNumericLiteral is called when exiting the numericLiteral production.
	ExitNumericLiteral(c *NumericLiteralContext)

	// ExitBooleanLiteral is called when exiting the booleanLiteral production.
	ExitBooleanLiteral(c *BooleanLiteralContext)

	// ExitBinaryComparisonPredicate is called when exiting the binaryComparisonPredicate production.
	ExitBinaryComparisonPredicate(c *BinaryComparisonPredicateContext)

	// ExitInPredicate is called when exiting the inPredicate production.
	ExitInPredicate(c *InPredicateContext)

	// ExitLikePredicate is called when exiting the likePredicate production.
	ExitLikePredicate(c *LikePredicateContext)

	// ExitRegexpPredicate is called when exiting the regexpPredicate production.
	ExitRegexpPredicate(c *RegexpPredicateContext)

	// ExitValueExpressionPredicate is called when exiting the valueExpressionPredicate production.
	ExitValueExpressionPredicate(c *ValueExpressionPredicateContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitProperties is called when exiting the properties production.
	ExitProperties(c *PropertiesContext)

	// ExitPropertyAssignments is called when exiting the propertyAssignments production.
	ExitPropertyAssignments(c *PropertyAssignmentsContext)

	// ExitProperty is called when exiting the property production.
	ExitProperty(c *PropertyContext)

	// ExitDefaultPropertyValue is called when exiting the defaultPropertyValue production.
	ExitDefaultPropertyValue(c *DefaultPropertyValueContext)

	// ExitNonDefaultPropertyValue is called when exiting the nonDefaultPropertyValue production.
	ExitNonDefaultPropertyValue(c *NonDefaultPropertyValueContext)

	// ExitBooleanValue is called when exiting the booleanValue production.
	ExitBooleanValue(c *BooleanValueContext)

	// ExitBasicStringLiteral is called when exiting the basicStringLiteral production.
	ExitBasicStringLiteral(c *BasicStringLiteralContext)

	// ExitUnquotedIdentifier is called when exiting the unquotedIdentifier production.
	ExitUnquotedIdentifier(c *UnquotedIdentifierContext)

	// ExitQuotedIdentifier is called when exiting the quotedIdentifier production.
	ExitQuotedIdentifier(c *QuotedIdentifierContext)

	// ExitBackQuotedIdentifier is called when exiting the backQuotedIdentifier production.
	ExitBackQuotedIdentifier(c *BackQuotedIdentifierContext)

	// ExitDigitIdentifier is called when exiting the digitIdentifier production.
	ExitDigitIdentifier(c *DigitIdentifierContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitDoubleLiteral is called when exiting the doubleLiteral production.
	ExitDoubleLiteral(c *DoubleLiteralContext)

	// ExitIntegerLiteral is called when exiting the integerLiteral production.
	ExitIntegerLiteral(c *IntegerLiteralContext)

	// ExitNonReserved is called when exiting the nonReserved production.
	ExitNonReserved(c *NonReservedContext)
}
