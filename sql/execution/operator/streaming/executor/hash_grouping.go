package executor

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator/streaming/aggregation"
	"github.com/lindb/lindb/sql/planner/plan"
)

type aggregators []aggregation.Aggregator

type HashGrouping struct {
	node         *plan.AggregationNode
	colIdxOfKeys []int // coloumn index of grouping keys
	aggregators  []aggregation.Aggregator

	assignments []*plan.Assignment // assignments of projection

	grouping map[string]aggregators
}

func NewHashGrouping(node *plan.AggregationNode, assignments []*plan.Assignment) Executor {
	sourceLayout := node.Source.GetOutputSymbols()
	groupingKeys := node.GetGroupingKeys()
	colIdxOfKeys := make([]int, len(groupingKeys))
	for index, groupingKey := range groupingKeys {
		_, colIdx, ok := lo.FindIndexOf(sourceLayout, func(item *plan.Symbol) bool {
			return item.Name == groupingKey.Name
		})
		if ok {
			colIdxOfKeys[index] = colIdx
		} else {
			panic("grouping keys not match")
		}
	}

	fmt.Printf("output=======> %v\n", node.Outputs)

	return &HashGrouping{
		node:         node,
		colIdxOfKeys: colIdxOfKeys,
		assignments:  assignments,
		aggregators:  createAggregators(node),
		grouping:     make(map[string]aggregators),
	}
}

func (g *HashGrouping) Enter(page *types.Page) {
	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		for _, colIdx := range g.colIdxOfKeys {
			_ = row.Get(colIdx)
			// fmt.Printf("column value===> %v=%v %v\n", key.Name, key.DataType, column)
		}

		for _, agg := range g.aggregators {
			agg.Enter(row)
		}
	}
	// fmt.Println("enter....")
}

func (g *HashGrouping) Leave(output chan<- *types.Page) {
	// fmt.Println("leave....")
	newPage := types.NewPage()
	outputColumns := make([]*types.Column, len(g.node.Outputs))
	for i, assign := range g.node.Outputs {
		outputColumns[i] = types.NewColumn()
		newPage.AppendColumn(types.NewColumnInfo(assign.Name, assign.DataType), outputColumns[i])
	}

	for _, agg := range g.aggregators {
		agg.Flush(outputColumns[len(g.node.Outputs)-1])
	}

	output <- newPage
}

func createAggregators(node *plan.AggregationNode) []aggregation.Aggregator {
	var aggregators []aggregation.Aggregator
	for _, agg := range node.Aggregations {
		aggregator, err := aggregation.CreateAggregator(agg.Aggregation.Function, agg.Aggregation.Arguments)
		if err != nil {
			fmt.Printf("create agg error=%v\n", err)
			continue
		}
		aggregators = append(aggregators, aggregator)
	}
	// FIXME: check keys==grouping keys
	return aggregators
}
