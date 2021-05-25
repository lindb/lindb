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

package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HTTPHandler represents amock http handler
type HTTPHandler struct {
	Method      string
	URL         string
	RequestBody interface{}

	HandlerFunc http.HandlerFunc

	ExpectHTTPCode int
	ExpectResponse interface{}
}

// DoRequest helps test that whether the http response with the given
// method,url,requestJSON,handlerFunc matches with the expectHTTPCode and
// expectResponse. if false the test will be not failed
func DoRequest(t *testing.T, httpHandler *HTTPHandler) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	requestBodyBytes, err := json.Marshal(httpHandler.RequestBody)
	if err != nil {
		t.Fatal(err)
		return
	}
	reader := bytes.NewReader(requestBodyBytes)
	req, err := http.NewRequest(httpHandler.Method, httpHandler.URL, reader)
	req.RequestURI = httpHandler.URL
	if err != nil {
		t.Fatal(err)
		return
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to
	// record the response.
	rr := httptest.NewRecorder()

	handler := httpHandler.HandlerFunc

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	assert.Equal(t, httpHandler.ExpectHTTPCode, rr.Code)

	// verify that the response matches the response
	if httpHandler.ExpectResponse != nil {
		readCloser := rr.Result().Body
		defer readCloser.Close()
		responseBytes, err := ioutil.ReadAll(readCloser)
		if err != nil {
			t.Fatal(err)
			return
		}
		expect, _ := json.Marshal(httpHandler.ExpectResponse)
		assert.Equal(t, expect, responseBytes)
	}
}
