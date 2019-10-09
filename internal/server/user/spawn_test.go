package user

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
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
	r := &model.Route{
		Entrypoint: "/bin/this_executable_is_not_likely_to_exist",
	}

	h := &model.Handler{
		Route: r,
	}

	err := spawn(h, nil)
	if err == nil {
		t.Error("Bad executable not reported")
	}
}

func TestSpawnReturnsNilWhenEntrypointIsGood(t *testing.T) {
	r := &model.Route{
		Entrypoint: locateJailLover(),
	}

	h := &model.Handler{
		Route: r,
	}

	err := spawn(h, nil)
	if err != nil {
		t.Error("Good executable reported")
	}
}

func TestSpawnWritesToStdout(t *testing.T) {
	r := &model.Route{
		Entrypoint: locateJailLover(),
	}

	h := &model.Handler{
		Route: r,
	}

	out := &bytes.Buffer{}

	_ = spawn(h, out)

	jldata := decodeJailLover(out.Bytes())

	if jldata.Cmdline[0] != locateJailLover() {
		t.Error("Ouput does not match jaillover's")
	}

}

func TestSpawnSetsKapowURLEnvVar(t *testing.T) {
	r := &model.Route{
		Entrypoint: locateJailLover(),
	}

	h := &model.Handler{
		Route: r,
	}

	out := &bytes.Buffer{}

	_ = spawn(h, out)

	jldata := decodeJailLover(out.Bytes())

	if v, ok := jldata.Env["KAPOW_URL"]; !ok || v != "http://localhost:8081" {
		t.Error("KAPOW_URL is not set properly")
	}
}

func TestSpawnSetsKapowHandlerIDEnvVar(t *testing.T) {
	r := &model.Route{
		Entrypoint: locateJailLover(),
	}

	h := &model.Handler{
		ID:    "HANDLER_ID_FOO",
		Route: r,
	}

	out := &bytes.Buffer{}

	_ = spawn(h, out)

	jldata := decodeJailLover(out.Bytes())

	if v, ok := jldata.Env["KAPOW_HANDLER_ID"]; !ok || v != "HANDLER_ID_FOO" {
		t.Error("KAPOW_HANDLER_ID is not set properly")
	}
}
