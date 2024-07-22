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
	"github.com/go-http-utils/headers"
	"github.com/lindb/common/pkg/ltoml"
	"github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/series/metric"
)

func TestWrite_Flat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replica.NewMockChannelManager(ctrl)
	api := NewWrite(&deps.HTTPDeps{
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
			metrics.NewLimitStatistics("flat_write_test", linmetric.BrokerRegistry)),
	})
	r := gin.New()
	api.Register(r)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, WritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	// parse err
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&enrich_tag=a=b", "error")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	var brokerRow metric.BrokerRow
	err := converter.ConvertTo(&protoMetricsV1.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
	}, &brokerRow)
	assert.NoError(t, err)
	var buf bytes.Buffer
	_, _ = brokerRow.WriteTo(&buf)
	body := buf.String()

	header := make(http.Header)
	header.Set(headers.ContentType, constants.ContentTypeFlat)

	// write error
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test3&enrich_tag=a=b", body, header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns4&enrich_tag=a=b", body, header)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestWrite_Influx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replica.NewMockChannelManager(ctrl)
	api := NewWrite(&deps.HTTPDeps{
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
			metrics.NewLimitStatistics("test", linmetric.BrokerRegistry)),
	})
	r := gin.New()
	api.Register(r)

	header := make(http.Header)
	header.Set(headers.ContentType, constants.ContentTypeInflux)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, WritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// influx line format without timestamp
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns3&enrich_tag=a=b", `
# bad line
a,v=c,d=f a=2 b=3 c=4
`, header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// write error
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test3&enrich_tag=a=b", `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`, header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns4&enrich_tag=a=b", `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`, header)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestWrite_Proto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replica.NewMockChannelManager(ctrl)
	limits := models.NewDefaultLimits()
	limits.MaxNamespaceLength = 5
	limits.MaxTagNameLength = 5
	limits.MaxTagValueLength = 5
	models.SetDatabaseLimits("test", limits)
	models.SetDatabaseLimits("test3", limits)
	api := NewWrite(&deps.HTTPDeps{
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
			metrics.NewLimitStatistics("test", linmetric.BrokerRegistry)),
	})
	r := gin.New()
	api.Register(r)

	// missing db param
	resp := mock.DoRequest(t, r, http.MethodPut, WritePath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// enrich_tag bad format
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns2&enrich_tag=a", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// namespace too lang
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=namespace3&enrich_tag=a=b", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// tag key too lang
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test3&enrich_tag=system=b", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// tag value too lang
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test3&enrich_tag=ip=127.0.0.1", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	header := make(http.Header)
	header.Set(headers.ContentType, constants.ContentTypeProto)

	// bad format
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test&ns=ns3&enrich_tag=a=b", `xxxx`, header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// write error
	resp = mock.DoRequest(t, r, http.MethodPut, WritePath+"?db=test3&enrich_tag=a=b", `ok`, header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	metricList := protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{
		{Name: "1", Namespace: "ns", SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "counter", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 23},
		}},
	}}
	data, _ := metricList.Marshal()
	resp = mock.DoRequest(t, r, http.MethodPost, WritePath+"?db=test&ns=ns4&enrich_tag=a=b", string(data), header)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	cm.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	resp = mock.DoRequest(t, r, http.MethodPost, WritePath+"?db=test&ns=ns4&enrich_tag=a=b", string(data), header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
