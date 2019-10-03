package http

import (
	"net/http"
	"strings"
)

func GetReason(r *http.Response) string {
	if i := strings.IndexByte(r.Status, ' '); i != -1 {
		return r.Status[i+1:]
	} else {
		return ""
	}
}
