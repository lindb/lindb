package api

import (
	"errors"
	"net/http"

	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/models"
)

// LoginAPI represents login param
type LoginAPI struct {
	user models.User
}

// NewLoginAPI creates login api instance
func NewLoginAPI(user models.User) *LoginAPI {
	return &LoginAPI{
		user: user,
	}
}

// Login responses unique token
// if use name or password is empty will responses error msg
// if use name or password is error also will responses error msg
func (l *LoginAPI) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := GetJSONBodyFromRequest(r, &user)
	// login request is error
	if err != nil {
		OK(w, err)
		return
	}
	// user name is empty
	if len(user.UserName) == 0 {
		err = errors.New("user name is empty")
		Error(w, err)
		return
	}
	// password is empty
	if len(user.Password) == 0 {
		err = errors.New("password is empty")
		Error(w, err)
		return
	}
	// user name is error
	if l.user.UserName != user.UserName {
		err = errors.New("user name is error")
		Error(w, err)
		return
	}
	// password is error
	if l.user.Password != user.Password {
		err = errors.New("password is error")
		Error(w, err)
		return
	}
	token, err := middleware.CreateToken(user)
	if err != nil {
		Error(w, err)
	}
	OK(w, token)
}

// Check responses use msg
// this method use for test
func (l *LoginAPI) Check(w http.ResponseWriter, r *http.Request) {
	OK(w, l.user)
}
