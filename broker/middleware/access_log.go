package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
)

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
			unescapedPath, err := url.PathUnescape(path)
			if err != nil {
				unescapedPath = path
			}
			// http://httpd.apache.org/docs/1.3/logs.html?PHPSESSID=026558d61a93eafd6da3438bb9605d4d#common
			logger.AccessLog.Info(realIP(r) + " " + strconv.Itoa(int(timeutil.Now()-start)) + "ms" +
				" \"" + r.Method + " " + unescapedPath + " " + r.Proto + "\" " +
				strconv.Itoa(log.status) + " " + strconv.Itoa(log.size))
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
