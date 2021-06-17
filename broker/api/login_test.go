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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	httppkg "github.com/lindb/lindb/pkg/http"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createTokenFn = httppkg.CreateToken
		ctrl.Finish()
	}()

	user := config.User{UserName: "admin", Password: "admin123"}
	api := NewLoginAPI(user)
	r := gin.New()
	api.Register(r)

	resp := mock.DoRequest(t, r, http.MethodPut, LoginPath, "")
	assert.Equal(t, http.StatusOK, resp.Code)

	//create success
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, "")
	assert.Equal(t, http.StatusOK, resp.Code)

	//user failure error password
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "admin", "password": "admin1234"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	//user failure error user name
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "123", "password": "admin1234"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	//user failure error password
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "123", "password": "admin12dd34"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	//user login failure  password is empty
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "123"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	// token create fail
	createTokenFn = func(user config.User) (string, error) {
		return "", fmt.Errorf("err")
	}
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "admin", "password": "admin123"}`)
	assert.Equal(t, http.StatusOK, resp.Code)

	// token create ok
	createTokenFn = httppkg.CreateToken
	resp = mock.DoRequest(t, r, http.MethodPut, LoginPath, `{"username": "admin", "password": "admin123"}`)
	assert.Equal(t, http.StatusOK, resp.Code)
}
