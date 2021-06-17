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
	"encoding/json"
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/lindb/lindb/config"
	httppkg "github.com/lindb/lindb/pkg/http"
)

//go:generate mockgen -source=./authentication.go -destination=./authentication_mock.go -package=middleware

type Authentication interface {
	// Validate validates the token
	Validate(next http.Handler) http.Handler
}

// userAuthentication represents user authentication using jwt
type userAuthentication struct {
	user config.User
}

// NewAuthentication creates authentication api instance
func NewAuthentication(user config.User) Authentication {
	return &userAuthentication{
		user: user,
	}
}

// Validate creates middleware for user permissions validation by request header Authorization
// if not authorization throw error
// else perform the next action
func (u *userAuthentication) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 0 {
			claims := parseToken(token, u.user)
			if claims.UserName == u.user.UserName && claims.Password == u.user.Password {
				next.ServeHTTP(w, r)
				return
			}
		}
		err := errors.New("authorization token invalid")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(err.Error())
		_, _ = w.Write(b)
	})
}

// ParseToken returns jwt claims by token
// get secret key use Md5Encrypt method with username and password
// then jwt parse token by secret key
func parseToken(tokenString string, user config.User) *httppkg.CustomClaims {
	claims := httppkg.CustomClaims{}
	cid := httppkg.Md5Encrypt(user)
	_, _ = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return cid, nil
	})
	return &claims
}
