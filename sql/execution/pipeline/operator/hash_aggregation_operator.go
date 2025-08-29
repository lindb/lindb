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

package operator

import (
	"context"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
)

type HashAggregationOperator struct {
	node  *plan.AggregationNode
	child Operator

	inbound *Queue
}

func NewHashAggregationOperator(node *plan.AggregationNode, child Operator) Operator {
	return &HashAggregationOperator{
		child:   child,
		node:    node,
		inbound: NewQueue(make(chan *types.Page)),
	}
}

func (h *HashAggregationOperator) Run(ctx context.Context, output chan<- *types.Page) {
	for {
		source, ok := h.inbound.Consume(ctx)
		if !ok {
			return
		}
		output <- source
	}
}

func (h *HashAggregationOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

func (h *HashAggregationOperator) Children() []Operator {
	return []Operator{h.child}
}

func (h *HashAggregationOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{h.inbound.GetInbound()}
}
