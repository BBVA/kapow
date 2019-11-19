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

package data

import (
	"reflect"
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestNewReturnAnEmptyStruct(t *testing.T) {
	shm := New()

	if len(shm.hs) != 0 {
		t.Error("Unexpected member in map")
	}
}

func TestPackageHaveASingletonEmptyHandlersList(t *testing.T) {
	t.Skip("Remove later")
	if !reflect.DeepEqual(Handlers, New()) {
		t.Error("Handlers is not an empty safeHandlerMap")
	}
}

func TestAddAddsANewHandlerToTheMap(t *testing.T) {
	shm := New()

	shm.Add(&model.Handler{ID: "FOO"})

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler not added to the map")
	}
}

func TestAddAdquiresMutexBeforeAdding(t *testing.T) {
	shm := New()

	shm.m.Lock()
	defer shm.m.Unlock()
	go shm.Add(&model.Handler{ID: "FOO"})

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; ok {
		t.Error("Handler added while mutex was acquired")
	}
}

func TestAddAddsHandlerAfterMutexIsReleased(t *testing.T) {
	shm := New()

	shm.m.Lock()
	go shm.Add(&model.Handler{ID: "FOO"})
	shm.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler not added after mutex release")
	}
}

func TestRemoveRemovesAHandlerFromTheMap(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.Remove("FOO")

	if _, ok := shm.hs["FOO"]; ok {
		t.Error("Handler not removed from the map")
	}
}

func TestRemoveAdquiresMutexBeforeRemoving(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.Lock()
	defer shm.m.Unlock()

	go shm.Remove("FOO")

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler was remove while mutex was acquired")
	}
}

func TestRemoveRemovesHandlerAfterMutexIsReleased(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.Lock()
	go shm.Remove("FOO")
	shm.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; ok {
		t.Error("Handler was not removed after mutex release")
	}
}

func TestGetReturnFalseWhenHandlerDoesNotExist(t *testing.T) {
	shm := New()

	if _, exists := shm.Get("FOO"); exists {
		t.Error("Get should return false when handler does not exist")
	}
}

func TestGetReturnTrueWhenHandlerExists(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	if _, exists := shm.Get("FOO"); !exists {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetReturnExistingHandler(t *testing.T) {
	shm := New()
	expected := &model.Handler{ID: "FOO"}
	shm.Add(expected)

	if current, _ := shm.Get("FOO"); current != expected {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetWaitsForTheWriterToFinish(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.Lock()
	defer shm.m.Unlock()

	c := make(chan *model.Handler)
	go func() { h, _ := shm.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Handler read while mutex was acquired")
	default: // This default prevents the select from being blocking
	}
}

func TestGetNonBlockingReadWithOtherReaders(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.RLock()
	defer shm.m.RUnlock()

	c := make(chan *model.Handler)
	go func() { h, _ := shm.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Handler couldn't read while mutex was acquired for read")
	}
}

func TestListIDsReturnsTheListOfHandlerIDs(t *testing.T) {
	shm := New()
	shm.hs["FOO"] = nil
	shm.hs["BAR"] = nil
	shm.hs["BAZ"] = nil

	ids := make(map[string]bool)
	for _, id := range shm.ListIDs() {
		ids[id] = true
	}

	_, okFoo := ids["FOO"]
	_, okBar := ids["BAR"]
	_, okBaz := ids["BAZ"]
	if !okFoo || !okBar || !okBaz {
		t.Error("Some IDs not returned")
	}
}

func TestListIDsWaitsForTheWriterToFinish(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.Lock()
	defer shm.m.Unlock()

	c := make(chan []string)
	go func() { c <- shm.ListIDs() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Handler read while mutex was acquired")
	default: // This default prevents the select from being blocking
	}
}

func TestListIDsNonBlockingReadWithOtherReaders(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{ID: "FOO"})

	shm.m.RLock()
	defer shm.m.RUnlock()

	c := make(chan []string)
	go func() { c <- shm.ListIDs() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Handler couldn't read while mutex was acquired for read")
	}
}
