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
	"net"
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/logger"
)

// Run Starts the control server listening in bindAddr
func Run(bindAddr string, wg *sync.WaitGroup, serverCert, serverKey, clientCert []byte) {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(clientCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on bindAddr with the TLS config
	server := &http.Server{
		Addr:      bindAddr,
		TLSConfig: tlsConfig,
		Handler:   configRouter(),
	}

	// Signal startup
	logger.L.Printf("ControlServer listening at %s\n", bindAddr)
	wg.Done()

	// Listen to HTTPS connections with the server certificate and wait
	logger.L.Fatal(server.ListenAndServeTLS(serverCert, serverKey))
}
