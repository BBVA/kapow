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
	"io/ioutil"
	nethttp "net/http"
	"strings"
	"testing"
)

func TestReasonExtractsReasonFromJSON(t *testing.T) {
	r := &nethttp.Response{
		Status: "200 OK",
		Body: ioutil.NopCloser(
			strings.NewReader(
				`{"reason": "Because reasons", "foo": "bar"}`,
			),
		),
	}

	reason, _ := Reason(r)

	if reason != "Because reasons" {
		t.Errorf(`reason mismatch, want "Because reasons", got %q`, reason)
	}
}

func TestReasonErrorsOnJSONWithNoReason(t *testing.T) {
	r := &nethttp.Response{
		Status: "200 OK",
		Body: ioutil.NopCloser(
			strings.NewReader(
				`{"madness": "Because madness", "foo": "bar"}`,
			),
		),
	}

	_, err := Reason(r)

	if err == nil {
		t.Error("error not reported")
	}
}

func TestReasonErrorsOnJSONWithEmptyReason(t *testing.T) {
	r := &nethttp.Response{
		Body: ioutil.NopCloser(
			strings.NewReader(
				`{"reason": "", "foo": "bar"}`,
			),
		),
	}

	_, err := Reason(r)

	if err == nil {
		t.Error("error not reported")
	}
}

func TestReasonErrorsOnNoJSON(t *testing.T) {
	r := &nethttp.Response{
		Body: ioutil.NopCloser(
			strings.NewReader(""),
		),
	}

	_, err := Reason(r)

	if err == nil {
		t.Error("error not reported")
	}
}

func TestReasonErrorsOnInvalidJSON(t *testing.T) {
	r := &nethttp.Response{
		Body: ioutil.NopCloser(
			strings.NewReader(
				`{"reason": "Because reasons", "cliffhanger...`,
			),
		),
	}

	_, err := Reason(r)

	if err == nil {
		t.Error("error not reported")
	}
}
