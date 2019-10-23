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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func handlerStatusOK(route model.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func handleRouteIDToBody(route model.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, route.ID)
	})
}

func TestGorillizeReturnsAnEmptyMuxWhenAnEmptyRouteList(t *testing.T) {
	m := gorillize([]model.Route{}, handlerStatusOK)

	if !reflect.DeepEqual(*m, *mux.NewRouter()) {
		t.Error("Returned mux not empty")
	}
}

func TestGorillizeReturnsAMuxThat404sWhenEmptyRouteList(t *testing.T) {
	m := *gorillize([]model.Route{}, handlerStatusOK)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status mismatch, got %d, want 404", res.StatusCode)
	}

}

func TestGorillizeReturnsAMuxThatMatchesByRoute(t *testing.T) {
	var rs []model.Route
	rs = append(rs, model.Route{
		Pattern: "/foo",
		Method:  "GET",
	})

	m := *gorillize(rs, handlerStatusOK)

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status mismatch, got %d, want 200", res.StatusCode)
	}
}

func TestGorillizeReturnsAMuxThat405sWhenMethodMismatch(t *testing.T) {
	var rs []model.Route
	rs = append(rs, model.Route{
		Pattern: "/foo",
		Method:  "GET",
	})

	m := *gorillize(rs, handlerStatusOK)

	req := httptest.NewRequest("POST", "/foo", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status mismatch, got %d, want 405", res.StatusCode)
	}
}

func TestGorillizeReturnsAMuxThatMatchesByMethod(t *testing.T) {
	var rs []model.Route
	rs = append(rs, model.Route{
		Pattern: "/foo",
		Method:  "UNORTHODOX",
	})

	m := *gorillize(rs, handlerStatusOK)

	req := httptest.NewRequest("UNORTHODOX", "/foo", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status mismatch, got %d, want 200", res.StatusCode)
	}
}

func TestGorillizeReturnsAMuxThatRespectsRouteOrderAB(t *testing.T) {
	var rs []model.Route
	rs = append(rs,
		model.Route{
			ID:      "routeA",
			Pattern: "/foo",
			Method:  "GET",
		},
		model.Route{
			ID:      "routeB",
			Pattern: "/foo",
			Method:  "GET",
		},
	)
	m := *gorillize(rs, handleRouteIDToBody)

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if body, _ := ioutil.ReadAll(res.Body); string(body) != "routeA" {
		t.Errorf("Mux did not respect route order %q", body)
	}
}

func TestGorillizeReturnsAMuxThatRespectsRouteOrderBA(t *testing.T) {
	var rs []model.Route
	rs = append(rs,
		model.Route{
			ID:      "routeB",
			Pattern: "/foo",
			Method:  "GET",
		},
		model.Route{
			ID:      "routeA",
			Pattern: "/foo",
			Method:  "GET",
		},
	)
	m := *gorillize(rs, handleRouteIDToBody)

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	res := w.Result()

	if body, _ := ioutil.ReadAll(res.Body); string(body) != "routeB" {
		t.Errorf("Mux did not respect route order %q", body)
	}
}
