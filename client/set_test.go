package client_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/BBVA/kapow/client"

	gock "gopkg.in/h2non/gock.v1"
)

// Test that no content errors are detected as non-existent resource
func TestNoContent(t *testing.T) {
	expectedErr := "Resource Item Not Found"
	host := "http://localhost:8080"
	hid := "xxxxxxxxxxxxxx"
	path := "/unpath"
	reader := strings.NewReader("Esto es un peacho de dato pa repartir")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusNoContent)

	if err := client.SetData(host, hid, path, reader); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != expectedErr {
		t.Errorf("Error don't match: expected \"%s\", got \"%s\"", expectedErr, err.Error())
	}
}

// Test that bad request errors are detected as invalid resource
func TestBadRequest(t *testing.T) {
	expectedErr := "Invalid Resource Path"
	host := "http://localhost:8080"
	hid := "xxxxxxxxxxxxxx"
	path := "/unpath"
	reader := strings.NewReader("Esto es un peacho de dato pa repartir")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusBadRequest)

	if err := client.SetData(host, hid, path, reader); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != expectedErr {
		t.Errorf("Error don't match: expected \"%s\", got \"%s\"", expectedErr, err.Error())
	}
}

// Test that not found errors are detected as invalid handler id
func TestNotFound(t *testing.T) {
	expectedErr := "Not Found"
	host := "http://localhost:8080"
	hid := "xxxxxxxxxxxxxx"
	path := "/unpath"
	reader := strings.NewReader("Esto es un peacho de dato pa repartir")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusNotFound)

	if err := client.SetData(host, hid, path, reader); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != expectedErr {
		t.Errorf("Error don't match: expected \"%s\", got \"%s\"", expectedErr, err.Error())
	}
}

// Test that internal server errors are detected correctly
func TestInternalServerError(t *testing.T) {
	expectedErr := "Internal Server Error"
	host := "http://localhost:8080"
	hid := "xxxxxxxxxxxxxx"
	path := "/unpath"
	reader := strings.NewReader("Esto es un peacho de dato pa repartir")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusInternalServerError)

	if err := client.SetData(host, hid, path, reader); err == nil {
		t.Error("Expected error not present")
	} else if err.Error() != expectedErr {
		t.Errorf("Error don't match: expected \"%s\", got \"%s\"", expectedErr, err.Error())
	}
}

// Test a http ok request
func TestOkRequest(t *testing.T) {
	host := "http://localhost:8080"
	hid := "xxxxxxxxxxxxxx"
	path := "/response/status/code"
	reader := strings.NewReader("200")

	defer gock.Off()

	gock.New(host).Put("/" + hid + path).Reply(http.StatusOK)

	if err := client.SetData(host, hid, path, reader); err != nil {
		t.Error("Unexpected error")
	}
}
