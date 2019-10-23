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
package http

import (
	"net/http"
	"strings"
)

// GetReason returns the reason phrase part of an HTTP response
func GetReason(r *http.Response) string {
	if i := strings.IndexByte(r.Status, ' '); i != -1 {
		return r.Status[i+1:]
	}
	return ""
}
