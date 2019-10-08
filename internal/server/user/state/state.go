package state

import (
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
)

type safeRouteList struct {
	rs []model.Route
	m  *sync.RWMutex
}

func New() safeRouteList {
	return safeRouteList{
		rs: []model.Route{},
		m:  &sync.RWMutex{},
	}
}

func (srl *safeRouteList) Append(r model.Route) model.Route {
	srl.m.Lock()
	srl.rs = append(srl.rs, r)
	l := len(srl.rs)
	srl.m.Unlock()

	return model.Route{Index: l - 1}
}

func (srl *safeRouteList) Snapshot() []model.Route {
	srl.m.RLock()
	defer srl.m.RUnlock()

	rs := make([]model.Route, len(srl.rs))
	copy(rs, srl.rs)
	return rs
}
