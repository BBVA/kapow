package server_test

import (
	"testing"

	"github.com/BBVA/kapow/internal/server"
)

func TestStartServerWhenInvalidBindAddrReturnsError(t *testing.T) {

	err := server.StartServer("foo;bar", "", "", true)
	if err == nil {
		t.Errorf("Expected error not found")
	}
}

func TestStartServerWhenInvalidPortNumberReturnsError(t *testing.T) {

	err := server.StartServer("0.0.0.0:bar", "", "", true)
	if err == nil {
		t.Errorf("Expected error not found")
	}
}
