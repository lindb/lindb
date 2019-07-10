package signature

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/eleme/lindb/config"

	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/eleme/lindb/pkg/signature/apierrors"

	"net/http"
)

var (
	userPath = ""
)

var WhiteRoute = map[string]bool{"/login": true, "/health": true}

type UserInfo struct {
	Cid      string `json:"cid"`
	Token    string `json:"token"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

const (
	linDBUser          = "user.toml"
	defaultUserCfgFile = "/etc/lindb/" + linDBUser
)

type TokenStatus int

var userConfig *config.LinDbUser

func init() {
	if userPath == "" {
		userPath = defaultUserCfgFile
	}
	userConfig = new(config.LinDbUser)
	if _, err := toml.DecodeFile(userPath, userConfig); err != nil {
		return
	}
}

// BuildNewUserInfo create cid and token by username and password
func BuildNewUserInfo(userName string, password string) (string, string) {
	key := userName + "/" + password
	token := jwt.New(jwt.SigningMethodHS256)
	cid, _ := AesEncrypt(key)
	enc, _ := token.SignedString([]byte(cid))
	return cid, enc
}

// LoginLinDB login lindb
// user info query from config file
func LoginLinDB(w http.ResponseWriter, r *http.Request) {
	var user UserInfo
	er := json.NewDecoder(r.Body).Decode(&user)
	if er != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, "Error in request")
		return
	}
	if len(user.UserName) == 0 {
		writeErrorResponse(w, apierrors.APIErrUserNameIsEmpty, r.URL)
		return
	}

	if len(user.Password) == 0 {
		writeErrorResponse(w, apierrors.APIErrPasswordIsEmpty, r.URL)
		return
	}

	if userConfig == nil {
		writeErrorResponse(w, apierrors.APIErrAccessDenied, r.URL)
		return
	}

	if userConfig.USER.UserName != user.UserName {
		writeErrorResponse(w, apierrors.APIErrUserName, r.URL)
		return
	}

	if userConfig.USER.Password != user.Password {
		writeErrorResponse(w, apierrors.APIErrPassword, r.URL)
		return
	}

	userInfo := UserInfo{}
	cid, token := BuildNewUserInfo(userConfig.USER.UserName, userConfig.USER.Password)
	userInfo.Cid = cid
	userInfo.Token = token

	response, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, response)
}

// ValidateTokenMiddleware lindb user info validate token middleware
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	path := r.URL.Path
	if WhiteRoute[path] {
		next(w, r)
	}
	cid := r.Header.Get("Cid")
	if len(cid) == 0 {
		writeErrorResponse(w, apierrors.APIErrUnsignedHeaders, r.URL)
	} else {
		tokenString, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(cid), nil
			})
		if err == nil {
			if tokenString.Valid {
				next(w, r)
			} else {
				writeErrorResponse(w, apierrors.APIErrInvalidToken, r.URL)
			}
		} else {
			writeErrorResponse(w, apierrors.APIErrInvalidUnauthorized, r.URL)
		}
	}
}

// HealthCheck lindb health check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(`{"alive": true}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeSuccessResponse(w, response)
}
