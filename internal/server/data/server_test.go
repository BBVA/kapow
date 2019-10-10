package data

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestConfigRouterHasRoutesWellConfigured(t *testing.T) {
	testCases := []struct {
		pattern, method string
		handler         uintptr
		mustMatch       bool
		vars            []struct{ k, v string }
	}{
		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodGet, reflect.ValueOf(readResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "request"}, {"resource", "params/name"}}},
		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodPut, reflect.ValueOf(updateResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "request"}, {"resource", "params/name"}}},
		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/cookies/name", http.MethodGet, reflect.ValueOf(readResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "response"}, {"resource", "cookies/name"}}},
		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/cookies/name", http.MethodPut, reflect.ValueOf(updateResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "response"}, {"resource", "cookies/name"}}},
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
				for _, v := range tc.vars {
					if value, exists := rm.Vars[v.k]; !exists {
						t.Errorf("Variable not present: %s", v.k)
					} else if v.v != value {
						t.Errorf("Variable value mismatch. Expected: %s, got: %s", v.v, value)
					}
				}
			}
		} else {
			t.Errorf("Route mismatch: %+v", tc)
		}
	}
}
