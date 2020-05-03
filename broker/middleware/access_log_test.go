package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessLogMiddleware(t *testing.T) {
	defer func() {
		pathUnescapeFunc = url.PathUnescape
	}()
	pathUnescapeFunc = func(s string) (string, error) {
		return "err-path", fmt.Errorf("err")
	}
	req, err := http.NewRequest("GET", "/health-check", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	accessLogHandler := AccessLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "ok")
	}))

	accessLogHandler.ServeHTTP(rr, req)

	accessLogHandler = AccessLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "ok")
	}))
	accessLogHandler.ServeHTTP(rr, req)
}

func Test_real_ip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health-check", nil)
	req.Header.Add("X-Real-Ip", "real-ip")
	assert.Equal(t, "real-ip", realIP(req))

	req, _ = http.NewRequest("GET", "/health-check", nil)
	req.Header.Add("X-Forwarded-For", "forward-ip")
	assert.Equal(t, "forward-ip", realIP(req))
	req, _ = http.NewRequest("GET", "/health-check", nil)
	req.RemoteAddr = "1.1.1.1:1023"
	assert.Equal(t, "1.1.1.1", realIP(req))
}
