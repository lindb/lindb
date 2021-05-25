// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package api

import (
	"errors"
	"net/http"
	"regexp"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"

	"github.com/lindb/lindb/pkg/logger"
)

var staticPath = "./../../web/build"

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

// AddMiddleware adds middleware func base on url path pattern
func AddMiddleware(middleware mux.MiddlewareFunc, regexp *regexp.Regexp) {
	middlewareHandlers = append(middlewareHandlers, middlewareHandler{middleware: middleware, regexp: regexp})
}

// AddRoute adds http route handle func for urp pattern
func AddRoute(name, method, pattern string, handler http.HandlerFunc) {
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
				handler = md.Middleware(route.handler)
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
	box, err := rice.FindBox(staticPath)
	if err != nil {
		log.Error("cannot find static resource", logger.Error(err))
	} else {
		router.Path("/").Handler(http.HandlerFunc(redirectToConsole))
		router.PathPrefix("/console/").
			Handler(http.StripPrefix("/console/",
				http.FileServer(box.HTTPBox())))
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

// redirectToConsole redirects to admin console
func redirectToConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/console/", http.StatusFound)
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
				log.Error("serve http func panic", logger.Error(err), logger.Stack())
			}
		}()
		h.ServeHTTP(w, r)
	})
}
