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

package httperror

import (
	"encoding/json"
	"net/http"
)

// A ServerErrMessage represents the reason why the error happened
type ServerErrMessage struct {
	Reason string `json:"reason"`
}

// ErrorJSON writes the provided error as a JSON body to the provided
// http.ResponseWriter, after setting the appropriate Content-Type header
func ErrorJSON(w http.ResponseWriter, error string, code int) {
	body, _ := json.Marshal(
		ServerErrMessage{
			Reason: error,
		},
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}
