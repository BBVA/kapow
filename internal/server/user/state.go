/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
	r.Index = len(srl.rs)
	srl.rs = append(srl.rs, r)
	srl.m.Unlock()

	Server.Handler.(*mux.SwappableMux).Update(srl.Snapshot())

	return r
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
	// TODO: Refactor with `defer` if applicable
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

func (srl *safeRouteList) Get(ID string) (r model.Route, err error) {
	srl.m.RLock()
	defer srl.m.RUnlock()
	for _, r = range srl.rs {
		if r.ID == ID {
			return
		}
	}

	err = errors.New("Route not found")
	return
}
