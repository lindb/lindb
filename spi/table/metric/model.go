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

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
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
				result.ColumnAssignments = append(result.ColumnAssignments,
					&spi.ColumnAssignment{Column: agg.Column, Handler: &ColumnHandle{Downsampling: tree.Max, Aggregation: agg.AggFuncName}},
				)
			}
			return result
		})
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
	return fmt.Sprintf("(downsampling=%s,aggregation=%s)", c.Downsampling, c.Aggregation)
}

type DataSplit struct {
	partition       *Partition
	groupingContext flow.GroupingContext
	groupingAgg     grouping

	seriesIDHighKey uint16
	lowSeriesIDs    roaring.Container
}

type Partition struct {
	tableScan *TableScan
	shard     *metric.Shard
	segments  []*metric.Segment

	fieldsData []flow.FilterResultSet
}

type TimeSeries struct {
	timestamps []int64
	values     []float64
}

func newTimeSeries(capacity int) *TimeSeries {
	return &TimeSeries{
		timestamps: make([]int64, 0, capacity),
		values:     make([]float64, 0, capacity),
	}
}

func (ts *TimeSeries) Append(timestamp int64, value float64) {
	ts.timestamps = append(ts.timestamps, timestamp)
	ts.values = append(ts.values, value)
}

func (ts *TimeSeries) Reset() {
	ts.timestamps = ts.timestamps[:0]
	ts.values = ts.values[:0]
}
