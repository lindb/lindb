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
	"fmt"
	"sort"

	"github.com/lindb/lindb/pkg/collections"
)

type bucket struct {
	value      *collections.FloatArray
	upperBound float64
	count      float64
}

type buckets []bucket

func (bkt buckets) Len() int { return len(bkt) }

func (bkt buckets) Less(i, j int) bool { return bkt[i].upperBound < bkt[j].upperBound }

func (bkt buckets) Swap(i, j int) { bkt[i], bkt[j] = bkt[j], bkt[i] }

// EnsureCountFieldCumulative ensures count in buckets are cumulative for quantile function
func (bkt buckets) EnsureCountFieldCumulative() {
	if bkt.Len() == 0 {
		return
	}
	last := bkt[0].count
	for i := 1; i < bkt.Len(); i++ {
		last += bkt[i].count
		bkt[i].count = last
	}
}

// QuantileCall references to prometheus implementation.
// https://github.com/prometheus/prometheus/blob/39d79c3cfb86c47d6bc06a9e9317af582f1833bb/promql/quantile.go
// 0 <= q <= 1
// q = 0, returns 0
// q = 1, last UpperBound before Inf is returned
func QuantileCall(q float64, histogramFields map[float64][]*collections.FloatArray) (*collections.FloatArray, error) {
	if q < 0 || q > 1 {
		return nil, fmt.Errorf("QuantileCall with illegal value: %f", q)
	}
	var histogramBuckets buckets
	for upperBound, arrays := range histogramFields {
		if len(arrays) >= 1 {
			histogramBuckets = append(histogramBuckets, bucket{upperBound: upperBound, value: arrays[0]})
		}
	}
	if len(histogramBuckets) == 0 {
		return nil, nil
	}
	sort.Sort(histogramBuckets)

	capacity := histogramFields[histogramBuckets[0].upperBound][0].Capacity()
	targetFloatArray := collections.NewFloatArray(capacity)
	for pos := 0; pos < capacity; pos++ {
		for bucketIdx := 0; bucketIdx < len(histogramBuckets); bucketIdx++ {
			if histogramBuckets[bucketIdx].value.HasValue(pos) {
				v := histogramBuckets[bucketIdx].value.GetValue(pos)
				histogramBuckets[bucketIdx].count = v
			} else {
				histogramBuckets[bucketIdx].count = 0
			}
		}
		histogramBuckets.EnsureCountFieldCumulative()

		// last bucket count = total
		observations := histogramBuckets[len(histogramBuckets)-1].count
		if observations == 0 {
			targetFloatArray.SetValue(pos, 0)
			continue
		}
		if len(histogramBuckets) == 1 {
			targetFloatArray.SetValue(pos, histogramBuckets[0].upperBound)
			continue
		}

		rank := q * observations
		b := sort.Search(len(histogramBuckets)-1, func(i int) bool { return histogramBuckets[i].count >= rank })
		if b == len(histogramBuckets)-1 {
			targetFloatArray.SetValue(pos, histogramBuckets[len(histogramBuckets)-2].upperBound)
			continue
		} else if b == 0 && histogramBuckets[0].upperBound <= 0 {
			targetFloatArray.SetValue(pos, histogramBuckets[0].upperBound)
			continue
		}

		var (
			bucketStart float64
			bucketEnd   = histogramBuckets[b].upperBound
			count       = histogramBuckets[b].count
		)
		if b > 0 {
			bucketStart = histogramBuckets[b-1].upperBound
			count -= histogramBuckets[b-1].count
			rank -= histogramBuckets[b-1].count
		}
		targetFloatArray.SetValue(pos, bucketStart+(bucketEnd-bucketStart)*(rank/count))
	}

	return targetFloatArray, nil
}
