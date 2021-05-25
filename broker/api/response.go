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
	"encoding/json"
	"net/http"
)

// OK responses with content and set the http status code 200
func OK(w http.ResponseWriter, a interface{}) {
	b, _ := json.Marshal(a)
	response(w, http.StatusOK, b)
}

// NoContent responses with empty content and set the http status code 204
func NoContent(w http.ResponseWriter) {
	response(w, http.StatusNoContent, nil)
}

// NotFound responses resource not found
func NotFound(w http.ResponseWriter) {
	response(w, http.StatusNotFound, nil)
}

// Error responses error message and set the http status code 500
func Error(w http.ResponseWriter, err error) {
	b, _ := json.Marshal(err.Error())
	response(w, http.StatusInternalServerError, b)
}

// response responses json body for http restful api
func response(w http.ResponseWriter, httpCode int, content []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	if len(content) > 0 {
		_, _ = w.Write(content)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
