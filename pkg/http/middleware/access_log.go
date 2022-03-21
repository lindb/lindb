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
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/pkg/logger"
)

// for testing
var (
	pathUnescapeFunc = url.PathUnescape
)

// AccessLog returns access log middleware
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		r := c.Request
		defer func() {
			// add access log
			path := r.RequestURI
			unescapedPath, err := pathUnescapeFunc(path)
			if err != nil {
				unescapedPath = path
			}
			status := c.Writer.Status()
			// http://httpd.apache.org/docs/1.3/logs.html?PHPSESSID=026558d61a93eafd6da3438bb9605d4d#common
			requestInfo := realIP(r) + " " + strconv.Itoa(int(time.Since(start).Microseconds())) + "us" +
				" \"" + r.Method + " " + unescapedPath + " " + r.Proto + "\" " +
				strconv.Itoa(status) + " " + strconv.Itoa(c.Writer.Size())
			if status >= 400 {
				logger.AccessLog.Error(requestInfo)
			} else {
				logger.AccessLog.Debug(requestInfo)
			}
		}()
		c.Next()
	}
}

// realIP return the real ip
func realIP(r *http.Request) string {
	xRealIP := r.Header.Get("X-Real-Ip")
	if xRealIP != "" {
		return xRealIP
	}
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			return address
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
