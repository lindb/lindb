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

package context

import (
	"context"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/sql/tree"
)

// MetadataDeps represents metric metadata search dependency.
type MetadataDeps struct {
	Ctx     context.Context
	Request *models.Request

	Database     string
	Statement    *tree.MetricMetadata
	CurrentNode  models.StatelessNode
	Choose       flow.NodeChoose
	TransportMgr rpc.TransportManager
}

// MetadataContext represents metric metadata search context.
type MetadataContext struct {
	baseTaskContext

	Deps *MetadataDeps
	// handle response
	results []string
}

// NewMetadataContext creates metric metadata search context.
func NewMetadataContext(deps *MetadataDeps) *MetadataContext {
	if deps.Statement.Limit == 0 || deps.Statement.Limit > constants.MaxSuggestions {
		// if limit =0 or > max suggestion items, need reset limit
		deps.Statement.Limit = constants.MaxSuggestions
	}
	return &MetadataContext{
		baseTaskContext: newBaseTaskContext(deps.Ctx, deps.TransportMgr),
		Deps:            deps,
	}
}

// WaitResponse waits metric metadata search task completed and returns metric data.
func (ctx *MetadataContext) WaitResponse() (any, error) {
	select {
	case <-ctx.doneCh:
		// received all data, break for loop
		return ctx.results, ctx.err
	case <-ctx.Deps.Ctx.Done():
		return nil, constants.ErrTimeout
	}
}

// HandleResponse handles metric metadata task response.
func (ctx *MetadataContext) HandleResponse(resp *protoCommonV1.TaskResponse, fromNode string) {
	ctx.handleResponse(resp, fromNode)
	ctx.tryClose()
}

// MakePlan makes the metric metadata physical plan.
func (ctx *MetadataContext) MakePlan() error {
	physicalPlans, err := ctx.Deps.Choose.Choose(ctx.Deps.Database, 1)
	if err != nil {
		return err
	}
	if len(physicalPlans) == 0 {
		return constants.ErrTargetNodesNotFound
	}

	suggestMarshalData, _ := ctx.Deps.Statement.MarshalJSON()
	for _, physicalPlan := range physicalPlans {
		physicalPlan.AddReceiver(ctx.Deps.CurrentNode.Indicator())
		if err := physicalPlan.Validate(); err != nil {
			return err
		}
		ctx.addRequests(
			&protoCommonV1.TaskRequest{
				RequestID:    ctx.Deps.Request.RequestID,
				RequestType:  protoCommonV1.RequestType_Metadata,
				PhysicalPlan: encoding.JSONMarshal(physicalPlan),
				Payload:      suggestMarshalData,
			}, physicalPlan)
	}

	return nil
}

// handleResponse handles metric metadata search task response with lock.
func (ctx *MetadataContext) handleResponse(resp *protoCommonV1.TaskResponse, fromNode string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.handleTaskState(resp, fromNode)
	ctx.expectResults--

	result := &models.SuggestResult{}
	if err := encoding.JSONUnmarshal(resp.Payload, result); err != nil {
		ctx.err = err
	}
	ctx.results = append(ctx.results, result.Values...)
}
