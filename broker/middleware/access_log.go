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

// accessStats represents http access stats record
type accessStats struct {
	status int
	size   int
}

// loggingWriter represents logging stats for http response
type loggingWriter struct {
	http.ResponseWriter
	accessStats accessStats
}

// Write writes the data to the connection as part of an HTTP reply,
// records http status and response size for access log.
func (r *loggingWriter) Write(p []byte) (int, error) {
	if r.accessStats.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		r.accessStats.status = http.StatusOK
	}
	written, err := r.ResponseWriter.Write(p)
	r.accessStats.size += written
	return written, err
}

// WriteHeader sends an HTTP response header with the provided status code,
// records http status for access log.
func (r *loggingWriter) WriteHeader(status int) {
	r.accessStats.status = status
	r.ResponseWriter.WriteHeader(status)
}

// AccessLogMiddleware returns access log middleware
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := timeutil.Now()
		writer := &loggingWriter{
			ResponseWriter: w,
			accessStats:    accessStats{},
		}
		defer func() {
			// add access log
			log := writer.accessStats
			path := r.RequestURI
			unescapedPath, err := pathUnescapeFunc(path)
			if err != nil {
				unescapedPath = path
			}
			// http://httpd.apache.org/docs/1.3/logs.html?PHPSESSID=026558d61a93eafd6da3438bb9605d4d#common
			logger.AccessLog.Info(realIP(r) + " " + strconv.Itoa(int(timeutil.Now()-start)) + "ms" +
				" \"" + r.Method + " " + unescapedPath + " " + r.Proto + "\" " +
				strconv.Itoa(log.status) + " " + strconv.Itoa(log.size))
			paths := strings.Split(unescapedPath, "?")
			if len(paths) > 0 {
				path = paths[0]
			}
			httpHandleTimer.WithLabelValues(path, strconv.Itoa(log.status)).Observe(float64(timeutil.Now() - start))
		}()
		next.ServeHTTP(writer, r)
	})
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
