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

package plan

type SymbolMap struct {
	From *Symbol
	To   *Symbol
}

type SymbolMapper struct {
	fn func(symbol *Symbol) *Symbol
}

func NewSymbolMapper(mapping map[string]*Symbol) *SymbolMapper {
	return &SymbolMapper{
		fn: func(symbol *Symbol) *Symbol {
			for {
				val, ok := mapping[symbol.Name]
				if ok && val.Name != symbol.Name {
					symbol = val
				} else {
					break
				}
			}
			return symbol
		},
	}
}

func (m *SymbolMapper) MapAggregation(node *AggregationNode, source PlanNode, newNodeID PlanNodeID) *AggregationNode {
	return NewAggregationNode(newNodeID, source, node.Aggregations, node.GroupingSets, node.Step)
}

func (m *SymbolMapper) MapSymbol(symobl *Symbol) *Symbol {
	return m.fn(symobl)
}
