package user

import (
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
)

type safeRouteList struct {
	rs []model.Route
	m  sync.RWMutex
}

var Routes = New()

func New() safeRouteList {
	return safeRouteList{}
}

func (srl *safeRouteList) Append(r model.Route) {
	srl.m.Lock()
	srl.rs = append(srl.rs, r)
	srl.m.Unlock()
}

func (srl *safeRouteList) Snapshot() []model.Route {
	srl.m.RLock()
	defer srl.m.RUnlock()

	if srl.rs == nil {
		return nil
	} else {
		rs := make([]model.Route, len(srl.rs))
		copy(rs, srl.rs)
		return rs
	}
}
