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

package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/arrow/builder"
	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/execution/grouping"
	"github.com/lindb/lindb/sql/execution/operator/streaming/aggregation"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

type aggregators []aggregation.Aggregator

// HashGrouping implements hash-based grouping aggregation for streaming data.
// It maintains a hash map of groups and their associated aggregators.
type HashGrouping struct {
	node                 *plan.AggregationNode
	sourceLayout         []*plan.Symbol
	colIdxOfGroupingKeys []int // coloumn index of grouping keys

	assignments []*plan.Assignment // assignments of projection

	grouping map[string]aggregators
	rules    []grouping.Rule

	buf    *grouping.Buffer
	mapper *grouping.StringMapper
}

// NewHashGrouping creates a new HashGrouping executor.
// It initializes grouping rules for each grouping key based on their data types,
// and sets up the internal structures for hash-based aggregation.
func NewHashGrouping(node *plan.AggregationNode, assignments []*plan.Assignment) Executor {
	sourceLayout := node.Source.GetOutputSymbols()
	groupingKeys := node.GetGroupingKeys()
	colIdxOfKeys := make([]int, len(groupingKeys))
	rules := make([]grouping.Rule, len(groupingKeys))
	mapper := grouping.NewStringMapper()
	for index, groupingKey := range groupingKeys {
		column, colIdx, ok := lo.FindIndexOf(sourceLayout, func(item *plan.Symbol) bool {
			return item.Name == groupingKey.Name
		})
		if ok {
			colIdxOfKeys[index] = colIdx
			rule, err := grouping.CreateRule(column.DataType, mapper)
			if err != nil {
				panic(err)
			}
			rules[index] = rule
		} else {
			panic("grouping keys not match")
		}
	}

	exec := &HashGrouping{
		node:                 node,
		sourceLayout:         sourceLayout,
		colIdxOfGroupingKeys: colIdxOfKeys,
		assignments:          assignments,
		grouping:             make(map[string]aggregators),
		rules:                rules,
		buf:                  grouping.NewBuffer(),
		mapper:               mapper,
	}
	return exec
}

// Enter processes an incoming record of data by extracting grouping keys,
// computing hash keys, and feeding rows into the appropriate aggregators.
// Each row is assigned to a group based on its grouping key values.
func (g *HashGrouping) Enter(record arrow.RecordBatch) {
	fmt.Println(record)
	numOfRows := int(record.NumRows())
	for row := range numOfRows {
		for i, colIdx := range g.colIdxOfGroupingKeys {
			column := record.Column(colIdx)
			g.rules[i].Map(column, row, g.buf)
		}

		data := encoding.U32SliceToBytes(g.buf.GetData())
		key := strutil.ByteSlice2String(data)
		aggregators, ok := g.grouping[key]
		if !ok {
			aggregators = g.createAggregators(record)
			// need clone string(key reuse)
			g.grouping[strings.Clone(key)] = aggregators
		}

		g.buf.Reset()

		for _, agg := range aggregators {
			agg.Enter(record, row)
		}
	}
}

// Leave finalizes the aggregation and outputs the results.
// It iterates through all groups, reconstructs the grouping key values,
// flushes the aggregation results, and sends the output page to the channel.
func (g *HashGrouping) Leave(output chan<- arrow.RecordBatch) {
	// TODO: create new grouping map???
	fields := lo.Map(g.node.GetOutputSymbols(), func(symbol *plan.Symbol, _ int) arrow.Field {
		return arrow.Field{Name: symbol.Name, Type: symbol.DataType}
	})
	rb := builder.NewRecordBuilder(memory.DefaultAllocator, arrow.NewSchema(fields, nil))
	defer rb.Release()

	builders := rb.Fields()

	for key, aggs := range g.grouping {
		data := encoding.BytesToU32Slice(strutil.String2ByteSlice(key))
		g.buf.ResetWithData(data)

		outputIndex := 0
		for i := range g.colIdxOfGroupingKeys {
			g.rules[i].Unmap(builders[outputIndex], g.buf)
			outputIndex++
		}

		for _, agg := range aggs {
			agg.Flush(builders[outputIndex])
			outputIndex++
		}
	}

	output <- larrow.NewFilterableRecord(rb.NewRecord(), nil)
}

// createAggregators creates aggregator instances for each aggregation function
// defined in the aggregation node. Each aggregator is responsible for computing
// one aggregate function (e.g., SUM, COUNT, AVG).
func (g *HashGrouping) createAggregators(record arrow.RecordBatch) []aggregation.Aggregator {
	var aggregators []aggregation.Aggregator
	ctx := expression.NewEvalContext(context.TODO())
	for _, agg := range g.node.Aggregations {
		args := make([]expression.Expression, len(agg.Aggregation.Arguments))
		for i, arg := range agg.Aggregation.Arguments {
			args[i] = expression.Rewrite(&expression.RewriteContext{
				SourceLayout: g.node.Source.GetOutputSymbols(),
				EvalContext:  ctx,
			}, arg)
		}
		// TODO: extract eval context
		aggregator, err := aggregation.CreateAggregator(ctx,
			agg.Aggregation.Function, args)
		if err != nil {
			panic(err)
		}
		aggregator.Initialize(record)
		aggregators = append(aggregators, aggregator)
	}
	// FIXME: check keys==grouping keys
	return aggregators
}
