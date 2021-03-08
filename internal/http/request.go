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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/BBVA/kapow/internal/logger"
)

var ControlClientGenerator = GenControlHTTPSClient

func AsJSON(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

// Get perform a request using Request with the GET method
func Get(url string, r io.Reader, w io.Writer, clientGenerator func() *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("GET", url, r, w, clientGenerator, reqTuner...)
}

// Post perform a request using Request with the POST method
func Post(url string, r io.Reader, w io.Writer, clientGenerator func() *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("POST", url, r, w, clientGenerator, reqTuner...)
}

// Put perform a request using Request with the PUT method
func Put(url string, r io.Reader, w io.Writer, clientGenerator func() *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("PUT", url, r, w, clientGenerator, reqTuner...)
}

// Delete perform a request using Request with the DELETE method
func Delete(url string, r io.Reader, w io.Writer, clientGenerator func() *http.Client, reqTuner ...func(*http.Request)) error {
	return Request("DELETE", url, r, w, clientGenerator, reqTuner...)
}

var devnull = ioutil.Discard

// Request will perform the request to the given url and method sending the
// content of the given reader as the body and writing all the contents
// of the response to the given writer. The reader and writer are
// optional.
func Request(method string, url string, r io.Reader, w io.Writer, clientGenerator func() *http.Client, reqTuners ...func(*http.Request)) error {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return err
	}

	for _, reqTuner := range reqTuners {
		reqTuner(req)
	}

	var client *http.Client
	if clientGenerator == nil {
		client = new(http.Client)
	} else {
		client = clientGenerator()
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

func GenControlHTTPSClient() *http.Client {

	serverCert, exists := os.LookupEnv("KAPOW_CONTROL_SERVER_CERT")
	if !exists {
		logger.L.Fatal("KAPOW_CONTROL_SERVER_CERT not in the environment")
	}

	clientCert, exists := os.LookupEnv("KAPOW_CONTROL_CLIENT_CERT")
	if !exists {
		logger.L.Fatal("KAPOW_CONTROL_CLIENT_CERT not in the environment")
	}

	clientKey, exists := os.LookupEnv("KAPOW_CONTROL_CLIENT_KEY")
	if !exists {
		logger.L.Fatal("KAPOW_CONTROL_CLIENT_KEY not in the environment")
	}

	// Load client cert
	clientTLSCert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		logger.L.Fatal(err)
	}

	// Load Server cert
	serverCertPool := x509.NewCertPool()
	serverCertPool.AppendCertsFromPEM([]byte(serverCert))

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientTLSCert},
		RootCAs:      serverCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// The client is always right!
	return client
}
