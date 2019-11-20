package srverrors

import (
	"encoding/json"
	"net/http"
)

// A ServerErrMessage represents the reason why the error happened
type ServerErrMessage struct {
	Reason string `json:"reason"`
}

// WriteErrorResponse writes the error JSON body to the provided http.ResponseWriter,
// after setting the appropiate Content-Type header
func WriteErrorResponse(statusCode int, reasonMsg string, res http.ResponseWriter) {
	respBody := ServerErrMessage{}
	respBody.Reason = reasonMsg
	bb, _ := json.Marshal(respBody)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	_, _ = res.Write(bb)
}
