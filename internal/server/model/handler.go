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

package model

import (
	"net/http"
	"sync"
)

// Handler represents an open HTTP connection in the User Server.
//
// This struct contains the connection Writer and Request to be managed
// by endpoints of the Data Server.
type Handler struct {
	// ID is unique identifier of the request.
	ID string

	// Route is the original route that matched this request.
	Route

	// Writing is a mutex that prevents two goroutines from writing at
	// the same time in the response.
	Writing sync.Mutex

	// Request is a pointer to the in-progress request.
	Request *http.Request

	// Writer is the original http.ResponseWriter of the request.
	Writer http.ResponseWriter

	// Status is the returned status code
	Status int

	// SentBytes is the number of sent bytes
	SentBytes int64
}
