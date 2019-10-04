// +build !race

package data

import (
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
)

func TestNewShouldReturnAnEmptyStruct(t *testing.T) {
	hs := New()

	if len(hs.h) > 0 {
		t.Error("Unexpected member in map")
	}
}

func TestAddAddsANewHandlerToTheMap(t *testing.T) {
	hs := New()

	hs.Add(&model.Handler{Id: "FOO"})

	if _, ok := hs.h["FOO"]; !ok {
		t.Error("Handler not added to the map")
	}
}

func TestAddAdquiresMutexBeforeAdding(t *testing.T) {
	hs := New()

	hs.m.Lock()
	defer hs.m.Unlock()
	go hs.Add(&model.Handler{Id: "FOO"})

	time.Sleep(10 * time.Millisecond)

	if _, ok := hs.h["FOO"]; ok {
		t.Error("Handler added while mutex was adquired")
	}
}

func TestAddAddsHandlerAfterMutexIsReleased(t *testing.T) {
	hs := New()

	hs.m.Lock()
	go hs.Add(&model.Handler{Id: "FOO"})
	hs.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if _, ok := hs.h["FOO"]; !ok {
		t.Error("Handler not added after mutex release")
	}
}

func TestRemoveRemovesAHandlerFromTheMap(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	hs.Remove("FOO")

	if _, ok := hs.h["FOO"]; ok {
		t.Error("Handler not removed from the map")
	}
}

func TestRemoveAdquiresMutexBeforeRemoving(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	hs.m.Lock()
	defer hs.m.Unlock()

	go hs.Remove("FOO")

	time.Sleep(10 * time.Millisecond)

	if _, ok := hs.h["FOO"]; !ok {
		t.Error("Handler was remove while mutex was adquired")
	}
}

func TestRemoveRemovesHandlerAfterMutexIsReleased(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	hs.m.Lock()
	go hs.Remove("FOO")
	hs.m.Unlock()

	time.Sleep(10 * time.Millisecond)

	if _, ok := hs.h["FOO"]; ok {
		t.Error("Handler was not removed after mutex release")
	}
}

func TestGetReturnFalseWhenHandlerDoesNotExist(t *testing.T) {
	hs := New()

	if _, exists := hs.Get("FOO"); exists {
		t.Error("Get should return false when handler does not exist")
	}

}

func TestGetReturnTrueWhenHandlerExists(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	if _, exists := hs.Get("FOO"); !exists {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetReturnExistingHandler(t *testing.T) {
	hs := New()
	expected := &model.Handler{Id: "FOO"}
	hs.Add(expected)

	if current, _ := hs.Get("FOO"); current != expected {
		t.Error("Get should return true when handler do exist")
	}
}

func TestGetWaitsForTheWriterToFinish(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	hs.m.Lock()
	defer hs.m.Unlock()

	c := make(chan *model.Handler)
	go func() { h, _ := hs.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Handler readed while mutex was adquired")
	default: // This default prevents the select from being blocking
	}
}

func TestGetNonBlockingReadWithOtherReaders(t *testing.T) {
	hs := New()
	hs.Add(&model.Handler{Id: "FOO"})

	hs.m.RLock()
	defer hs.m.RUnlock()

	c := make(chan *model.Handler)
	go func() { h, _ := hs.Get("FOO"); c <- h }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default: // This default prevents the select from being blocking
		t.Error("Handler couldn't read while mutex was adquired for read")
	}
}
