package user

import (
	"reflect"
	"testing"

	"github.com/BBVA/kapow/internal/server/user/state"
)

func TestPackageHaveASingletonEmptyRouteList(t *testing.T) {
	if !reflect.DeepEqual(Routes, state.New()) {
		t.Error("Routes is not an empty safeRouteList")
	}
}
