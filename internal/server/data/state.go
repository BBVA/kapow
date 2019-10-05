package data

import (
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
)

type safeHandlerMap struct {
	hs map[string]*model.Handler
	m  sync.RWMutex
}

var Handlers = New()

func New() safeHandlerMap {
	return safeHandlerMap{
		hs: make(map[string]*model.Handler),
		m:  sync.RWMutex{},
	}
}

func (shm *safeHandlerMap) Add(h *model.Handler) {
	shm.m.Lock()
	shm.hs[h.Id] = h
	shm.m.Unlock()
}

func (shm *safeHandlerMap) Remove(id string) {
	shm.m.Lock()
	delete(shm.hs, id)
	shm.m.Unlock()
}

func (shm *safeHandlerMap) Get(id string) (*model.Handler, bool) {
	shm.m.RLock()
	h, ok := shm.hs[id]
	shm.m.RUnlock()
	return h, ok
}
