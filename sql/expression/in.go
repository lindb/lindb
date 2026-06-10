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

package expression

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
)

// In evaluates `value IN (candidate1, candidate2, …)` element-wise,
// producing a *array.Boolean with one entry per row in the RecordBatch.
// Each candidate is evaluated as a scalar constant once per Eval call.
// Supported value types: int64, float64, string.
type In struct {
	value      Expression
	candidates []Expression
}

// NewIn constructs an In expression.
func NewIn(_ EvalContext, value Expression, candidates []Expression) Expression {
	return &In{value: value, candidates: candidates}
}

func (in *In) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("in: scalar evaluation not supported")
}

func (in *In) ResultType() ResultType {
	return Array
}

func (in *In) String() string {
	parts := make([]string, len(in.candidates))
	for i, c := range in.candidates {
		parts[i] = c.String()
	}
	return fmt.Sprintf("(%s IN (%s))", in.value, strings.Join(parts, ", "))
}

// Eval checks, for every row in record, whether the value column is contained
// in the candidates list. Candidates are expected to be scalar constants;
// they are evaluated once against the batch to obtain their values.
func (in *In) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	valueArr, err := in.value.Eval(record)
	if err != nil {
		return nil, fmt.Errorf("in: value eval: %w", err)
	}

	n := int(record.NumRows())
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(n)

	switch col := valueArr.(type) {
	case *array.String:
		set, err := in.buildStringSet(record)
		if err != nil {
			return nil, err
		}
		for i := range n {
			_, ok := set[col.Value(i)]
			b.Append(ok)
		}
	case *array.Int64:
		set, err := in.buildInt64Set(record)
		if err != nil {
			return nil, err
		}
		for i := range n {
			_, ok := set[col.Value(i)]
			b.Append(ok)
		}
	case *array.Float64:
		set, err := in.buildFloat64Set(record)
		if err != nil {
			return nil, err
		}
		for i := range n {
			_, ok := set[col.Value(i)]
			b.Append(ok)
		}
	default:
		return nil, fmt.Errorf("in: unsupported value type %T", valueArr)
	}

	return b.NewArray(), nil
}

// buildStringSet evaluates each candidate against the batch and collects
// their string values into a set for O(1) lookup.
func (in *In) buildStringSet(record arrow.RecordBatch) (map[string]struct{}, error) {
	set := make(map[string]struct{}, len(in.candidates))
	for _, c := range in.candidates {
		arr, err := c.Eval(record)
		if err != nil {
			return nil, fmt.Errorf("in: candidate eval: %w", err)
		}
		switch a := arr.(type) {
		case *array.String:
			for i := range a.Len() {
				set[a.Value(i)] = struct{}{}
			}
		default:
			return nil, fmt.Errorf("in: candidate type %T cannot be compared to string column", arr)
		}
	}
	return set, nil
}

// buildInt64Set evaluates each candidate and collects int64 values.
func (in *In) buildInt64Set(record arrow.RecordBatch) (map[int64]struct{}, error) {
	set := make(map[int64]struct{}, len(in.candidates))
	for _, c := range in.candidates {
		arr, err := c.Eval(record)
		if err != nil {
			return nil, fmt.Errorf("in: candidate eval: %w", err)
		}
		switch a := arr.(type) {
		case *array.Int64:
			for i := range a.Len() {
				set[a.Value(i)] = struct{}{}
			}
		default:
			return nil, fmt.Errorf("in: candidate type %T cannot be compared to int64 column", arr)
		}
	}
	return set, nil
}

// buildFloat64Set evaluates each candidate and collects float64 values.
func (in *In) buildFloat64Set(record arrow.RecordBatch) (map[float64]struct{}, error) {
	set := make(map[float64]struct{}, len(in.candidates))
	for _, c := range in.candidates {
		arr, err := c.Eval(record)
		if err != nil {
			return nil, fmt.Errorf("in: candidate eval: %w", err)
		}
		switch a := arr.(type) {
		case *array.Float64:
			for i := range a.Len() {
				set[a.Value(i)] = struct{}{}
			}
		case *array.Int64:
			// Allow int64 literals (e.g. IN (1, 2)) to match float64 columns.
			for i := range a.Len() {
				set[float64(a.Value(i))] = struct{}{}
			}
		default:
			return nil, fmt.Errorf("in: candidate type %T cannot be compared to float64 column", arr)
		}
	}
	return set, nil
}
