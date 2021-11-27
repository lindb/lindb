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

package api

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/mock"
)

func TestReverseProxy_Proxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewReverseProxy()
	r := gin.New()
	api.Register(r)

	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	}))
	// hack
	_ = backend.Listener.Close()
	l, err := net.Listen("tcp", "127.0.0.1:8089")
	assert.NoError(t, err)
	backend.Listener = l
	// Start the server.
	backend.Start()
	resp := mock.DoRequest(t, r, http.MethodGet, ProxyPath+"?target=127.0.0.1:8089&path="+ProxyPath, "")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Body.String())

	resp = mock.DoRequest(t, r, http.MethodGet, ProxyPath, "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
