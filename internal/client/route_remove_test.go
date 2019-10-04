package client

import (
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestRemoveRouteExistent(t *testing.T) {
	const (
		host    = "http://localhost:8080"
		routeID = "ROUTE_ID_OLDER_BUT_IT_CHECKS_OUT"
	)

	defer gock.Off()
	gock.New(host).Delete("/routes/" + routeID).Reply(http.StatusNoContent)

	err := RemoveRoute(host, routeID)
	if err != nil {
		t.Errorf("unexpected error: ‘%s’", err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

func TestRemoveRouteNonExistent(t *testing.T) {
	const (
		host    = "http://localhost:8080"
		routeID = "ROUTE_THIS_ONE_WONT_WORK_BUDDY"
	)
	expected := http.StatusText(http.StatusNotFound)

	defer gock.Off()
	gock.New(host).Delete("/routes/" + routeID).Reply(http.StatusNotFound)

	err := RemoveRoute(host, routeID)
	if err == nil {
		t.Errorf("error not reported for nonexistent route")
	} else if err.Error() != expected {
		t.Errorf("error mismatch: expected ‘%s’, got ‘%s’", expected, err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
