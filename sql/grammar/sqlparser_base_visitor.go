// Code generated from ./sql/grammar/SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // SQLParser
import "github.com/antlr4-go/antlr/v4"

type BaseSQLParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDdlStatement(ctx *DdlStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitStatementDefault(ctx *StatementDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitExplain(ctx *ExplainContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitExplainAnalyze(ctx *ExplainAnalyzeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitAdminStatement(ctx *AdminStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitUtilityStatement(ctx *UtilityStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitExplainType(ctx *ExplainTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitCreateDatabase(ctx *CreateDatabaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitEngineOption(ctx *EngineOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitRollupOptions(ctx *RollupOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDropDatabase(ctx *DropDatabaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitCreateBroker(ctx *CreateBrokerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitFlushDatabase(ctx *FlushDatabaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitCompactDatabase(ctx *CompactDatabaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowMaster(ctx *ShowMasterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowBrokers(ctx *ShowBrokersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowRequests(ctx *ShowRequestsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowLimit(ctx *ShowLimitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowMetadataTypes(ctx *ShowMetadataTypesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowMetadatas(ctx *ShowMetadatasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowState(ctx *ShowStateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowDatabases(ctx *ShowDatabasesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowMemoryDatabases(ctx *ShowMemoryDatabasesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowReplications(ctx *ShowReplicationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowNamespaces(ctx *ShowNamespacesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowTableNames(ctx *ShowTableNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitShowColumns(ctx *ShowColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitUseStatement(ctx *UseStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitWith(ctx *WithContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitNamedQuery(ctx *NamedQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQueryNoWith(ctx *QueryNoWithContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQueryTermDefault(ctx *QueryTermDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSubquery(ctx *SubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQuerySpecification(ctx *QuerySpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSelectSingle(ctx *SelectSingleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSelectAll(ctx *SelectAllContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitRelationDefault(ctx *RelationDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitJoinRelation(ctx *JoinRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitJoinType(ctx *JoinTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitJoinCriteria(ctx *JoinCriteriaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitAliasedRelation(ctx *AliasedRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitTableName(ctx *TableNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSubQueryRelation(ctx *SubQueryRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSingleGroupingSet(ctx *SingleGroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitGroupByAllColumns(ctx *GroupByAllColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitGroupingSet(ctx *GroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitHaving(ctx *HavingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitOrderBy(ctx *OrderByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitSortItem(ctx *SortItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitLimitRowCount(ctx *LimitRowCountContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitPredicatedExpression(ctx *PredicatedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitOr(ctx *OrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitAnd(ctx *AndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitValueExpressionDefault(ctx *ValueExpressionDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitArithmeticBinary(ctx *ArithmeticBinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDereference(ctx *DereferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitColumnReference(ctx *ColumnReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitParenExpression(ctx *ParenExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitNumericLiteral(ctx *NumericLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitIntervalLiteral(ctx *IntervalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitTimestampPredicate(ctx *TimestampPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBinaryComparisonPredicate(ctx *BinaryComparisonPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBetweenPredicate(ctx *BetweenPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitInPredicate(ctx *InPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitLikePredicate(ctx *LikePredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitRegexpPredicate(ctx *RegexpPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitValueExpressionPredicate(ctx *ValueExpressionPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitProperty(ctx *PropertyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDefaultPropertyValue(ctx *DefaultPropertyValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBooleanValue(ctx *BooleanValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBasicStringLiteral(ctx *BasicStringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitUnquotedIdentifier(ctx *UnquotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitQuotedIdentifier(ctx *QuotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDigitIdentifier(ctx *DigitIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitInterval(ctx *IntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitIntervalUnit(ctx *IntervalUnitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLParserVisitor) VisitNonReserved(ctx *NonReservedContext) interface{} {
	return v.VisitChildren(ctx)
}
