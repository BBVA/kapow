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

// Route contains the data needed to represent a Kapow! user route.
type Route struct {
	// ID is the unique identifier of the Route.
	ID string `json:"id,omitempty"`

	// Method is the HTTP method that will match this Route.
	Method string `json:"method"`

	// Pattern is the gorilla/mux path pattern that will match this
	// Route.
	Pattern string `json:"url_pattern"`

	// Entrypoint is the string that will be executed when the Route
	// match.
	//
	// This string will be split according to the shell parsing rules to
	// be passed as a list to exec.Command.
	Entrypoint string `json:"entrypoint,omitempty"`

	// Command is the last argument to be passed to exec.Command when
	// executing the Entrypoint
	Command string `json:"command"`

	// Index is this route position in the server's routes list.
	// It is an output field, its value is ignored as input.
	Index int `json:"index"`

	Debug bool `json:"debug,omitempty"`
}
