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

package builtin

import (
	"math/rand"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// ── RAND ──────────────────────────────────────────────────────────────────────

// randInstance is a per-query PRNG for RAND() / RAND(N).
// The seed is evaluated once at construction time.
type randInstance struct {
	rng *rand.Rand
}

// RandFactory creates a randInstance with the seed evaluated once at construction.
var RandFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	var seed int64
	if len(args) > 0 {
		if seedScalar, err := args[0].EvalScalar(); err == nil {
			seed = int64(scalarToFloat64(seedScalar))
		} else {
			seed = time.Now().UnixNano()
		}
	} else {
		seed = time.Now().UnixNano()
	}
	return &randInstance{
		// G404: math/rand is intentional — SQL RAND() does not require crypto-safe randomness.
		rng: rand.New(rand.NewSource(seed)), //nolint:gosec
	}
}

func (f *randInstance) EvalScalar() (scalar.Scalar, error) {
	return scalar.NewFloat64Scalar(f.rng.Float64()), nil
}

func (f *randInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	numRows := int(record.NumRows())
	builder := array.NewFloat64Builder(memory.DefaultAllocator)
	defer builder.Release()
	builder.Reserve(numRows)
	for range numRows {
		builder.UnsafeAppend(f.rng.Float64())
	}
	return builder.NewArray(), nil
}
