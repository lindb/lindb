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

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
)

// for testing
var (
	pathUnescapeFunc = url.PathUnescape
)

var (
	httpHandleTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_handle_duration",
			Help:    "HTTP handle duration(ms).",
			Buckets: monitoring.DefaultHistogramBuckets,
		},
		[]string{"path", "status"},
	)
)

func init() {
	monitoring.BrokerRegistry.MustRegister(httpHandleTimer)
}

// AccessLogMiddleware returns access log middleware
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := timeutil.Now()
		r := c.Request
		defer func() {
			// add access log
			path := r.RequestURI
			unescapedPath, err := pathUnescapeFunc(path)
			if err != nil {
				unescapedPath = path
			}
			// http://httpd.apache.org/docs/1.3/logs.html?PHPSESSID=026558d61a93eafd6da3438bb9605d4d#common
			requestInfo := realIP(r) + " " + strconv.Itoa(int(timeutil.Now()-start)) + "ms" +
				" \"" + r.Method + " " + unescapedPath + " " + r.Proto + "\" " +
				strconv.Itoa(c.Writer.Status()) + " " + strconv.Itoa(c.Writer.Size())
			if len(c.Errors) > 0 {
				logger.AccessLog.Error(requestInfo, logger.Error(c.Errors[0].Err))
			} else {
				logger.AccessLog.Info(requestInfo)
			}
			paths := strings.Split(unescapedPath, "?")
			if len(paths) > 0 {
				path = paths[0]
			}
			httpHandleTimer.WithLabelValues(path, strconv.Itoa(c.Writer.Status())).Observe(float64(timeutil.Now() - start))
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
