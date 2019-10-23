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
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestRemoveRouteOKExistent(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Delete("/routes/ROUTE_FOO").
		Reply(http.StatusNoContent)

	err := RemoveRoute("http://localhost:8080", "ROUTE_FOO")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}

func TestRemoveRouteErrorNonExistent(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8080").
		Delete("/routes/ROUTE_BAD").
		Reply(http.StatusNotFound)

	err := RemoveRoute("http://localhost:8080", "ROUTE_BAD")
	if err == nil {
		t.Errorf("Error not reported for nonexistent route")
	} else if err.Error() != "Not Found" {
		t.Errorf(`Error mismatch: got %q, want "Not Found"`, err)
	}

	if !gock.IsDone() {
		t.Errorf("No endpoint called")
	}
}
