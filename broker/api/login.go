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
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	createTokenFn = http.CreateToken
	LoginPath     = "/login"
)

// LoginAPI represents login param
type LoginAPI struct {
	user config.User

	logger *logger.Logger
}

// NewLoginAPI creates login api instance
func NewLoginAPI(user config.User) *LoginAPI {
	return &LoginAPI{
		user:   user,
		logger: logger.GetLogger("broker", "LoginAPI"),
	}
}

// Register adds login url route.
func (l *LoginAPI) Register(route gin.IRoutes) {
	route.PUT(LoginPath, l.Login)
}

// Login responses unique token
// if use name or password is empty will responses error msg
// if use name or password is error also will responses error msg
func (l *LoginAPI) Login(c *gin.Context) {
	user := config.User{}
	err := c.ShouldBind(&user)
	if err != nil {
		l.logger.Error("cannot get user info from request")
		http.OK(c, "")
		return
	}
	// user name is error
	if l.user.UserName != user.UserName {
		l.logger.Error("username is invalid")
		http.OK(c, "")
		return
	}
	// password is error
	if l.user.Password != user.Password {
		l.logger.Error("password is invalid")
		http.OK(c, "")
		return
	}
	token, err := createTokenFn(user)
	if err != nil {
		http.OK(c, "")
		return
	}
	http.OK(c, token)
}
