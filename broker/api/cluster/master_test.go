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

package cluster

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestMasterAPI_GetMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	master := coordinator.NewMockMaster(ctrl)

	api := NewMasterAPI(&deps.HTTPDeps{
		Master: master,
	})
	r := gin.New()
	api.Register(r)

	m := models.Master{ElectTime: timeutil.Now(), Node: models.Node{IP: "1.1.1.1", Port: 8000}}
	// get success
	master.EXPECT().GetMaster().Return(&m)
	resp := mock.DoRequest(t, r, http.MethodGet, MasterStatePath, "")
	assert.Equal(t, http.StatusOK, resp.Code)

	master.EXPECT().GetMaster().Return(nil)
	resp = mock.DoRequest(t, r, http.MethodGet, MasterStatePath, "")
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
