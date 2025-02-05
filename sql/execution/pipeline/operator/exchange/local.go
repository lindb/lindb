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
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type LocalExchangeOperator struct {
	inbound chan *types.Page
	node    *plan.ExchangeNode
}

// FIXME: source op
func NewLocalExchangeOperator(node *plan.ExchangeNode) operator.Operator {
	return &LocalExchangeOperator{
		node:    node,
		inbound: make(chan *types.Page),
	}
}

// Finish implements operator.Operator.
func (l *LocalExchangeOperator) Run(output chan<- *types.Page) {
	for page := range l.inbound {
		output <- page
	}
}

func (l *LocalExchangeOperator) Receive(page *types.Page) {
	l.inbound <- page
}

func (l *LocalExchangeOperator) GetLayout() []*plan.Symbol {
	return l.node.GetOutputSymbols()
}
