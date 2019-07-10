package middleware

import (
	/* #nosec */
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/eleme/lindb/models"

	"github.com/dgrijalva/jwt-go"
)

// UserAuthentication represents authentication param
type UserAuthentication struct {
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

// NewUserAuthentication creates authentication api instance
func NewUserAuthentication(user models.User) *UserAuthentication {
	return &UserAuthentication{
		user: user,
	}
}

// ValidateTokenMiddleware creates middleware for user permissions validation by request header Authorization
// if not authorization throw error
// else perform the next action
func (u *UserAuthentication) ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			err := errors.New("header cannot have authorization")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			b, _ := json.Marshal(err.Error())
			_, _ = w.Write(b)
			return
		}
		claims, _ := ParseToken(token, u.user)
		if claims.UserName == u.user.UserName && claims.Password == u.user.Password {
			next.ServeHTTP(w, r)
		}
	})
}

// ParseToken returns jwt claims by token
// get secret key use Md5Encrypt method with username and password
// then jwt parse token by secret key
func ParseToken(tokenString string, user models.User) (*CustomClaims, error) {
	claims := CustomClaims{}
	cid := Md5Encrypt(user)
	_, _ = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return cid, nil
	})
	return &claims, nil
}

// CreateLToken returns token use jwt with custom claims
func CreateToken(user models.User) (string, error) {
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
