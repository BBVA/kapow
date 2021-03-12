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

package client

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/BBVA/kapow/internal/http"
)

// AddRoute will add a new route in kapow
func AddRoute(host, path, method, entrypoint, command string, w io.Writer) error {
	url := host + "/routes"
	payload := map[string]string{
		"method":      method,
		"url_pattern": path,
		"command":     command,
	}
	if entrypoint != "" {
		payload["entrypoint"] = entrypoint
	}
	body, _ := json.Marshal(payload)
	return http.Post(url, bytes.NewReader(body), w, http.ControlClientGenerator, http.AsJSON)
}
