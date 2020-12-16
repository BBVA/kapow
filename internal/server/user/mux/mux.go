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
		root: gorillize([]model.Route{}, handlerBuilder),
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
