package api

import (
	"encoding/json"
	"net/http"
)

// OK responses with content and set the http status code 200
func OK(w http.ResponseWriter, a interface{}) {
	b, _ := json.Marshal(a)
	response(w, http.StatusOK, b)
}

// NoContent responses with empty content and set the http status code 204
func NoContent(w http.ResponseWriter) {
	response(w, http.StatusNoContent, nil)
}

// NotFound responses resource not found
func NotFound(w http.ResponseWriter) {
	response(w, http.StatusNotFound, nil)
}

// Error responses error message and set the http status code 500
func Error(w http.ResponseWriter, err error) {
	b, _ := json.Marshal(err.Error())
	response(w, http.StatusInternalServerError, b)
}

// response responses json body for http restful api
func response(w http.ResponseWriter, httpCode int, content []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	if len(content) > 0 {
		_, _ = w.Write(content)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
