/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Get perform a request using Request with the GET method
func Get(url string, contentType string, r io.Reader, w io.Writer) error {
	return Request("GET", url, contentType, r, w)
}

// Post perform a request using Request with the POST method
func Post(url string, contentType string, r io.Reader, w io.Writer) error {
	return Request("POST", url, contentType, r, w)
}

// Put perform a request using Request with the PUT method
func Put(url string, contentType string, r io.Reader, w io.Writer) error {
	return Request("PUT", url, contentType, r, w)
}

// Delete perform a request using Request with the DELETE method
func Delete(url string, contentType string, r io.Reader, w io.Writer) error {
	return Request("DELETE", url, contentType, r, w)
}

var devnull = ioutil.Discard

// Request will perform the request to the given url and method sending the
// content of the given reader as the body and writing all the contents
// of the response to the given writer. The reader and writer are
// optional.
func Request(method string, url string, contentType string, r io.Reader, w io.Writer) error {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return err
	}

	req.Header.Add("X-Kapow-Token", os.Getenv("KAPOW_CONTROL_TOKEN"))

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		reason, err := Reason(res)
		if err != nil {
			return err
		}
		return errors.New(reason)
	}

	if w == nil {
		_, err = io.Copy(devnull, res.Body)
	} else {
		_, err = io.Copy(w, res.Body)
	}

	return err
}
