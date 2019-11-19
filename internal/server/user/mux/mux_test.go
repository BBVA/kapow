// +build !race

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

package mux

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

func TestNewReturnsAProperlyInitializedMux(t *testing.T) {
	sm := New()
	sm.root.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	sm.root.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Error("mux not properly initialized")
	}
}

func TestSwappableMuxGetReturnsTheCurrentMux(t *testing.T) {
	sm := SwappableMux{}
	mux := sm.get()
	if !reflect.DeepEqual(mux, sm.root) {
		t.Errorf("Returned mux is not the same %#v", mux)
	}
}

func TestSwappableMuxGetReturnsADifferentInstance(t *testing.T) {
	sm := SwappableMux{}
	mux := sm.get()
	if &mux == &sm.root {
		t.Error("Returned mux is the same instance")
	}
}

func TestSwappableMuxGetWaitsForTheMutexToBeReleased(t *testing.T) {
	sm := SwappableMux{}

	sm.m.Lock()
	defer sm.m.Unlock()

	c := make(chan *mux.Router)
	go func() { c <- sm.get() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}

func TestSwappableMuxGetIsAbleToReadWhileOthersAreReading(t *testing.T) {
	sm := SwappableMux{}

	sm.m.RLock()
	defer sm.m.RUnlock()

	c := make(chan *mux.Router)
	go func() { c <- sm.get() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default:
		t.Error("The mutex cannot be acquired")
	}
}

func TestSwappableMuxSetSetsTheGivenMux(t *testing.T) {
	sm := SwappableMux{}
	m := mux.NewRouter()
	// nolint
	m.KeepContext = true

	sm.set(m)

	// nolint
	if !sm.root.KeepContext {
		t.Error("mux not set")
	}
}

func TestSwappableMuxSetSetsTheSameInstance(t *testing.T) {
	sm := SwappableMux{}
	m := mux.NewRouter()

	sm.set(m)

	if m != sm.root {
		t.Error("Set mux is not the same instance")
	}
}

func TestSwappableMuxSetWaitsForWriterToReleaseMutex(t *testing.T) {
	sm := SwappableMux{}

	sm.m.Lock()
	defer sm.m.Unlock()

	c := make(chan bool)
	go func() { sm.set(mux.NewRouter()); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}

func TestSwappableMuxSetWaitsForReadersToReleaseMutex(t *testing.T) {
	sm := SwappableMux{}

	sm.m.RLock()
	defer sm.m.RUnlock()

	c := make(chan bool)
	go func() { sm.set(mux.NewRouter()); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}

func TestServeHTTPCallsInnerMux(t *testing.T) {
	called := false

	m := mux.NewRouter()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { called = true })

	sm := SwappableMux{root: m}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sm.ServeHTTP(w, req)

	if !called {
		t.Error("Inner mux wasn't called")
	}
}

func TestServeHTTPCanServeWhenMuxIsReadLocked(t *testing.T) {
	called := false

	m := mux.NewRouter()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { called = true })

	sm := SwappableMux{root: m}
	sm.m.RLock()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	go sm.ServeHTTP(w, req)

	time.Sleep(10 * time.Millisecond)

	if !called {
		t.Error("Inner mux not called while mutex is read locked")
	}
}

func TestServeHTTPCallsInnerMuxAfterAcquiringLock(t *testing.T) {
	called := false

	m := mux.NewRouter()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { called = true })

	sm := SwappableMux{root: m}
	sm.m.Lock()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	go sm.ServeHTTP(w, req)

	time.Sleep(10 * time.Millisecond)

	if called {
		t.Fatal("Mutex not acquired")
	}

	sm.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if !called {
		t.Error("Inner mux wasn't called after mutex released")
	}
}

func TestUpdateUpdatesMuxWithProvideRouteList(t *testing.T) {
	sm := New()
	rs := []model.Route{
		model.Route{
			Method:     "GET",
			Pattern:    "/",
			Entrypoint: "/bin/sh -c",
			Command:    "jaillover > /tmp/kapow-test-update-mux",
		},
	}
	os.Remove("/tmp/kapow-test-update-mux")
	defer os.Remove("/tmp/kapow-test-update-mux")

	sm.Update(rs)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sm.ServeHTTP(w, req)

	if _, err := os.Stat("/tmp/kapow-test-update-mux"); os.IsNotExist(err) {
		t.Error("Routes not updated")
	} else if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
