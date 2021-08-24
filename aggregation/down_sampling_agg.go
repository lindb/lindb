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
	"encoding/json"
	"fmt"
	"math"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./down_sampling_agg.go -destination=./down_sampling_agg_mock.go -package=aggregation

// DownSamplingResult represents the result of down sampling aggregator.
type DownSamplingResult interface {
	// Append appends time and value.
	Append(slot bit.Bit, value float64)
}

type downSamplingMergeResult struct {
	agg FieldAggregator

	pos int
}

func NewDownSamplingMergeResult(agg FieldAggregator) DownSamplingResult {
	return &downSamplingMergeResult{
		agg: agg,
		pos: 0,
	}
}

func (d *downSamplingMergeResult) Append(slot bit.Bit, value float64) {
	if slot == bit.One {
		d.agg.AggregateBySlot(d.pos, value)
	}
	d.pos++
}

// TSDDownSamplingResult implements DownSamplingResult using encoding TSDEncoder.
type TSDDownSamplingResult struct {
	stream encoding.TSDEncoder
}

// NewTSDDownSamplingResult creates tsd down sampling result.
func NewTSDDownSamplingResult(stream encoding.TSDEncoder) DownSamplingResult {
	return &TSDDownSamplingResult{stream: stream}
}

// Append appends time and value into tsd encode stream.
func (rs *TSDDownSamplingResult) Append(slot bit.Bit, value float64) {
	rs.stream.AppendTime(slot)
	if slot == bit.One {
		rs.stream.AppendValue(math.Float64bits(value))
	}
}

// DownSamplingMultiSeriesInto merges field data from source time range => target time range,
// data will be merged into DownSamplingResult
// for example: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min.
func DownSamplingMultiSeriesInto(
	target timeutil.SlotRange, ratio uint16,
	aggFunc field.AggFunc, decoders []*encoding.TSDDecoder,
	rs DownSamplingResult,
) {
	// first loop: target slot range
	// todo:remove

	var m = make(map[uint16]float64)

	for j := target.Start; j <= target.End; j += ratio {
		hasValue := bit.Zero
		result := constants.EmptyValue
		// loop: source slot range and ratio(target interval/source interval)
		intervalEnd := ratio * (j + 1)

		// seek reads from the start slot
		// flushed data always starts from 0, but target is arbitrarily
		// decoders: 0-359
		for _, d := range decoders {
			if d != nil {
				d.Seek(j)
			}
		}
		for pos := j; pos < intervalEnd; pos++ {
			for _, decoder := range decoders {
				if decoder == nil {
					// if series id not exist, value maybe nil
					continue
				}
				if decoder.HasValueWithSlot(pos) {
					if !hasValue {
						// if target value not exist, set it
						result = math.Float64frombits(decoder.Value())
						hasValue = bit.One
					} else {
						// if target value exist, do aggregate
						result = aggFunc.Aggregate(result, math.Float64frombits(decoder.Value()))
					}
				}
			}
		}

		m[j] = result
		rs.Append(hasValue, result)
	}
	data, _ := json.Marshal(m)
	fmt.Println(string(data))
}
