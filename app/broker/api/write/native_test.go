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
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"io"
	"net/http"
	"testing"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/mock"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/replication"
)

func Test_NativeWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewNativeWriter(&deps.HTTPDeps{CM: cm})
	r := gin.New()
	api.Register(r)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, NativeWritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, NativeWritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// bad format
	resp = mock.DoRequest(t, r, http.MethodPut, NativeWritePath+"?db=test&ns=ns3&enrich_tag=a=b", `xxxx`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// write error
	resp = mock.DoRequest(t, r, http.MethodPut, NativeWritePath+"?db=test3&enrich_tag=a=b", `ok`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	var metricList = protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
		{Name: "1", Namespace: "ns", SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "counter", Type: protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM, Value: 23},
		}},
	}}
	data, _ := metricList.Marshal()
	resp = mock.DoRequest(t, r, http.MethodPost, NativeWritePath+"?db=test&ns=ns4&enrich_tag=a=b", string(data))
	assert.Equal(t, http.StatusNoContent, resp.Code)

	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPost, NativeWritePath+"?db=test&ns=ns4&enrich_tag=a=b", string(data))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

}
