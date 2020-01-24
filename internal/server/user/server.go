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
	"log"
	"net/http"

	"github.com/BBVA/kapow/internal/server/config"
	"github.com/BBVA/kapow/internal/server/user/mux"
)

// Server is a singleton that stores the http.Server for the user package
var Server = http.Server{
	Handler: mux.New(),
}

// Run finishes configuring Server and runs ListenAndServe on it
func Run(cfg config.ServerConfig) error {
	Server = http.Server{
		Addr:    cfg.UserBindAddr,
		Handler: mux.New(),
	}

	if (cfg.CertFile != "") && (cfg.KeyFile != "") {
		if cfg.ClientAuth {
			if Server.TLSConfig == nil {
				Server.TLSConfig = &tls.Config{}
			}

			var err error
			Server.TLSConfig.ClientCAs, err = loadCertificatesFromFile(cfg.ClientCaFile)
			if err != nil {
				log.Fatalf("UserServer failed to load CA certs: %s\n", err)
			} else {
				CAStore := "System store"
				if Server.TLSConfig.ClientCAs != nil {
					CAStore = cfg.ClientCaFile
				}
				log.Printf("UserServer using CA certs from %s\n", CAStore)
				Server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}
		}

		return Server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	} else {
		return Server.ListenAndServe()
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
