package server

import (
	"errors"
	"testing"

	"github.com/BBVA/kapow/internal/server/config"
)

func TestRunServerPushNilWhenNoError(t *testing.T) {
	errors := make(chan error)

	go runServer(func(cfg config.ServerConfig) error { return nil }, config.ServerConfig{}, errors)

	err := <-errors

	if err != nil {
		t.Error("Error should be nil")
	}
}

func TestRunServerPushServerErrors(t *testing.T) {
	expected := errors.New("foo")
	errors := make(chan error)
	go runServer(func(cfg config.ServerConfig) error { return expected }, config.ServerConfig{}, errors)

	current := <-errors

	if expected != current {
		t.Error("Not the expected error")
	}
}

func TestRunServerPushErrorIfServerPanics(t *testing.T) {
	errors := make(chan error)
	go runServer(func(cfg config.ServerConfig) error { panic("foo") }, config.ServerConfig{}, errors)

	err := <-errors

	if err == nil {
		t.Error("Error should not be null on panic")
	}
}
