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
)

func AsJSON(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

// Get perform a request using Request with the GET method
func Get(url string, r io.Reader, w io.Writer, client *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("GET", url, r, w, client, reqTuner...)
}

// Post perform a request using Request with the POST method
func Post(url string, r io.Reader, w io.Writer, client *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("POST", url, r, w, client, reqTuner...)
}

// Put perform a request using Request with the PUT method
func Put(url string, r io.Reader, w io.Writer, client *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("PUT", url, r, w, client, reqTuner...)
}

// Delete perform a request using Request with the DELETE method
func Delete(url string, r io.Reader, w io.Writer, client *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("DELETE", url, r, w, client, reqTuner...)
}

var devnull = ioutil.Discard

// Request will perform the request to the given url and method sending the
// content of the given reader as the body and writing all the contents
// of the response to the given writer. The reader and writer are
// optional.
func Request(method string, url string, r io.Reader, w io.Writer, client *http.Client, reqTuners ...func(*http.Request)) error {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return err
	}

	for _, reqTuner := range reqTuners {
		reqTuner(req)
	}

	if client == nil {
		client = new(http.Client)
	}

	res, err := client.Do(req)
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
