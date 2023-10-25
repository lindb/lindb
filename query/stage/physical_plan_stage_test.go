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

package stage

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
)

func TestPhysicalPlanStage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskCtx := context.NewMockTaskContext(ctrl)
	taskCtx.EXPECT().GetRequests().Return(map[string]*protoCommonV1.TaskRequest{
		"target": nil,
	})
	s := NewPhysicalPlanStage(taskCtx)
	assert.Equal(t, "Physical Plan", s.Identifier())
	assert.NotNil(t, s.Plan())
	assert.NotNil(t, s.NextStages())
}
