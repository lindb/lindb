package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eleme/lindb/broker"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type route struct {
	name    string
	method  string
	pattern string
	handler http.HandlerFunc
}

type routes []route

var rs = routes{
	route{"CreateOrUpdateDatabase", "POST", "/database", CreateOrUpdateDatabase},
	route{"GetDatabase", "Get", "/database", GetDatabase},
}

func NewRouter(config *broker.Config) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(PanicHandler(route.handler))
	}
	// static server path
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir(config.HTTP.Static))))
	return router
}

// Broker server handler
type APIHandler struct {
	// rpc gateway proxy handler
	Mux *runtime.ServeMux
	// base http handler
	Route *mux.Router
}

// Check the current http request type.if rpc proxy request,handler it with the runtime.ServeMux
func (handler *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO distinct the request type gracefully
	if p := strings.TrimPrefix(r.URL.Path, "/rpc/"); len(p) < len(r.URL.Path) {
		if authSuccess := authRPC(r); authSuccess {
			handler.Mux.ServeHTTP(w, r)
		} else {
			http.Error(w, "403 rpc auth failed", http.StatusForbidden)
		}
	} else {
		handler.Route.ServeHTTP(w, r)
	}
}

// Rpc request auth function.
func authRPC(r *http.Request) bool {
	s := r.Header.Get("Authorization")
	//todo Implement authentication and authorization
	return s == "auth"
}

// Parse json from request body into specified struct
func GetJSONBodyFromRequest(r *http.Request, t interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if decoder != nil {
		err := decoder.Decode(&t)
		return err
	}
	return fmt.Errorf("could parse request body")
}

// Get parameter value from the request。
// If there are multiple parameters with the same name, only the first value is returned.
// If there does not have the value and parameter has the required attribute it will return an error,
// otherwise it will return the defaultValue
func GetParamsFromRequest(paramsName string, r *http.Request, defaultValue string, required bool) (string, error) {
	if len(paramsName) == 0 {
		return "", fmt.Errorf("the params name must not be null")
	}
	var value string
	method := r.Method
	//Get request parameters according to different request methods
	switch method {
	case http.MethodGet:
		values := r.URL.Query()[paramsName]
		if len(values) > 0 {
			value = values[0]
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			return "", err
		}
		values := r.PostForm[paramsName]
		if len(values) > 0 {
			value = values[0]
		}
	default:
		return "", fmt.Errorf("only GET and POST methods are supported")
	}

	if len(value) > 0 {
		return value, nil
	}
	if !required {
		return defaultValue, nil
	}
	return "", fmt.Errorf("could not find the param;[%s] values ", paramsName)
}

// Represents a 200 handler and response message
func OKResponse(w http.ResponseWriter, a interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(a)
	_, _ = w.Write(b)
}

// Represents a 204 handler and response nothing
func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

// Represents a 500 handler and response error message
func ErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(err.Error())
	_, _ = w.Write(b)
}

// Panic handler for http request
func PanicHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("UnKnow ERROR")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
