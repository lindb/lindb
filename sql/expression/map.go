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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
)

type mapValuesFunc struct {
	ctx  EvalContext
	args []Expression

	targetKeys map[string]struct{}

	// FIXME: add expression clear func
	mb *array.MapBuilder
}

func newMapValuesFunc(ctx EvalContext, args []Expression) Func {
	targetKeys := make(map[string]struct{}, len(args)-1)
	for i := 1; i < len(args); i++ {
		key, err := args[i].EvalScalar()
		if err != nil {
			panic(fmt.Sprintf("failed to evaluate key argument in map_values function: %v", err))
		}
		targetKeys[scalar.ToString(key)] = struct{}{}
	}

	if len(targetKeys) == 0 {
		panic("map_values function requires at least one key argument")
	}
	mb := array.NewMapBuilder(memory.DefaultAllocator, arrow.BinaryTypes.String, arrow.BinaryTypes.String, false)

	return &mapValuesFunc{
		ctx:        ctx,
		args:       args,
		targetKeys: targetKeys,
		mb:         mb,
	}
}

func (n *mapValuesFunc) EvalScalar() (scalar.Scalar, error) {
	panic("map_values is not supported in scalar execution")
}

func (n *mapValuesFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	input, err := n.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release() // release the input array after processing

	inputMap, ok := input.(*array.Map)
	if !ok {
		return nil, fmt.Errorf("input of map_values should be map type, but got %T", input)
	}

	keys := inputMap.Keys().(*array.String)
	values := inputMap.Items().(*array.String)
	offsets := inputMap.Offsets()

	mb := n.mb
	// defer mb.Release()

	keysBuilder := mb.KeyBuilder().(*array.StringBuilder)
	valuesBuilder := mb.ItemBuilder().(*array.StringBuilder)

	mb.Reserve(inputMap.Len())

	for row := 0; row < inputMap.Len(); row++ {
		if inputMap.IsNull(row) {
			mb.AppendNull()
			continue
		}
		mb.Append(true)

		start, end := offsets[row], offsets[row+1]
		for i := int(start); i < int(end); i++ {
			// if keys.IsNull(i) {
			// 	continue
			// }
			key := keys.Value(i)

			// check if the key is in the target keys, if yes, append the key and value to the builder
			if _, exists := n.targetKeys[key]; exists {
				keysBuilder.Append(key)
				if values.IsNull(i) {
					valuesBuilder.AppendNull()
				} else {
					valuesBuilder.Append(values.Value(i))
				}
			}
		}
	}
	return mb.NewArray(), nil
}
