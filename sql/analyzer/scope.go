// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package analyzer

import (
	"fmt"

	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/samber/lo"

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

	Dynamic bool

	RelationID   *RelationID
	RelationType *Relation
	NamedQueries map[string]*tree.WithQuery
}

func createScope(parent *Scope) *Scope {
	return &Scope{
		Parent:       parent,
		RelationType: NewRelation(nil), // FIXME:???
	}
}

func (scope *Scope) getField(index tree.FieldIndex) *ResolvedField {
	parentFieldCount := 0
	parentScope := scope.getLocalParent()
	if parentScope != nil {
		parentFieldCount = parentScope.getLocalScopeFieldCount()
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
	return scope.findLocally(func(matchScope *Scope) bool {
		if matchScope == other {
			return true
		}
		// FIXME: check replation set
		return matchScope.RelationID != nil &&
			other.RelationID != nil &&
			matchScope.RelationID.SourceNode.GetID() == other.RelationID.SourceNode.GetID()
	}) != nil
}

func (scope *Scope) tryResolveField(node tree.Expression, name *tree.QualifiedName) *ResolvedField {
	return scope.resolveField(node, name, true)
}

func (scope *Scope) resolveField(node tree.Expression, name *tree.QualifiedName, local bool) *ResolvedField { //nolint
	fields := scope.RelationType.resolveFields(name)
	if len(fields) > 1 {
		panic(fmt.Sprintf("column '%s' is ambiguous", name.Name))
	}
	if len(fields) == 1 {
		// TODO: dup
		parentFieldCount := 0
		parentScope := scope.getLocalParent()
		if parentScope != nil {
			parentFieldCount = parentScope.getLocalScopeFieldCount()
		}

		return scope.asResolvedField(fields[0], parentFieldCount, local)
	}
	if scope.Dynamic {
		return &ResolvedField{
			Field: &tree.Field{
				Name:     name.Name,
				DataType: larrow.ExtensionTypes.Dynamic,
			},
			Scope: scope,
			// RelationFieldIndex:  relationFieldIndex,
			// HierarchyFieldIndex: relationFieldIndex + tree.FieldIndex(fieldIndexOffset),
			Local: local,
		}
	}
	// TODO: column ref
	if scope.Parent != nil {
		// TODO: query boundary
		return scope.Parent.resolveField(node, name, local)
	}
	return nil
}

func (scope *Scope) asResolvedField(field *tree.Field, fieldIndexOffset int, local bool) *ResolvedField {
	relationFieldIndex := scope.RelationType.IndexOf(field)
	return &ResolvedField{
		Field:               field,
		Scope:               scope,
		RelationFieldIndex:  relationFieldIndex,
		HierarchyFieldIndex: relationFieldIndex + tree.FieldIndex(fieldIndexOffset),
		Local:               local,
	}
}

func (scope *Scope) resolveAsteriskedIdentifierChain(
	identifierChain *tree.QualifiedName,
	selectItem *tree.AllColumns, //nolint
) *AsteriskedIdentifierChain {
	partsLen := len(identifierChain.Parts)
	var (
		scopeForTableRef *Scope
		scopeForFieldRef *Scope
	)
	find := func(scope *Scope, match func(field *tree.Field) bool) bool {
		fields := scope.RelationType.Fields
		lo.ContainsBy(fields, func(item *tree.Field) bool {
			return match(item)
		})
		return false
	}
	if partsLen <= 3 {
		scopeForTableRef = scope.findLocally(func(matchScope *Scope) bool {
			return find(scope, func(field *tree.Field) bool {
				return field.MatchesPrefix(identifierChain)
			})
		})
	}
	if partsLen >= 2 {
		part0 := identifierChain.Parts[0]
		part1 := identifierChain.Parts[1]
		scopeForFieldRef = scope.findLocally(func(matchScope *Scope) bool {
			return find(scope, func(field *tree.Field) bool {
				return field.Name != "" && field.Name == part1 && field.MatchesPrefix(
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
	if scope.Parent == nil {
		return nil
	}
	// FIXME:
	return scope.Parent.resolveAsteriskedIdentifierChain(identifierChain, selectItem)
}

func (scope *Scope) getLocalScopeFieldCount() int {
	parent := 0
	parentScope := scope.getLocalParent()
	if parentScope != nil {
		parent = parentScope.getLocalScopeFieldCount()
	}
	return parent + len(scope.RelationType.Fields) // TODO: all field??
}

func (scope *Scope) findLocally(match func(matchScope *Scope) bool) *Scope {
	s := scope
	for {
		if match(s) {
			return s
		}
		parent := s.getLocalParent()
		if parent == nil {
			break
		}
		s = parent
	}
	return nil
}

func (scope *Scope) getLocalParent() *Scope {
	// FIXME: check query boundary, out query
	return nil
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
