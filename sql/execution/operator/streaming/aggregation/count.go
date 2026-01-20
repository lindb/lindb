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

package aggregation

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
)

func newCountAggregator(ctx expression.EvalContext, args []expression.Expression) Aggregator {
	return &countAggregator{}
}

type countAggregator struct {
	value int64
}

func (c *countAggregator) Enter(row types.Row) {
	c.value++
}

func (c *countAggregator) Flush(column *types.Column) {
	column.Append(c.value)

	// need reset value after flush
	c.value = 0
}
