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
	"math"
	"sync"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

const infBlockSize = 360

var (
	infFilledBlock = make([]float64, infBlockSize)
)

func init() {
	for i := 0; i < infBlockSize; i++ {
		infFilledBlock[i] = math.Inf(1) + 1
	}
}

var float64Pool sync.Pool

func fillInfBlock(sl []float64) []float64 {
	length := len(sl)
	for i := 0; i <= length/infBlockSize; i++ {
		from := i * infBlockSize
		to := (i + 1) * infBlockSize
		if to > length {
			to = length
		}
		copy(sl[from:to], infFilledBlock)
	}
	return sl
}

func getFloat64Slice(size int) []float64 {
	item := float64Pool.Get()
	if item == nil {
		return make([]float64, size)
	}
	sl := item.(*[]float64)
	if cap(*sl) < size {
		return make([]float64, size)
	}
	return (*sl)[:size]
}

func putFloat64Slice(sl *[]float64) {
	float64Pool.Put(sl)
}

// DownSamplingMultiSeriesInto merges field data from source time range => target time range,
// data will be merged into DownSamplingResult
// for example: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min.
func DownSamplingMultiSeriesInto(
	target timeutil.SlotRange, ratio uint16,
	aggFunc field.AggFunc, decoders []*encoding.TSDDecoder,
	emitValue func(targetPos int, value float64),
) {
	targetValues := make([]float64, infBlockSize)
	length := int(target.End-target.Start) + 1
	if length <= infBlockSize {
		// on stack
		targetValues = targetValues[:length]
	} else {
		// on heap
		targetValues = getFloat64Slice(length)
		defer putFloat64Slice(&targetValues)
	}
	// first loop: filled target values with inf value
	// inf value is invalid, and won't be emitted after downsampling
	fillInfBlock(targetValues)

	// second loop: iterating tsd decoder
	for _, decoder := range decoders {
		if decoder == nil {
			continue
		}
		for movingSourceSlot := decoder.StartTime(); movingSourceSlot <= decoder.EndTime(); movingSourceSlot++ {
			if !decoder.HasValueWithSlot(movingSourceSlot) {
				continue
			}
			value := math.Float64frombits(decoder.Value())
			targetPos := int(movingSourceSlot/ratio) - int(target.Start)
			if targetPos < 0 {
				continue
			}
			// exhausted
			if targetPos >= length {
				break
			}
			// not set before
			if math.IsInf(targetValues[targetPos], 1) {
				targetValues[targetPos] = value
				// set before, aggregate
			} else {
				targetValues[targetPos] = aggFunc.Aggregate(targetValues[targetPos], value)
			}
		}
	}
	// third loop, emit downsampling data
	for offset, value := range targetValues {
		emitValue(offset, value)
	}
}
