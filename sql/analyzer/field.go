package analyzer

import (
	"strings"

	"github.com/lindb/lindb/sql/tree"
)

type FieldID struct {
	RelationID *RelationID
	FieldIndex int
}

type Field struct {
	Name          string
	RelationAlias *tree.QualifiedName
}

func (f *Field) matchesPrefix(prefix *tree.QualifiedName) bool {
	return prefix == nil || (f.RelationAlias != nil && f.RelationAlias.HasSuffix(prefix))
}

func (f *Field) canResolve(name *tree.QualifiedName) bool {
	if f.Name == "" {
		return false
	}
	// TODO: need to know whether the qualified name and the name of this field were quoted
	return f.matchesPrefix(name.Prefix) && strings.ToLower(f.Name) == strings.ToLower(name.Suffix)
}

func (f *Field) String() string {
	if f.RelationAlias == nil {
		return f.Name
	}
	return f.RelationAlias.Name + "." + f.Name
}

type ResolvedField struct {
	Scope               *Scope
	Field               *Field
	HierarchyFieldIndex int
	RelationFieldIndex  int
	Local               bool
}

func (rf *ResolvedField) FieldID() *FieldID {
	return &FieldID{
		RelationID: rf.Scope.RelationID,
		FieldIndex: rf.Scope.RelationType.IndexOf(rf.Field),
	}
}
