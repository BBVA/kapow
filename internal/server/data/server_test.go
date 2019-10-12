package data

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetRequestMethodReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestMethod(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "GET" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "GET", value)
	}
}

func TestGetRequestHostReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestHost(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "www.example.com" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "www.example.com", value)
	}
}

func TestGetRequestPathReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestPath(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "/this/is/a/test" {
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
