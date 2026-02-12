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

package operator

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ValuesOperator struct {
	node *plan.ValuesNode
}

func NewValuesOperator(node *plan.ValuesNode) Operator {
	return &ValuesOperator{
		node: node,
	}
}

func (op *ValuesOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	var page *types.Page
	node := op.node
	if node.Rows != nil {
		page = node.Rows
	} else if node.RowCount == 1 {
		page = types.RowWithEmptyValue
	}
	panic(page)
	// FIXME: output <- page
}

func (op *ValuesOperator) GetLayout() []*plan.Symbol {
	return op.node.GetOutputSymbols()
}

func (op *ValuesOperator) Children() []Operator {
	return nil
}

func (op *ValuesOperator) GetInbounds() []chan arrow.RecordBatch {
	return nil
}

func (op *ValuesOperator) String() string {
	return "ValuesOperator"
}
