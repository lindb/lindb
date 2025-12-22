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

package processor

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/executor/aggregation"
	"github.com/lindb/lindb/streaming/query/context"
)

type AggregatorProcessor struct {
	ctx *context.QueryContext

	aggregations []*plan.AggregationAssignment
	aggregators  []aggregation.Aggregator
}

func NewAggregatorProcessor(ctx *context.QueryContext,
	aggregations []*plan.AggregationAssignment,
) Processor {
	p := &AggregatorProcessor{
		ctx:          ctx,
		aggregations: aggregations,
	}
	for _, agg := range p.aggregations {
		fmt.Printf("func:%v\n", agg.Aggregation.Function)
		aggregator, err := aggregation.NewAggregator(agg.Aggregation.Function, agg.Aggregation.Arguments)
		if err != nil {
			panic(err)
		}
		p.aggregators = append(p.aggregators, aggregator)
	}
	return p
}

// Process implements query.Processor.
func (a *AggregatorProcessor) Process(event any) {
	fmt.Printf("current:%v, aggregation event=%v\n", a.ctx.Query, event)
	for _, agg := range a.aggregators {
		agg.Process(event)
	}

	values := make([]any, len(a.aggregators))
	for i := range a.aggregators {
		values[i] = a.aggregators[i].GetValue()
	}

	// a.Next(values)
	fmt.Println(values)
}

func (a *AggregatorProcessor) Run(output chan<- any) {}

func (a *AggregatorProcessor) Children() []Processor {
	return nil
}

func (a *AggregatorProcessor) GetInbounds() []chan any {
	return nil
}

func (a *AggregatorProcessor) String() string {
	return "AggregatorProcessor"
}
