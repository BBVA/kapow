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

package user

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/logger"
	"github.com/BBVA/kapow/internal/server/user/mux"
)

// Server is a singleton that stores the http.Server for the user package
var Server = http.Server{
	Handler: mux.New(),
}

var DebugEndpoints bool

// Run finishes configuring Server and runs ListenAndServe on it
func Run(bindAddr string, wg *sync.WaitGroup, certFile, keyFile, cliCaFile string, cliAuth, debug bool) {

	Server = http.Server{
		Addr:    bindAddr,
		Handler: mux.New(),
	}

	if debug {
		Routes.SetDebug()
	}

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		logger.L.Fatal(err)
	}

	if (certFile != "") && (keyFile != "") {
		if cliAuth {
			if Server.TLSConfig == nil {
				Server.TLSConfig = &tls.Config{}
			}

			var err error
			Server.TLSConfig.ClientCAs, err = loadCertificatesFromFile(cliCaFile)
			if err != nil {
				logger.L.Fatalf("UserServer failed to load CA certs: %s\n", err)
			} else {
				CAStore := "System store"
				if Server.TLSConfig.ClientCAs != nil {
					CAStore = cliCaFile
				}
				logger.L.Printf("UserServer using CA certs from %s\n", CAStore)
				Server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}
		}

		// Signal startup
		logger.L.Printf("UserServer listening at %s\n", bindAddr)
		wg.Done()

		logger.L.Fatal(Server.ServeTLS(listener, certFile, keyFile))
	} else {
		// Signal startup
		logger.L.Printf("UserServer listening at %s\n", bindAddr)
		wg.Done()

		logger.L.Fatal(Server.Serve(listener))
	}
}

func loadCertificatesFromFile(certFile string) (pool *x509.CertPool, err error) {
	if certFile != "" {
		var caCerts []byte
		caCerts, err = ioutil.ReadFile(certFile)
		if err == nil {
			pool = x509.NewCertPool()
			if !pool.AppendCertsFromPEM(caCerts) {
				err = fmt.Errorf("Invalid certificate file %s", certFile)
			}
		}
	}

	return
}
