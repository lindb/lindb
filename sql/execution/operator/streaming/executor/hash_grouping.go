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
	"strings"

	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/grouping"
	"github.com/lindb/lindb/sql/execution/operator/streaming/aggregation"
	"github.com/lindb/lindb/sql/planner/plan"
)

type aggregators []aggregation.Aggregator

type HashGrouping struct {
	node          *plan.AggregationNode
	outputColumns []types.ColumnMetadata
	sourceLayout  []*plan.Symbol
	colIdxOfKeys  []int // coloumn index of grouping keys

	assignments []*plan.Assignment // assignments of projection

	grouping map[string]aggregators
	rules    []grouping.Rule

	buf    *grouping.Buffer
	mapper *grouping.StringMapper
}

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
		node:          node,
		outputColumns: createOutputs(node),
		sourceLayout:  sourceLayout,
		colIdxOfKeys:  colIdxOfKeys,
		assignments:   assignments,
		grouping:      make(map[string]aggregators),
		rules:         rules,
		buf:           grouping.NewBuffer(),
		mapper:        mapper,
	}
	return exec
}

func (g *HashGrouping) Enter(page *types.Page) {
	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		for i, colIdx := range g.colIdxOfKeys {
			val := row.Get(colIdx)
			// fmt.Printf("column value===> %v=%d=%d\n", val, colIdx, i)
			g.rules[i].Map(val, g.buf)
		}

		data := encoding.U32SliceToBytes(g.buf.GetData())
		key := strutil.ByteSlice2String(data)
		aggregators, ok := g.grouping[key]
		if !ok {
			aggregators = createAggregators(g.node)
			// need clone string(key reuse)
			g.grouping[strings.Clone(key)] = aggregators
		}

		g.buf.Reset()

		for _, agg := range aggregators {
			agg.Enter(row)
		}

	}
}

func (g *HashGrouping) Leave(output chan<- *types.Page) {
	// TODO: create new grouping map???
	newPage := types.NewPage()
	outputColumns := make([]*types.Column, len(g.outputColumns))
	for i, column := range g.outputColumns {
		outputColumns[i] = types.NewColumn()
		newPage.AppendColumn(column, outputColumns[i])
	}

	for key, aggs := range g.grouping {
		data := encoding.BytesToU32Slice(strutil.String2ByteSlice(key))
		g.buf.ResetWithData(data)

		outputIndex := 0
		for i := range g.colIdxOfKeys {
			val := g.rules[i].Unmap(g.buf)
			outputColumns[outputIndex].Append(val)
			outputIndex++
		}

		for _, agg := range aggs {
			agg.Flush(outputColumns[outputIndex])
			outputIndex++
		}
	}

	output <- newPage
}

func createOutputs(node *plan.AggregationNode) (columns []types.ColumnMetadata) {
	for _, key := range node.GroupingSets.GroupingKeys {
		columns = append(columns, types.NewColumnInfo(key.Name, key.DataType))
	}

	for _, agg := range node.Aggregations {
		// TODO: add aggregation type
		columns = append(columns, types.NewColumnInfo(agg.Symbol.Name, agg.Symbol.DataType))
	}
	return columns
}

func createAggregators(node *plan.AggregationNode) []aggregation.Aggregator {
	var aggregators []aggregation.Aggregator
	for _, agg := range node.Aggregations {
		aggregator, err := aggregation.CreateAggregator(agg.Aggregation.Function, agg.Aggregation.Arguments)
		if err != nil {
			panic(err)
		}
		aggregators = append(aggregators, aggregator)
	}
	// FIXME: check keys==grouping keys
	return aggregators
}
