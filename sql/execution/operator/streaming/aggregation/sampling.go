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
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/arrow/pkg/model"

	"github.com/lindb/lindb/sql/expression"
)

func newSampllingAggregator(ctx expression.EvalContext, args []expression.Expression) Aggregator {
	sampling := &samplingAggregator{
		args: args,
		ctx:  ctx,
	}
	sampling.indexes.traceID = -1
	sampling.indexes.spanID = -1
	sampling.indexes.duration = -1
	return sampling
}

type samplingAggregator struct {
	ctx  expression.EvalContext
	args []expression.Expression

	initialized bool

	indexes struct {
		traceID  int
		spanID   int
		duration int
	}

	traceID  *larray.Generic[[]byte]
	spanID   *larray.Generic[[]byte]
	duration *larray.Generic[arrow.Duration]

	value *model.Exemplar
}

func (a *samplingAggregator) Initialize(record arrow.RecordBatch) {
	if !a.initialized {
		schema := record.Schema()
		for _, arg := range a.args {
			argStr := arg.String()
			fIndexes := schema.FieldIndices(argStr)
			if len(fIndexes) != 1 {
				// invalid argument, skip this argument
				panic(fmt.Sprintf("invalid argument %s, found %d fields", argStr, len(fIndexes)))
			}
			switch {
			case strings.Contains(argStr, "trace"):
				a.indexes.traceID = fIndexes[0]
			case strings.Contains(argStr, "span"):
				a.indexes.spanID = fIndexes[0]
			case strings.Contains(argStr, "duration"):
				a.indexes.duration = fIndexes[0]
			}
		}

		a.value = &model.Exemplar{}

		a.initialized = true
	}
	if a.indexes.traceID >= 0 {
		a.traceID = record.Column(a.indexes.traceID).(*larray.Generic[[]byte])
	}
	if a.indexes.spanID >= 0 {
		a.spanID = record.Column(a.indexes.spanID).(*larray.Generic[[]byte])
	}
	if a.indexes.duration >= 0 {
		a.duration = record.Column(a.indexes.duration).(*larray.Generic[arrow.Duration])
	}
}

func (a *samplingAggregator) Enter(record arrow.RecordBatch, row int) {
	if a.traceID == nil || a.spanID == nil {
		// invalid exemplar func
		return
	}
	var duration int64
	if a.duration != nil {
		duration = int64(a.duration.Value(row))
	}
	if a.value.TraceID != nil && duration <= a.value.Duration {
		return
	}

	a.value.Duration = duration
	a.value.TraceID = a.traceID.Value(row)
	a.value.SpanID = a.spanID.Value(row)
}

func (a *samplingAggregator) Flush(builder array.Builder) {
	if a.value == nil || a.value.TraceID == nil {
		// no valid exemplar, append null value
		builder.AppendNull()
		return
	}

	sb := builder.(*array.ExtensionBuilder)
	exemplarBuilder := larray.NewExemplarBuilder(sb)
	exemplarBuilder.Append(a.value)
	// need reset value after flush
	a.value.Reset()
}
