package client

import (
	"github.com/BBVA/kapow/internal/http"
)

func RemoveRoute(host, id string) error {
	url := host + "/routes/" + id
	return http.Delete(url, nil, nil)
}
