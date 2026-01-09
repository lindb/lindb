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

package join

import (
	"context"
	"fmt"
	"sync"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type rows struct {
	left, right types.Row
}

type HashJoinOperator struct {
	node        *plan.JoinNode
	left, right operator.Operator

	memTable              map[string]*rows
	leftKeys, rightKeys   []int
	leftScope, rightScope []*plan.Symbol

	leftInbound, rightInbound *operator.Queue
}

func NewHashJoinOperator(node *plan.JoinNode, left, right operator.Operator) operator.Operator {
	return &HashJoinOperator{
		node:       node,
		left:       left,
		right:      right,
		leftScope:  node.Left.GetOutputSymbols(),
		rightScope: node.Right.GetOutputSymbols(),

		leftInbound:  operator.NewQueue(make(chan *types.Page)),
		rightInbound: operator.NewQueue(make(chan *types.Page)),
	}
}

// Children implements operator.Operator.
func (h *HashJoinOperator) Children() []operator.Operator {
	return []operator.Operator{h.left, h.right}
}

// GetInbounds implements operator.Operator.
func (h *HashJoinOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{h.leftInbound.GetInbound(), h.rightInbound.GetInbound()}
}

// GetLayout implements operator.Operator.
func (h *HashJoinOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

func (h *HashJoinOperator) prepare() {
	h.memTable = make(map[string]*rows)

	joinCriteria := h.node.Criteria
	h.leftKeys = make([]int, len(joinCriteria))
	h.rightKeys = make([]int, len(joinCriteria))

	findIndex := func(keys []int, index int, key *plan.Symbol, scope []*plan.Symbol) {
		_, keySymbol, ok := lo.FindIndexOf(scope, func(item *plan.Symbol) bool {
			return item.Name == key.Name
		})
		if !ok {
			panic("not find left key")
		}
		keys[index] = keySymbol
	}

	for idx, expr := range joinCriteria {
		findIndex(h.leftKeys, idx, expr.Left, h.leftScope)
		findIndex(h.rightKeys, idx, expr.Right, h.rightScope)
	}
}

// Run implements operator.Operator.
func (h *HashJoinOperator) Run(ctx context.Context, output chan<- *types.Page) {
	h.prepare()

	var wg sync.WaitGroup
	wg.Add(2)

	consume := func(inbound *operator.Queue, keys []int, isLeft bool) {
		defer wg.Done()
		for {
			page, ok := inbound.Consume(ctx)
			if !ok {
				return
			}
			h.process(page, keys, isLeft)
			fmt.Printf("isLeft %v=>>>>>>>%v\n", isLeft, page)
		}
	}

	go consume(h.leftInbound, h.leftKeys, true)
	go consume(h.rightInbound, h.rightKeys, false)

	wg.Wait()

	newPage := types.NewPage()
	outputs := h.node.GetOutputSymbols()
	outputColumns := make([]*types.Column, len(outputs))
	for i, output := range outputs {
		outputColumns[i] = types.NewColumn()
		newPage.AppendColumn(
			types.NewColumnInfo(output.Name, output.DataType, output.Hidden, output.AggType),
			outputColumns[i])
	}

	for key, memRows := range h.memTable {
		fmt.Printf("join key=%v,memRows=%v........\n", key, memRows)
		if memRows.left != nil && memRows.right != nil {
			for i := range len(h.leftScope) {
				outputColumns[i].Append(memRows.left.Get(i))
			}

			for i := range len(h.rightScope) {
				outputColumns[i+len(h.leftScope)].Append(memRows.right.Get(i))
			}
		}
	}
	output <- newPage
}

func (h *HashJoinOperator) String() string {
	return "HashJoinOperator"
}

func (h *HashJoinOperator) process(page *types.Page, keys []int, isLeft bool) {
	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		column := row.GetString(keys[0])
		memRows, ok := h.memTable[string(*column)]
		if !ok {
			memRows = &rows{}
			h.memTable[string(*column)] = memRows
		}

		if isLeft {
			memRows.left = row
		} else {
			memRows.right = row
		}
	}
}
