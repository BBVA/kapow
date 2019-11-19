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

package user

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user/mux"
)

func TestNewReturnAnEmptyStruct(t *testing.T) {
	srl := New()

	if len(srl.rs) != 0 {
		t.Error("Unexpected member in slice")
	}
}

func TestAppendAppendsANewRouteToTheList(t *testing.T) {
	srl := New()

	srl.Append(model.Route{})

	if len(srl.rs) == 0 {
		t.Error("Route not added to the list")
	}
}

func TestAppendAdquiresMutexBeforeAdding(t *testing.T) {
	srl := New()

	srl.m.Lock()
	defer srl.m.Unlock()
	go srl.Append(model.Route{})

	time.Sleep(10 * time.Millisecond)

	if len(srl.rs) != 0 {
		t.Error("Route added while mutex was acquired")
	}
}

func TestAppendAddsRouteAfterMutexIsReleased(t *testing.T) {
	srl := New()

	srl.m.Lock()
	go srl.Append(model.Route{})
	srl.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if len(srl.rs) != 1 {
		t.Error("Route not added after mutex release")
	}
}

func TestSnapshotReturnTheCurrentListOfRoutes(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	rs := srl.Snapshot()

	if !reflect.DeepEqual(srl.rs, rs) {
		t.Error("Route list returned is not the current one")
	}
}

func TestSnapshotReturnADeepCopyOfTheListWhenIsNil(t *testing.T) {
	srl := New()
	srl.rs = nil

	rs := srl.Snapshot()

	if len(rs) != 0 {
		t.Fatal("Route list copy is not empty")
	}
}

func TestSnapshotReturnADeepCopyOfTheListWhenEmpty(t *testing.T) {
	srl := New()

	rs := srl.Snapshot()

	if len(rs) != 0 {
		t.Fatal("Route list copy is not empty")
	}
}

func TestSnapshotReturnADeepCopyOfTheListWhenNonEmpty(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	rs := srl.Snapshot()

	if &rs == &srl.rs {
		t.Fatal("Route list is not a copy")
	}

	for i := 0; i < len(rs); i++ {
		if &rs[i] == &srl.rs[i] {
			t.Errorf("Route %q is not a copy", i)
		}
	}
}

func TestSnapshotWaitsForTheWriterToFinish(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	srl.m.Lock()
	defer srl.m.Unlock()

	c := make(chan []model.Route)
	go func() { c <- srl.Snapshot() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Route list readed while mutex was acquired")
	default: // This default prevents the select from being blocking
	}
}

func TestSnapshotNonBlockingReadWithOtherReaders(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	srl.m.RLock()
	defer srl.m.RUnlock()

	c := make(chan []model.Route)
	go func() { c <- srl.Snapshot() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Route list couldn't be readed while mutex was acquired for read")
	}
}

func TestAppendReturnsTheInsertedRouted(t *testing.T) {
	srl := New()

	r := srl.Append(model.Route{ID: "FOO"})

	if r.ID != "FOO" {
		t.Errorf(`ID of the returned route is not "FOO", but %q`, r.ID)
	}
}

func TestAppendReturnsTheNumberedRoutesWhenEmpty(t *testing.T) {
	srl := New()

	r := srl.Append(model.Route{})

	if r.Index != 0 {
		t.Errorf("Index of the returned route is not 0, but %d", r.Index)
	}
}

func TestAppendReturnsTheInsertedRoutedWithTheActualIndexWhenPopulated(t *testing.T) {
	srl := New()

	var r model.Route

	for i := 0; i < 42; i++ {
		r = srl.Append(model.Route{})
	}
	if r.Index != 42-1 {
		t.Errorf("Index of the returned route is not the last one, i.e., 41, but %d", r.Index)
	}
}

func TestListReturnsTheSameNumberOfRoutesThanSnapshot(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	if len(srl.List()) != len(srl.Snapshot()) {
		t.Error("The number of routes returned is not correct")
	}
}

func TestListReturnsANumberedListOfRoutes(t *testing.T) {
	srl := New()

	for i := 0; i < 42; i++ {
		srl.Append(model.Route{})
	}

	l := srl.List()

	for i, r := range l {
		if i != r.Index {
			t.Fatalf("Route is correctly numbered. Got %v, expected %v", r.Index, i)
		}
	}
}

func TestDeleteReturnsAnErrorOnEmptyListOfRoutes(t *testing.T) {
	srl := New()

	err := srl.Delete("FOO")

	if err == nil {
		t.Error("Expected error not returned")
	}
}

func TestDeleteReturnsNilWhenTheRouteIsInTheList(t *testing.T) {
	srl := New()
	srl.rs = append(srl.rs, model.Route{ID: "FOO"})

	err := srl.Delete("FOO")

	if err != nil {
		t.Errorf("Nil was expected but an error was returned %q", err)
	}
}

func TestDeleteActuallyRemovesTheElementFromTheList(t *testing.T) {
	srl := New()
	srl.rs = append(srl.rs, model.Route{ID: "FOO"})

	_ = srl.Delete("FOO")

	if len(srl.rs) != 0 {
		t.Error("The route was not removed from the list")
	}
}

func TestDeleteRemovesARouteFromTheMiddleOfTheList(t *testing.T) {
	srl := New()
	srl.rs = append(srl.rs, model.Route{ID: "FOO"})
	srl.rs = append(srl.rs, model.Route{ID: "BAR"})
	srl.rs = append(srl.rs, model.Route{ID: "QUX"})

	_ = srl.Delete("BAR")

	if len(srl.rs) != 2 || srl.rs[0].ID != "FOO" || srl.rs[1].ID != "QUX" {
		t.Error("The route was not properly removed")
	}
}

func TestDeleteWaitsForWriterToFinishWriting(t *testing.T) {
	srl := New()

	srl.m.Lock()
	defer srl.m.Unlock()

	c := make(chan error)
	go func() { c <- srl.Delete("FOO") }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't wait for the writer to finish")
	default:
	}
}

func TestDeleteWaitsForReadersToFinishReading(t *testing.T) {
	srl := New()

	srl.m.RLock()
	defer srl.m.RUnlock()

	c := make(chan error)
	go func() { c <- srl.Delete("FOO") }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't wait for the reader to finish")
	default:
	}
}

func TestPackageHaveASingletonEmptyRouteList(t *testing.T) {
	if !reflect.DeepEqual(Routes, New()) {
		t.Error("Routes is not an empty safeRouteList")
	}
}

func TestAppendUpdatesMuxWithProvideRoute(t *testing.T) {
	Server = http.Server{
		Handler: mux.New(),
	}
	srl := New()
	route := model.Route{
		Method:     "GET",
		Pattern:    "/",
		Entrypoint: "/bin/sh -c",
		Command:    "jaillover > /tmp/kapow-test-append-updates-mux",
	}
	os.Remove("/tmp/kapow-test-append-updates-mux")
	defer os.Remove("/tmp/kapow-test-append-updates-mux")

	srl.Append(route)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Server.Handler.ServeHTTP(w, req)

	if _, err := os.Stat("/tmp/kapow-test-append-updates-mux"); os.IsNotExist(err) {
		t.Error("Routes not updated")
	} else if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeleteUpdatesMuxWithRemainingRoutes(t *testing.T) {
	Server = http.Server{
		Handler: mux.New(),
	}
	srl := New()
	route := srl.Append(
		model.Route{
			Method:     "GET",
			Pattern:    "/",
			Entrypoint: "/bin/sh -c",
			Command:    "jaillover > /tmp/kapow-test-remove-updates-mux",
		},
	)
	os.Remove("/tmp/kapow-test-remove-updates-mux")
	defer os.Remove("/tmp/kapow-test-remove-updates-mux")

	_ = srl.Delete(route.ID)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Server.Handler.ServeHTTP(w, req)

	if _, err := os.Stat("/tmp/kapow-test-remove-updates-mux"); err == nil {
		t.Error("Routes not updated")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGetReturnsAnErrorWhenEmptyList(t *testing.T) {
	srl := New()

	if _, err := srl.Get("FOO"); err == nil {
		t.Error("Expected error not returned")
	}
}

func TestGetReturnsAnErrorWhenRouteNotExists(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	if _, err := srl.Get("BAR"); err == nil {
		t.Error("Expected error not returned")
	}
}

func TestGetReturnsTheRequestedRoute(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	if r, err := srl.Get("FOO"); err != nil {
		t.Errorf("Unexpected error %+v", err)
	} else if r.ID != "FOO" {
		t.Errorf(`Route mismatch. Expected: "FOO". Got %q`, r.ID)
	}
}

func TestGetWaitsForTheWriterToFinish(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	srl.m.Lock()
	defer srl.m.Unlock()

	c := make(chan model.Route)
	go func() { r, _ := srl.Get("FOO"); c <- r }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Route list readed while mutex was acquired")
	default: // This default prevents the select from being blocking
	}
}

func TestGetNonBlockingReadWithOtherReaders(t *testing.T) {
	srl := New()
	srl.Append(model.Route{ID: "FOO"})

	srl.m.RLock()
	defer srl.m.RUnlock()

	c := make(chan model.Route)
	go func() { r, _ := srl.Get("FOO"); c <- r }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Route list couldn't be readed while mutex was acquired for read")
	}
}
