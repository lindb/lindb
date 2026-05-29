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
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// columnAggregator fills one output column for all rows of a batch.
// Implementations are chosen by buildColumnAgg based on the aggregation function name.
type columnAggregator interface {
	fill(b array.Builder, batch arrow.RecordBatch, numRows int)
}

// buildColumnAgg is the aggregation-function dispatch factory.
// It inspects agg.Aggregation.Function and returns the appropriate columnAggregator:
//   - histogram functions → histogramScalarAgg or histogramTimeSeriesAgg
//   - standard functions  → passThroughAgg (broker relays already-aggregated value)
//   - unknown / missing   → nullColumnAgg
func buildColumnAgg(
	agg *plan.AggregationAssignment,
	inputIdx map[string]int,
	isTimeSeries bool,
) columnAggregator {
	if tree.IsHistogramFunc(agg.Aggregation.Function) {
		phi := extractHistoPhi(agg.Aggregation.Arguments)
		histoName := extractHistoName(agg)

		var buckets []bucketCol
		for colName, ci := range inputIdx {
			if seriesmetric.IsBucketField(colName) && seriesmetric.HistoNameFromField(colName) == histoName {
				if bound, ok := seriesmetric.ParseBucketBound(colName); ok {
					buckets = append(buckets, bucketCol{bound: bound, colIdx: ci})
				}
			}
		}
		sort.Slice(buckets, func(i, j int) bool {
			if math.IsInf(buckets[i].bound, 1) {
				return false
			}
			if math.IsInf(buckets[j].bound, 1) {
				return true
			}
			return buckets[i].bound < buckets[j].bound
		})

		sumIdx, hasSum := inputIdx[seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatSum)]
		countIdx, hasCount := inputIdx[seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatCount)]

		p := rowFunc(agg.Aggregation.Function, phi, buckets, sumIdx, countIdx, hasSum, hasCount)
		if isTimeSeries {
			return &histogramTimeSeriesAgg{params: p}
		}
		return &histogramScalarAgg{params: p}
	}

	// Standard aggregation: storage already produced the final (or partial) value;
	// the broker copies it verbatim.  Column name equals the output symbol name.
	if ci, ok := inputIdx[agg.Symbol.Name]; ok {
		return &passThroughAgg{colIdx: ci}
	}
	return &nullColumnAgg{}
}
