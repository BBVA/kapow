package client

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestWriteContentToWriter(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Get("/handlers/THIS-IS-THE-HANDLER-ID/request/body").
		Reply(http.StatusOK).
		Body(bytes.NewReader([]byte("FOO")))

	var b bytes.Buffer
	buf := bufio.NewWriter(&b)
	err := GetData("http://localhost", "THIS-IS-THE-HANDLER-ID", "/request/body", buf)

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if !bytes.Equal(b.Bytes(), []byte("FOO")) {
		t.Errorf("Received content mismatch: %q != %q", b.Bytes(), []byte("FOO"))
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}

func TestPropagateHTTPError(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Get("/handlers/THIS-IS-THE-HANDLER-ID/request/body").
		Reply(http.StatusTeapot)

	err := GetData("http://localhost", "THIS-IS-THE-HANDLER-ID", "/request/body", nil)

	if err == nil {
		t.Errorf("Expected error not returned")
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}
