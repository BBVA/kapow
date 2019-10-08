package mux

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func handlerStatusOK(h *model.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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

// TODO: TestGorillizeReturnsAMuxThatRespectsRouteOrder
