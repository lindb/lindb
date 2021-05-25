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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
)

var tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvc" +
	"mQiOiJhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk"

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := middleware.NewMockAuthentication(ctrl)

	user := config.User{UserName: "admin", Password: "admin123"}
	api := NewLoginAPI(user, auth)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    []byte{1, 2, 3},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
	})

	//create success
	auth.EXPECT().CreateToken(gomock.Any()).Return(tokenStr, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    user,
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: tokenStr,
	})

	// token create fail
	auth.EXPECT().CreateToken(gomock.Any()).Return("", fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{UserName: "admin", Password: "admin123"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	//user failure error password
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{UserName: "admin", Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	//user failure error user name
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{UserName: "123", Password: "admin123"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	//user failure error password
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{UserName: "admin", Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	//user login failure  password is empty
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{UserName: "admin"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	//user login failure  user name is empty
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    config.User{Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
		ExpectResponse: "",
	})

	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/check",
		HandlerFunc:    api.Check,
		ExpectHTTPCode: 200,
		ExpectResponse: user,
	})
}

func Test_JWT(t *testing.T) {
	user := config.User{UserName: "admin", Password: "admin123"}
	claims := middleware.CustomClaims{
		UserName: user.UserName,
		Password: user.Password,
	}
	cid := middleware.Md5Encrypt(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte(cid))

	assert.Equal(t, tokenStr, tokenString)

	mapClaims := middleware.CustomClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, &mapClaims, func(token *jwt.Token) (i interface{}, e error) {
		return cid, nil
	})
	assert.Equal(t, user.Password, mapClaims.Password)
	assert.Equal(t, user.UserName, mapClaims.UserName)
}
