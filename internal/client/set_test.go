package client_test

import (
	"net/http"
	"strings"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/BBVA/kapow/internal/client"
)

// Test that not found errors are detected as invalid handler id
func TestNotFound(t *testing.T) {
	expectedErr := "Not Found"
	host := "http://localhost:8080"
	hid := "inventedID"
	path := "/response/status/code"
	reader := strings.NewReader("200")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusNotFound)

	if err := client.SetData(host, hid, path, reader); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != expectedErr {
		t.Errorf("Error don't match: expected \"%s\", got \"%s\"", expectedErr, err.Error())
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

// Test a http ok request
func TestOkRequest(t *testing.T) {
	host := "http://localhost:8080"
	hid := "HANDLER_XXXXXXXXXXXX"
	path := "/response/status/code"
	reader := strings.NewReader("200")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusOK)

	if err := client.SetData(host, hid, path, reader); err != nil {
		t.Error("Unexpected error")
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
