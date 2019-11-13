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

func TestListRoutesOKEmpty(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Get("/routes").
		Reply(http.StatusOK)

	err := ListRoutes("http://localhost:8080", nil)
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

func TestListRoutesOKSome(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Get("/routes").
		Reply(http.StatusOK).
		JSON([]map[string]string{
			{"foo": "bar"},
			{"bar": "foo"},
		})

	var b bytes.Buffer
	err := ListRoutes("http://localhost:8080", &b)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if !bytes.Equal(
		b.Bytes(), []byte(`[{"foo":"bar"},{"bar":"foo"}]`+"\n")) {
		t.Errorf("Unexpected error: got %q, want %q",
			b.String(), `[{"foo":"bar"},{"bar":"foo"}]`+"\n")
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
