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
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestSuccessOnCorrectRoute(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Post("/routes").
		MatchType("json").
		JSON(map[string]string{
			"method":      "GET",
			"url_pattern": "/hello",
			"entrypoint":  "",
			"command":     "echo Hello World | kapow set /response/body",
		}).
		Reply(http.StatusCreated).
		JSON(map[string]string{})

	err := AddRoute(
		"http://localhost",
		"/hello", "GET", "", "echo Hello World | kapow set /response/body", nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !gock.IsDone() {
		t.Error("Expected endpoint call not made")
	}
}
