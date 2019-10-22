package data

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/gorilla/mux"

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
