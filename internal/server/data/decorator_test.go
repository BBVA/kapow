/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package data

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestLockResponseWriterReturnsAFunctionsThatCallsTheGivenCallback(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	called := false

	fn := lockResponseWriter(func(http.ResponseWriter, *http.Request, *model.Handler) { called = true })

	fn(w, r, &h)
	if !called {
		t.Error("Callback not called")
	}
}

func TestLockResponseWriterReturnsAFunctionThatWaitsForTheLockToBeReleased(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	h.Writing.Lock()
	defer h.Writing.Unlock()

	fn := lockResponseWriter(func(http.ResponseWriter, *http.Request, *model.Handler) {})

	c := make(chan bool)
	go func() { fn(w, r, &h); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Lock not acquired during call")
	default: // This default prevents the select from being blocking
	}
}

func TestLockResponseWriterReturnsAFunctionReleaseTheLockAfterCallingTheGivenCallback(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	fn := lockResponseWriter(func(http.ResponseWriter, *http.Request, *model.Handler) {})

	fn(w, r, &h)

	c := make(chan bool)
	go func() { h.Writing.Lock(); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Lock not released after call")
	}
}

func TestLockResponseWriterReturnsAFunctionReleaseTheLockEvenIfTheGivenCallbackPanics(t *testing.T) {
	h := model.Handler{
		Request: httptest.NewRequest("POST", "/", nil),
		Writer:  httptest.NewRecorder(),
	}
	r := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	fn := lockResponseWriter(func(http.ResponseWriter, *http.Request, *model.Handler) { panic("BOOM!") })
	defer func() {
		_ = recover()

		c := make(chan bool)
		go func() { h.Writing.Lock(); c <- true }()

		time.Sleep(10 * time.Millisecond)

		select {
		case <-c:
		default: // This default prevents the select from being blocking
			t.Error("Lock not released after panic")
		}
	}()

	fn(w, r, &h)
}

func TestCheckHandlerReturnsAFunctionsThat404sWhenHandlerDoesNotExist(t *testing.T) {
	Handlers = New()
	r := createMuxRequest("/handlers/{handlerID}", "/handlers/BAZ", "GET", nil)
	w := httptest.NewRecorder()
	fn := checkHandler(func(http.ResponseWriter, *http.Request, *model.Handler) {})

	fn(w, r)

	for _, e := range checkErrorResponse(w.Result(), http.StatusNotFound, "Handler ID Not Found") {
		t.Error(e)
	}
}

func TestCheckHandlerReturnsAFunctionsThatCallsTheGivenCallbackWhenHandlerExists(t *testing.T) {
	Handlers = New()
	Handlers.Add(&model.Handler{ID: "BAZ"})
	r := createMuxRequest("/handlers/{handlerID}", "/handlers/BAZ", "GET", nil)
	w := httptest.NewRecorder()
	called := false

	fn := checkHandler(func(http.ResponseWriter, *http.Request, *model.Handler) { called = true })

	fn(w, r)
	if !called {
		t.Error("Callback not called")
	}
}

func TestCheckHandlerReturnsAFunctionsThatCallsTheGivenCallbackWithTheProperHandler(t *testing.T) {
	Handlers = New()
	Handlers.Add(&model.Handler{ID: "BAZ"})
	r := createMuxRequest("/handlers/{handlerID}", "/handlers/BAZ", "GET", nil)
	w := httptest.NewRecorder()
	var handlerID string

	fn := checkHandler(func(w http.ResponseWriter, r *http.Request, h *model.Handler) { handlerID = h.ID })

	fn(w, r)
	if handlerID != "BAZ" {
		t.Errorf(`Handler mismatch. Expected "BAZ". Got %q`, handlerID)
	}
}
