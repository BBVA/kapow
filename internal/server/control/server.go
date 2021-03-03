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
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net"
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/logger"
)

// Run Starts the control server listening in bindAddr
func Run(bindAddr string, wg *sync.WaitGroup, serverCert *x509.Certificate, serverCertPrivKey crypto.PrivateKey, serverCertBytes, clientCertBytes []byte) {
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(clientCert)

	// // Create the TLS Config with the CA pool and enable Client certificate validation
	// tlsConfig := &tls.Config{
	// 	ClientCAs:  caCertPool,
	// 	ClientAuth: tls.RequireAndVerifyClientCert,
	// }
	// tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on bindAddr with the TLS config

	clientCertPEM := new(bytes.Buffer)
	err := pem.Encode(clientCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertBytes,
	})
	if err != nil {
		logger.L.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(clientCertPEM.Bytes())

	ln, err := net.Listen("tcp", bindAddr)
	if err != nil {
		logger.L.Fatal(err)
	}

	server := &http.Server{
		Addr: bindAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				tls.Certificate{
					Certificate: [][]byte{serverCertBytes},
					PrivateKey:  serverCertPrivKey,
					Leaf:        serverCert,
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

// ListenAndServeTLSKeyPair start a server using in-memory TLS KeyPair
// func ListenAndServeTLSKeyPair(addr string, cert tls.Certificate, handler http.Handler) error {

// 	// as defined in https://github.com/golang/go/blob/c0547476f342665514904cf2581a62135d2366c3/src/net/http/server.go#L3034
// 	if addr == "" {
// 		addr = ":https"
// 	}
// 	// as defined in https://github.com/golang/go/blob/c0547476f342665514904cf2581a62135d2366c3/src/net/http/server.go#L3037
// 	ln, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return err
// 	}
// 	server := &http.Server{
// 		Addr:    addr,
// 		Handler: handler,
// 		TLSConfig: &tls.Config{
// 			// alternatifely we can use GetCertificate func(*ClientHelloInfo) (*Certificate, error)
// 			// for host-dependant certificates (possibly let's encrypt)
// 			Certificates: []tls.Certificate{cert},
// 		},
// 	}
// 	// if TLS config is defined, and no actual key path is provided, ServeTLS keeps the certificate
// 	// https://github.com/golang/go/blob/c0547476f342665514904cf2581a62135d2366c3/src/net/http/server.go#L2832
// 	return server.ServeTLS(ln, "", "")
// }
