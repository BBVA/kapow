// +build !race

package state

import (
	"reflect"
	"testing"
	"time"

	"github.com/BBVA/kapow/internal/server/model"
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

func TestAppendReturnsTheInsertedRoutedWithTheActualIndexWhenEmpty(t *testing.T) {
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
