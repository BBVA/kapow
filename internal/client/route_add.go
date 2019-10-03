package client

import (
	"bytes"
	"encoding/json"
	"github.com/BBVA/kapow/internal/http"
)

// AddRoute will add a new route in kapow
func AddRoute(host, path, method, entrypoint, command string) error {
	url := host + "/routes"
	body, _ := json.Marshal(map[string]string{
		"method":      method,
		"url_pattern": path,
		"entrypoint":  entrypoint,
		"command":     command})
	return http.Put(url, "application/json", bytes.NewReader(body), nil)
}
