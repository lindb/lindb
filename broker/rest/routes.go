package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/eleme/lindb/broker"
)

type route struct {
	name    string
	method  string
	pattern string
	handler http.HandlerFunc
}

type routes []route

var rs = routes{
	route{
		"Index",
		"GET",
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to my website!")
		},
	},
}

func NewRouter(config *broker.Config) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(route.handler)
	}
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(config.Http.Static))))

	return router
}
