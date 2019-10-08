package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestSwappableMuxGetReturnsTheCurrentMux(t *testing.T) {
	sm := swappableMux{}
	mux := sm.get()
	if !reflect.DeepEqual(mux, sm.root) {
		t.Errorf("Returned mux is not the same %#v", mux)
	}
}

func TestSwappableMuxGetReturnsADifferentInstance(t *testing.T) {
	sm := swappableMux{}
	mux := sm.get()
	if &mux == &sm.root {
		t.Error("Returned mux is the same instance")
	}
}

func TestSwappableMuxGetWaitsForTheMutexToBeReleased(t *testing.T) {
	sm := swappableMux{}

	sm.m.Lock()
	defer sm.m.Unlock()

	c := make(chan mux.Router)
	go func() { c <- sm.get() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}

func TestSwappableMuxGetIsAbleToReadWhileOthersAreReading(t *testing.T) {
	sm := swappableMux{}

	sm.m.RLock()
	defer sm.m.RUnlock()

	c := make(chan mux.Router)
	go func() { c <- sm.get() }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
	default:
		t.Error("The mutex cannot be acquired")
	}
}

func TestSwappableMuxSetSetsTheGivenMux(t *testing.T) {
	sm := swappableMux{}
	mux := mux.Router{
		KeepContext: true,
	}

	sm.set(mux)

	//nolint
	if !sm.root.KeepContext {
		t.Error("mux not set")
	}
}

func TestSwappableMuxSetSetsADifferentInstance(t *testing.T) {
	sm := swappableMux{}
	mux := mux.Router{}

	sm.set(mux)

	if &mux == &sm.root {
		t.Error("Set mux is the same instance")
	}
}

func TestSwappableMuxSetWaitsForWriterToReleaseMutex(t *testing.T) {
	sm := swappableMux{}

	sm.m.Lock()
	defer sm.m.Unlock()

	c := make(chan bool)
	go func() { sm.set(mux.Router{}); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}

func TestSwappableMuxSetWaitsForReadersToReleaseMutex(t *testing.T) {
	sm := swappableMux{}

	sm.m.RLock()
	defer sm.m.RUnlock()

	c := make(chan bool)
	go func() { sm.set(mux.Router{}); c <- true }()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-c:
		t.Error("Didn't acquire the mutex")
	default:
	}
}
