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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/mock"
)

func TestGetEnv(t *testing.T) {
	r := gin.New()
	api := NewEnvAPI(
		config.Monitor{URL: "http://localhost?db=_internal"},
		"Broker",
	)
	api.Register(r)

	testCases := []struct {
		desc    string
		prepare func()
		assert  func(resp *httptest.ResponseRecorder)
	}{
		{
			desc: "parse url fail",
			prepare: func() {
				urlParseFn = func(_ string) (*url.URL, error) {
					return nil, fmt.Errorf("err")
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			desc: "parse params fail",
			prepare: func() {
				urlParseQueryFn = func(_ string) (url.Values, error) {
					return url.Values{}, fmt.Errorf("err")
				}
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			desc: "get monitor env successfully",
			prepare: func() {
			},
			assert: func(resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(_ *testing.T) {
			defer func() {
				urlParseFn = url.Parse
				urlParseQueryFn = url.ParseQuery
			}()
			tC.prepare()
			resp := mock.DoRequest(t, r, http.MethodGet, EnvPath, "")
			tC.assert(resp)
		})
	}
}
