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
	"github.com/lindb/lindb/sql/tree"
)

type FieldID struct {
	RelationID *RelationID
	FieldIndex tree.FieldIndex
}

type ResolvedField struct {
	Scope               *Scope
	Field               *tree.Field
	HierarchyFieldIndex tree.FieldIndex
	RelationFieldIndex  tree.FieldIndex
	Local               bool
}

func (rf *ResolvedField) FieldID() *FieldID {
	return &FieldID{
		RelationID: rf.Scope.RelationID,
		FieldIndex: rf.Scope.RelationType.IndexOf(rf.Field),
	}
}
