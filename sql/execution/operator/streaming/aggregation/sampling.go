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
	"strings"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
)

func newSampllingAggregator(ctx expression.EvalContext, args []expression.Expression) Aggregator {
	sampling := &samplingAggregator{args: args}
	for _, arg := range args {
		argStr := arg.String()
		switch {
		case strings.Contains(argStr, "trace"):
			sampling.traceID = arg
		case strings.Contains(argStr, "span"):
			sampling.spanID = arg
		case strings.Contains(argStr, "duration"):
			sampling.duration = arg
		}
	}
	return sampling
}

type samplingAggregator struct {
	ctx  expression.EvalContext
	args []expression.Expression

	traceID  expression.Expression
	spanID   expression.Expression
	duration expression.Expression

	value *models.Exemplar
}

func (c *samplingAggregator) Enter(row types.Row) {
	if c.traceID == nil || c.spanID == nil {
		// invalid exemplar func
		return
	}
	var duration int64
	if c.duration != nil {
		duration, _, _ = c.duration.EvalInt(row)
	}
	if c.value != nil && duration <= c.value.Duration {
		return
	}
	traceID, _, _ := c.traceID.EvalString(row)
	spanID, _, _ := c.spanID.EvalString(row)

	if c.value == nil {
		c.value = &models.Exemplar{}
	}

	c.value.Duration = duration
	c.value.TraceID = traceID
	c.value.SpanID = spanID
}

func (c *samplingAggregator) Flush(column *types.Column) {
	column.Append(c.value)

	// need reset value after flush
	c.value = nil
}
