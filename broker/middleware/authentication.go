package middleware

import (
	/* #nosec */
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lindb/lindb/models"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=./authentication.go -destination=./authentication_mock.go -package=middleware

type Authentication interface {
	// CreateLToken returns the authentication token
	CreateToken(user models.User) (string, error)
	// Validate validates the token
	Validate(next http.Handler) http.Handler
}

// userAuthentication represents user authentication using jwt
type userAuthentication struct {
	user models.User
}

// CustomClaims represents jwt custom claims param
// need username and password and some standard claims
type CustomClaims struct {
	jwt.StandardClaims
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Valid rewrites jwt.Claims valid method return nil
func (*CustomClaims) Valid() error {
	return nil
}

// NewAuthentication creates authentication api instance
func NewAuthentication(user models.User) Authentication {
	return &userAuthentication{
		user: user,
	}
}

// Validate creates middleware for user permissions validation by request header Authorization
// if not authorization throw error
// else perform the next action
func (u *userAuthentication) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 0 {
			claims := parseToken(token, u.user)
			if claims.UserName == u.user.UserName && claims.Password == u.user.Password {
				next.ServeHTTP(w, r)
				return
			}
		}
		err := errors.New("authorization token invalid")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(err.Error())
		_, _ = w.Write(b)
	})
}

// ParseToken returns jwt claims by token
// get secret key use Md5Encrypt method with username and password
// then jwt parse token by secret key
func parseToken(tokenString string, user models.User) *CustomClaims {
	claims := CustomClaims{}
	cid := Md5Encrypt(user)
	_, _ = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return cid, nil
	})
	return &claims
}

// CreateLToken returns token use jwt with custom claims
func (u *userAuthentication) CreateToken(user models.User) (string, error) {
	claims := CustomClaims{
		UserName: user.UserName,
		Password: user.Password,
	}
	cid := Md5Encrypt(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(cid))
}

// Md5Encrypt returns secret key use Mk5 encryption with username and password
func Md5Encrypt(user models.User) string {
	/* #nosec */
	md5Encrypt := md5.New()
	key := fmt.Sprintf("%s/%s", user.UserName, user.Password)
	_, _ = md5Encrypt.Write([]byte(key))
	cipher := md5Encrypt.Sum(nil)
	return string(cipher)
}
