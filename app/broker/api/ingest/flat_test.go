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

package ingest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/series/metric"
)

func Test_Flat_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replica.NewMockChannelManager(ctrl)
	api := NewFlatWriter(&deps.HTTPDeps{
		BrokerCfg: &config.Broker{
			BrokerBase: config.BrokerBase{
				Ingestion: config.Ingestion{
					IngestTimeout: ltoml.Duration(time.Second * 2),
				},
			},
		},
		CM: cm,
		IngestLimiter: concurrent.NewLimiter(
			context.TODO(),
			32,
			time.Second,
			linmetric.NewScope("influx_write_test")),
	})
	r := gin.New()
	api.Register(r)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, FlatWritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, FlatWritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// parse err
	resp = mock.DoRequest(t, r, http.MethodPut, FlatWritePath+"?db=test&enrich_tag=a=b", "error")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	converter := metric.NewProtoConverter()
	var brokerRow metric.BrokerRow
	err := converter.ConvertTo(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
	}, &brokerRow)
	assert.NoError(t, err)
	var buf bytes.Buffer
	_, _ = brokerRow.WriteTo(&buf)
	body := buf.String()
	// write error
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPut, FlatWritePath+"?db=test3&enrich_tag=a=b", body)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, FlatWritePath+"?db=test&ns=ns4&enrich_tag=a=b", body)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}
