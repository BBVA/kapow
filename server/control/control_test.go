package control

import "testing"

func TestNewControlServerFillsTheStruct(t *testing.T) {

	server := NewControlServer("0.0.0.0", 8080, "/certfile.pem", "/keyfile.pem")

	if server.bindAddr != "0.0.0.0:8080" {
		t.Errorf("BindAddress incorrectly composed. Expected: %s, got: %s", "0.0.0.0:8080", server.bindAddr)
	}

	if server.certfile != "/certfile.pem" {
		t.Errorf("BindAddress incorrectly composed. Expected: %s, got: %s", "/certfile.pem", server.certfile)
	}

	if server.keyfile != "/keyfile.pem" {
		t.Errorf("BindAddress incorrectly composed. Expected: %s, got: %s", "/keyfile.pem", server.keyfile)
	}
}
