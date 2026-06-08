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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

// Comparison evaluates `left op right` element-wise, producing a boolean Arrow
// array with one entry per row in the input RecordBatch.
// Supported types: int64, float64, string (for the common SQL comparison operators).
type Comparison struct {
	left  Expression
	right Expression
	op    tree.ComparisonOperator
}

// NewComparison constructs a Comparison expression.
func NewComparison(_ EvalContext, op tree.ComparisonOperator, left, right Expression) Expression {
	return &Comparison{left: left, right: right, op: op}
}

func (c *Comparison) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("comparison: scalar evaluation not supported")
}

func (c *Comparison) ResultType() ResultType {
	return Array
}

func (c *Comparison) String() string {
	return fmt.Sprintf("(%s %s %s)", c.left, c.op, c.right)
}

// Eval evaluates the comparison for every row in the batch, returning a
// *array.Boolean whose i-th entry is true when row i satisfies the predicate.
func (c *Comparison) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	leftArr, err := c.left.Eval(record)
	if err != nil {
		return nil, fmt.Errorf("comparison: left eval: %w", err)
	}
	rightArr, err := c.right.Eval(record)
	if err != nil {
		return nil, fmt.Errorf("comparison: right eval: %w", err)
	}

	n := int(record.NumRows())
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(n)

	var buildErr error
	switch l := leftArr.(type) {
	case *array.Int64:
		switch r := rightArr.(type) {
		case *array.Int64:
			for i := range n {
				b.Append(compareInt64(l.Value(i), r.Value(i), c.op))
			}
		case *array.Float64:
			// Symmetric widening: int64 left vs float64 right.
			for i := range n {
				b.Append(compareFloat64(float64(l.Value(i)), r.Value(i), c.op))
			}
		default:
			buildErr = fmt.Errorf("comparison: unsupported right type %T for int64 left", rightArr)
		}
	case *array.Float64:
		switch r := rightArr.(type) {
		case *array.Float64:
			for i := range n {
				b.Append(compareFloat64(l.Value(i), r.Value(i), c.op))
			}
		case *array.Int64:
			// Numeric widening: literal integers (e.g. 2000 in HAVING count(*) > 2000)
			// are stored as Int64 even when the left side is Float64; convert implicitly.
			for i := range n {
				b.Append(compareFloat64(l.Value(i), float64(r.Value(i)), c.op))
			}
		default:
			buildErr = fmt.Errorf("comparison: unsupported right type %T for float64 left", rightArr)
		}
	case *array.String:
		switch r := rightArr.(type) {
		case *array.String:
			for i := range n {
				b.Append(compareString(l.Value(i), r.Value(i), c.op))
			}
		default:
			buildErr = fmt.Errorf("comparison: unsupported right type %T for string left", rightArr)
		}
	default:
		buildErr = fmt.Errorf("comparison: unsupported left type %T", leftArr)
	}

	if buildErr != nil {
		return nil, buildErr
	}
	return b.NewArray(), nil
}

func compareInt64(l, r int64, op tree.ComparisonOperator) bool {
	switch op {
	case tree.ComparisonEQ:
		return l == r
	case tree.ComparisonNEQ, tree.ComparisonOperator("<>"):
		return l != r
	case tree.ComparisonGT:
		return l > r
	case tree.ComparisonGTE:
		return l >= r
	case tree.ComparisonLT:
		return l < r
	case tree.ComparisonLTE:
		return l <= r
	default:
		return false
	}
}

func compareFloat64(l, r float64, op tree.ComparisonOperator) bool {
	switch op {
	case tree.ComparisonEQ:
		return l == r
	case tree.ComparisonNEQ, tree.ComparisonOperator("<>"):
		return l != r
	case tree.ComparisonGT:
		return l > r
	case tree.ComparisonGTE:
		return l >= r
	case tree.ComparisonLT:
		return l < r
	case tree.ComparisonLTE:
		return l <= r
	default:
		return false
	}
}

func compareString(l, r string, op tree.ComparisonOperator) bool {
	switch op {
	case tree.ComparisonEQ:
		return l == r
	case tree.ComparisonNEQ, tree.ComparisonOperator("<>"):
		return l != r
	case tree.ComparisonGT:
		return l > r
	case tree.ComparisonGTE:
		return l >= r
	case tree.ComparisonLT:
		return l < r
	case tree.ComparisonLTE:
		return l <= r
	default:
		return false
	}
}
