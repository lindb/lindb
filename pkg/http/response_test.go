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

package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	OK(c, "ok")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `"ok"`, resp.Body.String())
}

func TestNoContent(t *testing.T) {
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	NoContent(c)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	assert.Equal(t, 0, resp.Body.Len())
}

func TestNotFound(t *testing.T) {
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	NotFound(c)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, 4, resp.Body.Len())
}

func TestError(t *testing.T) {
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	Error(c, fmt.Errorf("err"))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, `"err"`, resp.Body.String())
}
