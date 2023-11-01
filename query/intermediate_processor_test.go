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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	queryctx "github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/sql/stmt"
)

func TestProcess_Fail(t *testing.T) {
	p := &intermediateTaskProcessor{}
	err := p.Process(nil, nil, &protoCommonV1.TaskRequest{PhysicalPlan: []byte("abc")})
	assert.Error(t, err)

	ip := NewIntermediateTaskProcessor(models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}, time.Second, nil, nil, nil)
	err = ip.Process(nil, nil, &protoCommonV1.TaskRequest{
		PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
			Targets: []*models.Target{{Indicator: "1.1.1.1:8000"}},
		}),
	})
	assert.Error(t, err)

	err = ip.Process(nil, nil, &protoCommonV1.TaskRequest{
		PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
			Targets: []*models.Target{{Indicator: "1.1.1.1:9000", ReceiveOnly: true}},
		}),
	})
	assert.NoError(t, err)

	err = ip.Process(nil, nil, &protoCommonV1.TaskRequest{
		RequestType: protoCommonV1.RequestType(10),
		PhysicalPlan: encoding.JSONMarshal(&models.PhysicalPlan{
			Targets: []*models.Target{{Indicator: "1.1.1.1:9000"}},
		}),
	})
	assert.NoError(t, err)
}

func TestProcessMetricDataSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		execFn = exec
		ctrl.Finish()
	}()

	physicalPlan := encoding.JSONMarshal(&models.PhysicalPlan{
		Targets: []*models.Target{{Indicator: "1.1.1.1:9000"}},
	})
	ip := NewIntermediateTaskProcessor(models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}, time.Second, nil, nil, nil)
	taskCtx := &flow.TaskContext{}
	err := ip.Process(taskCtx, nil, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Data,
		Payload:      []byte("abc"),
		PhysicalPlan: physicalPlan,
	})
	assert.Error(t, err)

	execFn = func(ctx queryctx.TaskContext, req *models.Request, mgr *SearchMgr) (any, error) {
		return nil, fmt.Errorf("err")
	}
	statement, _ := (&stmt.Query{}).MarshalJSON()
	err = ip.Process(taskCtx, nil, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Data,
		Payload:      statement,
		PhysicalPlan: physicalPlan,
	})
	assert.Error(t, err)

	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	stream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	execFn = func(ctx queryctx.TaskContext, req *models.Request, mgr *SearchMgr) (any, error) {
		return &protoCommonV1.TaskResponse{}, nil
	}
	err = ip.Process(taskCtx, stream, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Data,
		Payload:      statement,
		PhysicalPlan: physicalPlan,
	})
	assert.NoError(t, err)
}
func TestProcessMetricMetadataSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		metricMetadataSearchFn = MetricMetadataSearch
		ctrl.Finish()
	}()

	physicalPlan := encoding.JSONMarshal(&models.PhysicalPlan{
		Targets: []*models.Target{{Indicator: "1.1.1.1:9000"}},
	})
	ip := NewIntermediateTaskProcessor(models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}, time.Second, nil, nil, nil)
	taskCtx := &flow.TaskContext{}
	err := ip.Process(taskCtx, nil, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      []byte("abc"),
		PhysicalPlan: physicalPlan,
	})
	assert.Error(t, err)

	metricMetadataSearchFn = func(ctx context.Context, param *models.ExecuteParam,
		statement *stmt.MetricMetadata, mgr *SearchMgr) (any, error) {
		return nil, fmt.Errorf("err")
	}
	statement, _ := (&stmt.MetricMetadata{}).MarshalJSON()
	err = ip.Process(taskCtx, nil, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      statement,
		PhysicalPlan: physicalPlan,
	})
	assert.Error(t, err)

	metricMetadataSearchFn = func(ctx context.Context, param *models.ExecuteParam,
		statement *stmt.MetricMetadata, mgr *SearchMgr) (any, error) {
		return []string{}, nil
	}
	stream := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	stream.EXPECT().Send(gomock.Any()).Return(nil)
	err = ip.Process(taskCtx, stream, &protoCommonV1.TaskRequest{
		RequestType:  protoCommonV1.RequestType_Metadata,
		Payload:      statement,
		PhysicalPlan: physicalPlan,
	})
	assert.NoError(t, err)
}
