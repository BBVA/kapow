package srverrors_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BBVA/kapow/internal/server/srverrors"
)

func TestErrorJSONSetsAppJsonContentType(t *testing.T) {
	w := httptest.NewRecorder()

	srverrors.ErrorJSON(w, "Not Important Here", 0)

	if v := w.Result().Header.Get("Content-Type"); v != "application/json; charset=utf-8" {
		t.Errorf("Content-Type header mismatch. Expected: %q, got: %q", "application/json; charset=utf-8", v)
	}
}

func TestErrorJSONSetsRequestedStatusCode(t *testing.T) {
	w := httptest.NewRecorder()

	srverrors.ErrorJSON(w, "Not Important Here", http.StatusGone)

	if v := w.Result().StatusCode; v != http.StatusGone {
		t.Errorf("Status code mismatch. Expected: %d, got: %d", http.StatusGone, v)
	}
}

func TestErrorJSONSetsBodyCorrectly(t *testing.T) {
	expectedReason := "Something Not Found"
	w := httptest.NewRecorder()

	srverrors.ErrorJSON(w, expectedReason, http.StatusNotFound)

	errMsg := srverrors.ServerErrMessage{}
	if bodyBytes, err := ioutil.ReadAll(w.Result().Body); err != nil {
		t.Errorf("Unexpected error reading response body: %v", err)
	} else if err := json.Unmarshal(bodyBytes, &errMsg); err != nil {
		t.Errorf("Response body contains invalid JSON entity: %v", err)
	} else if errMsg.Reason != expectedReason {
		t.Errorf("Unexpected reason in response. Expected: %q, got: %q", expectedReason, errMsg.Reason)
	}
}
