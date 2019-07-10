package signature

import (
	"bytes"
	"encoding/json"

	"github.com/eleme/lindb/pkg/signature/apierrors"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"

	"testing"
)

const UserName = "etrace"
const Password = "etrace"

func Test_LoginLinDB(t *testing.T) {
	userInfo := UserInfo{}
	userInfo.UserName = UserName
	userInfo.Password = Password
	user, _ := json.Marshal(userInfo)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(user))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	LoginLinDB(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resultUser UserInfo
	er := json.Unmarshal(rr.Body.Bytes(), &resultUser)
	assert.Equal(t, true, er == nil)
	assert.True(t, true, len(resultUser.Token) > 0)
}

func TestLoginLinDBBad(t *testing.T) {
	userInfo := UserInfo{}
	userInfo.UserName = UserName
	user, _ := json.Marshal(userInfo)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(user))
	rr := httptest.NewRecorder()

	LoginLinDB(rr, req)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginLinDBBad2(t *testing.T) {
	userInfo := UserInfo{}
	userInfo.UserName = UserName
	userInfo.Password = "etrace2"
	user, _ := json.Marshal(userInfo)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(user))

	rr := httptest.NewRecorder()

	LoginLinDB(rr, req)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginLinDBBad3(t *testing.T) {
	userInfo := UserInfo{}
	userInfo.Password = Password
	user, _ := json.Marshal(userInfo)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(user))

	rr := httptest.NewRecorder()

	LoginLinDB(rr, req)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
}

func TestLoginLinDBBad4(t *testing.T) {
	userInfo := UserInfo{}
	userInfo.UserName = "etrace3"
	userInfo.Password = "etrace4"
	user, _ := json.Marshal(userInfo)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(user))
	rr := httptest.NewRecorder()

	LoginLinDB(rr, req)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestValidateTokenMiddleware(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test", nil)
	rr := httptest.NewRecorder()
	ValidateTokenMiddleware(rr, req, Validate)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
}

func TestValidateTokenMiddleware2(t *testing.T) {
	req, err := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Cid", "06e1b91ea72901f6b6f1e4b91cf7fa936cde548e420b727482108efa402be6")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.cM79qrxJwMCoObxWjL00OLOtGxDsFhnIkKkhB1y_3Ms")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ValidateTokenMiddleware(rr, req, Validate)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test me", rr.Body.String())
}

func Validate(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("test me"))
}

func TestValidateTokenMiddleware3(t *testing.T) {
	req, err := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Cid", "06e1b91ea72901f6b6f1e4b91cf7fa936cde548e420b727482108efa402be6")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.cM79qrxwMCoObxWjL00OLOtGxDsFhnIkKkhB1y_3Ms")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ValidateTokenMiddleware(rr, req, Validate)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var apiError apierrors.APIError
	er := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.Equal(t, true, er == nil)
	assert.Equal(t, apiError.HTTPStatusCode, rr.Code)
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `"{\"alive\": true}"`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
