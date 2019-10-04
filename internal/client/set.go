package client

import (
	"io"

	"github.com/BBVA/kapow/internal/http"
)

func SetData(host, handlerID, path string, r io.Reader) error {

	url := host + "/handlers/" + handlerID + path
	return http.Put(url, "", r, nil)
}
