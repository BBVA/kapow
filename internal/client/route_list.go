package client

import (
	"io"

	"github.com/BBVA/kapow/internal/http"
)

// ListRoutes list the routes registered on the kapow! instance
func ListRoutes(host string, w io.Writer) error {
	url := host + "/routes/"
	return http.Get(url, "application/json", nil, w)
}
