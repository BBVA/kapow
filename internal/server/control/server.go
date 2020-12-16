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
func Run(bindAddr string, wg *sync.WaitGroup) {

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		logger.L.Fatal(err)
	}

	// Signal startup
	logger.L.Printf("ControlServer listening at %s\n", bindAddr)
	wg.Done()

	logger.L.Fatal(http.Serve(listener, configRouter()))
}
