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
	"crypto/md5"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/lindb/lindb/config"
)

// CustomClaims represents jwt custom claims param
// need username and password and some standard claims
type CustomClaims struct {
	jwt.StandardClaims
	UserName string `json:"username"`
	Password string `json:"password"`
}

// CreateToken returns token use jwt with custom claims
func CreateToken(user config.User) (string, error) {
	claims := CustomClaims{
		UserName: user.UserName,
		Password: user.Password,
	}
	cid := Md5Encrypt(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(cid))
}

// Md5Encrypt returns secret key use Mk5 encryption with username and password
func Md5Encrypt(user config.User) string {
	/* #nosec */
	md5Encrypt := md5.New()
	key := fmt.Sprintf("%s/%s", user.UserName, user.Password)
	_, _ = md5Encrypt.Write([]byte(key))
	cipher := md5Encrypt.Sum(nil)
	return string(cipher)
}
