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
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
)

// randFunc implements MySQL-compatible RAND() / RAND(N) scalar function.
// RAND() returns a random float64 in [0, 1), re-seeded with the current time at construction.
// RAND(N) seeds the PRNG once with N, producing a deterministic sequence per query.
type randFunc struct {
	ctx     EvalContext // reserved; RAND() currently needs no query context
	rng     *rand.Rand
	builder *array.Float64Builder
}

func newRandFunc(ctx EvalContext, args []Expression) Func {
	var seed int64
	if len(args) > 0 {
		seedScalar, err := args[0].EvalScalar()
		if err != nil {
			panic(fmt.Sprintf("failed to evaluate seed argument in rand function: %v", err))
		}
		// Accept both integer and float seeds to match MySQL behavior (float is truncated).
		seed = int64(scalarToFloat64(seedScalar))
	} else {
		// Use current time as seed so each unseeded call produces distinct sequences.
		seed = time.Now().UnixNano()
	}
	return &randFunc{
		ctx: ctx,
		// Each randFunc instance owns its own PRNG to avoid global-state races.
		// G404: math/rand is intentional here — SQL RAND() does not require crypto-safe randomness.
		rng:     rand.New(rand.NewSource(seed)), //nolint:gosec
		builder: array.NewFloat64Builder(memory.DefaultAllocator),
	}
}

func (f *randFunc) EvalScalar() (scalar.Scalar, error) {
	return scalar.NewFloat64Scalar(f.rng.Float64()), nil
}

func (f *randFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	numRows := int(record.NumRows())
	builder := f.builder
	builder.Reserve(numRows)
	for i := 0; i < numRows; i++ {
		// UnsafeAppend skips the redundant Reserve(1) check inside Append — safe here
		// because Reserve(numRows) above guarantees sufficient capacity.
		builder.UnsafeAppend(f.rng.Float64())
	}
	return builder.NewArray(), nil
}
