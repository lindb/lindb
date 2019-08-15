package api

import (
	"errors"
	"net/http"
	"regexp"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"

	"github.com/lindb/lindb/pkg/fileutil"
)

type route struct {
	name    string
	method  string
	pattern string
	handler http.HandlerFunc
}

type middlewareHandler struct {
	regexp     *regexp.Regexp
	middleware mux.MiddlewareFunc
}

var routes []route

var middlewareHandlers []middlewareHandler

func AddMiddleware(middleware mux.MiddlewareFunc, regexp *regexp.Regexp) {
	middlewareHandlers = append(middlewareHandlers, middlewareHandler{middleware: middleware, regexp: regexp})
}

func AddRoutes(name, method, pattern string, handler http.HandlerFunc) {
	routes = append(routes, route{name: name, method: method, pattern: pattern, handler: handler})
}

// NewRouter returns a new router with a panic handler and a static server handler.
// middleware Method by method
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		mds := getMiddleware(route.pattern)
		var handler http.Handler
		// this route.pattern set middleware
		if len(mds) > 0 {
			for _, md := range mds {
				if handler != nil {
					handler = md.Middleware(handler)
				} else {
					handler = md.Middleware(route.handler)
				}
			}
		} else {
			handler = route.handler
		}
		router.
			Methods([]string{route.method, http.MethodOptions}...).
			Name(route.name).
			Handler(panicHandler(handler)).
			Path(route.pattern)
	}
	// static server path exist, serve web console
	webPath := "./web/build"
	if fileutil.Exist(webPath) {
		router.PathPrefix("/static/").
			Handler(http.StripPrefix("/static/",
				http.FileServer(rice.MustFindBox("./../../web/build").HTTPBox())))
	}
	// add cors support
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Method", "POST, OPTIONS, GET, HEAD, PUT, PATCH, DELETE")

				w.Header().Set("Access-Control-Allow-Headers",
					"Origin, X-Requested-With, X-HTTP-Method-Override,accept-charset,accept-encoding "+
						", Content-Type, Accept, Authorization")

				if r.Method == http.MethodOptions {
					return
				}
				next.ServeHTTP(w, r)
			})
		})
	return router
}

// getMiddleware returns suited middleware by pattern
func getMiddleware(pattern string) []mux.MiddlewareFunc {
	var ms []mux.MiddlewareFunc
	for _, middlewareHandler := range middlewareHandlers {
		if middlewareHandler.regexp.MatchString(pattern) {
			ms = append(ms, middlewareHandler.middleware)
		}
	}
	return ms
}

// panicHandler handles panics and returns a json response with error message
// and http code 500
func panicHandler(h http.Handler) http.Handler {
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
