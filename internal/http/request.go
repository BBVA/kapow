package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

//Get perform a request using Request with the GET method
func Get(url string, r io.Reader, w io.Writer) error {
	return Request("GET", url, r, w)
}

//Post perform a request using Request with the POST method
func Post(url string, r io.Reader, w io.Writer) error {
	return Request("POST", url, r, w)
}

//Put perform a request using Request with the PUT method
func Put(url string, r io.Reader, w io.Writer) error {
	return Request("PUT", url, r, w)
}

//Delete perform a request using Request with the DELETE method
func Delete(url string, r io.Reader, w io.Writer) error {
	return Request("DELETE", url, r, w)
}

var devnull = ioutil.Discard

//Request will perform the request to the given url and method sending the
//content of the given reader as the body and writing all the contents
//of the response to the given writer. The reader and writer are
//optional.
func Request(method string, url string, r io.Reader, w io.Writer) error {
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
