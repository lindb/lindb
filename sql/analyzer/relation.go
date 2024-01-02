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
	Fields       []*Field
	FieldIndexes map[*Field]int
}

func NewRelation(relationType RelationType, fields []*Field) *Relation {
	rt := &Relation{
		Type:         relationType,
		Fields:       fields,
		FieldIndexes: make(map[*Field]int),
	}
	for i, f := range fields {
		rt.FieldIndexes[f] = i
	}
	return rt
}

func (r *Relation) withAlias(relationAlias string) *Relation {
	var fields []*Field
	for i := range r.Fields {
		field := r.Fields[i]

		fields = append(fields, &Field{
			Name:          field.Name,
			RelationAlias: tree.NewQualifiedName([]*tree.Identifier{{Value: relationAlias}}),
		})
	}
	return NewRelation(AliasedRelation, fields)
}

func (r *Relation) joinWith(other *Relation) *Relation {
	var fields []*Field
	fields = append(fields, r.Fields...)
	fields = append(fields, other.Fields...)
	return NewRelation(JoinRelation, fields)
}

func (r *Relation) getFieldByIndex(fieldIndex int) *Field {
	return r.Fields[fieldIndex]
}

func (r *Relation) resolveFields(name *tree.QualifiedName) (result []*Field) {
	return lo.Filter(r.Fields, func(item *Field, _ int) bool {
		return item.canResolve(name)
	})
}

func (r *Relation) IndexOf(field *Field) int {
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
