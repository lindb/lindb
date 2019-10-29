package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	resp := httptest.NewRecorder()
	OK(resp, "ok")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `"ok"`, resp.Body.String())
}

func TestNoContent(t *testing.T) {
	resp := httptest.NewRecorder()
	NoContent(resp)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	assert.Equal(t, 0, resp.Body.Len())
}

func TestNotFound(t *testing.T) {
	resp := httptest.NewRecorder()
	NotFound(resp)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, 0, resp.Body.Len())
}

func TestError(t *testing.T) {
	resp := httptest.NewRecorder()
	Error(resp, fmt.Errorf("err"))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, `"err"`, resp.Body.String())
}
