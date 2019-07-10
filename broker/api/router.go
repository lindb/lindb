package api

import (
	"errors"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
)

type route struct {
	name    string
	method  string
	pattern string
	handler http.HandlerFunc
}

var routes []route

func AddRoute(name, method, pattern string, handler http.HandlerFunc) {
	routes = append(routes, route{name: name, method: method, pattern: pattern, handler: handler})
}

// NewRouter returns a new router with a panic handler and a static server handler.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(panicHandler(route.handler))
	}

	// static server path exist, serve web console
	webPath := "./web/build"
	if util.Exist(webPath) {
		router.PathPrefix("/static/").
			Handler(http.StripPrefix("/static/",
				http.FileServer(rice.MustFindBox("./../../web/build").HTTPBox())))
	}
	return router
}

// panicHandler handles panics and returns a json response with error message
// and http code 500
func panicHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.GetLogger().Info("TTTTT", zap.Any("fff", r))
		var err error
		defer func() {
			r := recover()
			logger.GetLogger().Info("errr", zap.Stack("dfsfds"), zap.Any("fff", r))
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
