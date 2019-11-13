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
	"net/http/httptest"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestConfigRouterReturnsRouterWithDecoratedRoutes(t *testing.T) {
	var handlerID string
	rs := []routeSpec{
		{
			"/handlers/{handlerID}/dummy",
			"GET",
			func(w http.ResponseWriter, r *http.Request, h *model.Handler) { handlerID = h.ID },
		},
	}
	Handlers = New()
	Handlers.Add(&model.Handler{ID: "FOO"})
	m := configRouter(rs)

	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/handlers/FOO/dummy", nil))

	if handlerID != "FOO" {
		t.Errorf(`Handler ID mismatch. Expected "FOO". Got %q`, handlerID)
	}
}

func TestConfigRouterReturnsRouterThat400sOnUnconfiguredResources(t *testing.T) {
	m := configRouter([]routeSpec{})
	w := httptest.NewRecorder()

	m.ServeHTTP(w, httptest.NewRequest("GET", "/handlers/FOO/dummy", nil))

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code mismatch. Expected 400. Got %d", res.StatusCode)
	}
}
