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

package streaming

import (
	"context"
	"time"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/execution/operator/streaming/executor"
	"github.com/lindb/lindb/sql/planner/plan"
)

type TimeWindowOperator struct {
	node  plan.PlanNode
	child operator.Operator

	executor executor.Executor

	inbound *operator.Queue

	ticker *time.Ticker
}

func NewTimeWindowOperator(node plan.PlanNode, executor executor.Executor, child operator.Operator) operator.Operator {
	ticker := time.NewTicker(time.Second * 5)
	return &TimeWindowOperator{
		node:     node,
		executor: executor,
		child:    child,
		inbound:  operator.NewQueue(make(chan arrow.RecordBatch, 1024)),
		ticker:   ticker,
	}
}

func (t *TimeWindowOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	defer t.ticker.Stop()

	for {
		select {
		case source := <-t.inbound.GetInbound():
			if source == nil {
				return
			}
			t.executor.Enter(source)
		case <-t.ticker.C:
			t.executor.Leave(output)
		}
	}
}

func (t *TimeWindowOperator) GetLayout() []*plan.Symbol {
	return t.node.GetOutputSymbols()
}

func (t *TimeWindowOperator) Children() []operator.Operator {
	return []operator.Operator{t.child}
}

func (t *TimeWindowOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{t.inbound.GetInbound()}
}

func (t *TimeWindowOperator) String() string {
	return "TimeWindowOperator"
}
