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
	nethttp "net/http"
	"testing"
)

func TestEmptyReasonWhenEmptyString(t *testing.T) {
	r := &nethttp.Response{Status: ""}
	if GetReason(r) != "" {
		t.Errorf("We consider an empty status line to have an empty reason")
	}
}

func TestEmptyReasonWhenOnlyCode(t *testing.T) {
	r := &nethttp.Response{Status: "200"}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestEmptyReasonWhenOnlyCodePlusSpace(t *testing.T) {
	r := &nethttp.Response{Status: "200 "}
	if GetReason(r) != "" {
		t.Errorf("We consider an status line with just the status code to have an empty reason")
	}
}

func TestReasonOfOneWord(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO"}
	if GetReason(r) != "FOO" {
		t.Errorf("Unexpected reason found")
	}
}

func TestReasonOfMultipleWords(t *testing.T) {
	r := &nethttp.Response{Status: "200 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}

func TestBehaveWithOddSizeStatusCode(t *testing.T) {
	r := &nethttp.Response{Status: "2 FOO BAR BAZ"}
	if GetReason(r) != "FOO BAR BAZ" {
		t.Errorf("Unexpected reason found")
	}
}
