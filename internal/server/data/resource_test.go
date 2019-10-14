package data

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
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
		Request: httptest.NewRequest("POST", "/foo", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/handlers/HANDLER_ID/request/body", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("Status code mismatch")
	}
}

func TestGetRequestsBodySetsOctectStreamContentType(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/handlers/HANDLER_ID/request/body", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type mismatch")
	}
}

func TestGetRequestBodyWritesHandlerRequestBodyToResponseWriter(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo", strings.NewReader("BAR")),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/handlers/HANDLER_ID/request/body", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if body, _ := ioutil.ReadAll(res.Body); string(body) != "BAR" {
		t.Error("Body mismatch")
	}
}

func TestGetRequestBody500sWhenHandlerRequestErrors(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo", BadReader("User closed the connection")),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/handlers/HANDLER_ID/request/body", nil)
	w := httptest.NewRecorder()

	getRequestBody(w, r, &h)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Error("status not 500")
	}
}

func TestGetRequestBodyClosesConnectionWhenReaderErrorsAfterWrite(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/foo", ErrorOnSecondReadReader(strings.NewReader("FOO"))),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("GET", "/handlers/HANDLER_ID/request/body", nil)
	w := httptest.NewRecorder()
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("Didn't panic")
		}
	}()

	getRequestBody(w, r, &h)
}
