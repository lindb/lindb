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

package function

import (
	"github.com/lindb/lindb/pkg/collections"
)

// FuncCall calls the function calc by function type and params
func FuncCall(funcType FuncType, params ...collections.FloatArray) collections.FloatArray {
	switch funcType {
	case Sum, Min, Max, Count:
		if len(params) == 0 {
			return nil
		}
		return params[0]
	case Avg:
		// params: 0=>sum, 1=>count
		if len(params) < 2 {
			return nil
		}
		result := collections.NewFloatArray(params[0].Capacity())
		it := params[0].Iterator()
		for it.HasNext() {
			idx, sum := it.Next()
			if params[1].HasValue(idx) {
				count := params[1].GetValue(idx)
				if count != 0 {
					result.SetValue(idx, sum/count)
				}
			}
		}
		return result
	default:
		return nil
	}
}
