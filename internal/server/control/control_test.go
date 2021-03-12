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

package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/httperror"
	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user"
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

func TestPathValidatorNoErrorWhenCorrectPath(t *testing.T) {
	err := pathValidator("/routes/{routeID}")

	if err != nil {
		t.Error(err)
	}
}

func TestPathValidatorErrorWhenInvalidPath(t *testing.T) {
	err := pathValidator("/routes/{routeID{")

	if err == nil {
		t.FailNow()
	}
}

func TestAddRouteReturnsBadRequestWhenMalformedJSONBody(t *testing.T) {
	reqPayload := `{
	method": "GET",
	url_pattern": "/hello",
	entrypoint": null,
	command": "echo Hello World | kapow set /response/body"
  }`

	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()

	addRoute(resp, req)

	for _, e := range checkErrorResponse(resp.Result(), http.StatusBadRequest, "Malformed JSON") {
		t.Error(e)
	}
}

func TestAddRouteReturns422ErrorWhenMandatoryFieldsMissing(t *testing.T) {
	tc := []struct {
		payload, testCase string
		testMustFail      bool
	}{
		{`{}`, "EmptyBody", true},
		{`{
	  "method": "GET"
	  }`,
			"Missing url_pattern",
			true,
		},
		{`{
	  "url_pattern": "/hello"
	  }`,
			"Missing method",
			true,
		},
		{`{
	  "method": "GET",
	  "url_pattern": "/hello"
	  }`,
			"",
			false,
		},
		{`{
	  "method": "GET",
	  "url_pattern": "/hello",
	  "entrypoint": ""
	  }`,
			"",
			false,
		},
		{`{
	  "method": "GET",
	  "url_pattern": "/hello",
	  "command": ""
	  }`,
			"",
			false,
		},
		{`{
	  "method": "GET",
	  "url_pattern": "/hello",
	  "entrypoint": "",
	  "command": ""
	  }`,
			"",
			false,
		},
	}

	for _, test := range tc {
		req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(test.payload))
		resp := httptest.NewRecorder()

		addRoute(resp, req)
		r := resp.Result()
		if test.testMustFail {
			for _, e := range checkErrorResponse(r, http.StatusUnprocessableEntity, "Invalid Route") {
				t.Error(e)
			}
		} else if !test.testMustFail {
			if r.StatusCode != http.StatusCreated {
				t.Errorf("HTTP status mismatch in case %s. Expected: %d, got: %d", test.testCase, http.StatusUnprocessableEntity, r.StatusCode)
			}

			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Incorrect content type in response. Expected: application/json, got: %q", ct)
			}
		}
	}
}

func TestAddRouteGeneratesRouteID(t *testing.T) {
	reqPayload := `{
	"method": "GET",
	"url_pattern": "/hello",
	"entrypoint": "/bin/sh -c",
	"command": "echo Hello World | kapow set /response/body"
  }`
	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()
	var genID string
	funcAdd = func(input model.Route) model.Route {
		genID = input.ID
		input.Index = 0
		return input
	}
	origPathValidator := pathValidator
	defer func() { pathValidator = origPathValidator }()
	pathValidator = func(path string) error { return nil }

	addRoute(resp, req)

	if _, err := uuid.Parse(genID); err != nil {
		t.Error("ID not generated properly")
	}
}

func TestAddRoute500sWhenIDGeneratorFails(t *testing.T) {
	reqPayload := `{
	"method": "GET",
	"url_pattern": "/hello",
	"entrypoint": "/bin/sh -c",
	"command": "echo Hello World | kapow set /response/body"
  }`
	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()

	origPathValidator := pathValidator
	defer func() { pathValidator = origPathValidator }()
	pathValidator = func(path string) error { return nil }

	idGenOrig := idGenerator
	defer func() { idGenerator = idGenOrig }()
	idGenerator = func() (uuid.UUID, error) {
		var uuid uuid.UUID
		return uuid, errors.New("End of Time reached; Try again before, or in the next Big Bang cycle")
	}

	addRoute(resp, req)

	for _, e := range checkErrorResponse(resp.Result(), http.StatusInternalServerError, "Internal Server Error") {
		t.Error(e)
	}
}

func TestAddRouteReturnsCreated(t *testing.T) {
	reqPayload := `{
	"method": "GET",
	"url_pattern": "/hello",
	"entrypoint": "/bin/sh -c",
	"command": "echo Hello World | kapow set /response/body"
  }`

	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()
	var genID string
	funcAdd = func(input model.Route) model.Route {
		expected := model.Route{ID: input.ID, Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body"}
		if input == expected {
			genID = input.ID
			input.Index = 0
			return input
		}

		return model.Route{}
	}
	origPathValidator := pathValidator
	defer func() { pathValidator = origPathValidator }()
	pathValidator = func(path string) error { return nil }

	addRoute(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusCreated, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}

	respJson := model.Route{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respJson); err != nil {
		t.Errorf("Invalid JSON response. %s", resp.Body.String())
	}

	expectedRouteSpec := model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 0, ID: genID}
	if respJson != expectedRouteSpec {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteSpec, respJson)
	}
}

func TestAddRoute422sWhenInvalidRoute(t *testing.T) {
	reqPayload := `{
	"method": "GET",
	"url_pattern": "/he{{o",
	"entrypoint": "/bin/sh -c",
	"command": "echo Hello World | kapow set /response/body"
}`
	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()
	origPathValidator := pathValidator
	defer func() { pathValidator = origPathValidator }()
	pathValidator = func(path string) error { return errors.New("Invalid route") }

	addRoute(resp, req)

	for _, e := range checkErrorResponse(resp.Result(), http.StatusUnprocessableEntity, "Invalid Route") {
		t.Error(e)
	}
}

func TestRemoveRouteReturnsNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/routes/{id}", removeRoute).
		Methods("DELETE")
	funcRemove = func(id string) error {
		if id == "ROUTE_XXXXXXXXXXXXXXXXXX" {
			return errors.New(id)
		}

		return nil
	}

	handler.ServeHTTP(resp, req)

	for _, e := range checkErrorResponse(resp.Result(), http.StatusNotFound, "Route Not Found") {
		t.Error(e)
	}
}

func TestRemoveRouteReturnsNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/routes/{id}", removeRoute).
		Methods("DELETE")

	funcRemove = func(id string) error {
		if id == "ROUTE_XXXXXXXXXXXXXXXXXX" {
			return nil
		}
		return errors.New(id)
	}

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusNoContent, resp.Code)
	}
}

func TestListRoutesReturnsEmptyList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/routes/", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(listRoutes)

	funcList = func() []model.Route {

		return []model.Route{}
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusOK, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}
}

func TestListRoutesReturnsTwoElementsList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/routes", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(listRoutes)

	funcList = func() []model.Route {
		return []model.Route{
			{Method: "GET", Pattern: "/hello1", Entrypoint: "/bin/sh -c", Command: "echo Hello World1 | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"},
			{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 1, ID: "ROUTE_YYYYYYYYYYYYYYYYYY"},
		}
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusOK, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}

	respJson := []model.Route{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respJson); err != nil {
		t.Errorf("Invalid JSON response. %s", resp.Body.String())
	}

	expectedRouteList := []model.Route{
		{Method: "GET", Pattern: "/hello1", Entrypoint: "/bin/sh -c", Command: "echo Hello World1 | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"},
		{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 1, ID: "ROUTE_YYYYYYYYYYYYYYYYYY"},
	}

	if !reflect.DeepEqual(respJson, expectedRouteList) {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteList, respJson)
	}
}

func TestGetRouteReturns404sWhenRouteDoesntExist(t *testing.T) {
	handler := mux.NewRouter()
	handler.HandleFunc("/routes/{id}", getRoute).
		Methods("GET")
	r := httptest.NewRequest(http.MethodGet, "/routes/FOO", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Route Not Found") {
		t.Error(e)
	}
}

func TestGetRouteSetsCorrectContentType(t *testing.T) {
	handler := mux.NewRouter()
	handler.HandleFunc("/routes/{id}", getRoute).
		Methods("GET")
	r := httptest.NewRequest(http.MethodGet, "/routes/FOO", nil)
	w := httptest.NewRecorder()
	user.Routes.Append(model.Route{ID: "FOO"})

	handler.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	if hVal := resp.Header.Get("Content-Type"); hVal != "application/json" {
		t.Errorf(`Route mismatch. Expected: "application/json". Got: %s`, hVal)
	}
}

func TestGetRouteReturnsTheRequestedRoute(t *testing.T) {
	handler := mux.NewRouter()
	handler.HandleFunc("/routes/{id}", getRoute).
		Methods("GET")
	r := httptest.NewRequest(http.MethodGet, "/routes/FOO", nil)
	w := httptest.NewRecorder()
	user.Routes.Append(model.Route{ID: "FOO"})

	handler.ServeHTTP(w, r)

	resp := w.Result()
	respJson := model.Route{}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	bBytes, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(bBytes, &respJson); err != nil {
		t.Errorf("Invalid JSON response. %s", string(bBytes))
	}

	if respJson.ID != "FOO" {
		t.Errorf(`Route mismatch. Expected: "FOO". Got: %s`, respJson.ID)
	}
}
