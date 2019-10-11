package user

import (
	"errors"
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user/mux"
)

type safeRouteList struct {
	rs []model.Route
	m  *sync.RWMutex
}

var Routes safeRouteList = New()

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

	Server.Handler.(*mux.SwappableMux).Update(srl.Snapshot())

	return model.Route{Index: l - 1}
}

func (srl *safeRouteList) Snapshot() []model.Route {
	srl.m.RLock()
	defer srl.m.RUnlock()

	rs := make([]model.Route, len(srl.rs))
	copy(rs, srl.rs)
	return rs
}

func (srl *safeRouteList) List() []model.Route {
	rs := srl.Snapshot()
	for i := 0; i < len(rs); i++ {
		rs[i].Index = i
	}
	return rs
}

func (srl *safeRouteList) Delete(ID string) error {
	srl.m.Lock()
	for i := 0; i < len(srl.rs); i++ {
		if srl.rs[i].ID == ID {
			srl.rs = append(srl.rs[:i], srl.rs[i+1:]...)
			srl.m.Unlock()
			Server.Handler.(*mux.SwappableMux).Update(srl.Snapshot())
			return nil

		}
	}
	srl.m.Unlock()
	return errors.New("Route not found")
}
