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

package printer

import (
	"fmt"
	"strings"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// collapseHistogramBuckets collapses histogram bucket entries in a parallel
// (names, formatted) slice into summary entries, leaving non-bucket entries
// unchanged.
//
//   - names[i]     — the raw field/column name used for bucket detection
//   - formatted[i] — the display string (e.g. "latency.__bucket_0.5:sum" or
//     "latency.__bucket_0.5->(downsampling=sum)"); the bucket name
//     portion is replaced with the summary in the output.
//
// All bucket fields that share the same histogram prefix (e.g. "latency") are
// collapsed into a single entry "latency.__bucket[N]<rest>" where N is the
// total count of buckets for that prefix.
func collapseHistogramBuckets(names, formatted []string) []string {
	// Pass 1: count bucket fields per histogram prefix.
	bucketCounts := make(map[string]int)
	for _, n := range names {
		if seriesmetric.IsBucketField(n) {
			bucketCounts[seriesmetric.HistoNameFromField(n)]++
		}
	}

	// Pass 2: emit one condensed summary per prefix; skip subsequent buckets.
	seenPrefix := make(map[string]bool)
	var parts []string
	for i, n := range names {
		if !seriesmetric.IsBucketField(n) {
			parts = append(parts, formatted[i])
			continue
		}
		prefix := seriesmetric.HistoNameFromField(n)
		if seenPrefix[prefix] {
			continue // already emitted the summary for this prefix
		}
		seenPrefix[prefix] = true
		// Replace the full bucket name inside the formatted string with the summary.
		// formatted[i] starts with the bucket name followed by a delimiter
		// (":" for symbols, "->" for assignments), so Replace with count=1 is safe.
		summary := fmt.Sprintf("%s.__bucket[%d]", prefix, bucketCounts[prefix])
		parts = append(parts, strings.Replace(formatted[i], n, summary, 1))
	}
	return parts
}

// formatSymbols formats a list of plan Symbols for display.
// Histogram bucket symbols (<prefix>.__bucket_<bound>) are collapsed into a
// single summary entry "<prefix>.__bucket[N]:<type>" to keep the output concise.
func formatSymbols(symbols []*plan.Symbol) string {
	names := make([]string, len(symbols))
	formatted := make([]string, len(symbols))
	for i, s := range symbols {
		names[i] = s.Name
		formatted[i] = s.String()
	}
	return "[" + strings.Join(collapseHistogramBuckets(names, formatted), ", ") + "]"
}

func formatAggregation(aggregation *plan.Aggregation) string {
	var args []string
	for _, arg := range aggregation.Arguments {
		args = append(args, fmt.Sprintf("%q", tree.FormatExpression(arg)))
	}
	return fmt.Sprintf("%s(%s)", aggregation.Function, strings.Join(args, ", "))
}
