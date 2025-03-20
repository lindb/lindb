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

package exchange

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type RemoteExchangeOperator struct {
	ctx     context.Context
	node    *plan.RemoteSourceNode
	inbound chan *types.Page
}

func NewRemoteExchangeOperator(ctx context.Context, node *plan.RemoteSourceNode, numOfChild int) operator.SourceOperator {
	return &RemoteExchangeOperator{
		ctx:     ctx,
		node:    node,
		inbound: make(chan *types.Page),
	}
}

func (op *RemoteExchangeOperator) GetSourceID() plan.PlanNodeID {
	return op.node.GetNodeID()
}

// GetOutput implements Operator
func (op *RemoteExchangeOperator) Run(output chan<- *types.Page) {
	for page := range op.inbound {
		if page == nil {
			fmt.Println("page nil")
			continue
		}
		// FIXME: do merge logic
		// it := page.Iterator()
		// groupingColumns := page.Grouping
		// for row := it.Begin(); row != it.End(); row = it.Next() {
		// 	fmt.Println("kkkk......")
		// 	op.mergedPage.AppendColumn(page.Layout[], page.Columns[])
		// }
		fmt.Printf("exchange merge,page=%v\n", string(encoding.JSONMarshal(page)))

		output <- page
	}
}

func (op *RemoteExchangeOperator) Receive(page *types.Page) {
	op.inbound <- page
}

func (op *RemoteExchangeOperator) GetLayout() []*plan.Symbol {
	return op.node.GetOutputSymbols()
}

func (op *RemoteExchangeOperator) Complete() {
	close(op.inbound)
}

func (op *RemoteExchangeOperator) Children() []operator.Operator {
	return nil
}
