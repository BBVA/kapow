package http

import (
	nethttp "net/http"
	"testing"
)

func TestEmptyReasonWhenEmptyString(t *testing.T) {
	r := &nethttp.Response{Status: ""}
	if GetReason(r) != "" {
		t.Errorf("We consider an empty status line to have an empty reason")
	}
}

func TestEmptyReasonWhenOnlyCode(t *testing.T) {
	r := &nethttp.Response{Status: "200"}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestEmptyReasonWhenOnlyCodePlusSpace(t *testing.T) {
	r := &nethttp.Response{Status: "200 "}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestReasonOfOneWord(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO"}
	if GetReason(r) != "FOO" {
		t.Errorf("Unexpected reason found")
	}
}

func TestReasonOfMultipleWords(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}

func TestBehaveWithOddSizeStatusCode(t *testing.T) {
	r := &nethttp.Response{Status: "2 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}
