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

package metricsdata

import (
	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/timeutil"
)

// downsampling merges field data from source time range => target time range,
// data will be merged into DownSamplingResult
// for example: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min.
func downsampling[V float64 | *models.Exemplar](
	context *mergerContext,
	timeRange timeutil.SlotRange,
	result *collections.Array[V], getter func(sloat uint16) (V, bool),
	aggFn func(old V, new V) V,
) {
	bs := int(context.baseSlot)
	target := context.targetRange
	ratio := context.ratio

	for movingSourceSlot := timeRange.Start; movingSourceSlot <= timeRange.End; movingSourceSlot++ {
		v, ok := getter(movingSourceSlot)
		if !ok {
			continue
		}
		targetPos := bs + int(movingSourceSlot/ratio) - int(target.Start)
		if targetPos < 0 {
			continue
		}

		// TODO: target pos exceed target range?

		if result.HasValue(targetPos) {
			result.SetValue(targetPos, aggFn(result.GetValue(targetPos), v))
		} else {
			result.SetValue(targetPos, v)
		}
	}
}
