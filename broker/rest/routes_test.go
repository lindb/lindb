package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eleme/lindb/pkg/option"

	"github.com/stretchr/testify/assert"
)

func TestGetJSONBodyFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/database", nil)
	_, err := GetParamsFromRequest("test", req, "defaultValue", true)
	assert.True(t, err != nil)

	value, _ := GetParamsFromRequest("test", req, "defaultValue", false)
	assert.Equal(t, "defaultValue", value)

	req2, _ := http.NewRequest("POST", "/database", nil)
	_, err2 := GetParamsFromRequest("test", req2, "defaultValue", true)
	assert.True(t, err2 != nil)

	value2, _ := GetParamsFromRequest("test", req2, "defaultValue", false)
	assert.Equal(t, "defaultValue", value2)

}

func TestGetParamsFromRequest(t *testing.T) {
	database := option.Database{Name: "test"}
	databaseByte, _ := json.Marshal(database)
	req, _ := http.NewRequest("POST", "/database", bytes.NewReader(databaseByte))
	newDatabase := &option.Database{}
	_ = GetJSONBodyFromRequest(req, newDatabase)
	assert.Equal(t, newDatabase.Name, database.Name)
}

//httpHandlerTest defines the http handler test struct
type httpHandlerTest struct {
	method         string
	url            string
	requestJSON    interface{}
	handlerFunc    http.HandlerFunc
	expectHTTPCode int
	expectResponse interface{}
}

//testHTTPHandler helps test that whether the http response with the given
// method,url,requestJSON,handlerFunc matches with the expectHTTPCode and
// expectResponse. if false the test will be not failed
func testHTTPHandler(t *testing.T, httpHandlerTest *httpHandlerTest) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	requestBodyBytes, err := json.Marshal(httpHandlerTest.requestJSON)
	if err != nil {
		t.Errorf("parse requestJSON  error :%s", err.Error())
	}
	reader := bytes.NewReader(requestBodyBytes)
	req, err := http.NewRequest(httpHandlerTest.method, httpHandlerTest.url, reader)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to
	// record the response.
	rr := httptest.NewRecorder()
	handler := httpHandlerTest.handlerFunc

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != httpHandlerTest.expectHTTPCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, httpHandlerTest.expectHTTPCode)
	}
	// verify that the response matches the  response
	if httpHandlerTest.expectResponse != nil {
		readCloser := rr.Result().Body
		defer readCloser.Close()
		responseBytes, err := ioutil.ReadAll(readCloser)
		if err != nil {
			t.Errorf("read response body error :%s", err.Error())
		}
		expectResponseBytes, err := json.Marshal(httpHandlerTest.expectResponse)
		if err != nil {
			t.Errorf("parse expectResponse error :%s", err.Error())
		}
		if !bytes.Equal(responseBytes, expectResponseBytes) {
			var respInterface interface{}
			err := json.Unmarshal(requestBodyBytes, respInterface)
			if err != nil {
				t.Error("handler returned wrong response:,parse error")
			}
			t.Errorf("handler returned wrong response: got %v want %v",
				respInterface, httpHandlerTest.expectResponse)
		}
	}
}
