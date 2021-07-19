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
	"math"
	"sort"

	"github.com/lindb/lindb/pkg/collections"
)

type bucket struct {
	upperBound float64
	count      float64
	itr        collections.FloatArrayIterator
}

type buckets []bucket

func (bkt buckets) Len() int           { return len(bkt) }
func (bkt buckets) Less(i, j int) bool { return bkt[i].upperBound < bkt[j].upperBound }
func (bkt buckets) Swap(i, j int)      { bkt[i], bkt[j] = bkt[j], bkt[i] }

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
func QuantileCall(q float64, histogramFields map[float64][]collections.FloatArray) (collections.FloatArray, error) {
	if q < 0 || q > 1 {
		return nil, fmt.Errorf("QuantileCall with illegal value: %f", q)
	}
	var histogramBuckets = make(buckets, len(histogramFields))

	var idx = 0
	for upperBound, arrays := range histogramFields {
		if len(arrays) != 1 {
			return nil, fmt.Errorf("QuantileCall buckets's floatArray count: %d not equals 1", len(arrays))
		}
		histogramBuckets[idx] = bucket{upperBound: upperBound, itr: arrays[0].Iterator()}
		idx++
	}
	sort.Sort(histogramBuckets)

	if len(histogramBuckets) < 2 {
		return nil, fmt.Errorf("QuantileCall with buckets count: %d less than 2", len(histogramBuckets))
	}
	if !math.IsInf(histogramBuckets[len(histogramBuckets)-1].upperBound, +1) {
		return nil, fmt.Errorf("QuantileCall's largest upper bound is not +Inf")
	}
	capacity := histogramFields[histogramBuckets[0].upperBound][0].Capacity()
	targetFloatArray := collections.NewFloatArray(capacity)

	itr := histogramBuckets[0].itr
	for itr.HasNext() {
		pos, v := itr.Next()
		histogramBuckets[0].count = v

		for bucketIdx := 1; bucketIdx < len(histogramBuckets); bucketIdx++ {
			if !histogramBuckets[bucketIdx].itr.HasNext() {
				return nil, fmt.Errorf("QuantileCall floatArray length")
			}
			_, v := histogramBuckets[bucketIdx].itr.Next()
			histogramBuckets[bucketIdx].count = v
		}
		histogramBuckets.EnsureCountFieldCumulative()

		observations := histogramBuckets[len(histogramFields)-1].count
		if observations == 0 {
			targetFloatArray.SetValue(pos, 0)
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
