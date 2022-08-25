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

package state

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/query"
)

func TestRequestAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	pipeline := query.NewMockPipeline(ctrl)
	api := NewRequestAPI()
	r := gin.New()
	api.Register(r)

	t.Run("get all requests", func(t *testing.T) {
		resp := mock.DoRequest(t, r, http.MethodGet, RequestsPath, "")
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("param invalid", func(t *testing.T) {
		resp := mock.DoRequest(t, r, http.MethodGet, RequestPath, "")
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("reuqest not found", func(t *testing.T) {
		resp := mock.DoRequest(t, r, http.MethodGet, RequestPath+"?requestId=id", "")
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("found reuqest", func(t *testing.T) {
		query.GetPipelineManager().AddPipeline("id", pipeline)
		pipeline.EXPECT().Stats().Return(nil)
		resp := mock.DoRequest(t, r, http.MethodGet, RequestPath+"?requestId=id", "")
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
