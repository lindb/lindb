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
	"github.com/lindb/lindb/pkg/timeutil"
)

// RateCall represents rate function call.
func RateCall(interval int64, params ...*collections.FloatArray) *collections.FloatArray {
	if len(params) == 0 {
		return nil
	}
	result := collections.NewFloatArray(params[0].Capacity())
	itr := params[0].NewIterator()
	for itr.HasNext() {
		idx, val := itr.Next()
		result.SetValue(idx, val/float64(interval/timeutil.OneSecond))
	}
	return result
}
