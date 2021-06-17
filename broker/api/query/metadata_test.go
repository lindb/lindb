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
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataAPI_Handle_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		parseSQLFunc = parseSQL

		ctrl.Finish()
	}()

	api := NewMetadataAPI(&deps.HTTPDeps{})
	r := gin.New()
	api.Register(r)
	// case 1: database name not input
	resp := mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// case 2: parse sql err
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=show d", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// case 3: wrong type
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=select f1 from cpu", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// case 4: unknown metadata type
	parseSQLFunc = func(ql string) (*stmt.Metadata, error) {
		return &stmt.Metadata{}, nil
	}
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=select f1 from cpu", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestMetadataAPI_ShowDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	databaseService := service.NewMockDatabaseService(ctrl)
	api := NewMetadataAPI(&deps.HTTPDeps{DatabaseSrv: databaseService})
	r := gin.New()
	api.Register(r)

	databaseService.EXPECT().List().Return(nil, nil)
	resp := mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?sql=show databases", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	databaseService.EXPECT().List().Return(
		[]*models.Database{
			{Name: "test1"},
			{Name: "test2"},
		},
		nil)
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?sql=show databases", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	databaseService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?sql=show databases", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestMetadataAPI_SuggestCommon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := parallel.NewMockExecutorFactory(ctrl)
	exec := parallel.NewMockMetadataExecutor(ctrl)

	factory.EXPECT().NewMetadataBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(exec).AnyTimes()

	api := NewMetadataAPI(&deps.HTTPDeps{ExecutorFct: factory, StateMachines: &coordinator.BrokerStateMachines{}})
	r := gin.New()
	api.Register(r)

	resp := mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?sql=show namespaces", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	exec.EXPECT().Execute().Return(nil, fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=show namespaces", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	exec.EXPECT().Execute().Return([]string{"a", "b"}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=show namespaces", "")
	assert.Equal(t, http.StatusOK, resp.Code)

	exec.EXPECT().Execute().Return([]string{"ddd"}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=show fields from cpu", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	exec.EXPECT().Execute().Return([]string{string(encoding.JSONMarshal(&[]field.Meta{{Name: "test", Type: field.SumField}}))}, nil)
	resp = mock.DoRequest(t, r, http.MethodGet, MetadataQueryPath+"?db=db&sql=show fields from cpu", "")
	assert.Equal(t, http.StatusOK, resp.Code)
}
