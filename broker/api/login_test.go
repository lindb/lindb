package api

import (
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"

	"github.com/eleme/lindb/broker/middleware"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
)

type testLoginAPISuite struct {
	mock.RepoTestSuite
}

var test *testing.T

func TestLoginApi(t *testing.T) {
	check.Suite(&testLoginAPISuite{})
	test = t
	check.TestingT(t)
}

func (tl *testLoginAPISuite) TestLogin(c *check.C) {
	user := models.User{UserName: "admin", Password: "admin123"}
	api := NewLoginAPI(user)

	//create success
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    user,
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 200,
	})

	//user failure error password
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    models.User{UserName: "admin", Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 500,
	})

	//user failure error user name
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    models.User{UserName: "123", Password: "admin123"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 500,
	})

	//user failure error password
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    models.User{UserName: "admin", Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 500,
	})

	//user login failure  password is empty
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    models.User{UserName: "admin"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 500,
	})

	//user login failure  user name is empty
	mock.DoRequest(test, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/user",
		RequestBody:    models.User{Password: "admin1234"},
		HandlerFunc:    api.Login,
		ExpectHTTPCode: 500,
	})
}

func Test_JWT(t *testing.T) {
	user := models.User{UserName: "admin", Password: "admin123"}
	claims := middleware.CustomClaims{
		UserName: user.UserName,
		Password: user.Password,
	}
	cid := middleware.Md5Encrypt(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte(cid))

	assert.Equal(t,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvc"+
			"mQiOiJhZG1pbjEyMyJ9.YbNGN0V-U5Y3xOIGNXcgbQkK2VV30UDDEZV19FN62hk", tokenString)

	mapClaims := middleware.CustomClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, &mapClaims, func(token *jwt.Token) (i interface{}, e error) {
		return cid, nil
	})
	assert.Equal(t, user.Password, mapClaims.Password)
	assert.Equal(t, user.UserName, mapClaims.UserName)
}
