package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var devnull io.Writer = ioutil.Discard

//Do will perform the request to the given url and method sending the
//content of the given reader as the body and writing all the contents
//of the response to the given writer. The reader and writer are
//optional.
func Do(method string, url string, r io.Reader, w io.Writer) error {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(GetReason(res))
	}

	if w == nil {
		_, err = io.Copy(devnull, res.Body)
	} else {
		_, err = io.Copy(w, res.Body)
	}
	return err
}
