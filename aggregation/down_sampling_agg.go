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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./down_sampling_agg.go -destination=./down_sampling_agg_mock.go -package=aggregation

// DownSamplingAggregator represents down sampling for field data.
type DownSamplingAggregator interface {
	// DownSampling merges fields' data by target interval and time range.
	DownSampling(aggFunc field.AggFunc, values []*encoding.TSDDecoder)
}

// DownSamplingResult represents the result of down sampling aggregator.
type DownSamplingResult interface {
	// Append appends time and value.
	Append(slot bit.Bit, value float64)
	Reset()
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

func (d *downSamplingMergeResult) Reset() {
	d.pos = 0
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

func (rs *TSDDownSamplingResult) Reset() {
	//do nothing
}

// downSamplingAggregator implements DownSamplingAggregator interface.
type downSamplingAggregator struct {
	source, target timeutil.SlotRange
	ratio          uint16

	rs DownSamplingResult
}

// NewDownSamplingAggregator creates DownSamplingAggregator.
func NewDownSamplingAggregator(source, target timeutil.SlotRange,
	ratio uint16, rs DownSamplingResult) DownSamplingAggregator {
	return &downSamplingAggregator{
		source: source,
		target: target,
		ratio:  ratio,
		rs:     rs,
	}
}

// DownSampling merges field data from source time range => target time range,
// for example: source range[5,182]=>target range[0,6], ratio:30, source interval:10s, target interval:5min.
func (ds *downSamplingAggregator) DownSampling(aggFunc field.AggFunc, values []*encoding.TSDDecoder) {
	hasValue := false
	pos := ds.source.Start
	end := ds.source.End
	result := 0.0
	rs := ds.rs
	// first loop: target slot range
	for j := ds.target.Start; j <= ds.target.End; j++ {
		// second loop: source slot range and ratio(target interval/source interval)
		intervalEnd := ds.ratio * (j + 1)
		for pos <= end && pos < intervalEnd {
			// 1. merge data by time slot
			for _, value := range values {
				if value == nil {
					// if series id not exist, value maybe nil
					continue
				}
				if value.HasValueWithSlot(pos) {
					if !hasValue {
						// if target value not exist, set it
						result = math.Float64frombits(value.Value())
						hasValue = true
					} else {
						// if target value exist, do aggregate
						result = aggFunc.Aggregate(result, math.Float64frombits(value.Value()))
					}
				}
			}
			pos++
		}
		// 2. add data into rs stream
		if hasValue {
			rs.Append(bit.One, result)
			// reset has value for next loop
			hasValue = false
			result = 0.0
		} else {
			rs.Append(bit.Zero, constants.EmptyValue)
		}
	}
}
