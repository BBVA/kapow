package client

import (
	"github.com/BBVA/kapow/internal/http"
)

// RemoveRoute removes a registered route in Kapow! server
func RemoveRoute(host, id string) error {

	url := host + "/routes/" + id
	return http.Delete(url, "", nil, nil)
}
