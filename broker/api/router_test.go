package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"

	"github.com/stretchr/testify/assert"
)

func TestAddMiddleware(t *testing.T) {
	reg, _ := regexp.Compile("/check/*")
	AddMiddleware(func(next http.Handler) http.Handler {
		return nil
	}, reg)

	assert.Equal(t, 1, len(middlewareHandlers))

	middleware := getMiddleware("/check/test")
	assert.Equal(t, 1, len(middleware))
}

func TestAddRoute(t *testing.T) {
	AddRoute("test", http.MethodGet, "/test", func(writer http.ResponseWriter, request *http.Request) {})
	assert.Equal(t, 1, len(routes))
}

func TestNewRouter(t *testing.T) {
	reg, _ := regexp.Compile("/test/*")
	AddMiddleware(func(next http.Handler) http.Handler {
		return nil
	}, reg)
	AddRoute("test", http.MethodGet, "/test/login", func(writer http.ResponseWriter, request *http.Request) {})
	AddRoute("GetUser", http.MethodGet, "/get/user", func(writer http.ResponseWriter, request *http.Request) {})
	router := NewRouter()
	assert.NotNil(t, router)
}

func TestNewRouter_Static(t *testing.T) {
	old := staticPath
	staticPath = "/test/static/path"
	defer func() {
		_ = fileutil.RemoveDir(staticPath)
		staticPath = old
	}()
	_ = NewRouter()
	_ = fileutil.MkDir(staticPath)
	_ = NewRouter()
}

func TestToConsole(t *testing.T) {
	req, _ := http.NewRequest("GET", "/database", nil)
	resp := httptest.NewRecorder()
	redirectToConsole(resp, req)
}

func TestRouteHandle(t *testing.T) {
	// add panic handle
	AddRoute("test panic", http.MethodGet, "/panic", func(writer http.ResponseWriter, request *http.Request) {
		param, _ := GetParamsFromRequest("type", request, "test", true)
		switch param {
		case "err":
			panic(fmt.Errorf("fff"))
		case "str":
			panic("fff")
		case "other":
			panic([]byte{1, 2, 3})
		}
	})
	AddRoute("test", http.MethodGet, "/test", func(writer http.ResponseWriter, request *http.Request) {
		OK(writer, request)
	})
	r := NewRouter()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	req, _ = http.NewRequest(http.MethodOptions, "/test", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	req, _ = http.NewRequest(http.MethodGet, "/panic?type=err", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	req, _ = http.NewRequest(http.MethodGet, "/panic?type=str", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	req, _ = http.NewRequest(http.MethodGet, "/panic?type=other", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
}
