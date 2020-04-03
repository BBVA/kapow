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

package spawn

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/BBVA/kapow/internal/server/model"
)

type Output struct {
	Cmdline []string          `json:"cmdline"`
	Env     map[string]string `json:"env"`
}

func decodeJailLover(out []byte) (jldata Output) {
	err := json.Unmarshal(out, &jldata)
	if err != nil {
		log.Fatal("jaillover output is malformed", err)
	}
	return
}

func locateJailLover() string {
	out, err := exec.Command("which", "jaillover").Output()
	if err != nil {
		log.Fatal("jaillover not found in PATH", err)
	}
	return strings.TrimRight(string(out), "\n")
}

func TestSpawnRetursErrorWhenEntrypointIsBad(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: "/bin/this_executable_is_not_likely_to_exist",
		},
	}

	err := Spawn(h, nil)

	if err == nil {
		t.Error("Bad executable not reported")
	}
}

func TestSpawnReturnsNilWhenEntrypointIsGood(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover(),
		},
	}

	err := Spawn(h, nil)

	if err != nil {
		t.Error("Good executable reported")
	}
}

func TestSpawnWritesToStdout(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover(),
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if jldata.Cmdline[0] != locateJailLover() {
		t.Error("Output does not match jaillover's")
	}
}

func TestSpawnSetsKapowURLEnvVar(t *testing.T) {
	t.Skip("Not neccessary as this variable is now set by server at start up")

	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover(),
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if v, ok := jldata.Env["KAPOW_DATA_URL"]; !ok || v != "http://localhost:8082" {
		t.Error("KAPOW_DATA_URL is not set properly")
	}
}

func TestSpawnSetsKapowHandlerIDEnvVar(t *testing.T) {
	h := &model.Handler{
		ID: "HANDLER_ID_FOO",
		Route: model.Route{
			Entrypoint: locateJailLover(),
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if v, ok := jldata.Env["KAPOW_HANDLER_ID"]; !ok || v != "HANDLER_ID_FOO" {
		t.Error("KAPOW_HANDLER_ID is not set properly")
	}
}

func TestSpawnRunsOKEntrypointsWithAParam(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover() + " -foo",
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if !reflect.DeepEqual(jldata.Cmdline, []string{locateJailLover(), "-foo"}) {
		t.Error("Args not as expected")
	}
}

func TestSpawnRunsOKEntrypointWithArgWithSpace(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover() + ` "foo bar"`,
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if !reflect.DeepEqual(jldata.Cmdline, []string{locateJailLover(), "foo bar"}) {
		t.Error("Args not parsed as expected")
	}
}

func TestSpawnErrorsWhenEntrypointIsInvalidShell(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover() + ` "`,
		},
	}
	out := &bytes.Buffer{}

	err := Spawn(h, out)

	if err == nil {
		t.Error("Invalid args not reported")
	}
}

func TestSpawnRunsOKEntrypointWithMultipleArgs(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover() + " foo bar",
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if !reflect.DeepEqual(jldata.Cmdline, []string{locateJailLover(), "foo", "bar"}) {
		t.Error("Args not parsed as expected")
	}
}

func TestSpawnRunsOKEntrypointAndCommand(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{
			Entrypoint: locateJailLover() + " foo bar",
			Command:    "baz qux",
		},
	}
	out := &bytes.Buffer{}

	_ = Spawn(h, out)

	jldata := decodeJailLover(out.Bytes())
	if !reflect.DeepEqual(jldata.Cmdline, []string{locateJailLover(), "foo", "bar", "baz qux"}) {
		t.Error("Malformed cmdline")
	}
}

func TestSpawnReturnsErrorIfEntrypointNotSet(t *testing.T) {
	h := &model.Handler{
		Route: model.Route{},
	}

	err := Spawn(h, nil)

	if err == nil {
		t.Error("Spawn() did not report entrypoint not set")
	}
}
