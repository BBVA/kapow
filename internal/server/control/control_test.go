package control

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestAddRouteWhenMalformedJSONBodyReturnsBadRequest(t *testing.T) {
	t.Skip("****** WIP ******")
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

func TestAddRouteWhenMandatoryFieldsMissingReturns422Error(t *testing.T) {
	t.Skip("****** WIP ******")
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

	respJson := model.Route{}
	if err := json.Unmarshal(resp.Body.Bytes(), &respJson); err == nil {
		t.Errorf("Invalid JSON response. %s", resp.Body.String())
	}

	expectedRouteSpec := model.Route{Method: "GET", Pattern: "/hello", Entrypoint: "/bin/sh -c", Command: "echo Hello World | kapow set /response/body", Index: 0, ID: "ROUTE_XXXXXXXXXXXXXXXXXX"}
	if respJson != expectedRouteSpec {
		t.Errorf("Response mismatch. Expected %#v, got: %#v", expectedRouteSpec, respJson)
	}
}

func TestAppendNewRouteFromRequest(t *testing.T) {
	t.Skip("****** WIP ******")
	reqPayload := `{
  "method": "GET",
  "url_pattern": "/hello",
  "entrypoint": null,
  "command": "echo Hello World | kapow set /response/body"
}`

	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(addRoute)

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusCreated, resp.Code)
	}

	//	expectedPayload := `{
	//  "method": "GET",
	//  "url_pattern": "/hello",
	//  "entrypoint": null,
	//  "command": "echo Hello World | kapow set /response/body"
	//}`
	respJsonRoute := model.Route{}
	_ = json.Unmarshal(resp.Body.Bytes(), &respJsonRoute)
	if respJsonRoute.Method != "GET" {
		t.Errorf("Method mismatch. Expected: %s, got: %s", "GET", respJsonRoute.Method)
	}

	if respJsonRoute.Entrypoint != "" {
		t.Errorf("Entrypoint mismatch. Expected: %s, got: %s", "", respJsonRoute.Entrypoint)
	}

	if respJsonRoute.Command != "echo Hello World | kapow set /response/body" {
		t.Errorf("Command mismatch. Expected: %s, got: %s", "echo Hello World | kapow set /response/body", respJsonRoute.Command)
	}

	if respJsonRoute.Pattern != "/hello" {
		t.Errorf("Pattern mismatch. Expected: %s, got: %s", "/hello", respJsonRoute.Pattern)
	}

	if respJsonRoute.Index > 0 {
		t.Errorf("Index mismatch. Expected: %d, got: %d", 0, respJsonRoute.Index)
	}
}
