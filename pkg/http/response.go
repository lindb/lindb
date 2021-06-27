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
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK responses with content and set the http status code 200.
func OK(c *gin.Context, content interface{}) {
	response(c, http.StatusOK, content)
}

// NoContent responses with empty content and set the http status code 204.
func NoContent(c *gin.Context) {
	response(c, http.StatusNoContent, nil)
}

// NotFound responses resource not found.
func NotFound(c *gin.Context) {
	_ = c.Error(errors.New("StatusNotFound"))
	response(c, http.StatusNotFound, nil)
}

// Error responses error message and set the http status code 500.
func Error(c *gin.Context, err error) {
	_ = c.Error(err)
	response(c, http.StatusInternalServerError, err.Error())
}

// response responses json body for http restful api
func response(c *gin.Context, httpCode int, content interface{}) {
	c.JSON(httpCode, content)
}
