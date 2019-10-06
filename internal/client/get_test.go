package client

import (
	"bytes"
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestWriteContentToWriter(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Get("/handlers/HANDLER_BAR/request/body").
		Reply(http.StatusOK).
		BodyString("FOO")

	var b bytes.Buffer
	err := GetData(
		"http://localhost", "HANDLER_BAR", "/request/body", &b)

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if !bytes.Equal(b.Bytes(), []byte("FOO")) {
		t.Errorf("Received content mismatch: %q != %q", b.Bytes(), []byte("FOO"))
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestPropagateHTTPError(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Get("/handlers/HANDLER_BAR/request/body").
		Reply(http.StatusTeapot)

	err := GetData(
		"http://localhost", "HANDLER_BAR", "/request/body", nil)

	if err == nil {
		t.Errorf("Expected error not returned")
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}
