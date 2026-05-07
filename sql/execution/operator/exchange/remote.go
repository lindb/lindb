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

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
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
		inbound: operator.NewQueue(make(chan arrow.RecordBatch)),
	}
}

func (op *RemoteExchangeOperator) GetSourceID() plan.PlanNodeID {
	return op.node.GetNodeID()
}

// Run runs the exchange operator, consuming the pages from inbound channel and merging the pages.
func (op *RemoteExchangeOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	var buffer []arrow.RecordBatch

	for {
		// consume the pages from inbound channel
		page, ok := op.inbound.Consume(ctx)
		if !ok {
			// Channel closed: check whether it was a normal Complete() or a Fail().
			if errMsg := op.inbound.ErrMsg(); errMsg != "" {
				// Remote storage task failed — panic so the pipeline captures the error
				// and propagates it back to the client via task_execution → dml → HTTP 500.
				panic(errMsg)
			}
			// TODO: merge pages (streaming)
			mergedPage := types.MergeRecords(buffer)
			if mergedPage != nil {
				output <- mergedPage
			}
			// return if inbound channel is closed
			return
		}
		if page == nil {
			continue
		}
		buffer = append(buffer, page)
	}
}

func (op *RemoteExchangeOperator) Receive(record arrow.RecordBatch) {
	op.inbound.Produce(record)
}

// Fail signals that the remote storage task failed. Run() will panic with
// errMsg so the pipeline error path surfaces the failure to the client.
func (op *RemoteExchangeOperator) Fail(errMsg string) {
	op.inbound.Fail(errMsg)
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

func (op *RemoteExchangeOperator) GetInbounds() []chan arrow.RecordBatch {
	return nil
}

func (op *RemoteExchangeOperator) String() string {
	return "RemoteExchangeOperator"
}
