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

package operator

import (
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

// metricAllSeries represents load all series ids operator.
type metricAllSeries struct {
	executeCtx *flow.ShardExecuteContext
	indexDB    indexdb.IndexDatabase
}

// NewMetricAllSeries creates a metricAllSeries instance.
func NewMetricAllSeries(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &metricAllSeries{
		executeCtx: executeCtx,
		indexDB:    shard.IndexDatabase(),
	}
}

// Execute executes load all series ids by given namespace/metric.
func (op *metricAllSeries) Execute() error {
	queryStmt := op.executeCtx.StorageExecuteCtx.Query
	// get series ids for metric level
	seriesIDs, err := op.indexDB.GetSeriesIDsForMetric(queryStmt.Namespace, queryStmt.MetricName)
	if err != nil {
		return err
	}
	if !queryStmt.HasGroupBy() {
		// add series id without tags, maybe metric has too many series, but one series without tags
		seriesIDs.Add(series.IDWithoutTags)
	}
	op.executeCtx.SeriesIDsAfterFiltering.Or(seriesIDs)
	return nil
}

// Identifier returns identifier string value of all series operator.
func (op *metricAllSeries) Identifier() string {
	return "All Series"
}

// Stats returns the stats of all series operator.
func (op *metricAllSeries) Stats() interface{} {
	return &models.SeriesStats{
		NumOfSeries: op.executeCtx.SeriesIDsAfterFiltering.GetCardinality(),
	}
}
