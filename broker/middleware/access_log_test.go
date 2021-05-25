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

package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessLogMiddleware(t *testing.T) {
	defer func() {
		pathUnescapeFunc = url.PathUnescape
	}()
	pathUnescapeFunc = func(s string) (string, error) {
		return "err-path", fmt.Errorf("err")
	}
	req, err := http.NewRequest("GET", "/health-check", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	accessLogHandler := AccessLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "ok")
	}))

	accessLogHandler.ServeHTTP(rr, req)

	accessLogHandler = AccessLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "ok")
	}))
	accessLogHandler.ServeHTTP(rr, req)
}

func Test_real_ip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health-check", nil)
	req.Header.Add("X-Real-Ip", "real-ip")
	assert.Equal(t, "real-ip", realIP(req))

	req, _ = http.NewRequest("GET", "/health-check", nil)
	req.Header.Add("X-Forwarded-For", "forward-ip")
	assert.Equal(t, "forward-ip", realIP(req))
	req, _ = http.NewRequest("GET", "/health-check", nil)
	req.RemoteAddr = "1.1.1.1:1023"
	assert.Equal(t, "1.1.1.1", realIP(req))
}
