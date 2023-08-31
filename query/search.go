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

package query

import (
	"context"
	"sort"
	"strings"
	"time"

	commonmodels "github.com/lindb/common/models"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/strutil"
	queryctx "github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/stage"
	trackerpkg "github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/field"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// SearchMgr represents the dependencies for searching.
type SearchMgr struct {
	// for intermediate processor set reqeust id, must keep using same request id
	RequestID    string
	Timeout      time.Duration
	CurNode      models.StatelessNode
	Choose       flow.NodeChoose
	TaskMgr      TaskManager
	TransportMgr rpc.TransportManager
}

// MetricMetadataSearchWithResult represents the metadata query executor and retruns the final result set.
func MetricMetadataSearchWithResult(ctx context.Context,
	param *models.ExecuteParam, statement *stmtpkg.MetricMetadata,
	mgr *SearchMgr,
) (any, error) {
	rs, err := MetricMetadataSearch(ctx, param, statement, mgr)
	if err != nil {
		return nil, err
	}
	return buildMetadataResultSet(statement, rs.([]string))
}

// MetricMetadata represents the metadata query executor, includes:
// 1. suggest metric name
// 2. suggest tag keys by spec metric name
// 3. suggest tag values by spec metric name and tag key
// 4. suggest fields by spec metric name
func MetricMetadataSearch(ctx context.Context,
	param *models.ExecuteParam, statement *stmtpkg.MetricMetadata,
	mgr *SearchMgr,
) (any, error) {
	req := models.NewRequest(mgr.CurNode.Indicator(), param.Database, param.SQL)
	taskCtx := queryctx.NewMetadataContext(&queryctx.MetadataDeps{
		Ctx:          ctx,
		Request:      req,
		Database:     param.Database,
		Statement:    statement,
		CurrentNode:  mgr.CurNode,
		Choose:       mgr.Choose,
		TransportMgr: mgr.TransportMgr,
	})
	return exec(taskCtx, req, mgr)
}

// MetricMetadata represents a query executor both storage/broker side.
// When returning query results the following is the order in which processing takes place:
// 1) filtering
// 2) Scanning
// 3) Grouping if it needs
// 4) Down sampling
// 5) Aggregation
// 6) Functions
// 7) Expressions
// Execute query
// 1) plan query language
// 2) aggregator data from time series(memory/file/network)

// BrokerExecutor represents the broker query executor,
// 1) chooses the storage nodes that the data is relatively complete
// 2) chooses broker nodes for root and intermediate computing from all available broker nodes
// 3) storage node as leaf computing node does filter and atomic compute
// 4) intermediate computing nodes are optional, only need if it has grouping query, does order by for grouping
// 4) root computing node does function and expression computing
// 5) finally returns result set to user
//
// NOTICE: there are some scenarios:
//  1. some assignment shards not in query replica shards,
//     maybe some expectant results are lost in data in offline shard, WHY can query not completely data,
//     because of for the system availability.
func MetricDataSearch(ctx context.Context,
	param *models.ExecuteParam, statement *stmtpkg.Query,
	mgr *SearchMgr,
) (any, error) {
	req := models.NewRequest(mgr.CurNode.Indicator(), param.Database, param.SQL)
	taskCtx := queryctx.NewRootMetricContext(
		&queryctx.RootMetricContextDeps{
			Ctx:          ctx,
			Request:      req,
			Database:     param.Database,
			CurrentNode:  mgr.CurNode,
			Statement:    statement,
			Choose:       mgr.Choose,
			TransportMgr: mgr.TransportMgr,
		})
	return exec(taskCtx, req, mgr)
}

// exec executes the query pipeline.
func exec(ctx queryctx.TaskContext, req *models.Request, mgr *SearchMgr) (any, error) {
	if strings.TrimSpace(req.DB) == "" {
		return nil, constants.ErrDatabaseNameRequired
	}
	if mgr.RequestID != "" {
		req.RequestID = mgr.RequestID
	}
	// set request id
	GetRequestManager().NewRequest(req)
	// execute metadata query pipeline
	tracker := trackerpkg.NewStageTracker(flow.NewTaskContextWithTimeout(ctx.Context(), mgr.Timeout))
	ctx.SetTracker(tracker)
	mgr.TaskMgr.AddTask(req.RequestID, ctx)

	defer func() {
		mgr.TaskMgr.RemoveTask(req.RequestID)
		GetRequestManager().CompleteRequest(req.RequestID)
	}()

	pipeline := newExecutePipelineFn(tracker, func(err error) {
		// remove pipeline from cache after execute completed
		defer GetPipelineManager().RemovePipeline(req.RequestID)

		ctx.Complete(err)
	})
	// cache pipeline
	GetPipelineManager().AddPipeline(req.RequestID, pipeline)
	pipeline.Execute(stage.NewPhysicalPlanStage(ctx))
	return ctx.WaitResponse()
}

// buildMetadataResultSet builds metric metadata result set.
func buildMetadataResultSet(statement *stmtpkg.MetricMetadata, result []string) (*commonmodels.Metadata, error) {
	values := strutil.DeDupStringSlice(result)
	sort.Strings(values)
	switch statement.Type {
	case stmtpkg.Field:
		// build field result model
		result := make(map[field.Name]field.Meta)
		fields := field.Metas{}
		for _, value := range values {
			err := encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				return nil, err
			}
			for _, f := range fields {
				result[f.Name] = f
			}
		}
		// HistogramSum(sum), HistogramCount(sum), HistogramMin(min), HistogramMax(max) is visible
		// __bucket_{id}(HistogramField) is not visible for api,
		// underlying histogram data is only restricted access by user via quantile function
		// furthermore, we suggest some quantile functions for user in field names, such as quantile(0.99)
		var (
			resultFields []commonmodels.Field
			hasHistogram bool
		)
		for _, f := range result {
			if f.Type != field.HistogramField {
				resultFields = append(resultFields, commonmodels.Field{
					Name: string(f.Name),
					Type: f.Type.String(),
				})
			} else {
				hasHistogram = true
			}
		}
		//
		if hasHistogram {
			resultFields = append(resultFields,
				commonmodels.Field{Name: "quantile(0.99)", Type: field.HistogramField.String()},
				commonmodels.Field{Name: "quantile(0.95)", Type: field.HistogramField.String()},
				commonmodels.Field{Name: "quantile(0.90)", Type: field.HistogramField.String()},
			)
		}
		sort.Slice(resultFields, func(i, j int) bool {
			return resultFields[i].Name < resultFields[j].Name
		})
		return &commonmodels.Metadata{
			Type:   statement.Type.String(),
			Values: resultFields,
		}, nil
	default:
		return &commonmodels.Metadata{
			Type:   statement.Type.String(),
			Values: values,
		}, nil
	}
}
