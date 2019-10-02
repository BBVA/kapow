package client

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

func getReason(r *http.Response) string {
	return strings.Join(strings.Split(r.Status, " ")[1:], " ")
}

//GetData will perform the request and write the results on the provided writer
func GetData(host, id, path string, wr io.Writer) error {
	url := host + "/handlers/" + id + path

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(getReason(res))
	}

	_, err = io.Copy(wr, res.Body)
	return err
}
