package client

import (
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestRemoveRouteOKExistent(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Delete("/routes/ROUTE_FOO").
		Reply(http.StatusNoContent)

	err := RemoveRoute("http://localhost:8080", "ROUTE_FOO")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

func TestRemoveRouteErrorNonExistent(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Delete("/routes/ROUTE_BAD").
		Reply(http.StatusNotFound)

	err := RemoveRoute("http://localhost:8080", "ROUTE_BAD")
	if err == nil {
		t.Errorf("Error not reported for nonexistent route")
	} else if err.Error() != "Not Found" {
		t.Errorf(`Error mismatch: got %q, want "Not Found"`, err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
