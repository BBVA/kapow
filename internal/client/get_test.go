package client

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestGetInvalidUrl(t *testing.T) {
	err := GetData("", "", "", os.Stdout)
	if err == nil {
		t.Error("Expected error with invalid url ''")
	}
}

func TestGetInvalidWriter(t *testing.T) {
	err := GetData("http://localhost:8081", "0000", "/", nil)
	if err == nil {
		t.Error("Expected error with no writer")
	}
}

func TestGetURLNotFoundWithUnknownID(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:8081").
		Get("/handlers/000/").Reply(http.StatusNotFound)

	err := GetData("http://localhost:8081", "000", "/", os.Stdout)

	if err == nil {
		t.Errorf("Expect not found error but get no error")
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}

func TestGetRetrieveRequestMethod(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:8081").
		Get("/handlers/000/request/method").
		Reply(http.StatusAccepted).
		BodyString("POST")

	rw := new(bytes.Buffer)

	err := GetData("http://localhost:8081", "000", "/request/method", rw)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	strRes := rw.String()

	if strRes != "POST" {
		t.Errorf("POST string expected but found: '%v'", strRes)
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}
