package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func Test_ParseToken(t *testing.T) {
	user := models.User{UserName: "admin", Password: "admin123"}
	claim, _ := ParseToken(
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOi"+
			"JhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk", user)
	assert.Equal(t, user.UserName, claim.UserName)
	assert.Equal(t, user.Password, claim.Password)
}

func Test_CreateToken(t *testing.T) {
	user := models.User{UserName: "admin", Password: "admin123"}
	u := NewUserAuthentication(user)
	token, err := CreateToken(u.user)
	assert.Equal(t, true, err == nil)
	assert.Equal(t,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGF"+
			"zc3dvcmQiOiJhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk", token)
}
