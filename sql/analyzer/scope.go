package analyzer

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type BasisType string

var (
	TABLE BasisType = "Table"
	FIELD BasisType = "Field"
)

type AsteriskedIdentifierChain struct {
	Type         BasisType
	Scope        *Scope
	RelationType *Relation
}

type Scope struct {
	Parent        *Scope
	QueryBoundary bool

	RelationID   *RelationID
	RelationType *Relation
	NamedQueries map[string]*tree.WithQuery
}

func createScope(parent *Scope) *Scope {
	return &Scope{
		Parent:       parent,
		RelationType: NewRelation(TableRelation, nil), // FIXME:???
	}
}

func (scope *Scope) getField(index int) *ResolvedField {
	parentFieldCount := 0
	parentScope := scope.getLocalParent()
	if parentScope != nil {
		parentFieldCount = parentScope.getlocalScopeFieldCount()
	}

	return scope.asResolvedField(scope.RelationType.getFieldByIndex(index), parentFieldCount, true)
}

func (scope *Scope) getNameQuery(name string) (withQuery *tree.WithQuery) {
	withQuery, ok := scope.NamedQueries[name]
	if ok {
		return withQuery
	}
	if scope.Parent != nil {
		return scope.Parent.getNameQuery(name)
	}
	return
}

func (scope *Scope) IsLocalScope(other *Scope) bool {
	return scope.findLocally(func(scope *Scope) bool {
		if scope == other {
			return true
		}
		// FIXME: check replation set
		fmt.Printf("scope=%v,other=%v\n", scope, other)
		return scope.RelationID != nil &&
			other.RelationID != nil &&
			scope.RelationID.SourceNode.GetID() == other.RelationID.SourceNode.GetID()
	}) != nil
}

func (scope *Scope) tryResolveField(node tree.Expression, name *tree.QualifiedName) *ResolvedField {
	return scope.resolveField(node, name, true)
}

func (scope *Scope) resolveField(node tree.Expression, name *tree.QualifiedName, local bool) *ResolvedField {
	fields := scope.RelationType.resolveFields(name)
	if len(fields) > 1 {
		panic(fmt.Sprintf("column '%s' is ambiguous", name.Name))
	}
	if len(fields) == 1 {
		// TODO: dup
		parentFieldCount := 0
		parentScope := scope.getLocalParent()
		if parentScope != nil {
			parentFieldCount = parentScope.getlocalScopeFieldCount()
		}

		return scope.asResolvedField(fields[0], parentFieldCount, local)
	}
	// TODO: column ref
	if scope.Parent != nil {
		// TODO: query boundary
		return scope.Parent.resolveField(node, name, local)
	}
	return nil
}

func (scope *Scope) asResolvedField(field *Field, fieldIndexOffset int, local bool) *ResolvedField {
	relationFieldIndex := scope.RelationType.IndexOf(field)
	return &ResolvedField{
		Field:               field,
		Scope:               scope,
		RelationFieldIndex:  relationFieldIndex,
		HierarchyFieldIndex: relationFieldIndex + fieldIndexOffset,
		Local:               local,
	}
}

func (scope *Scope) resolveAsteriskedIdentifierChain(identifierChain *tree.QualifiedName, selectItem *tree.AllColumns) *AsteriskedIdentifierChain {
	partsLen := len(identifierChain.Parts)
	var (
		scopeForTableRef *Scope
		scopeForFieldRef *Scope
	)
	find := func(scope *Scope, match func(field *Field) bool) bool {
		fmt.Println(scope.RelationType)
		fields := scope.RelationType.Fields
		for i := range fields {
			if match(fields[i]) {
				return true
			}
		}
		return false
	}
	if partsLen <= 3 {
		scopeForTableRef = scope.findLocally(func(scope *Scope) bool {
			return find(scope, func(field *Field) bool {
				return field.matchesPrefix(identifierChain)
			})
		})
	}
	if partsLen >= 2 {
		part0 := identifierChain.Parts[0]
		part1 := identifierChain.Parts[1]
		scopeForFieldRef = scope.findLocally(func(scope *Scope) bool {
			return find(scope, func(field *Field) bool {
				return field.Name != "" && field.Name == part1 && field.matchesPrefix(
					tree.NewQualifiedName([]*tree.Identifier{{Value: part0}}),
				)
			})
		})
	}
	if scopeForTableRef != nil && scopeForFieldRef != nil {
		panic("fix me.....")
	}
	if scopeForTableRef != nil {
		return &AsteriskedIdentifierChain{
			Type:         TABLE,
			Scope:        scopeForTableRef,
			RelationType: scopeForTableRef.RelationType,
		}
	}
	if scopeForFieldRef != nil {
		return &AsteriskedIdentifierChain{
			Type: FIELD,
		}
	}
	fmt.Println(scopeForFieldRef)
	fmt.Println(scopeForTableRef)
	if scope.Parent == nil {
		return nil
	}
	// FIXME:
	return scope.Parent.resolveAsteriskedIdentifierChain(identifierChain, selectItem)
}

func (scope *Scope) getlocalScopeFieldCount() int {
	parent := 0
	parentScope := scope.getLocalParent()
	if parentScope != nil {
		parent = parentScope.getlocalScopeFieldCount()
	}
	return parent + len(scope.RelationType.Fields) // TODO: all field??
}

func (scope *Scope) findLocally(match func(scope *Scope) bool) *Scope {
	s := scope
	for {
		if match(s) {
			return s
		}
		parent := s.getLocalParent()
		if parent == nil {
			break
		}
		if parent == scope {
			panic("===")
		}
		fmt.Println("partnn......")
		fmt.Println(parent)
		s = parent
	}
	return nil
}

func (scope *Scope) getOuterQueryParent() *Scope {
	s := scope
	if s.Parent != nil {
		// FIXME:MMMMMMMM
		return s.Parent
	}
	return nil
}

func (scope *Scope) getLocalParent() *Scope {
	// FIXME: check query boundary
	return scope.Parent
}

type ScopeBuilder struct {
	parent       *Scope
	relationID   *RelationID
	relationType *Relation
	namedQueries map[string]*tree.WithQuery
}

func NewScopeBuilder(parentScope *Scope) *ScopeBuilder {
	builder := &ScopeBuilder{
		namedQueries: make(map[string]*tree.WithQuery),
	}
	// FIXME:
	return builder.withParent(parentScope)
}

func (b *ScopeBuilder) withParent(parent *Scope) *ScopeBuilder {
	b.parent = parent
	return b
}

func (b *ScopeBuilder) withRelation(relationID *RelationID, relationType *Relation) *ScopeBuilder {
	b.relationID = relationID
	b.relationType = relationType
	return b
}

func (b *ScopeBuilder) withNameQuery(name string, withQuery *tree.WithQuery) *ScopeBuilder {
	b.namedQueries[name] = withQuery
	return b
}

func (b *ScopeBuilder) containsNamedQuery(name string) (exist bool) {
	_, exist = b.namedQueries[name]
	return
}

func (b *ScopeBuilder) build() *Scope {
	if b.relationType == nil {
		panic("2222....")
	}
	return &Scope{
		Parent:       b.parent,
		RelationID:   b.relationID,
		RelationType: b.relationType,
		NamedQueries: b.namedQueries, // FIXME: copy it?
	}
}
