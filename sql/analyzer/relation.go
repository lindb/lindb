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

	"github.com/lindb/common/pkg/encoding"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
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
	Fields []*tree.Field

	// cannot use field name, field name is empty when select item isn't Identifier/Dereference
	FieldIndexes map[*tree.Field]tree.FieldIndex

	fieldsMap map[tree.FieldIndex]*tree.Field

	maxIndex tree.FieldIndex
}

func NewRelation(fields []*tree.Field) *Relation {
	rt := &Relation{
		Fields:       fields,
		FieldIndexes: make(map[*tree.Field]tree.FieldIndex),
		fieldsMap:    make(map[tree.FieldIndex]*tree.Field),
	}
	for _, f := range fields {
		rt.FieldIndexes[f] = f.Index
		rt.fieldsMap[f.Index] = f

		if f.Index > rt.maxIndex {
			rt.maxIndex = f.Index
		}
	}
	fmt.Printf("new relation fields=%v\n", rt.FieldIndexes)
	return rt
}

func (r *Relation) withAlias(relationAlias string, columnAliases []string) *Relation {
	var fields []*tree.Field
	for i, field := range r.Fields {
		columnAlias := field.Name
		if len(columnAliases) != 0 {
			columnAlias = columnAliases[i]
		}
		fmt.Printf("columnAlias=%s,columnAliases=%v", columnAlias, columnAliases)

		fields = append(fields, &tree.Field{
			Name:          columnAlias,
			DataType:      field.DataType,
			AggType:       field.AggType,
			Index:         field.Index,
			Hidden:        field.Hidden,
			RelationAlias: relationAlias,
		})
	}
	return NewRelation(fields)
}

func (r *Relation) joinWith(other *Relation) *Relation {
	var fields []*tree.Field
	index := tree.FieldIndex(0)
	add := func(fieldList []*tree.Field) {
		for _, field := range fieldList {
			// create new field of join relation
			newField := field.Clone()
			newField.Index = index // NOTE: need reset index, maybe table a/b have same column name
			fields = append(fields, newField)
			index++
		}
	}

	add(r.Fields)
	add(other.Fields)
	return NewRelation(fields)
}

func (r *Relation) getFieldByIndex(fieldIndex tree.FieldIndex) *tree.Field {
	return r.fieldsMap[fieldIndex]
}

func (r *Relation) resolveFields(name *tree.QualifiedName) (result []*tree.Field) {
	return lo.Filter(r.Fields, func(item *tree.Field, _ int) bool {
		return item.CanResolve(name)
	})
}

func (r *Relation) IndexOf(field *tree.Field) tree.FieldIndex {
	index, ok := r.FieldIndexes[field]
	fmt.Printf("relation index of %v,%v,%v,%v\n", r.FieldIndexes, ok, r.maxIndex, string(encoding.JSONMarshal(field)))
	if ok {
		return index
	}
	if field.DataType == types.DTDynamic {
		r.Fields = append(r.Fields, field)

		r.maxIndex++
		field.Index = r.maxIndex
		r.FieldIndexes[field] = field.Index
		r.fieldsMap[field.Index] = field
		return field.Index
	}
	return 0
}

type RelationID struct {
	SourceNode tree.Node
}

func NewRelationID(sourceNode tree.Node) *RelationID {
	return &RelationID{
		SourceNode: sourceNode,
	}
}
