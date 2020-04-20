package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// An Output represents the execution context,
// meaning the command line and the environment
type Output struct {
	Cmdline []string          `json:"cmdline"`
	Env     map[string]string `json:"env"`
}

func getEnvMap() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		s := strings.SplitN(e, "=", 2)
		env[s[0]] = s[1]
	}
	return env
}

func main() {
	o := Output{
		Cmdline: os.Args,
		Env:     getEnvMap(),
	}
	res, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("JSON marshal failed %+v", err)
	}
	fmt.Println(string(res))
	if len(os.Args) > 1 && os.Args[1] == "--miserably-fail" {
		fmt.Fprintln(os.Stderr, "jailover miserably failed")
		os.Exit(1)
	}
}
