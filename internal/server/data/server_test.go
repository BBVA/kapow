package data

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

//func TestConfigRouterHasRoutesWellConfigured(t *testing.T) {
//	testCases := []struct {
//		pattern, method string
//		handler         uintptr
//		mustMatch       bool
//		vars            []struct{ k, v string }
//	}{
//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodGet, reflect.ValueOf(readResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "request"}, {"resource", "params/name"}}},
//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/request/params/name", http.MethodPut, reflect.ValueOf(updateResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "request"}, {"resource", "params/name"}}},
//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/cookies/name", http.MethodGet, reflect.ValueOf(readResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"root", "response"}, {"resource", "cookies/name"}}},
//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/headers/", http.MethodPut, reflect.ValueOf(updateResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}}},
//		{"/handlers/HANDLER_ZZZZZZZZZZZZZZZZ/response/headers/name", http.MethodPut, reflect.ValueOf(updateResource).Pointer(), true, []struct{ k, v string }{{"handler_id", "HANDLER_ZZZZZZZZZZZZZZZZ"}, {"key", "name"}}},
//	}
//	r := configRouter()
//
//	for _, tc := range testCases {
//		rm := mux.RouteMatch{}
//		rq, _ := http.NewRequest(tc.method, tc.pattern, nil)
//		if matched := r.Match(rq, &rm); tc.mustMatch != matched {
//			t.Errorf("Route mismatch: Expected: %+v\n\t\t\t\t\t\t got: %+v", tc, rm)
//		} else {
//			if tc.mustMatch {
//				// Check for Handler match.
//				realHandler := reflect.ValueOf(rm.Handler).Pointer()
//				if realHandler != tc.handler {
//					t.Errorf("Handler mismatch. Expected: %X, got: %X", tc.handler, realHandler)
//				}
//
//				// Check for variables
//				for _, v := range tc.vars {
//					if value, exists := rm.Vars[v.k]; !exists {
//						t.Errorf("Variable not present: %s", v.k)
//					} else if v.v != value {
//						t.Errorf("Variable value mismatch. Expected: %s, got: %s", v.v, value)
//					}
//				}
//			}
//		}
//	}
//}

func TestUpdateResourceNotFoundWhenInvalidHandlerID(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return nil, false
		}

		return nil, true
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusNotFound, response.Code)
	}
}

func TestUpdateResourceBadRequestWhenInvalidUrl(t *testing.T) {

	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return nil, true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusBadRequest, response.Code)
	}
}

func TestUpdateResourceOkWhenValidHandlerID(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXX/response/headers/name", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXX" {
			return nil, true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusOK, response.Code)
	}
}

func TestUpdateResourceBadRequestWhenInvalidCookiesUrl(t *testing.T) {
	t.Skip("**** WIP ****")
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/cookies", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return nil, true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusBadRequest, response.Code)
	}
}

// FIXME: Fails because URL doesn't match
//func TestUpdateResourceNotFoundWhenInvalidHandlerID(t *testing.T) {
//	request := httptest.NewRequest(http.MethodPut, "/handlers/response/headers/language", strings.NewReader("ES"))
//	response := httptest.NewRecorder()
//	handler := configRouter()
//
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusNotFound {
//		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusNotFound, response.Code)
//	}
//}

//func TestUpdateResourceBadRequestWhenIncompletedResourceURL(t *testing.T) {
//	request := httptest.NewRequest(http.MethodPut, "/handlers/xxxxxxxxx/response/headers/", strings.NewReader("ES"))
//	response := httptest.NewRecorder()
//	handler := configRouter()
//
//	getHandlerId = func(id string) (*model.Handler, bool) {
//		if id == "xxxxxxxxx" {
//			return nil, true
//		}
//		return nil, false
//	}
//
//	handler.ServeHTTP(response, request)
//	// TODO: We need to assure that an invalid resource path returns 400 (Bad Request)
//	if response.Code != http.StatusBadRequest {
//		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusBadRequest, response.Code)
//	}
//}

//func TestUpdateResourceSetHeaderWhenPutReceived(t *testing.T) {
//	request := httptest.NewRequest(http.MethodPut, "/handlers/xxxxxxxxxx/response/headers/language", strings.NewReader("ES"))
//	response := httptest.NewRecorder()
//	handler := configRouter()
//
//	getHandlerId = func(id string) (*model.Handler, bool) {
//		if id == "xxxxxxxxxx" {
//			return nil, true
//		}
//		return nil, false
//	}
//
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusOK {
//		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusOK, response.Code)
//	}
//}
//
