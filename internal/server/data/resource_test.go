/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package data

import (
	"bytes"
	"errors"
	"io"
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

type badReader struct {
	errorMessage string
}

func (r *badReader) Read(p []byte) (int, error) {
	return 0, errors.New(r.errorMessage)
}

func BadReader(m string) io.Reader {
	return &badReader{errorMessage: m}
}

type errorOnSecondReadReader struct {
	r    io.Reader
	last bool
}

func (r *errorOnSecondReadReader) Read(p []byte) (int, error) {
	if r.last {
		return 0, errors.New("Second read failed by design")
	} else {
		r.last = true
		return r.r.Read(p)
	}
}

func ErrorOnSecondReadReader(r io.Reader) io.Reader {
	return &errorOnSecondReadReader{r: r}
}

func TestGetRequestBody200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestBodySetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestBodyWritesHandlerRequestBodyToResponseWriter(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader("BAR")),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAR" {
		t.Error("Body mismatch")
	}
}

func TestGetRequestBody500sWhenHandlerRequestErrors(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", BadReader("User closed the connection")),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

func TestGetRequestBodyClosesConnectionWhenReaderErrorsAfterWrite(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", ErrorOnSecondReadReader(strings.NewReader("FOO"))),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("Didn't panic")
		}
	}()

	getRequestBody(w, r, &h)
}

func TestGetRequestMethod200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestMethod(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestMethodSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestMethod(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestMethodReturnsTheCorrectMethod(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("FOO", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestMethod(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "FOO" {
		t.Error("Body mismatch")
	}
}

func TestGetRequestHost200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestHost(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestHostReturnsTheCorrectHostname(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "http://www.foo.bar:8080/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestHost(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "www.foo.bar:8080" {
		t.Error("Body mismatch")
	}
}

func TestGetRequestHostSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestHost(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestPath200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestPath(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestPathSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestPath(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestPathReturnsPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestPath(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "/foo" {
		t.Error("Body mismatch")
	}
}

func TestGetRequestPathDoesntReturnQueryStringParams(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo?bar=1", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/not-important-here", nil)
	w := httptest.NewRecorder()

	getRequestPath(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "/foo" {
		t.Errorf("Body mismatch. Expected: /foo. Got: %v", string(body))
	}
}

func createMuxRequest(pattern, url, method string, content io.Reader) (req *http.Request) {
	m := mux.NewRouter()
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) { req = r })
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(method, url, content))
	return
}

func TestGetRequestMatches200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestMatchesSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestMatchesReturnsTheCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}

}

func TestGetRequestMatchesReturnsNotFoundWhenMatchDoesntExists(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/", "/", "GET", nil),
		Writer:  httptest.NewRecorder(),
	}

	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/foo", "GET", nil)
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func TestGetRequestParams200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestParamsSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestParamsReturnsTheCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestParams404sWhenParamDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

// FIXME: Discuss how return multiple values
func TestGetRequestParamsReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ&bar=QUX", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestHeaders200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "BAZ")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestHeadersSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "BAZ")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if v := res.Header.Get("Content-Type"); v != "application/octet-stream" {
		t.Errorf("Content Type mismatch. Expected: application/octet-stream. Got: %q", v)
	}
}

func TestGetRequestHeadersReturnsTheCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Bar", "BAZ")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestHeaders200sWhenHeaderIsEmptyString(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestHeadersReturnsEmptyBodyWhenHeaderIsEmptyString(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "" {
		t.Errorf(`Body mismatch. Expected "". Got: %q`, string(body))
	}
}

func TestGetRequestHeadersReturnsTheCorrectInsensitiveMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "BAZ")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestHeaders404sWhenHeaderDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func TestGetRequestHeadersReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "BAZ")
	h.Request.Header.Add("bar", "QUX")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestCookies200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "BAZ"})
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestCookiesSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "BAZ"})
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestCookiesReturnsMatchedCookieValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "BAZ"})
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestCookies404sIfCookieDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func TestGetRequestCookiesReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "BAZ"})
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "QUX"})
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

// NOTE: The current implementation doesn't allow us to decode
// form encoded data sent in a request with an arbitrary method. This is
// needed for Kapow! semantic so it MUST be changed in the future

// FIXME: Test form decoding with GET method
// FIXME: Test form decoding without Content-Type:
// application/x-www-form-urlencoded header

func TestGetRequestForm200sOnHappyPath(t *testing.T) {
	form := url.Values{}
	form.Add("bar", "BAZ")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestFormSetsOctectStreamContentType(t *testing.T) {
	form := url.Values{}
	form.Add("bar", "BAZ")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestFormReturnsTheCorrectMatchValue(t *testing.T) {
	form := url.Values{}
	form.Add("bar", "BAZ")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}
}

func TestGetRequestForm404sWhenFieldDoesntExist(t *testing.T) {
	form := url.Values{}
	form.Add("foo", "BAZ")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func TestGetRequestForm200sWhenFieldIsEmptyString(t *testing.T) {
	form := url.Values{}
	form.Add("bar", "")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestFormReturnsEmptyBodyWhenFieldIsEmptyString(t *testing.T) {
	form := url.Values{}
	form.Add("bar", "")
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", strings.NewReader(form.Encode())),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "" {
		t.Errorf(`Body mismatch. Expected: "". Got: %q`, string(body))
	}
}

// TODO: Discuss how to manage this use case, Not Found, Bad Request, ...
func TestGetRequestForm404sWhenFormDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET", nil)
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func createMultipartFileRequest(fieldName, filename, fileContent string) *http.Request {
	multPartBody := bytes.Buffer{}
	multPartWriter := multipart.NewWriter(&multPartBody)
	part, _ := multPartWriter.CreateFormFile(fieldName, filename)
	_, _ = part.Write([]byte(fileContent))
	_ = multPartWriter.Close()
	r := httptest.NewRequest("POST", "/", &multPartBody)
	r.Header.Add("Content-Type", multPartWriter.FormDataContentType())
	return r
}

func TestGetRequestFileName200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "foo", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/filename", "/handlers/HANDLERID/request/files/bar/filename", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileName(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestFileNameSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "foo", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/filename", "/handlers/HANDLERID/request/files/bar/filename", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileName(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestFileNameReturnsTheCorrectFilename(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "BAZ", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/filename", "/handlers/HANDLERID/request/files/bar/filename", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileName(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf(`Body mismatch. Expected: "BAZ". Got: %q`, string(body))
	}
}

func TestGetRequestFileName404sWhenFileDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("foo", "qux", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/filename", "/handlers/HANDLERID/request/files/bar/filename", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileName(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

// TODO: Discuss which error is appropiate when the form doesn't exists
func TestGetRequestFileName404sWhenFormDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/filename", "/handlers/HANDLERID/request/files/bar/filename", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileName(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

func TestGetRequestFileContent200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "foo", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestGetRequestFileContentSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "foo", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestFileContentReturnsTheCorrectFileContent(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("bar", "foo", "BAZ"),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf(`Body mismatch. Expected: "BAZ". Got: %q`, string(body))
	}
}

func TestGetRequestFileContent404sWhenFileDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: createMultipartFileRequest("foo", "qux", ""),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

// TODO: Discuss which error is appropiate when the form doesn't exists
func TestGetRequestFileContent404sWhenFormDoesntExist(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Resource Item Not Found") {
		t.Error(e)
	}
}

// TODO: Discuss what happens when request is interrupted
func TestGetRequestFileContent500sWhenHandlerRequestErrors(t *testing.T) {
	t.Skip("Undefined behavior")
	multPartBody := bytes.Buffer{}
	multPartWriter := multipart.NewWriter(&multPartBody)
	part, _ := multPartWriter.CreateFormFile("bar", "BAZ")
	_, _ = part.Write([]byte("qux"))
	_ = multPartWriter.Close()
	buf := bytes.NewBuffer(multPartBody.Bytes()[:len(multPartBody.Bytes())-1])

	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", ErrorOnSecondReadReader(buf)),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Add("Content-Type", multPartWriter.FormDataContentType())
	r := createMuxRequest("/handlers/HANDLERID/request/files/{name}/content", "/handlers/HANDLERID/request/files/bar/content", "GET", nil)
	w := httptest.NewRecorder()

	getRequestFileContent(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

func TestSetResponseStatus200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader("200"))
	w := httptest.NewRecorder()

	setResponseStatus(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestSetResponseStatusSetsGivenStatus(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader("418"))
	w := httptest.NewRecorder()

	setResponseStatus(w, r, &h)

	res := hw.Result()
	if res.StatusCode != http.StatusTeapot {
		t.Errorf("Status code mismatch. Expected: 418, Got: %d", res.StatusCode)
	}
}

func TestSetResponseStatus400sWhenNonparseableStatusCode(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader("foo"))
	w := httptest.NewRecorder()

	setResponseStatus(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusUnprocessableEntity, "Non Integer Value") {
		t.Error(e)
	}
}

func TestSetResponseStatus500sWhenErrorReadingRequest(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", BadReader("Failed by design"))
	w := httptest.NewRecorder()

	setResponseStatus(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

// FIXME: This is not the spec behavior but Go checks too many things to
// be sure. Discuss how to fix this.
func TestSetResponseStatus400sWhenStatusCodeNotSupportedByGo(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", strings.NewReader("99"))
	w := httptest.NewRecorder()

	setResponseStatus(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusBadRequest, "Invalid Status Code") {
		t.Error(e)
	}
}

func TestSetResponseHeaders200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/response/headers/{name}", "/handlers/HANDLERID/response/headers/bar", "PUT", strings.NewReader("BAZ"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestSetResponseHeadersSetsGivenHeader(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/headers/{name}", "/handlers/HANDLERID/response/headers/bar", "PUT", strings.NewReader("BAZ"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	res := hw.Result()
	// nolint
	if values, ok := res.Header["bar"]; !ok || values[0] != "BAZ" {
		t.Errorf(`Header mismatch. Expected "BAZ". Contents %v`, res.Header)
	}
}

// TODO: Decide if setResponseHeader should ADD or OVERWRITE headers
func TestSetResponseHeadersAddsGivenHeaderWhenAlreadySet(t *testing.T) {
	hw := httptest.NewRecorder()
	hw.Header()["bar"] = []string{"BAZ"}
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/headers/{name}", "/handlers/HANDLERID/response/headers/bar", "PUT", strings.NewReader("QUX"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	res := hw.Result()
	// nolint
	if values, ok := res.Header["bar"]; !ok || !reflect.DeepEqual(values, []string{"BAZ", "QUX"}) {
		t.Errorf(`Header mismatch. Expected ["BAZ", "QUX"]. Contents %v`, res.Header)
	}
}

func TestSetResponseHeaders500sWhenErrorReadingRequest(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", BadReader("Failed by design"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

// TODO: Validate Header Key encoding
func TestSetResponseHeaders400sOnInvalidHeaderKey(t *testing.T) {
	t.Skip("Somebody has to validate header key, but net/http doesn't give us any facility (yet).")
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/headers/{name}", "/handlers/HANDLERID/response/headers/%0a", "PUT", strings.NewReader("BAZ"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusBadRequest, "Invalid Header Name") {
		t.Error(e)
	}
}

// TODO: Validate Header Value encoding
func TestSetResponseHeaders400sOnInvalidHeaderValue(t *testing.T) {
	t.Skip("Somebody has to validate header value, but net/http doesn't give us any facility (yet).")
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/headers/{name}", "/handlers/HANDLERID/response/headers/foo", "PUT", strings.NewReader("\n"))
	w := httptest.NewRecorder()

	setResponseHeaders(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusBadRequest, "Invalid Header Value") {
		t.Error(e)
	}
}

func TestSetResponseCookies200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	setResponseCookies(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestSetResponseCookiesSetsGivenCookie(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/cookies/{name}", "/handlers/HANDLERID/response/cookies/bar", "PUT", strings.NewReader("BAZ"))
	w := httptest.NewRecorder()

	setResponseCookies(w, r, &h)

	res := hw.Result()
	// nolint
	if c := res.Cookies(); len(c) != 1 || c[0].Name != "bar" || c[0].Value != "BAZ" {
		t.Errorf(`Header mismatch. Expected "bar=BAZ". Contents %+v`, c)
	}
}

func TestSetResponseCookies500sWhenErrorReadingRequest(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", BadReader("Failed by design"))
	w := httptest.NewRecorder()

	setResponseCookies(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

func TestSetResponseCookiesAddsValueToExistingCookie(t *testing.T) {
	hw := httptest.NewRecorder()
	c := &http.Cookie{Name: "bar", Value: "BAZ"}
	http.SetCookie(hw, c)
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/cookies/{name}", "/handlers/HANDLERID/response/cookies/bar", "PUT", strings.NewReader("QUX"))
	w := httptest.NewRecorder()

	setResponseCookies(w, r, &h)

	res := hw.Result()
	// nolint
	if c := res.Cookies(); len(c) != 2 || c[1].Name != "bar" || c[1].Value != "QUX" {
		t.Errorf(`Header mismatch. Expected "bar=QUX". Contents %+v`, c)
	}
}

// TODO: Validate Cookie Name&Value encoding

func TestSetResponseBody200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	setResponseBody(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code mismatch. Expected: 200, Got: %d", res.StatusCode)
	}
}

func TestSetResponseBodySetsTheResponseBody(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/body", "/handlers/HANDLERID/response/body", "PUT", strings.NewReader("BAZ"))
	w := httptest.NewRecorder()

	setResponseBody(w, r, &h)

	res := hw.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf(`Body mismatch. Expected: "BAZ". Got: %q`, string(body))
	}
}

func TestSetResponseBody500sWhenReaderFailsInFirstRead(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/body", "/handlers/HANDLERID/response/body", "PUT", BadReader("Fail by design"))
	w := httptest.NewRecorder()

	setResponseBody(w, r, &h)

	for _, e := range checkErrorResponse(w.Result(), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)) {
		t.Error(e)
	}
}

func TestSetResponseBodyPanicsIfReaderFailsAfterFirstWrite(t *testing.T) {
	hw := httptest.NewRecorder()
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  hw,
	}
	r := createMuxRequest("/handlers/HANDLERID/response/body", "/handlers/HANDLERID/response/body", "PUT", ErrorOnSecondReadReader(strings.NewReader("FOO")))
	w := httptest.NewRecorder()

	defer func() {
		if rec := recover(); rec == nil {
			t.Error("Didn't panic")
		}
	}()
	setResponseBody(w, r, &h)
}
