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
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/broker/api/exec/command"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
	brokerQuery "github.com/lindb/lindb/query/broker"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

func TestExecuteAPI_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare
	ok := "ok"
	repo := state.NewMockRepository(ctrl)
	repoFct := state.NewMockRepositoryFactory(ctrl)
	master := coordinator.NewMockMasterController(ctrl)
	masterStateMgr := masterpkg.NewMockStateManager(ctrl)
	master.EXPECT().GetStateManager().Return(masterStateMgr).AnyTimes()
	queryFactory := brokerQuery.NewMockFactory(ctrl)
	stateMgr := broker.NewMockStateManager(ctrl)
	opt := &option.DatabaseOption{}
	api := NewExecuteAPI(&deps.HTTPDeps{
		Ctx:          context.Background(),
		Repo:         repo,
		RepoFactory:  repoFct,
		Master:       master,
		StateMgr:     stateMgr,
		QueryFactory: queryFactory,
		BrokerCfg: &config.Broker{BrokerBase: config.BrokerBase{
			HTTP: config.HTTP{ReadTimeout: ltoml.Duration(time.Second * 10)},
		}},
		QueryLimiter: concurrent.NewLimiter(
			context.TODO(),
			2,
			time.Second*5,
			metrics.NewLimitStatistics("exec", linmetric.BrokerRegistry),
		),
	})
	cfg := `{\"config\":{\"namespace\":\"test\",\"timeout\":10,\"dialTimeout\":10,`
	cfg += `\"leaseTTL\":10,\"endpoints\":[\"http://localhost:2379\"]}}`
	databaseCfg := `{\"name\":\"test\",\"storage\":\"cluster-test\",\"numOfShard\":12,`
	databaseCfg += `\"replicaFactor\":3,\"option\":{\"intervals\":[{\"interval\":\"10s\"}]}}`
	r := gin.New()
	api.Register(r)

	var backend *httptest.Server

	databaseCfgData := encoding.JSONMarshal(map[string]models.DatabaseConfig{
		"test": {
			ShardIDs: []models.ShardID{1, 2, 3},
			Option:   &option.DatabaseOption{},
		},
	})
	mockSrv := func(data []byte) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Add("content-type", "application/json")
			_, _ = rw.Write(data)
		}))
		u, err := url.Parse(server.URL)
		assert.NoError(t, err)
		p, err := strconv.Atoi(u.Port())
		assert.NoError(t, err)
		stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
			LiveNodes: map[models.NodeID]models.StatefulNode{1: {
				StatelessNode: models.StatelessNode{
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				},
				ID: 1,
			}}}, true)
	}

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
			name:    "unknown metadata statement type",
			reqBody: `{"sql":"show master"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.State{}, nil
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
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Use{}, nil
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
					Option:        opt,
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
			name:    "get database successfully, with one wrong data",
			reqBody: `{"sql":"show databases"}`,
			prepare: func() {
				// get ok
				database := models.Database{
					Name:          "test",
					Storage:       "cluster-test",
					NumOfShard:    12,
					ReplicaFactor: 3,
					Option:        opt,
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
			name:    "schema query, unknown metadata type",
			reqBody: `{"sql":"show database"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Schema{}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
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
					Option:        opt,
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
				queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
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
				queryFactory.EXPECT().NewMetricQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(metricQuery)
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
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.MetricMetadata{Type: stmtpkg.Field, Limit: 0}, nil
				}
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
				metricQuery.EXPECT().WaitResponse().
					Return([]string{string(encoding.JSONMarshal(&[]field.Meta{{Name: "test", Type: field.SumField}}))}, nil)
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
		{
			name:    "unknown storage op type",
			reqBody: `{"sql":"show storages"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Storage{Type: stmtpkg.StorageOpUnknown}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storages, get storages failure",
			reqBody: `{"sql":"show storages"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show storages, list storage successfully, but unmarshal failure",
			reqBody: `{"sql":"show storages"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte("[]")}}, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show storages successfully",
			reqBody: `{"sql":"show storages"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetStorage("xxx").Return(nil, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show storages successfully,but state not found",
			reqBody: `{"sql":"show storages"}`,
			prepare: func() {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(
					[]state.KeyValue{{Key: "", Value: []byte(`{ "config": {"namespace":"xxx"}}`)}}, nil)
				stateMgr.EXPECT().GetStorage("xxx").Return(nil, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "create storage json err",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Storage{Type: stmtpkg.StorageOpCreate, Value: "xx"}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create storage, config validate failure",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Storage{Type: stmtpkg.StorageOpCreate, Value: `{"config":{}}`}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create storage successfully, storage not exist",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						if err := check([]byte{1, 2, 3}); err != nil {
							return false, err
						}
						return true, nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "create storage successfully, storage exist",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ []byte, check func([]byte) error) (bool, error) {
						cfg1 := strings.ReplaceAll(cfg, `\"`, `"`)
						data := []byte(cfg1)
						storage := &config.StorageCluster{}
						err := encoding.JSONUnmarshal(data, storage)
						assert.NoError(t, err)
						data = encoding.JSONMarshal(storage)
						if err := check(data); err != nil {
							return false, err
						}
						return true, nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "create storage failure with err",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create storage failure",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(nil)
				repo.EXPECT().PutWithTX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create storage repo failure",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create storage, close repo failure",
			reqBody: `{"sql":"create storage ` + cfg + `"}`,
			prepare: func() {
				repoFct.EXPECT().CreateStorageRepo(gomock.Any()).Return(repo, nil)
				repo.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show broker alive node",
			reqBody: `{"sql":"show broker alive"}`,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "1.1.1.1",
					HTTPPort: 8080,
				}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show storage alive node",
			reqBody: `{"sql":"show storage alive"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorageList().Return([]*models.StorageState{})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show memory database state, but storage not found",
			reqBody: `{"sql":"show memory database where storage=a and database=b"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show replication state, but storage not found",
			reqBody: `{"sql":"show replication where storage=a and database=b"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show replication state, alive node empty",
			reqBody: `{"sql":"show replication where storage=a and database=b"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
					LiveNodes: nil}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show replication state, but fetch state failure",
			reqBody: `{"sql":"show replication where storage=a and database=b"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
					LiveNodes: map[models.NodeID]models.StatefulNode{1: {
						StatelessNode: models.StatelessNode{
							HostIP:   "127.0.01", // mock host err
							HTTPPort: 8080,
						},
						ID: 1,
					}}}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show replication state, but fetch state failure",
			reqBody: `{"sql":"show replication where storage=a and database=b"}`,
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("[]"))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{
					LiveNodes: map[models.NodeID]models.StatefulNode{1: {
						StatelessNode: models.StatelessNode{
							HostIP:   u.Hostname(),
							HTTPPort: uint16(p),
						},
						ID: 1,
					}}}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show broker metric, no alive node",
			reqBody: `{"sql":"show broker metric where metric in (a,b)"}`,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show broker metric, fetch metric failure",
			reqBody: `{"sql":"show broker metric where metric in (a,b)"}`,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "127.0.0.1",
					HTTPPort: 8080,
				}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show broker metric successfully",
			reqBody: `{"sql":"show broker metric where metric in (a,b)"}`,
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`{"cpu":[{"fields":[{"value":1}]},{"fields":[{"value":1}]}]}`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}, {
					HostIP:   u.Hostname(),
					HTTPPort: uint16(p),
				}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show storage metric, storage name empty",
			reqBody: `{"sql":"show storage metric where storage='' and metric in (a,b)"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show storage metric, storage not found",
			reqBody: `{"sql":"show storage metric where storage='a' and metric in (a,b)"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storage metric, storage no alive node",
			reqBody: `{"sql":"show storage metric where storage='a' and metric in (a,b)"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(&models.StorageState{}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storage metric successfully",
			reqBody: `{"sql":"show storage metric where storage='a' and metric in (a,b)"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).
					Return(&models.StorageState{LiveNodes: map[models.NodeID]models.StatefulNode{1: {}, 2: {}}}, true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show metadata path successfully",
			reqBody: `{"sql":"show metadata types"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "state from state machine, but type not found",
			reqBody: `{"sql":"show broker metadata from state_machine where type=abc"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Metadata{MetadataType: stmtpkg.MetadataType(1000),
						Source: stmtpkg.StateMachineSource}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "state from state machine, but source not found",
			reqBody: `{"sql":"show broker metadata from state_machine where type=abc"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Metadata{MetadataType: stmtpkg.BrokerMetadata,
						Source: stmtpkg.SourceType(100)}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "state from state machine, broker state",
			reqBody: `{"sql":"show broker metadata from state_machine where type=abc"}`,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{}})
				cli := client.NewMockStateMachineCli(ctrl)
				command.NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNodes(gomock.Any(), gomock.Any()).Return(&ok)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "state from state machine, master state",
			reqBody: `{"sql":"show master metadata from state_machine where type=abc"}`,
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{}})
				cli := client.NewMockStateMachineCli(ctrl)
				command.NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNodes(gomock.Any(), gomock.Any()).Return(&ok)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "state from state machine, storage state",
			reqBody: `{"sql":"show storage metadata from state_machine where type=abc and storage=xx"}`,
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{}})
				cli := client.NewMockStateMachineCli(ctrl)
				command.NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNode(gomock.Any(), gomock.Any()).Return(&ok, nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show broker metadata, but type not found",
			reqBody: `{"sql":"show broker metadata from state_repo where type=abc"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show broker metadata, but walk entry repo failure",
			reqBody: `{"sql":"show broker metadata from state_repo where type=LiveNode"}`,
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show broker metadata, but walk entry unmarshal failure",
			reqBody: `{"sql":"show broker metadata from state_repo where type=LiveNode"}`,
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), []byte("value"))
						return nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show broker metadata, get live node successfully",
			reqBody: `{"sql":"show broker metadata from state_repo where type=DatabaseConfig"}`,
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), encoding.JSONMarshal(&models.Database{Name: "1.1.1.1"}))
						return nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show master metadata, get master successfully",
			reqBody: `{"sql":"show master metadata from state_repo where type=Master"}`,
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), encoding.JSONMarshal(&models.Master{ElectTime: 11}))
						return nil
					})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show storage metadata, but storage name empty",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode and storage=''"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show storage metadata, but type not found",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode1 and storage='abc'"}`,
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storage metadata, but storage state not found",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode and storage='test'"}`,
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
				masterStateMgr.EXPECT().GetStorageCluster("test").Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storage metadata, no data",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode and storage='test'"}`,
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
				storage := masterpkg.NewMockStorageCluster(ctrl)
				masterStateMgr.EXPECT().GetStorageCluster("test").Return(storage)
				storage.EXPECT().GetRepo().Return(repo)
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show storage metadata, forward request failure",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode and storage='test'"}`,
			prepare: func() {
				master.EXPECT().IsMaster().Return(false)
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: 8089}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "show storage metadata, forward request successfully",
			reqBody: `{"sql":"show storage metadata from state_repo where type=LiveNode and storage='test'"}`,
			prepare: func() {
				port := uint16(8789)
				master.EXPECT().IsMaster().Return(false)
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: port}})
				backend = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("test"))
				}))
				// hack
				_ = backend.Listener.Close()
				l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				assert.NoError(t, err)
				backend.Listener = l
				// Start the server.
				backend.Start()
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "create database config unmarshal failure",
			reqBody: `{"sql":"create database {\"name\":\"name\"}"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Schema{Type: stmtpkg.CreateDatabaseSchemaType, Value: "err"}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create database validation failure",
			reqBody: `{"sql":"create database {\"name\":\"name\"}"}`,
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create database, persist failure",
			reqBody: `{"sql":"create database ` + databaseCfg + `"}`,
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create database, option validation failure",
			reqBody: `{"sql":"create database ` + databaseCfg + `"}`,
			prepare: func() {
				sqlParseFn = func(sql string) (stmt stmtpkg.Statement, err error) {
					return &stmtpkg.Schema{
						Type: stmtpkg.CreateDatabaseSchemaType,
						Value: string(encoding.JSONMarshal(&models.Database{
							Name:          "test",
							Storage:       "cluster-test",
							NumOfShard:    12,
							ReplicaFactor: 3,
							Option: &option.DatabaseOption{
								Intervals: option.Intervals{{Interval: 10}},
								Ahead:     "10",
							},
						})),
					}, nil
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "create database successfully",
			reqBody: `{"sql":"create database ` + databaseCfg + `"}`,
			prepare: func() {
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "drop database, but delete cfg failure",
			reqBody: `{"sql":"drop database test"}`,
			prepare: func() {
				// delete database cfg failure
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "drop database, but delete shard assignment failure",
			reqBody: `{"sql":"drop database test"}`,
			prepare: func() {
				// delete database cfg ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				// delete database shard assignment failure
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "drop database successfully",
			reqBody: `{"sql":"drop database test"}`,
			prepare: func() {
				// delete database cfg ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				// delete database shard assignment ok
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "show all requests, but no alive broker",
			reqBody: `{"sql":"show requests"}`,
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show all requests, but get err from broker",
			reqBody: `{"sql":"show requests"}`,
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("{}"))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "127.0.0.1",
					HTTPPort: uint16(p),
				}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name:    "show all requests successfully",
			reqBody: `{"sql":"show requests"}`,
			prepare: func() {
				svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("content-type", "application/json")
					_, _ = w.Write([]byte(`[{"start":12314}]`))
				}))
				u, err := url.Parse(svr.URL)
				assert.NoError(t, err)
				p, err := strconv.Atoi(u.Port())
				assert.NoError(t, err)
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
					HostIP:   "127.0.0.1",
					HTTPPort: uint16(p),
				}})
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "recover storage, but storage not found",
			reqBody: `{"sql":"recover storage test"}`,
			prepare: func() {
				stateMgr.EXPECT().GetStorage(gomock.Any()).Return(nil, false)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "recover storage, but get database config failure",
			reqBody: `{"sql":"recover storage test"}`,
			prepare: func() {
				mockSrv([]byte("abc"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name:    "recover storage, but recover shard assignment failure",
			reqBody: `{"sql":"recover storage test"}`,
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "recover storage, but recover database schema failure",
			reqBody: `{"sql":"recover storage test"}`,
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name:    "recover storage successfully",
			reqBody: `{"sql":"recover storage test"}`,
			prepare: func() {
				mockSrv(databaseCfgData)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				sqlParseFn = sql.Parse
				command.NewRestyFn = resty.New
				command.NewStateMachineCliFn = client.NewStateMachineCli
			}()
			command.NewRestyFn = func() *resty.Client {
				c := resty.New()
				c.SetTimeout(time.Second)
				return c
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			resp := mock.DoRequest(t, r, http.MethodPut, ExecutePath, tt.reqBody)
			if tt.assert != nil {
				tt.assert(resp)
			}
			if backend != nil {
				backend.Close()
				backend = nil
			}
		})
	}
}
