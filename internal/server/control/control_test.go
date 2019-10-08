package control

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestAppendNewRouteFromRequest(t *testing.T) {
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
	json.Unmarshal(resp.Body.Bytes(), &respJsonRoute)
	if respJsonRoute.Method != "GET" {
		t.Errorf("Method missmatch. Expected: %s, got: %s", "GET", respJsonRoute.Method)
	}

	if respJsonRoute.Entrypoint != "" {
		t.Errorf("Entrypoint missmatch. Expected: %s, got: %s", "", respJsonRoute.Entrypoint)
	}

	if respJsonRoute.Command != "echo Hello World | kapow set /response/body" {
		t.Errorf("Command missmatch. Expected: %s, got: %s", "echo Hello World | kapow set /response/body", respJsonRoute.Command)
	}

	if respJsonRoute.Pattern != "/hello" {
		t.Errorf("Pattern missmatch. Expected: %s, got: %s", "/hello", respJsonRoute.Pattern)
	}
}
