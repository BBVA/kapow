package data

import (
	"errors"
	// "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Error("status not 500")
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

func createMuxRequest(pattern, url, method string) (req *http.Request) {
	m := mux.NewRouter()
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) { req = r })
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(method, url, nil))
	return
}

func TestGetRequestMatches200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET"),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET")
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestMatchesSetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET"),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET")
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestMatchesReturnsTheCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/foo/{bar}", "/foo/BAZ", "GET"),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/bar", "GET")
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAZ" {
		t.Errorf("Body mismatch. Expected: BAZ. Got: %v", string(body))
	}

}

func TestGetRequestMatchesReturnsNotFoundWhenMatchDoesntExists(t *testing.T) {
	h := model.Handler{
		Request: createMuxRequest("/", "/", "GET"),
		Writer:  httptest.NewRecorder(),
	}

	r := createMuxRequest("/handlers/HANDLERID/request/matches/{name}", "/handlers/HANDLERID/request/matches/foo", "GET")
	w := httptest.NewRecorder()

	getRequestMatches(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code mismatch. Expected: 404. Got: %d", res.StatusCode)
	}
}

func TestGetRequestParams200sOnHappyPath(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET")
	w := httptest.NewRecorder()

	getRequestParams(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code mismatch. Expected: 404. Got: %d", res.StatusCode)
	}
}

// FIXME: Discuss how return multiple values
func TestGetRequestParamsReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/foo?bar=BAZ&bar=QUX", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := createMuxRequest("/handlers/HANDLERID/request/params/{name}", "/handlers/HANDLERID/request/params/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestHeadersReturnsTheCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("Bar", "BAZ")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
	w := httptest.NewRecorder()

	getRequestHeaders(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestHeadersReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.Header.Set("bar", "BAZ")
	h.Request.Header.Add("bar", "QUX")
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/cookies/{name}", "/handlers/HANDLERID/request/cookies/bar", "GET")
	w := httptest.NewRecorder()

	getRequestCookies(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code mismatch. Expected: 404, Got: %d", res.StatusCode)
	}
}

func TestGetRequestCookiesReturnsTheFirstCorrectMatchValue(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("GET", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "BAZ"})
	h.Request.AddCookie(&http.Cookie{Name: "bar", Value: "QUX"})
	r := createMuxRequest("/handlers/HANDLERID/request/headers/{name}", "/handlers/HANDLERID/request/headers/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code mismatch. Expected: 404, Got: %d", res.StatusCode)
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
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
	r := createMuxRequest("/handlers/HANDLERID/request/form/{name}", "/handlers/HANDLERID/request/form/bar", "GET")
	w := httptest.NewRecorder()

	getRequestForm(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code mismatch. Expected: 404, Got: %d", res.StatusCode)
	}
}
