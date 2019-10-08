package control

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAppendNerRouteFromRequest(t *testing.T) {
	reqPayload := "{}"

	req := httptest.NewRequest(http.MethodPost, "/routes", strings.NewReader(reqPayload))
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(addRoute)

	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusCreated, resp.Code)
	}

  expectedPayload := "{}"
  if resp.Body.String() != expectedPayload {
    t.Errorf("HTTP status mistmacht. Expected: %d, got: %d", http.StatusCreated, resp.Code)
  }
}
