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
	"fmt"
	"net/http"
)

// GetJSONBodyFromRequest gets json from request body and then parses into specified struct
func GetJSONBodyFromRequest(r *http.Request, t interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&t)
}

// GetParamsFromRequest gets parameter value from the request。
// If the request method is neither GET or POST,returns an error.
// If there are multiple parameters with the same paramsName,returns the first value.
// If there does not have the value and required is false, returns an error, otherwise returns the defaultValue
func GetParamsFromRequest(paramsName string, r *http.Request, defaultValue string, required bool) (string, error) {
	if len(paramsName) == 0 {
		return "", fmt.Errorf("the params name must not be null")
	}
	var value string
	method := r.Method
	//get value from different object according to different request methods
	switch method {
	// the parameter value need to be parsed from the form when the request method is POST
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			return "", err
		}
		values := r.PostForm[paramsName]
		if len(values) > 0 {
			value = values[0]
		}
	// the parameter value need to be parsed from the url when the request method is GET
	case http.MethodGet, http.MethodDelete, http.MethodPut:
		values := r.URL.Query()[paramsName]
		if len(values) > 0 {
			value = values[0]
		}
	// default return error
	default:
		return "", fmt.Errorf("only GET/POST/DELETE/PUT methods are supported")
	}
	if len(value) > 0 {
		return value, nil
	}
	if !required {
		return defaultValue, nil
	}
	return "", fmt.Errorf("please input %s", paramsName)
}
