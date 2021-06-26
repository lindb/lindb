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

func TestPrometheusWrite_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewPrometheusWriter(&deps.HTTPDeps{CM: cm})
	r := gin.New()
	api.Register(r)
	// case 1: param error
	resp := mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 2: read request body err
	resp = mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath+"?db=dal", "#$$#@#")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	//	// case 3: write wal err
	input := `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# 	TYPE go_gc_duration_seconds summary
go_gc_duration_seconds { quantile = "0.9999" } NaN
go_gc_duration_seconds_count 9
go_gc_duration_seconds_sum 90
//`
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(errors.New("err"))
	resp = mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath+"?db=dal", input)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// case 4: write wal success
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath+"?db=dal", input)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	// case 5: parse prometheus data err
	input = "# HELP go_gc_duration_seconds A summary of the GC invocation durations"
	resp = mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath+"?db=dal", input)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// case 6: enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, PrometheusWritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
