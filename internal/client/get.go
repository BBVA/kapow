package client

import (
	"io"

	"github.com/BBVA/kapow/internal/http"
)

//GetData will perform the request and write the results on the provided writer
func GetData(host, id, path string, wr io.Writer) error {
	url := host + "/handlers/" + id + path
	return http.Get(url, "", nil, wr)
}
