package mux

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/model"
)

type SwappableMux struct {
	m    sync.RWMutex
	root *mux.Router
}

func New() *SwappableMux {
	return &SwappableMux{
		root: mux.NewRouter(),
	}
}

func (sm *SwappableMux) get() *mux.Router {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.root
}

func (sm *SwappableMux) set(mux *mux.Router) {
	sm.m.Lock()
	sm.root = mux
	sm.m.Unlock()
}

func (sm *SwappableMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sm.get().ServeHTTP(w, r)
}

func (sm *SwappableMux) Update(rs []model.Route) {
	sm.set(gorillize(rs, handlerBuilder))
}
