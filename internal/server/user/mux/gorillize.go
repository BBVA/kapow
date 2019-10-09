package mux

import (
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func gorillize(rs []model.Route, buildHandler func(*model.Route) http.Handler) *mux.Router {
	m := mux.NewRouter()

	for _, r := range rs {
		m.Handle(r.Pattern, buildHandler(nil)).Methods(r.Method)
	}

	return m
}
