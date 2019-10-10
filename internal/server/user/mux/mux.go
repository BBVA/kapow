package mux

import (
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

type swappableMux struct {
	m    sync.RWMutex
	root *mux.Router
}

func New() *swappableMux {
	return &swappableMux{
		root: mux.NewRouter(),
	}
}

func (sm *swappableMux) get() *mux.Router {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.root
}

func (sm *swappableMux) set(mux *mux.Router) {
	sm.m.Lock()
	sm.root = mux
	sm.m.Unlock()
}

func (sm *swappableMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sm.get().ServeHTTP(w, r)
}

func (sm *swappableMux) Update(rs []model.Route) {
	sm.set(gorillize(rs, handlerBuilder))
}
