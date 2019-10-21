package data

import (
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

type resourceHandler func(http.ResponseWriter, *http.Request, *model.Handler)

func lockResponseWriter(fn resourceHandler) resourceHandler {
	return func(w http.ResponseWriter, r *http.Request, h *model.Handler) {
		h.Writing.Lock()
		defer h.Writing.Unlock()
		fn(w, r, h)
	}
}

func checkHandler(fn resourceHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerID := mux.Vars(r)["handlerID"]
		if h, ok := Handlers.Get(handlerID); ok {
			fn(w, r, h)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
