package client

import (
	"errors"
	"net/http"
	"strings"
)

// AddRoute will add a new route in kapow
func AddRoute(host, path, method, entrypoint, command string) error {
	reqData, err := http.NewRequest(
		"PUT",
		host+"/routes",
		strings.NewReader(command),
	)
	if err != nil {
		return err
	}

	var client = new(http.Client)
	res, err := client.Do(reqData)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(res.Status)
	}

	return nil
}
