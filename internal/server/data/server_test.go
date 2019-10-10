package data

import (
	"testing"
)

func TestConfigRouterHasRoutesWellConfigured(t *testing.T) {
	t.Skip("****** WIP ******")
	//	testCases := []struct {
	//		pattern, method string
	//		handler         uintptr
	//		mustMatch       bool
	//		vars            []string
	//	}{
	//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodGet, reflect.ValueOf().Pointer(), true, []struct{ key, value string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"resource", "params/name"}}},
	//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodPut, reflect.ValueOf().Pointer(), true, []struct{ key, value string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"resource", "params/name"}}},
	//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/cookies/name", http.MethodGet, reflect.ValueOf().Pointer(), true, []struct{ key, value string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"resource", "cookies/name"}}},
	//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/cookies/name", http.MethodPut, reflect.ValueOf().Pointer(), true, []struct{ key, value string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"resource", "cookies/name"}}},
	//	}
	//	r := configRouter()
	//
	//	for _, tc := range testCases {
	//		rm := mux.RouteMatch{}
	//		rq, _ := http.NewRequest(tc.method, tc.pattern, nil)
	//		if matched := r.Match(rq, &rm); tc.mustMatch == matched {
	//			if tc.mustMatch {
	//				// Check for Handler match.
	//				realHandler := reflect.ValueOf(rm.Handler).Pointer()
	//				if realHandler != tc.handler {
	//					t.Errorf("Handler mismatch. Expected: %X, got: %X", tc.handler, realHandler)
	//				}
	//
	//				// Check for variables
	//				for _, vn := range tc.vars {
	//					if _, exists := rm.Vars[vn]; !exists {
	//						t.Errorf("Variable not present: %s", vn)
	//					}
	//				}
	//			}
	//		} else {
	//			t.Errorf("Route mismatch: %+v", tc)
	//		}
	//	}
}
