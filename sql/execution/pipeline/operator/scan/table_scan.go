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

package scan

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type TableScanOperator struct {
	connector spi.SourceConnector
	node      *plan.TableScanNode
}

func NewTableScanOperator(connector spi.SourceConnector, node *plan.TableScanNode) operator.Operator {
	return &TableScanOperator{
		connector: connector,
		node:      node,
	}
}

func (op *TableScanOperator) Run(output chan<- *types.Page) {
	op.connector.Run(output)
}

func (op *TableScanOperator) GetLayout() []*plan.Symbol {
	return op.node.GetOutputSymbols()
}

func (op *TableScanOperator) Children() []operator.Operator {
	return nil
}
