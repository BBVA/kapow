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
package data

import (
	"sync"

	"github.com/BBVA/kapow/internal/server/model"
)

type safeHandlerMap struct {
	hs map[string]*model.Handler
	m  *sync.RWMutex
}

var Handlers = New()

func New() safeHandlerMap {
	return safeHandlerMap{
		hs: make(map[string]*model.Handler),
		m:  &sync.RWMutex{},
	}
}

func (shm *safeHandlerMap) Add(h *model.Handler) {
	shm.m.Lock()
	shm.hs[h.ID] = h
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

func (shm *safeHandlerMap) ListIDs() (ids []string) {
	shm.m.RLock()
	defer shm.m.RUnlock()
	for id := range shm.hs {
		ids = append(ids, id)
	}
	return
}
