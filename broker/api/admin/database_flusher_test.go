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
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
)

type mockIOReader struct {
}

func (m *mockIOReader) Close() error {
	return fmt.Errorf("err")
}

func (m *mockIOReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("err")
}

func TestNewDatabaseFlusherAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	master := coordinator.NewMockMaster(ctrl)
	flushAPI := NewDatabaseFlusherAPI(&deps.HTTPDeps{
		Master: master,
	})
	r := gin.New()
	flushAPI.Register(r)

	// no cluster
	resp := mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, "{}")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// no database name
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, ``)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// submit err
	master.EXPECT().IsMaster().Return(true)
	master.EXPECT().FlushDatabase(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	// submit ok
	master.EXPECT().IsMaster().Return(true)
	master.EXPECT().FlushDatabase(gomock.Any(), gomock.Any()).Return(nil)
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	defer func() {
		httpGet = http.Get
	}()

	// forward master
	master.EXPECT().IsMaster().Return(false)
	master.EXPECT().GetMaster().Return(&models.Master{
		Node: models.Node{
			IP:   "127.0.0.1",
			Port: 12345,
		},
	})
	httpGet = func(url string) (resp *http.Response, err error) {
		return nil, fmt.Errorf("err")
	}
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	httpGet = func(url string) (resp *http.Response, err error) {
		return nil, nil
	}
	master.EXPECT().IsMaster().Return(false)
	master.EXPECT().GetMaster().Return(&models.Master{
		Node: models.Node{
			IP:   "127.0.0.1",
			Port: 12345,
		},
	})
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusOK, resp.Code)
	httpGet = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	master.EXPECT().IsMaster().Return(false)
	master.EXPECT().GetMaster().Return(&models.Master{
		Node: models.Node{
			IP:   "127.0.0.1",
			Port: 12346,
		},
	})
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	httpGet = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &mockIOReader{},
		}, nil
	}
	master.EXPECT().IsMaster().Return(false)
	master.EXPECT().GetMaster().Return(&models.Master{
		Node: models.Node{
			IP:   "127.0.0.1",
			Port: 12346,
		},
	})
	resp = mock.DoRequest(t, r, http.MethodPut, FlushDatabasePath, `{"cluster":"test","database":"db"}`)
	assert.Equal(t, http.StatusOK, resp.Code)
}
