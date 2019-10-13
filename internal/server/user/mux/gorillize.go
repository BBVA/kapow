package mux

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/model"
)

func gorillize(rs []model.Route, buildHandler func(model.Route) http.Handler) *mux.Router {
	m := mux.NewRouter()

	for _, r := range rs {
		m.Handle(r.Pattern, buildHandler(r)).Methods(r.Method)
	}

	return m
}
