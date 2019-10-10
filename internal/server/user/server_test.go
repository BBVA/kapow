package user

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user/mux"
)

// TODO TestRunRunsAnHTTPServer(t *testing.T) {}

func TestAppendUpdatesMuxWithProvideRoute(t *testing.T) {
	Server = http.Server{
		Handler: mux.New(),
	}
	srl := New()
	route := model.Route{
		Method:     "GET",
		Pattern:    "/",
		Entrypoint: "/bin/sh -c",
		Command:    "jaillover > /tmp/kapow-test-append-updates-mux",
	}
	os.Remove("/tmp/kapow-test-append-updates-mux")
	defer os.Remove("/tmp/kapow-test-append-updates-mux")

	srl.Append(route)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Server.Handler.ServeHTTP(w, req)

	if _, err := os.Stat("/tmp/kapow-test-append-updates-mux"); os.IsNotExist(err) {
		t.Error("Routes not updated")
	} else if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
