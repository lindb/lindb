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
	"errors"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// ── MAP_VALUES ────────────────────────────────────────────────────────────────

// mapValuesInstance holds the target key set, computed once at construction time.
type mapValuesInstance struct {
	arg        function.Expr
	targetKeys map[string]struct{}
}

// MapValuesFactory creates a mapValuesInstance with constant key args cached.
var MapValuesFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	targetKeys := make(map[string]struct{}, len(args)-1)
	for i := 1; i < len(args); i++ {
		if keyScalar, err := args[i].EvalScalar(); err == nil {
			targetKeys[scalar.ToString(keyScalar)] = struct{}{}
		}
	}
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &mapValuesInstance{arg: arg, targetKeys: targetKeys}
}

func (f *mapValuesInstance) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("map_values is not supported in scalar execution")
}

func (f *mapValuesInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	if f.arg == nil {
		return nil, errors.New("map_values requires at least 1 argument")
	}
	input, err := f.arg.Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	inputMap, ok := input.(*larray.Map)
	if !ok {
		return nil, fmt.Errorf("map_values: input must be a map type, got %T", input)
	}

	keys := inputMap.Keys().(*array.String)
	values := inputMap.Items().(*array.String)
	offsets := inputMap.Offsets()

	mb := array.NewMapBuilder(memory.DefaultAllocator, arrow.BinaryTypes.String, arrow.BinaryTypes.String, false)
	defer mb.Release()
	keysBuilder := mb.KeyBuilder().(*array.StringBuilder)
	valuesBuilder := mb.ItemBuilder().(*array.StringBuilder)

	mb.Reserve(inputMap.Len())
	for i := range inputMap.Len() {
		if inputMap.IsNull(i) {
			mb.AppendNull()
			continue
		}
		mb.Append(true)
		row := inputMap.Row(i)
		start, end := offsets[row], offsets[row+1]
		for j := int(start); j < int(end); j++ {
			key := keys.Value(j)
			if _, exists := f.targetKeys[key]; exists {
				keysBuilder.Append(key)
				if values.IsNull(j) {
					valuesBuilder.AppendNull()
				} else {
					valuesBuilder.Append(values.Value(j))
				}
			}
		}
	}
	return mb.NewArray(), nil
}
