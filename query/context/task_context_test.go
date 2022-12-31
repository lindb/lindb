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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

func TestTaskContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transportMgr := rpc.NewMockTransportManager(ctrl)
	transportMgr.EXPECT().SendRequest(gomock.Any(), gomock.Any()).Return(nil)
	ctx := newBaseTaskContext(context.TODO(), transportMgr)
	assert.NotNil(t, ctx.Context())

	req := protoCommonV1.TaskRequest{}
	assert.NoError(t, ctx.SendRequest("target", &req))

	ctx.addRequests(&req, &models.PhysicalPlan{
		Targets: []*models.Target{
			{Indicator: "target-1"},
			{Indicator: "target-2"},
			{Indicator: "target-1"},
		},
	})
	assert.Len(t, ctx.GetRequests(), 2)

	ctx.Complete(fmt.Errorf("err"))
	ctx.Complete(fmt.Errorf("err"))
	ctx.handleTaskState(&protoCommonV1.TaskResponse{Completed: false}, "leaf")
	ctx.handleTaskState(&protoCommonV1.TaskResponse{Completed: true}, "leaf")
}
