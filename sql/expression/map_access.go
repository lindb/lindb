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
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/scalar"
)

// MapAccess evaluates base[key] where base is an Arrow map column and key is a string scalar.
// For each row it looks up the key in the map and returns the corresponding value (or null).
type MapAccess struct {
	base Expression
	key  Expression
}

func NewMapAccess(ctx EvalContext, base, key Expression) Expression {
	return &MapAccess{base: base, key: key}
}

func (m *MapAccess) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("map access is not supported in scalar execution")
}

func (m *MapAccess) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	baseArr, err := m.base.Eval(record)
	if err != nil {
		return nil, err
	}
	defer baseArr.Release()

	inputMap, ok := baseArr.(*larray.Map)
	if !ok {
		return nil, fmt.Errorf("map access: base must be a map type, got %T", baseArr)
	}

	keyScalar, err := m.key.EvalScalar()
	if err != nil {
		return nil, fmt.Errorf("map access: failed to evaluate key: %w", err)
	}
	targetKey := scalar.ToString(keyScalar)

	keys := inputMap.Keys().(*array.String)
	values := inputMap.Items().(*array.String)
	offsets := inputMap.Offsets()

	sb := array.NewStringBuilder(memory.DefaultAllocator)
	defer sb.Release()
	sb.Reserve(inputMap.Len())

	for i := 0; i < inputMap.Len(); i++ {
		if inputMap.IsNull(i) {
			sb.AppendNull()
			continue
		}
		row := inputMap.Row(i)
		start, end := offsets[row], offsets[row+1]
		found := false
		for j := int(start); j < int(end); j++ {
			if keys.Value(j) == targetKey {
				if values.IsNull(j) {
					sb.AppendNull()
				} else {
					sb.Append(values.Value(j))
				}
				found = true
				break
			}
		}
		if !found {
			sb.AppendNull()
		}
	}
	return sb.NewArray(), nil
}

func (m *MapAccess) ResultType() ResultType {
	return Array
}

func (m *MapAccess) String() string {
	return fmt.Sprintf("%s[%s]", m.base.String(), m.key.String())
}
