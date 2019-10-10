package control

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestConfigRouterHasRoutesWellConfigured(t *testing.T) {
	testCases := []struct {
		pattern, method string
		handler         uintptr
		mustMatch       bool
		vars            []string
	}{
		{"/routes/ROUTE_YYYYYYYYYYYYYYY", http.MethodGet, 0, false, []string{}},
		{"/routes/ROUTE_YYYYYYYYYYYYYYY", http.MethodPut, 0, false, []string{}},
		{"/routes/ROUTE_YYYYYYYYYYYYYYY", http.MethodPost, 0, false, []string{}},
		{"/routes/ROUTE_YYYYYYYYYYYYYYY", http.MethodDelete, reflect.ValueOf(removeRoute).Pointer(), true, []string{"id"}},
		{"/routes", http.MethodGet, reflect.ValueOf(listRoutes).Pointer(), true, []string{}},
		{"/routes", http.MethodPut, 0, false, []string{}},
		{"/routes", http.MethodPost, reflect.ValueOf(addRoute).Pointer(), true, []string{}},
		{"/routes", http.MethodDelete, 0, false, []string{}},
	}
	r := configRouter()

	for _, tc := range testCases {
		rm := mux.RouteMatch{}
		rq, _ := http.NewRequest(tc.method, tc.pattern, nil)
		if matched := r.Match(rq, &rm); tc.mustMatch == matched {
			if tc.mustMatch {
				// Check for Handler match.
				realHandler := reflect.ValueOf(rm.Handler).Pointer()
				if realHandler != tc.handler {
					t.Errorf("Handler mismatch. Expected: %X, got: %X", tc.handler, realHandler)
				}

				// Check for variables
				for _, vn := range tc.vars {
					if _, exists := rm.Vars[vn]; !exists {
						t.Errorf("Variable not present: %s", vn)
					}
				}
			}
		} else {
			t.Errorf("Route mismatch: %+v", tc)
		}
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
	handler := http.HandlerFunc(addRoute)

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusBadRequest, resp.Code)
	}

}

func TestAddRouteReturns422ErrorWhenMandatoryFieldsMissing(t *testing.T) {
	handler := http.HandlerFunc(addRoute)
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

		handler.ServeHTTP(resp, req)
		if test.testMustFail {
			if resp.Code != http.StatusUnprocessableEntity {
				t.Errorf("HTTP status mismatch in case %s. Expected: %d, got: %d", test.testCase, http.StatusUnprocessableEntity, resp.Code)
			}
		} else if !test.testMustFail {
			if resp.Code != http.StatusCreated {
				t.Errorf("HTTP status mismatch in case %s. Expected: %d, got: %d", test.testCase, http.StatusUnprocessableEntity, resp.Code)
			}

			if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
			}
		}
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
	handler := http.HandlerFunc(addRoute)

	funcAdd = func(input model.Route) model.Route {
		expected := model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body"}
		if input == expected {
			input.Index = 0
			input.ID = "ROUTE_XXXXXXXXXXXXXXXXXX"
			return input
		}

		return model.Route{}
	}

	handler.ServeHTTP(resp, req)
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

	expectedRouteSpec := model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"}
	if respJson != expectedRouteSpec {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteSpec, respJson)
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
	if resp.Code != http.StatusNotFound {
		t.Errorf("HTTP status mismatch. Expected: %d, got: %d", http.StatusNotFound, resp.Code)
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
			model.Route{Method: "GET", Pattern: "/hello1", Entrypoint: "/bin/sh -c", Command: "echo Hello World1 | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"},
			model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 1, ID: "ROUTE_YYYYYYYYYYYYYYYYYY"},
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
		model.Route{Method: "GET", Pattern: "/hello1", Entrypoint: "/bin/sh -c", Command: "echo Hello World1 | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"},
		model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 1, ID: "ROUTE_YYYYYYYYYYYYYYYYYY"},
	}

	if !reflect.DeepEqual(respJson, expectedRouteList) {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteList, respJson)
	}
}
