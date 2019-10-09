package control

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
)

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
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusBadRequest, resp.Code)
	}

}

func TestAddRouteReturns422ErrorWhenMandatoryFieldsMissing(t *testing.T) {
	reqPayload := `{
  }`

	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(addRoute)

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusUnprocessableEntity, resp.Code)
	}

}

func TestAddRouteReturnsCreated(t *testing.T) {
	t.Skip("****** WIP ******")
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
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusCreated, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}

	respJson := model.Route{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respJson); err == nil {
		t.Errorf("Invalid JSON response. %s", resp.Body.String()) //FIXME: String comparsion not working, comparing against itself?
	}

	expectedRouteSpec := model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"}
	if respJson != expectedRouteSpec {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteSpec, respJson)
	}
}

func TestRemoveRouteReturnsNotFound(t *testing.T) {
	t.Skip("****** WIP ******")
	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(removeRoute)

	funcRemove = func(id string) error {
		if id == "ROUTE_XXXXXXXXXXXXXXXXXX" {
			return errors.New(id)
		}

		return nil
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusNotFound, resp.Code)
	}
}

func TestRemoveRouteReturnsNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(removeRoute)

	funcRemove = func(id string) error {
		return nil
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusNoContent, resp.Code)
	}
}

// FIXME: ListRoutes is a get, no path params, call
func TestListRoutesReturnsEmptyList(t *testing.T) {
	t.Skip("****** WIP ******")

	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(removeRoute)

	funcList = func() []model.Route {

		return []model.Route{}
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusNotFound, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}
}

func TestListRoutesReturnsTwoElementsList(t *testing.T) {
	t.Skip("****** WIP ******")

	req := httptest.NewRequest(http.MethodDelete, "/routes/ROUTE_XXXXXXXXXXXXXXXXXX", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(removeRoute)

	funcList = func() []model.Route {
		return []model.Route{
			model.Route{Method: "GET", Pattern: "/hello1", Entrypoint: "/bin/sh -c", Command: "echo Hello World1 | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"},
			model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 1, ID: "ROUTE_YYYYYYYYYYYYYYYYYY"},
		}
	}

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusOK, resp.Code)
	}

	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Incorrect content type in response. Expected: application/json, got: %s", ct)
	}

	respJson := []model.Route{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respJson); err == nil {
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
