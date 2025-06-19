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

package output

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ResultSetOutputOperator struct {
	node  *plan.OutputNode
	child operator.Operator
}

func NewRSOutputOperator(node *plan.OutputNode, child operator.Operator) operator.Operator {
	return &ResultSetOutputOperator{
		node:  node,
		child: child,
	}
}

// AddInput implements operator.Operator
func (op *ResultSetOutputOperator) Run(ctx context.Context, output chan<- *types.Page) {
	inbound := make(chan *types.Page)

	go func() {
		defer close(inbound)

		op.child.Run(ctx, inbound)
	}()

	rebuildPage := false
	layout := op.node.GetOutputSymbols()

	sourceLayout := make(map[string]int)
	lo.ForEach(op.child.GetLayout(), func(symbol *plan.Symbol, index int) {
		sourceLayout[symbol.Name] = index
	})
	columnNames := op.node.ColumnNames

	for idx, symbol := range layout {
		sourceIdx, ok := sourceLayout[symbol.Name]
		if ok && (idx != sourceIdx || (len(columnNames) > 0 && columnNames[idx] != symbol.Name)) {
			rebuildPage = true
			break
		}
	}

	// process child output
	for page := range inbound {
		if page == nil || page.NumRows() == 0 {
			fmt.Printf("add empty page====%v\n", string(encoding.JSONMarshal(page)))
			return
		}
		if rebuildPage {
			targetPage := types.NewPage()
			for colIdx, col := range layout {
				if idx, ok := sourceLayout[col.Name]; ok {
					column := page.Layout[idx]
					if len(columnNames) > 0 {
						column.Name = columnNames[colIdx]
					}
					targetPage.AppendColumn(column, page.Columns[idx])
				}
			}
			output <- targetPage
		} else {
			output <- page
		}
	}
}

func (op *ResultSetOutputOperator) GetLayout() []*plan.Symbol {
	panic("result set output operator should not get layout")
}

func (op *ResultSetOutputOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}
