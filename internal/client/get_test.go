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
	err := GetData("http://localhost", "HANDLER_BAR", "/request/body", &b)

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
