package mock

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HTTPHandler represents amock http handler
type HTTPHandler struct {
	Method      string
	URL         string
	RequestBody interface{}

	HandlerFunc http.HandlerFunc

	ExpectHTTPCode int
	ExpectResponse interface{}
}

// DoRequest helps test that whether the http response with the given
// method,url,requestJSON,handlerFunc matches with the expectHTTPCode and
// expectResponse. if false the test will be not failed
func DoRequest(t *testing.T, httpHandler *HTTPHandler) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	requestBodyBytes, err := json.Marshal(httpHandler.RequestBody)
	if err != nil {
		t.Fatal(err)
		return
	}
	reader := bytes.NewReader(requestBodyBytes)
	req, err := http.NewRequest(httpHandler.Method, httpHandler.URL, reader)
	if err != nil {
		t.Fatal(err)
		return
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to
	// record the response.
	rr := httptest.NewRecorder()

	handler := httpHandler.HandlerFunc

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	assert.Equal(t, httpHandler.ExpectHTTPCode, rr.Code)

	// verify that the response matches the response
	if httpHandler.ExpectResponse != nil {
		readCloser := rr.Result().Body
		defer readCloser.Close()
		responseBytes, err := ioutil.ReadAll(readCloser)
		if err != nil {
			t.Fatal(err)
			return
		}
		expect, _ := json.Marshal(httpHandler.ExpectResponse)
		assert.Equal(t, expect, responseBytes)
	}
}
