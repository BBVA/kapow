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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BBVA/kapow/internal/server/user/mux"
)

// Server is a singleton that stores the http.Server for the user package
var Server = http.Server{
	Handler: mux.New(),
}

// Run finishes configuring Server and runs ListenAndServe on it
func Run(bindAddr, certFile, keyFile, cliCaFile string, cliAuth bool) {
	Server = http.Server{
		Addr:    bindAddr,
		Handler: mux.New(),
	}

	if (certFile != "") && (keyFile != "") {
		if cliAuth {
			if Server.TLSConfig == nil {
				Server.TLSConfig = &tls.Config{}
			}

			var err error
			Server.TLSConfig.ClientCAs, err = loadCertificatesFromFile(cliCaFile)
			if err != nil {
				log.Printf("UserServer failed to load CA certs: %s\nDefault to system CA store.", err)
			} else {
				CAStore := "System store"
				if Server.TLSConfig.ClientCAs != nil {
					CAStore = cliCaFile
				}
				log.Printf("UserServer using CA certs from %s\n", CAStore)
				Server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}
		}

		if err := Server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
			log.Fatalf("UserServer failed: %s", err)
		}
	} else {
		if err := Server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("UserServer failed: %s", err)
		}
	}
}

func loadCertificatesFromFile(certFile string) (pool *x509.CertPool, err error) {
	if certFile != "" {
		caCerts, err := ioutil.ReadFile(certFile)
		if err == nil {
			pool = x509.NewCertPool()
			pool.AppendCertsFromPEM(caCerts)
		}
	}

	return
}
