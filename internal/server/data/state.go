package data

import (
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
)

type safeHandlerMap struct {
	h map[string]*model.Handler
	m sync.RWMutex
}

var Handlers = New()

func New() safeHandlerMap {
	return safeHandlerMap{
		h: make(map[string]*model.Handler),
		m: sync.RWMutex{},
	}
}

func (hs *safeHandlerMap) Add(handler *model.Handler) {
	hs.m.Lock()
	hs.h[handler.Id] = handler
	hs.m.Unlock()
}

func (hs *safeHandlerMap) Remove(id string) {
	hs.m.Lock()
	delete(hs.h, id)
	hs.m.Unlock()
}

func (hs *safeHandlerMap) Get(id string) (*model.Handler, bool) {
	hs.m.RLock()
	hndl, ok := hs.h[id]
	hs.m.RUnlock()
	return hndl, ok
}
