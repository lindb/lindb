package analyzer

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

type GroupingSetAnalysis struct {
	originalExpressions []tree.Expression
	complexExpressions  []tree.Expression
	ordinarySets        [][]*FieldID
}

func NewGroupingSetAnalysis(
	originalExpressions []tree.Expression,
	ordinarySets [][]*FieldID,
	complexExpressions []tree.Expression,
) *GroupingSetAnalysis {
	return &GroupingSetAnalysis{
		originalExpressions: originalExpressions,
		ordinarySets:        ordinarySets,
		complexExpressions:  complexExpressions,
	}
}

func (gsa *GroupingSetAnalysis) GetComplexExpressions() []tree.Expression {
	return gsa.complexExpressions
}

func (gsa *GroupingSetAnalysis) GetAllFields() (rs []*FieldID) {
	for _, fields := range gsa.ordinarySets {
		rs = append(rs, fields...)
	}
	return
}

func (gsa *GroupingSetAnalysis) GetOriginalExpression() []tree.Expression {
	return gsa.originalExpressions
}

type Analysis struct {
	root                  tree.Statement
	scopes                map[tree.NodeID]*Scope      // TODO: node ref?
	namedQueries          map[tree.NodeID]*tree.Query // table reference to with query
	selectAllResultFields map[tree.NodeID][]*Field
	selectExpressions     map[tree.NodeID][]*SelectExpression
	aliasedRelations      map[*tree.QualifiedName]tree.Relation
	tableMetadatas        map[tree.NodeID]*spi.TableMetadata
	relationNames         map[tree.NodeID]*tree.QualifiedName
	joins                 map[tree.NodeID]tree.Expression
	where                 map[tree.NodeID]tree.Expression
	groupingSets          map[tree.NodeID]*GroupingSetAnalysis
	having                map[tree.NodeID]tree.Expression
	orderByExpressions    map[tree.NodeID][]tree.Expression
	limit                 map[tree.NodeID]int64

	columnReferences map[tree.NodeID]*ResolvedField
}

func NewAnalysis(root tree.Statement) *Analysis {
	return &Analysis{
		root:                  root,
		scopes:                make(map[tree.NodeID]*Scope),
		namedQueries:          make(map[tree.NodeID]*tree.Query),
		selectAllResultFields: make(map[tree.NodeID][]*Field),
		selectExpressions:     make(map[tree.NodeID][]*SelectExpression),
		tableMetadatas:        make(map[tree.NodeID]*spi.TableMetadata),
		relationNames:         make(map[tree.NodeID]*tree.QualifiedName),
		aliasedRelations:      make(map[*tree.QualifiedName]tree.Relation),
		joins:                 make(map[tree.NodeID]tree.Expression),
		where:                 make(map[tree.NodeID]tree.Expression),
		groupingSets:          make(map[tree.NodeID]*GroupingSetAnalysis),
		having:                make(map[tree.NodeID]tree.Expression),
		orderByExpressions:    make(map[tree.NodeID][]tree.Expression),
		limit:                 make(map[tree.NodeID]int64),

		columnReferences: make(map[tree.NodeID]*ResolvedField),
	}
}

func (a *Analysis) GetRoot() tree.Node {
	return a.root
}

func (a *Analysis) SetScope(node tree.Node, scope *Scope) {
	a.scopes[node.GetID()] = scope
}

func (a *Analysis) GetScope(node tree.Node) (scope *Scope) {
	scope = a.scopes[node.GetID()]
	return
}

func (a *Analysis) SetSelectAllResultFields(node *tree.AllColumns, fields []*Field) {
	a.selectAllResultFields[node.GetID()] = fields
}

func (a *Analysis) GetSelectAllResultFields(node *tree.AllColumns) (fields []*Field) {
	fields = a.selectAllResultFields[node.GetID()]
	return
}

func (a *Analysis) GetOutputDescriptor(node tree.Node) *Relation {
	return a.GetScope(node).RelationType
}

func (a *Analysis) RegisterNamedQuery(tableReference *tree.Table, query *tree.Query) {
	a.namedQueries[tableReference.GetID()] = query
}

func (a *Analysis) GetNamedQuery(tableReference *tree.Table) (query *tree.Query) {
	query = a.namedQueries[tableReference.GetID()]
	return
}

func (a *Analysis) SetSelectExpressions(node tree.Node, expressions []*SelectExpression) {
	a.selectExpressions[node.GetID()] = expressions
}

func (a *Analysis) GetSelectExpressions(node tree.Node) (expressions []*SelectExpression) {
	expressions = a.selectExpressions[node.GetID()]
	return
}

func (a *Analysis) AddAliased(relation tree.Relation, aliased *tree.QualifiedName) {
	a.aliasedRelations[aliased] = relation
}

func (a *Analysis) GetRelationByAliased(aliased *tree.QualifiedName) (relation tree.Relation) {
	relation = a.aliasedRelations[aliased]
	return
}

func (a *Analysis) RegisterTableMetadata(table *tree.Table, tableMetadata *spi.TableMetadata) {
	a.tableMetadatas[table.GetID()] = tableMetadata
}

func (a Analysis) GetTableMetadata(table *tree.Table) (tableMetadata *spi.TableMetadata) {
	tableMetadata = a.tableMetadatas[table.GetID()]
	return
}

func (a *Analysis) SetRelationName(relation tree.Relation, name *tree.QualifiedName) {
	a.relationNames[relation.GetID()] = name
}

func (a *Analysis) SetJoinCriteria(node *tree.Join, criteria tree.Expression) {
	a.joins[node.GetID()] = criteria
}

func (a *Analysis) GetJoinCriteria(node *tree.Join) (criteria tree.Expression) {
	criteria = a.joins[node.GetID()]
	return
}

func (a *Analysis) SetWhere(node tree.Node, expression tree.Expression) {
	a.where[node.GetID()] = expression
}

func (a *Analysis) GetWhere(node tree.Node) tree.Expression {
	return a.where[node.GetID()]
}

func (a *Analysis) SetGroupingSets(node tree.Node, groupingSets *GroupingSetAnalysis) {
	a.groupingSets[node.GetID()] = groupingSets
}

func (a *Analysis) IsAggregation(node tree.Node) (ok bool) {
	_, ok = a.groupingSets[node.GetID()]
	return
}

func (a *Analysis) GetGroupingSets(node tree.Node) *GroupingSetAnalysis {
	return a.groupingSets[node.GetID()]
}

func (a *Analysis) SetHaving(query *tree.QuerySpecification, expression tree.Expression) {
	a.having[query.GetID()] = expression
}

func (a *Analysis) GetHaving(query *tree.QuerySpecification) tree.Expression {
	return a.having[query.GetID()]
}

func (a *Analysis) SetOrderByExpressions(node tree.Node, orderByExpressions []tree.Expression) {
	a.orderByExpressions[node.GetID()] = orderByExpressions // FIXME: copy it?
}

func (a *Analysis) GetOrderByExpressions(node tree.Node) []tree.Expression {
	return a.orderByExpressions[node.GetID()]
}

func (a *Analysis) SetLimit(node tree.Node, rowCount int64) {
	a.limit[node.GetID()] = rowCount
}

func (a *Analysis) AddColumnReference(expression tree.Expression, field *ResolvedField) {
	a.columnReferences[expression.GetID()] = field
}

func (a *Analysis) IsColumnReference(node tree.Expression) bool {
	_, ok := a.columnReferences[node.GetID()]
	return ok
}

func (a *Analysis) GetColumnReferenceField(node tree.Expression) (field *ResolvedField) {
	field = a.columnReferences[node.GetID()]
	return
}

func (a *Analysis) RecordSubQueries(node tree.Node, expressionAnalysis *ExpressionAnalysis) {
	// panic("impl it")
}

func (a *Analysis) GetStatement() tree.Statement {
	return a.root
}
