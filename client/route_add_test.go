package client

import "testing"

func TestInvalidURL(t *testing.T) {
	err := AddRoute("http://localhost;8080", "/hi", "GET", "bash -c", "echo 'Hi' | kapow set /response/body")
	if err == nil {
		t.Error("expect to fail due invalid url")
		t.Fail()
	}
}
