/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	if !gock.IsDone() {
		t.Errorf("Expected request not performed")
	}
}

func TestReturnHTTPErrorAsIs(t *testing.T) {
	defer gock.Off()
	customError := errors.New("FOO")
	gock.New("http://localhost").ReplyError(customError)

	err := Request("GET", "http://localhost", nil, nil)
	if errors.Unwrap(err) != customError {
		t.Errorf("Returned error is not the expected error: '%v'", err)
	}

	if !gock.IsDone() {
		t.Errorf("Expected request not performed")
	}
}

func TestReturnHTTPReasonAsErrorWhenUnsuccessful(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Reply(http.StatusTeapot).
		BodyString(`{"reason": "I'm a teapot"}`)

	err := Request("GET", "http://localhost", nil, nil)
	if err == nil || err.Error() != http.StatusText(http.StatusTeapot) {
		t.Errorf("Reason should be returned as an error")
	}

	if !gock.IsDone() {
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

	if !gock.IsDone() {
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

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestSendContentTypeJSON(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		MatchHeader("Content-Type", "application/json").
		Reply(http.StatusOK)

	err := Request("GET", "http://localhost", nil, nil, AsJSON)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err.Error())
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestGetRequestsWithMethodGet(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Get("/").
		Reply(http.StatusOK)

	err := Get("http://localhost/", nil, nil)

	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestPostRequestsWithMethodPost(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Post("/").
		Reply(http.StatusOK)

	err := Post("http://localhost/", nil, nil)

	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestPutRequestsWithMethodPut(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Put("/").
		Reply(http.StatusOK)

	err := Put("http://localhost/", nil, nil)

	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}

func TestDeleteRequestsWithMethodDelete(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Delete("/").
		Reply(http.StatusOK)

	err := Delete("http://localhost/", nil, nil)

	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	if !gock.IsDone() {
		t.Error("No expected endpoint called")
	}
}
