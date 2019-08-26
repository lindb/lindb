package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lindb/lindb/config"

	"github.com/stretchr/testify/assert"
)

var tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGF" +
	"zc3dvcmQiOiJhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk"

func Test_ParseToken(t *testing.T) {
	user := config.User{UserName: "admin", Password: "admin123"}
	claim := parseToken(tokenStr, user)
	assert.Equal(t, user.UserName, claim.UserName)
	assert.Equal(t, user.Password, claim.Password)
}

func Test_CreateToken(t *testing.T) {
	user := config.User{UserName: "admin", Password: "admin123"}
	u := NewAuthentication(user)
	token, err := u.CreateToken(user)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, tokenStr, token)
}

func TestUserAuthentication_Validate(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer abc123")

	rr := httptest.NewRecorder()
	user := config.User{UserName: "admin", Password: "admin123"}
	auth := NewAuthentication(user)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "ok")
	})
	authHandler := auth.Validate(handler)

	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	req, err = http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", tokenStr)
	rr = httptest.NewRecorder()

	authHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}
