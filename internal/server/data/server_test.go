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
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/headers/name", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handlerId}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return createMockHandler(id, httptest.NewRecorder()), true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusNotFound, response.Code)
	}
}

func TestUpdateResourceOkWhenValidHandlerID(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYY/response/headers/name", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handlerId}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYY" {
			return createMockHandler(id, httptest.NewRecorder()), true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusOK, response.Code)
	}
}

func TestUpdateResourceBadRequestWhenInvalidUrl(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handlerId}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYYYYYY" {
			return createMockHandler(id, httptest.NewRecorder()), true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusBadRequest, response.Code)
	}
}

func TestUpdateResourceBadRequestWhenInvalidCookiesUrl(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYY/response/cookies", strings.NewReader("value"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handlerId}/{resource:.*$}", updateResource).
		Methods("PUT")

	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYY" {
			return createMockHandler(id, httptest.NewRecorder()), true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusBadRequest, response.Code)
	}
}

func TestUpdateResourceAddHeaderWhenRecieved(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYY/response/headers/pepe", strings.NewReader("mola"))
	response := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handlerId}/{resource:.*$}", updateResource).
		Methods("PUT")

	handlerInResponse := httptest.NewRecorder()
	getHandlerId = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_YYYYYYYYYYYY" {
			return createMockHandler(id, handlerInResponse), true
		}

		return nil, false
	}

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("HTTP Status mismatch. Expected: %d, got: %d", http.StatusOK, response.Code)
	}

	headerValue := handlerInResponse.Result().Header.Get("pepe")
	if headerValue != "mola" {
		t.Errorf("Invalid Header value. Expected: %s, got: %s", "mola", headerValue)
	}
}

func createMockHandler(id string, writer http.ResponseWriter) *model.Handler {
	return &model.Handler{
		ID:     id,
		Writer: writer,
	}
}
