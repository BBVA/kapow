package client

import (
	"fmt"
	"io"

	"github.com/BBVA/kapow/internal/http"
)

// TODO: Review spec: Data API > Error responses > 204 Resource Item Not Found should not be 2xx
func SetData(kapowURL, handlerId, path string, r io.Reader) error {

	return http.Put(fmt.Sprintf("%s/%s%s", kapowURL, handlerId, path), r, nil)

}
