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

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
)

//go:generate mockgen -source=./order_by.go -destination=./order_by_mock.go -package=aggregation

// aggResult represents the result of order by aggregate.
type aggResult struct {
	count                                   int
	min, max, sum, last, first, avg, stddev float64
}

// OrderByRow represents row for order by, implements Row interface.
type OrderByRow struct {
	tags   string
	fields map[string]*collections.FloatArray

	// cache order by point value
	points map[string]*aggResult
}

// NewOrderByRow creates a OrderByRow instance.
func NewOrderByRow(tags string, fields map[string]*collections.FloatArray) Row {
	return &OrderByRow{
		tags:   tags,
		fields: fields,
		points: make(map[string]*aggResult),
	}
}

// aggregate returns the value by given function type for this series data.
func (r *OrderByRow) aggregate(it *collections.FloatArrayIterator) *aggResult {
	var count int
	var min, max, sum, last, first, avg, stddev float64
	// stddev
	var value, mean float64

	for it.HasNext() {
		// get value
		_, val := it.Next()

		last = val

		if count == 0 {
			count = 1
			first = val
			min = val
			max = val
			sum = val
			mean = val
			value = val
			continue
		}

		if val < min {
			min = val
		} else if val > max {
			max = val
		}
		count++
		sum += val
		delta := val - mean
		mean += delta / float64(count)
		value += delta * (val - mean)
	}

	// finally calc avg/stddev
	if count > 0 {
		avg = sum / float64(count)
		stddev = math.Sqrt(value / float64(count))
	}
	return &aggResult{
		count:  count,
		min:    min,
		max:    max,
		sum:    sum,
		last:   last,
		first:  first,
		avg:    avg,
		stddev: stddev,
	}
}

// GetValue returns the value of aggregation for this series based on given field name/function type.
func (r *OrderByRow) GetValue(fieldName string, funcType function.FuncType) float64 {
	val, ok := r.points[fieldName]
	if !ok {
		field, ok := r.fields[fieldName]
		if !ok {
			return 0.0
		}
		if field == nil {
			return 0.0
		}
		val = r.aggregate(field.NewIterator())
		r.points[fieldName] = val
	}
	// return val base function type
	switch funcType {
	case function.Count:
		return float64(val.count)
	case function.Sum:
		return val.sum
	case function.Max:
		return val.max
	case function.Min:
		return val.min
	case function.First:
		return val.first
	case function.Last:
		return val.last
	case function.Avg:
		return val.avg
	case function.Stddev:
		return val.stddev
	}
	return 0.0
}

// ResultSet returns the resutl set of series(tags/fields).
func (r *OrderByRow) ResultSet() (tags string, fields map[string]*collections.FloatArray) {
	return r.tags, r.fields
}

// OrderBy represents order by container.
type OrderBy interface {
	// Push pushes row into container.
	Push(row Row)
	// ResultSet returns result set of order by.
	ResultSet() []Row
}

// resultLimiter represents a size limit container, implements OrderBy interface.
type resultLimiter struct {
	rows  []Row
	limit int
}

// NewResultLimiter creates a size limit container.
func NewResultLimiter(limit int) OrderBy {
	return &resultLimiter{
		limit: limit,
	}
}

// Push pushes row into limit container.
func (r *resultLimiter) Push(row Row) {
	if len(r.rows) < r.limit {
		r.rows = append(r.rows, row)
	}
}

// ResultSet returns result set of limiter.
func (r *resultLimiter) ResultSet() []Row {
	return r.rows
}

// optNOrderBy implements OrderBy interface(top n).
type topNOrderBy struct {
	topn *topNHeap
}

// NewTopNOrderBy creates a topNOrderBy container instance.
func NewTopNOrderBy(orderByItems []*OrderByItem, topN int) OrderBy {
	return &topNOrderBy{
		topn: newTopNHeap(orderByItems, topN),
	}
}

// Push pushes row into topN container.
func (o *topNOrderBy) Push(row Row) {
	o.topn.Add(row)
}

// ResultSet returns result set for topN.
func (o *topNOrderBy) ResultSet() []Row {
	return o.topn.ResultSet()
}
