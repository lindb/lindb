package analyzer

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/tree"
)

type RelationType string

var (
	TableRelation   RelationType = "Table"
	AliasedRelation RelationType = "Aliased"
	JoinRelation    RelationType = "Join"
	UnknownRelation RelationType = "Unknown"
)

type Relation struct {
	Type         RelationType
	Fields       []*tree.Field
	FieldIndexes map[*tree.Field]int
}

func NewRelation(relationType RelationType, fields []*tree.Field) *Relation {
	rt := &Relation{
		Type:         relationType,
		Fields:       fields,
		FieldIndexes: make(map[*tree.Field]int),
	}
	for i, f := range fields {
		rt.FieldIndexes[f] = i
	}
	return rt
}

func (r *Relation) withAlias(relationAlias string) *Relation {
	var fields []*tree.Field
	for i := range r.Fields {
		field := r.Fields[i]

		fields = append(fields, &tree.Field{
			Name:          field.Name,
			DataType:      field.DataType,
			RelationAlias: tree.NewQualifiedName([]*tree.Identifier{{Value: relationAlias}}),
		})
	}
	return NewRelation(AliasedRelation, fields)
}

func (r *Relation) joinWith(other *Relation) *Relation {
	var fields []*tree.Field
	fields = append(fields, r.Fields...)
	fields = append(fields, other.Fields...)
	return NewRelation(JoinRelation, fields)
}

func (r *Relation) getFieldByIndex(fieldIndex int) *tree.Field {
	return r.Fields[fieldIndex]
}

func (r *Relation) resolveFields(name *tree.QualifiedName) (result []*tree.Field) {
	return lo.Filter(r.Fields, func(item *tree.Field, _ int) bool {
		return item.CanResolve(name)
	})
}

func (r *Relation) IndexOf(field *tree.Field) int {
	return r.FieldIndexes[field]
}

type RelationID struct {
	SourceNode tree.Node
}

func NewRelationID(sourceNode tree.Node) *RelationID {
	return &RelationID{
		SourceNode: sourceNode,
	}
}
