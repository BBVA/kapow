package client_test

import (
	"net/http"
	"strings"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/BBVA/kapow/internal/client"
)

// Test an HTTP OK request
func TestSetDataSuccessOnCorrectRequest(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Put("/HANDLER_FOO/response/status/code").
		Reply(http.StatusOK)

	if err := client.SetData(
		"http://localhost:8080",
		"HANDLER_FOO",
		"/response/status/code",
		strings.NewReader("200"),
	); err != nil {
		t.Error("Unexpected error")
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

// Test that Not Found errors are detected when an invalid handler id is sent
func TestSetDataErrIfBadHandlerID(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Put("/HANDLER_BAD/response/status/code").
		Reply(http.StatusNotFound)

	if err := client.SetData(
		"http://localhost:8080",
		"HANDLER_BAD",
		"/response/status/code",
		strings.NewReader("200"),
	); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != "Not Found" {
		t.Errorf(`Error mismatch: expected "Not Found", got %q`, err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
