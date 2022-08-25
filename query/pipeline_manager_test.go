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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetPipelineManager(t *testing.T) {
	assert.NotNil(t, GetPipelineManager())
	assert.NotNil(t, GetPipelineManager())
}

func TestPipelineManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pipeline := NewMockPipeline(ctrl)
	mgr := newPipelineManager()

	req := "requestID"
	assert.Empty(t, mgr.GetAllAlivePipelines())

	mgr.AddPipeline(req, pipeline)
	assert.Len(t, mgr.GetAllAlivePipelines(), 1)
	assert.NotNil(t, mgr.GetPipeline(req))

	mgr.RemovePipeline(req)
	assert.Empty(t, mgr.GetAllAlivePipelines())
}
