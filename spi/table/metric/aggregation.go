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

package metric

import (
	"fmt"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series/field"
)

func aggregate(fn field.AggType, start, step int64, dst *collections.FloatArray, src *TimeSeries) {
	for index, timestamp := range src.timestamps {
		value := src.values[index]
		offset := 0
		if step > 0 {
			offset = int((timestamp - start) / step)
		}
		// TODO: add log
		if offset < 0 {
			panic(fmt.Sprintf("warn offset < 0, offset=%v\n", offset))
		}

		if dst.HasValue(offset) {
			dst.SetValue(offset, fn.Aggregate(dst.GetValue(offset), value))
		} else {
			dst.SetValue(offset, value)
		}
	}
}
