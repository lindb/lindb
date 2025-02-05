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

package tree

import (
	"strings"

	"github.com/lindb/lindb/spi/types"
)

type FieldIndex int

type Field struct {
	RelationAlias *QualifiedName
	Name          string
	DataType      types.DataType
	AggType       types.AggregateType
	Index         FieldIndex // set field index when statement analyzer(visit relation)
	Hidden        bool
}

func (f *Field) MatchesPrefix(prefix *QualifiedName) bool {
	return prefix == nil || (f.RelationAlias != nil && f.RelationAlias.HasSuffix(prefix))
}

func (f *Field) CanResolve(name *QualifiedName) bool {
	if f.Name == "" {
		return false
	}
	// TODO: need to know whether the qualified name and the name of this field were quoted
	return f.MatchesPrefix(name.Prefix) && strings.EqualFold(f.Name, name.Suffix)
}

func (f *Field) String() string {
	if f.RelationAlias == nil {
		return f.Name
	}
	return f.RelationAlias.Name + "." + f.Name
}
