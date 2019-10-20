package client

import (
	"io"

	"github.com/BBVA/kapow/internal/http"
)

// ListRoutes queries the kapow! instance for the routes that are registered
func ListRoutes(host string, w io.Writer) error {
	url := host + "/routes"
	return http.Get(url, "", nil, w)
}
