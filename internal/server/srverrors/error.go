package srverrors

import (
	"encoding/json"
	"net/http"
)

type ServerErrMessage struct {
	Reason string `json:"reason"`
}

func WriteErrorResponse(statusCode int, reasonMsg string, res http.ResponseWriter) {
	respBody := ServerErrMessage{}
	respBody.Reason = reasonMsg
	bb, _ := json.Marshal(respBody)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	_, _ = res.Write(bb)
}
