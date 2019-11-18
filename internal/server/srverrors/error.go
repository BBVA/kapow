package srverrors

import (
	"encoding/json"
	"net/http"
)

type ServerErrMessage struct {
	Reason string
}

func WriteErrorResponse(statusCode int, reasonMsg string, res http.ResponseWriter) {
	respBody := ServerErrMessage{}
	respBody.Reason = reasonMsg
	bb, _ := json.Marshal(respBody)
	res.Header().Add("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	_, _ = res.Write(bb)
}
