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

package admin

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/service"
)

func TestDatabaseAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := gin.New()

	databaseService := service.NewMockDatabaseService(ctrl)

	api := NewDatabaseAPI(&deps.HTTPDeps{
		DatabaseSrv: databaseService,
	})
	api.Register(r)

	db := models.Database{
		Name:          "test",
		Cluster:       "test",
		NumOfShard:    12,
		ReplicaFactor: 3,
		Option:        option.DatabaseOption{Interval: "10s"},
	}

	// get request error
	reps := mock.DoRequest(t, r, http.MethodPost, DatabasePath, "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// create success
	databaseService.EXPECT().Save(gomock.Any()).Return(nil)
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, `{"name":"test"}`)
	assert.Equal(t, http.StatusNoContent, reps.Code)
	// create err
	databaseService.EXPECT().Save(gomock.Any()).Return(fmt.Errorf("err"))
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, `{"name":"test"}`)
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// get success
	databaseService.EXPECT().Get(gomock.Any()).Return(&db, nil)
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=test", "")
	assert.Equal(t, http.StatusOK, reps.Code)
	// no database name
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath, "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	databaseService.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=test", "")
	assert.Equal(t, http.StatusNotFound, reps.Code)

	databaseService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	reps = mock.DoRequest(t, r, http.MethodGet, ListDatabasePath, "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	databaseService.EXPECT().List().Return([]*models.Database{&db}, nil)
	reps = mock.DoRequest(t, r, http.MethodGet, ListDatabasePath, "")
	assert.Equal(t, http.StatusOK, reps.Code)
}
