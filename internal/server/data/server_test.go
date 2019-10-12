package data

import (
	"net/http/httptest"
	"testing"
)

func TestGetRequestMethodReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/this/is/a/test?with=params&that=works", nil)

	req.Host = "www.example.com"
	value, err := getRequestMethod(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "GET" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "GET", value)
	}
}

func TestGetRequestHostReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/this/is/a/test?with=params&that=works", nil)

	req.Host = "www.example.com"
	value, err := getRequestHost(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "www.example.com" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "www.example.com", value)
	}
}

func TestGetRequestPathReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/this/is/a/test?with=params&that=works", nil)

	req.Host = "www.example.com"
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

	err := setResponseStatus(res, 500)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if res.Result().StatusCode != 500 {
		t.Errorf("Unexpected value. Expected: %d, got: %d", 500, res.Result().StatusCode)
	}
}

func TestSetResponseStatusFailsWhenNotInteger(t *testing.T) {
	res := httptest.NewRecorder()

	err := setResponseStatus(res, "200")
	if err == nil {
		t.Errorf("Expecting an error but got OK")
	}
}
