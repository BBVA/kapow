package http

import (
	nethttp "net/http"
	"testing"
)

func TestEmptyStringIsEmptyReason(t *testing.T) {
	r := &nethttp.Response{Status: ""}
	if GetReason(r) != "" {
		t.Errorf("We consider an empty status line to have an empty reason")
	}
}

func TestOnlyCodeIsEmptyReason(t *testing.T) {
	r := &nethttp.Response{Status: "200"}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestOnlyCodePlusSpaceIsEmptyReason(t *testing.T) {
	r := &nethttp.Response{Status: "200 "}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestOneWordReason(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO"}
	if GetReason(r) != "FOO" {
		t.Errorf("Unexpected reason found")
	}
}

func TestMultiWordReason(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}

func TestOddSizeStatusCode(t *testing.T) {
	r := &nethttp.Response{Status: "2 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}
