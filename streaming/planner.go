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
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/executor/aggregation"
	"github.com/lindb/lindb/streaming/query/input"
	"github.com/lindb/lindb/streaming/query/output"
	"github.com/lindb/lindb/streaming/query/processor"
	"github.com/lindb/lindb/streaming/stream"
)

type StreamingPlanner struct {
	planContext *PlanContext
}

func NewStreamingPlanner(planContext *PlanContext) *StreamingPlanner {
	return &StreamingPlanner{
		planContext: planContext,
	}
}

func (p *StreamingPlanner) Plan(node plan.PlanNode) processor.Processor {
	return node.Accept(nil, p).(processor.Processor)
}

func (p *StreamingPlanner) Visit(context any, n plan.PlanNode) (r any) {
	fmt.Printf("visit node:%T\n", n)
	switch node := n.(type) {
	case *plan.OutputNode:
		return node.Source.Accept(context, p)
	case *plan.InsertNode:
		return p.visitInsert(node)
	case *plan.FilterNode:
		return p.visitFilter(node)
	case *plan.AggregationNode:
		return p.visitAggregation(node)
	case *plan.ProjectionNode:
		return node.Source.Accept(context, p)
	case *plan.TableScanNode:
		return p.visitTable(node)
	default:
		panic(fmt.Sprintf("%T not supported for streaming", n))
	}
}

func (p *StreamingPlanner) visitInsert(node *plan.InsertNode) processor.Processor {
	child := node.Source.Accept(nil, p).(processor.Processor)
	return output.NewInsertInto(p.planContext.QueryContext, node, child)
}

func (p *StreamingPlanner) visitAggregation(node *plan.AggregationNode) processor.Processor {
	child := node.Source.Accept(nil, p).(processor.Processor)
	return processor.NewTimeWindowProcessor(p.planContext.QueryContext,
		aggregation.NewExecutor(p.planContext.QueryContext, node), child)
}

func (p *StreamingPlanner) visitFilter(node *plan.FilterNode) processor.Processor {
	child := node.Source.Accept(nil, p).(processor.Processor)
	return processor.NewFilterProcessor(p.planContext.QueryContext, node, child)
}

func (p *StreamingPlanner) visitTable(node *plan.TableScanNode) processor.Processor {
	table := node.Table.(*stream.TableHandle)
	streamName := table.Stream
	handle := p.planContext.QueryContext.InputManager.GetInputHandler(streamName)
	// create single stream receiver
	singleStream := input.NewSingleStreamReceiver(p.planContext.QueryContext, streamName)

	// subscribe event from stream
	handle.Subscribe(singleStream)
	return singleStream
}
