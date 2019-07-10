package api

import (
	"encoding/json"
	"net/http"
)

// OK responses with content and set the http status code 200
func OK(w http.ResponseWriter, a interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(a)
	_, _ = w.Write(b)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// NoContent responses with empty content and set the http status code 204
func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

// NotFound responses resource not found
func NotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
}

// Error responses error message and set the http status code 500
func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(err.Error())
	_, _ = w.Write(b)
}
