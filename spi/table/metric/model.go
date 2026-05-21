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

	"github.com/lindb/common/models"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage/metric"
)

func init() {
	// register table/column handle
	encoding.RegisterNodeType(TableHandle{})
	encoding.RegisterNodeType(ColumnHandle{})

	spi.RegisterCreateTableFn(spi.Metric, func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
			Database:  db,
			Namespace: ns,
			Metric:    name,
		}
	})

	spi.RegisterApplyAggregationFn(spi.Metric,
		func(table spi.TableHandle, tableMeta *types.TableMetadata,
			aggregations []spi.ColumnAggregation,
		) *spi.ApplyAggregationResult {
			result := &spi.ApplyAggregationResult{}
			// FIXME: find downSampling agg
			for _, agg := range aggregations {
				if isHistogramFunc(agg.AggFuncName) {
					// Expand histogram function into physical column assignments (sum/sum per field).
					// Histogram computation is handled at the broker aggregation layer.
					assignments := expandHistogramColumns(agg, tableMeta)
					result.ColumnAssignments = append(result.ColumnAssignments, assignments...)
				} else {
					result.ColumnAssignments = append(result.ColumnAssignments,
						&spi.ColumnAssignment{Column: agg.Column, Handler: &ColumnHandle{Downsampling: tree.Max, Aggregation: agg.AggFuncName}},
					)
				}
			}
			return result
		})
}

// isHistogramFunc reports whether the given function name is a histogram aggregation.
func isHistogramFunc(name tree.FuncName) bool {
	return tree.IsHistogramFunc(name)
}

// expandHistogramColumns maps a histogram ColumnAggregation to the physical storage columns
// that need to be read and the ColumnHandle that controls how each is aggregated.
//
// Physical naming convention (set by AppendFields in internal/linmetric/histogram_delta.go):
//   - <histoName>.__sum, <histoName>.__count, <histoName>.__min, <histoName>.__max  — scalar statistics
//   - <histoName>.__bucket_<bound>                                                  — per-bucket counts
//
// For histogram_sum / histogram_count we only need one stat field.
// For histogram_quantile / histogram_avg we need all bucket fields + sum + count.
//
// All physical field handles use plain (sum, sum) — histogram computation is handled
// at the broker aggregation layer, not inside the storage source connector.
func expandHistogramColumns(agg spi.ColumnAggregation, tableMeta *types.TableMetadata) []*spi.ColumnAssignment {
	sumHandle := &ColumnHandle{Downsampling: tree.Sum, Aggregation: tree.Sum}

	switch agg.AggFuncName {
	case tree.HistogramSum:
		return []*spi.ColumnAssignment{
			{Column: seriesmetric.HistoStatFieldName(agg.Column, seriesmetric.HistoStatSum), Handler: sumHandle},
		}
	case tree.HistogramCount:
		return []*spi.ColumnAssignment{
			{Column: seriesmetric.HistoStatFieldName(agg.Column, seriesmetric.HistoStatCount), Handler: sumHandle},
		}
	default:
		// histogram_quantile and histogram_avg need all bucket fields and the stat fields.
		// Collect them from the Arrow schema so we only request fields that actually exist.
		var assignments []*spi.ColumnAssignment
		if tableMeta != nil && tableMeta.Schema != nil {
			for _, f := range tableMeta.Schema.Fields() {
				name := f.Name
				if seriesmetric.IsBucketField(name) && seriesmetric.HistoNameFromField(name) == agg.Column {
					assignments = append(assignments, &spi.ColumnAssignment{
						Column:  name,
						Handler: sumHandle,
					})
				}
			}
		}
		// Always include sum and count so histogram_avg can divide them.
		assignments = append(assignments,
			&spi.ColumnAssignment{Column: seriesmetric.HistoStatFieldName(agg.Column, seriesmetric.HistoStatSum), Handler: sumHandle},
			&spi.ColumnAssignment{Column: seriesmetric.HistoStatFieldName(agg.Column, seriesmetric.HistoStatCount), Handler: sumHandle},
		)
		return assignments
	}
}

type TableHandle struct {
	Database  string `json:"database"`
	Namespace string `json:"namespace"`
	Metric    string `json:"metric"`

	TimeRange timeutil.TimeRange `json:"timeRange"`
	Interval  timeutil.Interval  `json:"interval"`
}

func (t *TableHandle) SetTimeRange(timeRange timeutil.TimeRange) {
	t.TimeRange = timeRange
}

func (t *TableHandle) GetTimeRange() timeutil.TimeRange {
	return t.TimeRange
}

func (t *TableHandle) SetInterval(interval timeutil.Interval) {
	t.Interval = interval
}

func (t *TableHandle) GetInterval() timeutil.Interval {
	return t.Interval
}

func (t *TableHandle) Kind() spi.DatasourceKind {
	return spi.Metric
}

func (t *TableHandle) String() string {
	return fmt.Sprintf("%s:%s:%s", t.Database, t.Namespace, t.Metric)
}

type ColumnHandle struct {
	Downsampling tree.FuncName `json:"downsampling"`
	Aggregation  tree.FuncName `json:"aggregation"`
}

func (c *ColumnHandle) String() string {
	// When downsampling and aggregation use the same function, omit the redundant
	// aggregation field to keep plan output concise.
	if c.Downsampling == c.Aggregation {
		return fmt.Sprintf("(downsampling=%s)", c.Downsampling)
	}
	return fmt.Sprintf("(downsampling=%s,aggregation=%s)", c.Downsampling, c.Aggregation)
}

type DataSplit struct {
	partition       *Partition
	groupingContext flow.GroupingContext
	groupingAgg     grouping

	numOfPoints int // num. of points per series(storage level)

	seriesIDHighKey uint16
	lowSeriesIDs    roaring.Container
}

type Partition struct {
	tableScan *TableScan
	shard     *metric.Shard
	segments  []*metric.Segment

	fieldsData []flow.FilterResultSet
}

type TimeSeries[V float64 | *models.Exemplar] struct {
	timestamps []int64
	values     []V
}

func newTimeSeries[V float64 | *models.Exemplar](capacity int) *TimeSeries[V] {
	return &TimeSeries[V]{
		timestamps: make([]int64, 0, capacity),
		values:     make([]V, 0, capacity),
	}
}

func (ts *TimeSeries[V]) Append(timestamp int64, value V) {
	ts.timestamps = append(ts.timestamps, timestamp)
	ts.values = append(ts.values, value)
}

func (ts *TimeSeries[V]) Reset() {
	ts.timestamps = ts.timestamps[:0]
	ts.values = ts.values[:0]
}
