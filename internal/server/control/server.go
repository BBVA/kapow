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

package control

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/certs"
	"github.com/BBVA/kapow/internal/logger"
)

// Run Starts the control server listening in bindAddr
func Run(bindAddr string, wg *sync.WaitGroup, serverCert, clientCert certs.Cert) {

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(clientCert.SignedCertPEMBytes())

	ln, err := net.Listen("tcp", bindAddr)
	if err != nil {
		logger.L.Fatal(err)
	}

	server := &http.Server{
		Addr: bindAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				tls.Certificate{
					Certificate: [][]byte{serverCert.SignedCert},
					PrivateKey:  serverCert.PrivKey,
					Leaf:        serverCert.X509Cert,
				},
			},
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
		},
		Handler: configRouter(),
	}

	// Signal startup
	logger.L.Printf("ControlServer listening at %s\n", bindAddr)
	wg.Done()

	// Listen to HTTPS connections with the server certificate and wait
	logger.L.Fatal(server.ServeTLS(ln, "", ""))
}
