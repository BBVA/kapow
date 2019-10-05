// +build !race

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
	if !reflect.DeepEqual(Handlers, New()) {
		t.Error("Handlers is not an empty safeHandlerMap")
	}
}

func TestAddAddsANewHandlerToTheMap(t *testing.T) {
	shm := New()

	shm.Add(&model.Handler{Id: "FOO"})

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler not added to the map")
	}
}

func TestAddAdquiresMutexBeforeAdding(t *testing.T) {
	shm := New()

	shm.m.Lock()
	defer shm.m.Unlock()
	go shm.Add(&model.Handler{Id: "FOO"})

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; ok {
		t.Error("Handler added while mutex was adquired")
	}
}

func TestAddAddsHandlerAfterMutexIsReleased(t *testing.T) {
	shm := New()

	shm.m.Lock()
	go shm.Add(&model.Handler{Id: "FOO"})
	shm.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler not added after mutex release")
	}
}

func TestRemoveRemovesAHandlerFromTheMap(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{Id: "FOO"})

	shm.Remove("FOO")

	if _, ok := shm.hs["FOO"]; ok {
		t.Error("Handler not removed from the map")
	}
}

func TestRemoveAdquiresMutexBeforeRemoving(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{Id: "FOO"})

	shm.m.Lock()
	defer shm.m.Unlock()

	go shm.Remove("FOO")

	time.Sleep(10 * time.Millisecond)

	if _, ok := shm.hs["FOO"]; !ok {
		t.Error("Handler was remove while mutex was adquired")
	}
}

func TestRemoveRemovesHandlerAfterMutexIsReleased(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{Id: "FOO"})

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
	shm.Add(&model.Handler{Id: "FOO"})

	if _, exists := shm.Get("FOO"); !exists {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetReturnExistingHandler(t *testing.T) {
	shm := New()
	expected := &model.Handler{Id: "FOO"}
	shm.Add(expected)

	if current, _ := shm.Get("FOO"); current != expected {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetWaitsForTheWriterToFinish(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{Id: "FOO"})

	shm.m.Lock()
	defer shm.m.Unlock()

	c := make(chan *model.Handler)
	go func() { h, _ := shm.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Handler readed while mutex was adquired")
	default: // This default prevents the select from being blocking
	}
}

func TestGetNonBlockingReadWithOtherReaders(t *testing.T) {
	shm := New()
	shm.Add(&model.Handler{Id: "FOO"})

	shm.m.RLock()
	defer shm.m.RUnlock()

	c := make(chan *model.Handler)
	go func() { h, _ := shm.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Handler couldn't read while mutex was adquired for read")
	}
}
