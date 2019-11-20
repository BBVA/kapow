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

package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BBVA/kapow/internal/server/httperror"
	"github.com/BBVA/kapow/internal/server/model"
)

func checkErrorResponse(r *http.Response, expectedErrcode int, expectedReason string) []error {
	errList := make([]error, 0)

	if r.StatusCode != expectedErrcode {
		errList = append(errList, fmt.Errorf("HTTP status mismatch. Expected: %d, got: %d", expectedErrcode, r.StatusCode))
	}

	if v := r.Header.Get("Content-Type"); v != "application/json; charset=utf-8" {
		errList = append(errList, fmt.Errorf("Content-Type header mismatch. Expected: %q, got: %q", "application/json; charset=utf-8", v))
	}

	errMsg := httperror.ServerErrMessage{}
	if bodyBytes, err := ioutil.ReadAll(r.Body); err != nil {
		errList = append(errList, fmt.Errorf("Unexpected error reading response body: %v", err))
	} else if err := json.Unmarshal(bodyBytes, &errMsg); err != nil {
		errList = append(errList, fmt.Errorf("Response body contains invalid JSON entity: %v", err))
	} else if errMsg.Reason != expectedReason {
		errList = append(errList, fmt.Errorf("Unexpected reason in response. Expected: %q, got: %q", expectedReason, errMsg.Reason))
	}

	return errList
}

func TestConfigRouterReturnsRouterWithDecoratedRoutes(t *testing.T) {
	var handlerID string
	rs := []routeSpec{
		{
			"/handlers/{handlerID}/dummy",
			"GET",
			func(w http.ResponseWriter, r *http.Request, h *model.Handler) { handlerID = h.ID },
		},
	}
	Handlers = New()
	Handlers.Add(&model.Handler{ID: "FOO"})
	m := configRouter(rs)

	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/handlers/FOO/dummy", nil))

	if handlerID != "FOO" {
		t.Errorf(`Handler ID mismatch. Expected "FOO". Got %q`, handlerID)
	}
}

func TestConfigRouterReturnsRouterThat400sOnUnconfiguredResources(t *testing.T) {
	m := configRouter([]routeSpec{})
	w := httptest.NewRecorder()

	m.ServeHTTP(w, httptest.NewRequest("GET", "/handlers/FOO/dummy", nil))

	for _, e := range checkErrorResponse(w.Result(), http.StatusBadRequest, "Invalid Resource Path") {
		t.Error(e)
	}
}
