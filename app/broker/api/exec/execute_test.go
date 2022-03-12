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

package exec

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	brokerQuery "github.com/lindb/lindb/query/broker"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

func TestExecuteAPI_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare
	repo := state.NewMockRepository(ctrl)
	master := coordinator.NewMockMasterController(ctrl)
	queryFactory := brokerQuery.NewMockFactory(ctrl)
	api := NewExecuteAPI(&deps.HTTPDeps{
		Ctx:          context.Background(),
		Repo:         repo,
		Master:       master,
		QueryFactory: queryFactory,
		BrokerCfg: &config.Broker{BrokerBase: config.BrokerBase{
			HTTP: config.HTTP{ReadTimeout: ltoml.Duration(time.Second * 10)},
		}},
		QueryLimiter: concurrent.NewLimiter(
			context.TODO(),
			2,
			time.Second*5,
			linmetric.NewScope("metric_data_search"),
		),
	})
	r := gin.New()
	api.Register(r)

	cases := []struct {
		name    string
		reqBody string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			name: "param invalid",
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "parse sql failure",
			reqBody: `{"sql":"show a"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "parse sql failure",
			reqBody: `{"sql":"show a"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "unknown state statement type",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt.Statement, error) {
					return &stmt.State{}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "unknown metadata statement type",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt.Statement, error) {
					return &stmt.Metadata{}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "unknown lin query language statement",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmt.Statement, err error) {
					return nil, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "master not found",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				master.EXPECT().GetMaster().Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "found master",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "get database list err",
			reqBody: `{"sql":"show databases"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "get database successfully, with one wrong data",
			reqBody: `{"sql":"show databases"}`,
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "get all database schemas",
			reqBody: `{"sql":"show schemas"}`,
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "database name cannot be empty when query metric",
			reqBody: `{"sql":"select f from cpu"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "query metric failure",
			reqBody: `{"sql":"select f from mem","db":"test"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetricQuery(ctrl)
				queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "query metric successfully",
			reqBody: `{"sql":"select f from mem","db":"test"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetricQuery(ctrl)
				queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return(&models.ResultSet{}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "get database list err",
			reqBody: `{"sql":"show databases"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "get database successfully, with one wrong data",
			reqBody: `{"sql":"show databases"}`,
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
				}
				database.Desc = database.String()
				data := encoding.JSONMarshal(&database)
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
					{Key: "db", Value: data},
					{Key: "err", Value: []byte{1, 2, 4}},
				}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "metadata query need input database",
			reqBody: `{"sql":"show namespaces"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "metadata query failure",
			reqBody: `{"sql":"show namespaces","db":"db"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetaDataQuery(ctrl)
				queryFactory.EXPECT().NewMetadataQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "metadata query successfully",
			reqBody: `{"sql":"show namespaces","db":"db"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetaDataQuery(ctrl)
				queryFactory.EXPECT().NewMetadataQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return([]string{"ns"}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show fields failure",
			reqBody: `{"sql":"show fields from cp","db":"db"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetaDataQuery(ctrl)
				queryFactory.EXPECT().NewMetadataQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return([]string{"ns"}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show fields successfully",
			reqBody: `{"sql":"show fields from cp","db":"db"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetaDataQuery(ctrl)
				queryFactory.EXPECT().NewMetadataQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				metricQuery.EXPECT().WaitResponse().Return([]string{string(encoding.JSONMarshal(&[]field.Meta{{Name: "test", Type: field.SumField}}))}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show histogram fields successfully",
			reqBody: `{"sql":"show fields from cp","db":"db"}`,
			prepare: func() {
				metricQuery := brokerQuery.NewMockMetaDataQuery(ctrl)
				queryFactory.EXPECT().NewMetadataQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
				// histogram
				metricQuery.EXPECT().WaitResponse().Return([]string{string(encoding.JSONMarshal(&[]field.Meta{
					{Name: "test", Type: field.SumField},
					{Name: "__bucket_0", Type: field.HistogramField},
					{Name: "__bucket_2", Type: field.HistogramField},
					{Name: "__bucket_3", Type: field.HistogramField},
					{Name: "__bucket_4", Type: field.HistogramField},
					{Name: "__bucket_99", Type: field.HistogramField},
					{Name: "histogram_sum", Type: field.SumField},
					{Name: "histogram_count", Type: field.SumField},
					{Name: "histogram_min", Type: field.MinField},
					{Name: "histogram_max", Type: field.MaxField},
				}))}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		//show fields from cp
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				sqlParseFn = sql.Parse
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			resp := mock.DoRequest(t, r, http.MethodPut, ExecutePath, tt.reqBody)
			if tt.assert != nil {
				tt.assert(resp)
			}
		})
	}
}
