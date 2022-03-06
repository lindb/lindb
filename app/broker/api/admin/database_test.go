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
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseAPI_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := gin.New()
	repo := state.NewMockRepository(ctrl)
	api := NewDatabaseAPI(&deps.HTTPDeps{
		Ctx:  context.Background(),
		Repo: repo,
		BrokerCfg: &config.Broker{BrokerBase: config.BrokerBase{
			HTTP: config.HTTP{ReadTimeout: ltoml.Duration(time.Second * 10)},
		}},
	})
	api.Register(r)

	// bind error
	reps := mock.DoRequest(t, r, http.MethodPost, DatabasePath, "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// database name empty
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, `{"name":""}`)
	assert.Equal(t, http.StatusInternalServerError, reps.Code)
	// cluster empty
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, `{"name":"23"}`)
	assert.Equal(t, http.StatusInternalServerError, reps.Code)
	// num shards < 0
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath,
		`{"name":"23", "storage": "xxx", "numOfShard": -2}`)
	assert.Equal(t, http.StatusInternalServerError, reps.Code)
	// ReplicaFactor < 0
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath,
		`{"name":"23", "storage": "xxx", "numOfShard": 2, "replicaFactor": -1}`)
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// validate error
	database := models.Database{
		Name:          "test",
		Storage:       "cluster-test",
		NumOfShard:    12,
		ReplicaFactor: 3,
		Option:        option.DatabaseOption{},
	}
	data := encoding.JSONMarshal(&database)
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, string(data))
	assert.Equal(t, http.StatusInternalServerError, reps.Code)
	// put
	database.Option = option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}
	data = encoding.JSONMarshal(&database)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	reps = mock.DoRequest(t, r, http.MethodPost, DatabasePath, string(data))
	assert.Equal(t, http.StatusNoContent, reps.Code)
}

func TestDatabaseAPI_GetByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := gin.New()
	repo := state.NewMockRepository(ctrl)
	api := NewDatabaseAPI(&deps.HTTPDeps{
		Ctx:  context.Background(),
		Repo: repo,
		BrokerCfg: &config.Broker{BrokerBase: config.BrokerBase{
			HTTP: config.HTTP{ReadTimeout: ltoml.Duration(time.Second * 10)}}},
	})
	api.Register(r)

	reps := mock.DoRequest(t, r, http.MethodGet, DatabasePath, "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// name empty
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=", "")
	assert.Equal(t, http.StatusInternalServerError, reps.Code)

	// get bad content
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte("bad-data"), nil)
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=xxx", "")
	assert.Equal(t, http.StatusNotFound, reps.Code)

	// get error
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, io.ErrClosedPipe)
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=xxx", "")
	assert.Equal(t, http.StatusNotFound, reps.Code)

	// get ok
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte(`{"name":"xxx"}`), nil)
	reps = mock.DoRequest(t, r, http.MethodGet, DatabasePath+"?name=xxx", "")
	assert.Equal(t, http.StatusOK, reps.Code)
}
