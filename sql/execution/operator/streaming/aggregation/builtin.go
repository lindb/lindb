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

	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/function"
	_ "github.com/lindb/lindb/sql/function/builtin" // register built-in incremental agg
	"github.com/lindb/lindb/sql/tree"
)

type NewAggregator func(ctx expression.EvalContext, args []expression.Expression) Aggregator

// legacyFuncs holds aggregators not yet migrated to the unified registry (e.g. Sampling).
var legacyFuncs = map[tree.FuncName]NewAggregator{
	tree.Sampling: newSampllingAggregator,
}

// CreateAggregator returns an Aggregator for the named function.
// Functions registered in function.DefaultRegistry (sum, count, min, max, first, last)
// are resolved through the unified IncrementalAgg interface;
// all others fall back to the legacy factory map.
func CreateAggregator(ctx expression.EvalContext, name tree.FuncName, args []expression.Expression) (Aggregator, error) {
	if incremental, ok := function.DefaultRegistry.Incremental(name); ok {
		// expression.Expression satisfies function.Expr via structural typing.
		fArgs := make([]function.Expr, len(args))
		for i, a := range args {
			fArgs[i] = a
		}
		return &accumulatorAdapter{
			acc: incremental.NewAccumulator(ctx, fArgs),
		}, nil
	}
	factory, ok := legacyFuncs[name]
	if !ok {
		return nil, fmt.Errorf("not support %s", name)
	}
	return factory(ctx, args), nil
}
