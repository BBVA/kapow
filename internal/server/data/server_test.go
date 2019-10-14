package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func TestGetRequestMethodReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	if value := getRequestMethod(req); value != "GET" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "GET", value)
	}
}

func TestGetRequestHostReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	if value := getRequestHost(req); value != "www.example.com" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "www.example.com", value)
	}
}

func TestGetRequestPathReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	if value := getRequestPath(req); value != "/this/is/a/test" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "/this/is/a/test", value)
	}
}

func TestSetResponseStatusSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseStatus(res, 500)
	if val := res.Result().StatusCode; val != 500 {
		t.Errorf("Unexpected value. Expected: %d, got: %d", 500, val)
	}
}

func TestGetRequestHeaderReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")

	value, err := getRequestHeader(req, "A-Header")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestHeaderReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")

	if _, err := getRequestHeader(req, "Other-Header"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestSetResponseHeaderSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseHeader(res, "A-Header", "With-Value")
	if val := res.Result().Header.Get("A-Header"); val != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", val)
	}
}

func TestGetRequestCookieReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	value, err := getRequestCookie(req, "A-Cookie")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestCookieReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	if _, err := getRequestCookie(req, "Other-Cookie"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestSetResponseCookieSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseCookie(res, "A-Cookie", "With-Value")
	cookies := res.Result().Cookies()
	val := ""
	for _, v := range cookies {
		if v.Name == "A-Cookie" {
			val = v.Value
		}
	}
	if val != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", val)
	}
}

func TestGetRequestParamReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	value, err := getRequestParam(req, "with")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "params" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "params", value)
	}
}

func TestGetRequestParamReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	if _, err := getRequestParam(req, "Other-Param"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestGetRequestFormReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.PostForm = url.Values{}
	req.PostForm.Set("A-Field", "With-Value")
	req.PostForm.Set("Another-Field", "With-AnotherValue")

	value, err := getRequestForm(req, "A-Field")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestFormReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.PostForm = url.Values{}
	req.PostForm.Set("A-Field", "With-Value")
	req.PostForm.Set("Another-Field", "With-AnotherValue")

	if _, err := getRequestForm(req, "Other-Field"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestGetRequestFileNameReturnsCorrectValue(t *testing.T) {
	multPartBody := bytes.Buffer{}
	multPartWriter := multipart.NewWriter(&multPartBody)
	part, _ := multPartWriter.CreateFormFile("A-File", "filename.txt")
	_, _ = part.Write([]byte("This is the file content\n"))
	_ = multPartWriter.Close()

	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", &multPartBody)
	req.Header.Add("Content-Type", multPartWriter.FormDataContentType())
	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	value, err := getRequestFileName(req, "A-File")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "filename.txt" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "filename.txt", value)
	}
}

func TestGetRequestFileNameReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	if _, err := getRequestFileName(req, "Other-File"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestCopyRequestFileReturnsOK(t *testing.T) {
	multPartBody := bytes.Buffer{}
	multPartWriter := multipart.NewWriter(&multPartBody)
	part, _ := multPartWriter.CreateFormFile("A-File", "filename.txt")
	_, _ = part.Write([]byte("This is the file content\n"))
	_ = multPartWriter.Close()

	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", &multPartBody)
	req.Header.Add("Content-Type", multPartWriter.FormDataContentType())
	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	result := strings.Builder{}
	err := copyRequestFile(req, "A-File", &result)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value := result.String(); value != "This is the file content\n" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestCopyRequestFileReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	result := strings.Builder{}
	if err := copyRequestFile(req, "Other-File", &result); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func generateTargetRequestForMatch() *http.Request {
	var targetRequest *http.Request

	h := mux.NewRouter()
	h.HandleFunc("/a/{foo}", func(res http.ResponseWriter, req *http.Request) { targetRequest = req }).Methods("GET")
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/a/bar", nil))

	return targetRequest
}

func TestGetRequestMatchReturnsCorrectValue(t *testing.T) {
	req := generateTargetRequestForMatch()

	value, err := getRequestMatch(req, "foo")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "bar" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "bar", value)
	}
}

func TestGetRequestMatchReturnsErrorWhenNotExists(t *testing.T) {
	req := generateTargetRequestForMatch()

	if _, err := getRequestMatch(req, "bar"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestCopyFromRequestBodyReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", strings.NewReader("This is a body content for testing purposes"))

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	result := strings.Builder{}
	if err := copyFromRequestBody(req, &result); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value := result.String(); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestCopyToResponseBodyReturnsOK(t *testing.T) {
	res := httptest.NewRecorder()

	if err := copyToResponseBody(res, strings.NewReader("This is a body content for testing purposes")); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Unexpected error while reading result body: %+v", err)
	}

	if value := string(bodyBytes); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestCopyToResponseStreamReturnsOK(t *testing.T) {
	res := httptest.NewRecorder()

	if err := copyToResponseStream(res, strings.NewReader("This is a body content for testing purposes")); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Unexpected error while reading result body: %+v", err)
	}

	if value := string(bodyBytes); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestRouterIsWellConfigured(t *testing.T) {
	testCases := []struct {
		pattern, method string
		handler         func(http.ResponseWriter, *http.Request)
		mustMatch       bool
		vars            []struct{ k, v string }
	}{
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/method", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "method"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/host", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "host"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/path", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "path"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/matches/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "matches/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/params/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "params/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "headers/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/cookies/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "cookies/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/form/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "form/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/files/name/filename", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "files/name/filename"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/files/name/content", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "files/name/content"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/body", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "body"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodPost, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodPut, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodDelete, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodGet, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodPost, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/status", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "status"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "headers/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/cookies/name", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "cookies/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/body", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "body"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/stream", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "stream"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodDelete, nil, false, []struct{ k, v string }{}},
	}

	r := configRouter()

	for _, tc := range testCases {
		rm := mux.RouteMatch{}
		rq, _ := http.NewRequest(tc.method, tc.pattern, nil)
		if matched := r.Match(rq, &rm); tc.mustMatch != matched {
			t.Errorf("Route mismatch: Expected: %+v\n\t\t\t\t\t\t got: %+v", tc, rm)
		} else {
			if tc.mustMatch {
				// Check for Handler match.
				realHandler := reflect.ValueOf(rm.Handler).Pointer()
				expectedHandler := reflect.ValueOf(tc.handler).Pointer()
				if realHandler != expectedHandler {
					t.Errorf("Handler mismatch. Expected: %X, got: %X", expectedHandler, realHandler)
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
		}
	}
}

func TestReadRequestResourcesReturnsNotFoundWhenHandlerIDNotExist(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/method", nil)
	resp := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/request/{resource_path:.*$}", readRequestResources).
		Methods(http.MethodGet)

	getHandler = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXXXXXX" {
			return &model.Handler{ID: id /*Request *http.Request, Writer http.ResponseWriter*/}, true
		}

		return nil, false
	}

	handler.ServeHTTP(resp, req)
	if got := resp.Result().StatusCode; got != http.StatusNotFound {
		t.Errorf("Unexpected status code. Expected: %d, got: %d", http.StatusNotFound, got)
	}
}

func TestWriteResponseResourcesReturnsNotFoundWhenHandlerIDNotExist(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/status", nil)
	resp := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/response/{resource_path:.*$}", writeResponseResources).
		Methods(http.MethodPut)

	getHandler = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXXXXXX" {
			return &model.Handler{ID: id /*Request *http.Request, Writer http.ResponseWriter*/}, true
		}

		return nil, false
	}

	handler.ServeHTTP(resp, req)
	if got := resp.Result().StatusCode; got != http.StatusNotFound {
		t.Errorf("Unexpected status code. Expected: %d, got: %d", http.StatusNotFound, got)
	}
}

func TestReadRequestResourcesReturnsBadRequestWhenInvalidResource(t *testing.T) {
	testCases := []string{
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/foo",
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/bar",
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/poor",
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/request/{resource_path:.*$}", readRequestResources).
		Methods(http.MethodGet)

	getHandler = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXXXXXX" {
			return &model.Handler{ID: id /*Request *http.Request, Writer http.ResponseWriter*/}, true
		}

		return nil, false
	}

	for _, testURL := range testCases {
		req := httptest.NewRequest(http.MethodGet, testURL, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if got := resp.Result().StatusCode; got != http.StatusBadRequest {
			t.Errorf("Unexpected status code. Expected: %d, got: %d", http.StatusBadRequest, got)
		}
	}
}

func TestWriteResponseResourcesReturnsBadRequestWhenInvalidResource(t *testing.T) {
	testCases := []string{
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/foo",
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/bar",
		"/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/poor",
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/response/{resource_path:.*$}", writeResponseResources).
		Methods(http.MethodPut)

	getHandler = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXXXXXX" {
			return &model.Handler{ID: id /*Request *http.Request, Writer http.ResponseWriter*/}, true
		}

		return nil, false
	}

	for _, testURL := range testCases {
		req := httptest.NewRequest(http.MethodPut, testURL, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if got := resp.Result().StatusCode; got != http.StatusBadRequest {
			t.Errorf("Unexpected status code. Expected: %d, got: %d", http.StatusBadRequest, got)
		}
	}
}

func TestReadRequestResourcesReturns(t *testing.T) {
	testCases := []struct {
		name, method, url string
		statusCode        int
		expectedBody      string
	}{
		{"Get method", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/method", 200, http.MethodPut},
		{"Get host", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/host", 200, "www.example.com"},
		{"Get path", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/path", 200, "/this/is/a/test"},
		{"Get body", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/body", 200, "bar for testing purposes"},
		{"Get param", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/params/with", 200, "params"},
		{"Get unexistent param", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/params/other", 404, ""},
		{"Get invalid param", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/params", 400, ""},
		{"Get header", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/headers/A-Header", 200, "With-Value"},
		{"Get unexistent header", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/headers/Other-Header", 404, ""},
		{"Get invalid header", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/headers", 400, ""},
		{"Get cookie", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/cookies/A-Cookie", 200, "With-Value"},
		{"Get unexistent cookie", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/cookies/Other-Cookie", 404, ""},
		{"Get invalid cookie", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/cookies", 400, ""},
		{"Get form field", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/form/A-Field", 200, "With-Value"},
		{"Get unexistent form field", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/form/Other-Field", 404, ""},
		{"Get invalid form field", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/form", 400, ""},
		{"Get match", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/matches/what", 200, "test"},
		{"Get unexistent match", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/matches/that", 404, ""},
		{"Get invalid match", http.MethodGet, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/request/matches", 400, ""},
	}

	getHandler = func(id string) (*model.Handler, bool) {
		if id == "HANDLER_XXXXXXXXXXXXXXXX" {
			var targetRequest *http.Request
			h := mux.NewRouter()
			h.HandleFunc("/this/is/a/{what}", func(res http.ResponseWriter, req *http.Request) { targetRequest = req }).Methods(http.MethodPut)
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPut, "http://www.example.com/this/is/a/test?with=params", strings.NewReader("bar for testing purposes")))

			targetRequest.Header.Add("A-Header", "With-Value")
			targetRequest.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})
			targetRequest.PostForm = url.Values{}
			targetRequest.PostForm.Set("A-Field", "With-Value")

			return &model.Handler{
					ID:      id,
					Request: targetRequest, /*Writer http.ResponseWriter*/
				},
				true
		}

		return nil, false
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/request/{resource_path:.*$}", readRequestResources).
		Methods(http.MethodGet)

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.url, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if got := resp.Result().StatusCode; got != tc.statusCode {
			t.Errorf("Unexpected status code for request %q. Expected: %d, got: %d", tc.name, tc.statusCode, got)
		} else {
			if tc.expectedBody != "" {
				if bodyBytes, err := ioutil.ReadAll(resp.Result().Body); err == nil {
					if got := string(bodyBytes); tc.expectedBody != got {
						t.Errorf("Unexpected response body for request %q. Expected: %s, got: %s", tc.name, tc.expectedBody, got)
					}
				} else {
					t.Errorf("Unexpected error reading response body: %v", err)
				}
			}
		}
	}
}

type testCase struct {
	name, method, url string
	statusCode        int
	payload           string
	validate          func(*http.Response, testCase) error
}

func TestWriteResponseResourcesReturns(t *testing.T) {
	testCases := []testCase{
		{"Set invalid status", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/status", 400, "hola", nil},
		{"Set status", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/status", 200, "300",
			func(res *http.Response, tc testCase) error {
				if res.StatusCode != 300 {
					return fmt.Errorf("Unexpected status code for request %q. Expected: %s, got: %d", tc.name, tc.payload, res.StatusCode)
				}
				return nil
			},
		},
		{"Set body", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/body", 200, "bar for testing purposes",
			func(res *http.Response, tc testCase) error {
				if bodyBytes, _ := ioutil.ReadAll(res.Body); tc.payload != string(bodyBytes) {
					return fmt.Errorf("Unexpected response body for request %q. Expected: %s, got: %s", tc.name, tc.payload, string(bodyBytes))
				}
				return nil
			},
		},
		{"Set stream", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/stream", 200, "bar for testing purposes",
			func(res *http.Response, tc testCase) error {
				if bodyBytes, _ := ioutil.ReadAll(res.Body); tc.payload != string(bodyBytes) {
					return fmt.Errorf("Unexpected response body for request %q. Expected: %s, got: %s", tc.name, tc.payload, string(bodyBytes))
				}
				return nil
			},
		},
		{"Set header", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/headers/A-Header", 200, "With-Value",
			func(res *http.Response, tc testCase) error {
				if res.Header.Get("A-Header") != tc.payload {
					return fmt.Errorf("Unexpected header value for request %q. Expected: %s, got: %s", tc.name, tc.payload, res.Header.Get("A-Header"))
				}
				return nil
			},
		},
		{"Set invalid header", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/headers", 400, "", nil},
		{"Set cookie", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/cookies/A-Cookie", 200, "With-Value",
			func(res *http.Response, tc testCase) error {
				if len(res.Cookies()) < 1 {
					return fmt.Errorf("Unexpected result for request %q. No cookies found", tc.name)
				} else if res.Cookies()[0].Name != "A-Cookie" || res.Cookies()[0].Value != tc.payload {
					return fmt.Errorf("Unexpected cookie value for request %q. Expected: %s, got: %s", tc.name, tc.payload, res.Cookies()[0].Name+"="+res.Cookies()[0].Value)
				}
				return nil
			},
		},
		{"Set invalid cookie", http.MethodPut, "/handlers/HANDLER_XXXXXXXXXXXXXXXX/response/cookies", 400, "", nil},
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/handlers/{handler_id}/response/{resource_path:.*$}", writeResponseResources).
		Methods(http.MethodPut)

	for _, tc := range testCases {
		targetResponse := httptest.NewRecorder()
		getHandler = func(id string) (*model.Handler, bool) {
			if id == "HANDLER_XXXXXXXXXXXXXXXX" {
				var targetRequest *http.Request
				h := mux.NewRouter()
				h.HandleFunc("/this/is/a/{what}", func(res http.ResponseWriter, req *http.Request) { targetRequest = req }).Methods(http.MethodPut)
				h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPut, "http://www.example.com/this/is/a/test?with=params", strings.NewReader("bar for testing purposes")))

				targetRequest.Header.Add("A-Header", "With-Value")
				targetRequest.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})
				targetRequest.PostForm = url.Values{}
				targetRequest.PostForm.Set("A-Field", "With-Value")

				return &model.Handler{
						ID:      id,
						Request: targetRequest,
						Writer:  targetResponse,
					},
					true
			}

			return nil, false
		}
		req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.payload))
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if got := resp.Result().StatusCode; got != tc.statusCode {
			t.Errorf("Unexpected status code for request %q. Expected: %d, got: %d", tc.name, tc.statusCode, got)
		} else {
			if tc.validate != nil {
				if err := tc.validate(targetResponse.Result(), tc); err != nil {
					t.Error(err.Error())
				}
			}
		}
	}
}
