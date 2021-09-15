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
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	brokerQuery "github.com/lindb/lindb/query/broker"
)

func TestMetricAPI_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFactory := brokerQuery.NewMockFactory(ctrl)
	metricQuery := brokerQuery.NewMockMetricQuery(ctrl)
	stateMgr := broker.NewMockStateManager(ctrl)

	queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)

	api := NewMetricAPI(&deps.HTTPDeps{
		BrokerCfg:    &config.Broker{Query: config.Query{Timeout: ltoml.Duration(time.Second)}},
		StateMgr:     stateMgr,
		QueryFactory: queryFactory,
		QueryLimiter: concurrent.NewLimiter(
			context.TODO(),
			2,
			time.Second*5,
			linmetric.NewScope("metric_data_search"),
		),
	})
	r := gin.New()
	api.Register(r)

	metricQuery.EXPECT().WaitResponse().Return(&models.ResultSet{}, nil)
	resp := mock.DoRequest(t, r, http.MethodGet, MetricQueryPath+"?db=test&sql=select f from cpu", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestNewMetricAPI_Search_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFactory := brokerQuery.NewMockFactory(ctrl)
	stateMgr := broker.NewMockStateManager(ctrl)

	api := NewMetricAPI(&deps.HTTPDeps{
		BrokerCfg:    &config.Broker{Query: config.Query{Timeout: ltoml.Duration(time.Second)}},
		QueryFactory: queryFactory,
		StateMgr:     stateMgr,
		QueryLimiter: concurrent.NewLimiter(
			context.TODO(),
			2,
			time.Second*5,
			linmetric.NewScope("metric_data_search"),
		),
	})
	r := gin.New()
	api.Register(r)

	// param error
	resp := mock.DoRequest(t, r, http.MethodGet, MetricQueryPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	metricQuery := brokerQuery.NewMockMetricQuery(ctrl)
	queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
	metricQuery.EXPECT().WaitResponse().Return(&models.ResultSet{}, fmt.Errorf("err"))

	resp = mock.DoRequest(t, r, http.MethodGet, MetricQueryPath+"?db=test&sql=select f from cpu", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
