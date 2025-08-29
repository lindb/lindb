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

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type RemoteExchangeOperator struct {
	ctx     context.Context
	node    *plan.RemoteSourceNode
	inbound *operator.Queue
}

func NewRemoteExchangeOperator(ctx context.Context, node *plan.RemoteSourceNode, numOfChild int) operator.SourceOperator {
	return &RemoteExchangeOperator{
		ctx:     ctx,
		node:    node,
		inbound: operator.NewQueue(make(chan *types.Page)),
	}
}

func (op *RemoteExchangeOperator) GetSourceID() plan.PlanNodeID {
	return op.node.GetNodeID()
}

// Run runs the exchange operator, consuming the pages from inbound channel and merging the pages.
func (op *RemoteExchangeOperator) Run(ctx context.Context, output chan<- *types.Page) {
	var buffer []*types.Page

	for {
		// consume the pages from inbound channel
		page, ok := op.inbound.Consume(ctx)
		fmt.Printf("receive pagll...e=%v\n", page)
		if !ok {
			// TODO: merge pages (streaming)
			mergedPage := types.MergePages(buffer)
			if mergedPage != nil {
				output <- mergedPage
			}
			// return if inbound channel is closed
			return
		}
		if page == nil {
			continue
		}
		if page.Error != "" {
			panic(page.Error)
		}
		buffer = append(buffer, page)
	}
}

func (op *RemoteExchangeOperator) Receive(page *types.Page) {
	op.inbound.Produce(page)
}

func (op *RemoteExchangeOperator) GetLayout() []*plan.Symbol {
	return op.node.GetOutputSymbols()
}

func (op *RemoteExchangeOperator) Complete() {
	op.inbound.Close()
}

func (op *RemoteExchangeOperator) Children() []operator.Operator {
	return nil
}

func (op *RemoteExchangeOperator) GetInbounds() []chan *types.Page {
	return nil
}
