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

package write

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/replication"
)

func TestWriteAPI_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewWriteAPI(&deps.HTTPDeps{CM: cm})
	r := gin.New()
	api.Register(r)
	// param error
	resp := mock.DoRequest(t, r, http.MethodPut, SumWritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, SumWritePath+"?db=dal&cluster=dal&count=1", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(errors.New("err")).AnyTimes()
	resp = mock.DoRequest(t, r, http.MethodPut, SumWritePath+"?db=dal&cluster=dal", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
