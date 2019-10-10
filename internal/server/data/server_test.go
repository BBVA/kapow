package data

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

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

func TestUpdateResourceAddHeaderWhenRecieved(t *testing.T) {
	t.Skip("**** WIP ****")
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/header/pepe", strings.NewReader("mola"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/{resource:.*$}", updateResource).
		Methods("PUT")

	handlerResponse := httptest.NewRecorder()
	myHandler := &model.Handler{
		ID:     "HANDLER_YYYYYYYYYYYYYYYY",
		Writer: handlerResponse,
	}
	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return myHandler, true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusOK, response.Code)
	}

	headerValue := handlerResponse.Header().Get("pepe")
	if headerValue != "mola" {
		t.Errorf("Invalid Cookie value. Expected: %s, got: %s", "mola", headerValue)
	}
}
