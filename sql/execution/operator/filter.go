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
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"

	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// FilterOperator applies a boolean predicate to each RecordBatch coming from
// child, forwarding only the rows for which the predicate evaluates to true.
// It is used to implement the HAVING clause (post-aggregation filter).
type FilterOperator struct {
	ctx     context.Context
	child   Operator
	rawPred tree.Expression      // AST predicate, rewritten lazily on first Run()
	pred    expression.Expression // compiled expression after rewrite
	inbound *Queue
}

// NewFilterOperator creates a FilterOperator that wraps child and filters
// every RecordBatch with the given AST predicate.
func NewFilterOperator(ctx context.Context, child Operator, pred tree.Expression) Operator {
	return &FilterOperator{
		ctx:     ctx,
		child:   child,
		rawPred: pred,
		inbound: NewQueue(make(chan arrow.RecordBatch, 1024)),
	}
}

func (op *FilterOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	if op.pred == nil {
		op.prepare()
	}

	for {
		record, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		filtered, err := op.filter(record)
		if err != nil {
			panic(fmt.Sprintf("FilterOperator: %v", err))
		}
		// Only forward non-empty batches to avoid polluting downstream operators.
		if filtered.NumRows() > 0 {
			output <- filtered
		}
	}
}

// filter evaluates the predicate against the batch and returns a new
// RecordBatch containing only the rows where the predicate is true.
func (op *FilterOperator) filter(record arrow.RecordBatch) (arrow.RecordBatch, error) {
	predArr, err := op.pred.Eval(record)
	if err != nil {
		return nil, fmt.Errorf("predicate eval: %w", err)
	}
	boolArr, ok := predArr.(*array.Boolean)
	if !ok {
		return nil, fmt.Errorf("predicate must return *array.Boolean, got %T", predArr)
	}

	// Collect row indices that pass the predicate.
	// Per SQL three-valued logic, NULL evaluates to UNKNOWN which is treated
	// as false — a row with a NULL predicate result is excluded.
	rows := int(record.NumRows())
	passing := make([]int, 0, rows)
	for i := range rows {
		if !boolArr.IsNull(i) && boolArr.Value(i) {
			passing = append(passing, i)
		}
	}

	if len(passing) == 0 {
		// Return an empty batch with the same schema.
		empty := make([]arrow.Array, record.NumCols())
		for i := range int(record.NumCols()) {
			empty[i] = BuildEmptyColumn(record.Column(i))
		}
		return larrow.NewFilterableRecord(
			array.NewRecordBatch(record.Schema(), empty, 0), nil), nil
	}

	// Reconstruct each column by selecting only the passing rows.
	cols := make([]arrow.Array, record.NumCols())
	for i := range int(record.NumCols()) {
		cols[i], err = SelectRows(record.Column(i), passing)
		if err != nil {
			return nil, fmt.Errorf("column %d row selection: %w", i, err)
		}
	}
	result := array.NewRecordBatch(record.Schema(), cols, int64(len(passing)))
	return larrow.NewFilterableRecord(result, nil), nil
}

// GetLayout proxies the child's output layout — FilterOperator does not change
// the schema.
func (op *FilterOperator) GetLayout() []*plan.Symbol {
	return op.child.GetLayout()
}

func (op *FilterOperator) Children() []Operator {
	return []Operator{op.child}
}

func (op *FilterOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{op.inbound.GetInbound()}
}

func (op *FilterOperator) String() string {
	return "FilterOperator"
}

// prepare compiles the AST predicate into an executable expression using the
// child's output layout as the column mapping.
func (op *FilterOperator) prepare() {
	op.pred = expression.Rewrite(&expression.RewriteContext{
		SourceLayout: op.child.GetLayout(),
		EvalContext:  expression.NewEvalContext(op.ctx),
	}, op.rawPred)
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers for row selection (no arrow/compute dependency)
// ─────────────────────────────────────────────────────────────────────────────

// SelectRows builds a new Arrow array that contains only the rows at the given
// indices. Supports the common types emitted by aggregations (int64, float64,
// string, bool).
func SelectRows(col arrow.Array, rows []int) (arrow.Array, error) {
	alloc := memory.DefaultAllocator
	switch src := col.(type) {
	case *array.Int64:
		b := array.NewInt64Builder(alloc)
		defer b.Release()
		b.Reserve(len(rows))
		for _, i := range rows {
			if src.IsNull(i) {
				b.AppendNull()
			} else {
				b.UnsafeAppend(src.Value(i))
			}
		}
		return b.NewArray(), nil
	case *array.Float64:
		b := array.NewFloat64Builder(alloc)
		defer b.Release()
		b.Reserve(len(rows))
		for _, i := range rows {
			if src.IsNull(i) {
				b.AppendNull()
			} else {
				b.UnsafeAppend(src.Value(i))
			}
		}
		return b.NewArray(), nil
	case *array.String:
		b := array.NewStringBuilder(alloc)
		defer b.Release()
		b.Reserve(len(rows))
		for _, i := range rows {
			if src.IsNull(i) {
				b.AppendNull()
			} else {
				b.Append(src.Value(i))
			}
		}
		return b.NewArray(), nil
	case *array.Boolean:
		b := array.NewBooleanBuilder(alloc)
		defer b.Release()
		b.Reserve(len(rows))
		for _, i := range rows {
			if src.IsNull(i) {
				b.AppendNull()
			} else {
				b.Append(src.Value(i))
			}
		}
		return b.NewArray(), nil
	default:
		return nil, fmt.Errorf("selectRows: unsupported column type %T", col)
	}
}

// BuildEmptyColumn creates a zero-length array of the same type as the source.
func BuildEmptyColumn(col arrow.Array) arrow.Array {
	result, _ := SelectRows(col, nil)
	if result == nil {
		// Fallback for unknown types: re-use the column itself (a 0-row batch
		// is discarded downstream anyway).
		return col
	}
	return result
}
