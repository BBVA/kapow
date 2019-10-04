package client

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

const (
	host = "http://localhost:8080"
)

func TestListRoutesEmpty(t *testing.T) {
	descr := fmt.Sprintf("ListRoutes(%q, nil)", host)

	defer gock.Off()
	gock.New(host).
		Get("/routes").
		Reply(http.StatusOK)

	err := ListRoutes(host, nil)
	if err != nil {
		t.Errorf("%s: unexpected error %q", descr, err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

func TestListRoutesSome(t *testing.T) {
	descr := fmt.Sprintf("ListRoutes(%q, buf)", host)

	const want = "JSON array of some routes..."

	defer gock.Off()
	gock.New(host).
		Get("/routes").
		Reply(http.StatusOK).
		BodyString(want)

	buf := new(bytes.Buffer)
	err := ListRoutes(host, buf)
	if err != nil {
		t.Errorf("%s: unexpected error: %q", descr, err)
	} else if got := buf.String(); got != want {
		t.Errorf("%s: got %q, expected %q", descr, buf, want)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
