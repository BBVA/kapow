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
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

type resourceHandler func(http.ResponseWriter, *http.Request, *model.Handler)

func lockResponseWriter(fn resourceHandler) resourceHandler {
	return func(w http.ResponseWriter, r *http.Request, h *model.Handler) {
		h.Writing.Lock()
		defer h.Writing.Unlock()
		fn(w, r, h)
	}
}

func checkHandler(fn resourceHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerID := mux.Vars(r)["handlerID"]
		if h, ok := Handlers.Get(handlerID); ok {
			fn(w, r, h)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
