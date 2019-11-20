package srverrors

import (
	"encoding/json"
	"net/http"
)

// A ServerErrMessage represents the reason why the error happened
type ServerErrMessage struct {
	Reason string `json:"reason"`
}

// ErrorJSON writes the provided error as a JSON body to the provided
// http.ResponseWriter, after setting the appropriate Content-Type header
func ErrorJSON(w http.ResponseWriter, error string, code int) {
	body, _ := json.Marshal(
		ServerErrMessage{
			Reason: error,
		},
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}
