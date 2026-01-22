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

package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type RelationPlan struct {
	Root          plan.PlanNode
	Scope         *analyzer.Scope
	OutContext    *TranslationMap
	FieldMappings []*plan.Symbol
}

func (r *RelationPlan) getSymbol(fieldIndex tree.FieldIndex) *plan.Symbol {
	fieldIdx := int(fieldIndex)
	if fieldIdx < 0 || fieldIdx >= len(r.FieldMappings) {
		panic(fmt.Sprintf("no field->symbol mapping for field %d", fieldIdx))
	}
	return r.FieldMappings[fieldIdx]
}
