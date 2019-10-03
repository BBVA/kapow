package http

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestReturnErrorOnInvalidURL(t *testing.T) {
	defer gock.Off()
	gock.New("").Reply(200)

	err := Request("GET", "://", nil, nil)
	if err == nil {
		t.Errorf("Expected error not returned")
	}

	if gock.IsDone() {
		t.Errorf("Request was performed anyway")
	}
}

func TestRequestGivenMethod(t *testing.T) {
	defer gock.Off()
	mock := gock.New("http://localhost")
	mock.Method = "FOO"
	mock.Reply(200)

	err := Request("FOO", "http://localhost", nil, nil)
	if err != nil {
		t.Errorf("Unexpected error on request")
	}

	if gock.IsDone() == false {
		t.Errorf("Expected request not performed")
	}
}

func TestReturnHTTPErrorAsIs(t *testing.T) {
	defer gock.Off()
	customError := errors.New("FOO")
	gock.New("http://localhost").ReplyError(customError)

	err := Request("GET", "http://localhost", nil, nil)
	if errors.Unwrap(err) != customError {
		t.Errorf("Returned error is not the expected error")
	}

	if gock.IsDone() == false {
		t.Errorf("Expected request not performed")
	}
}

func TestReturnHTTPReasonAsErrorWhenUnsuccessful(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").Reply(http.StatusTeapot)

	err := Request("GET", "http://localhost", nil, nil)
	if err == nil || err.Error() != http.StatusText(http.StatusTeapot) {
		t.Errorf("Reason should be returned as an error")
	}

	if gock.IsDone() == false {
		t.Errorf("Expected request not performed")
	}
}

func TestCopyResponseBodyToWriter(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost").Reply(200).BodyString("FOO")

	rw := new(bytes.Buffer)

	err := Request("GET", "http://localhost", nil, rw)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	res := rw.String()

	if res != "FOO" {
		t.Errorf("Unexpected output %v", res)
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}

func TestWriteToDevNullWhenNoWriter(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost").Reply(200).BodyString("FOO")

	original := devnull
	devnull = new(bytes.Buffer)

	defer func() { devnull = original }()

	err := Request("GET", "http://localhost", nil, nil)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	res := devnull.(*bytes.Buffer).String()

	if res != "FOO" {
		t.Errorf("Unexpected output %v", res)
	}

	if gock.IsDone() == false {
		t.Error("No expected endpoint called")
	}
}
