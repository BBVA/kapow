package data

import (
	"net/http"

	"github.com/gorilla/mux"
)

type routeSpec struct {
	route  string
	method string
	rh     resourceHandler
}

func configRouter(rs []routeSpec) (r *mux.Router) {
	r = mux.NewRouter()
	for _, s := range rs {
		r.HandleFunc(s.route, checkHandler(s.rh)).Methods(s.method)
	}
	r.HandleFunc(
		"/handlers/{handlerID}/{resource:.*}",
		func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) })
	return r
}
