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
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/replication"
)

func Test_Influx_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewInfluxWriter(&deps.HTTPDeps{CM: cm})
	r := gin.New()
	api.Register(r)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, InfluxWritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, InfluxWritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// bad influx line format
	resp = mock.DoRequest(t, r, http.MethodPut, InfluxWritePath+"?db=test&ns=ns3&enrich_tag=a=b", `
# bad line
a,v=c,d=f a=2 b=3 c=4
`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// write error
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPut, InfluxWritePath+"?db=test3&enrich_tag=a=b", `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, InfluxWritePath+"?db=test&ns=ns4&enrich_tag=a=b", `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}
