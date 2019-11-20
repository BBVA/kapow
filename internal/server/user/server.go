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
	"log"
	"net/http"

	"github.com/BBVA/kapow/internal/server/user/mux"
)

// Server is a singleton that stores the http.Server for the user package
var Server = http.Server{
	Handler: mux.New(),
}

// Run finishes configuring Server and runs ListenAndServe on it
func Run(bindAddr string) {
	Server = http.Server{
		Addr:    bindAddr,
		Handler: mux.New(),
	}
	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("UserServer failed: %s", err)
	}
}
