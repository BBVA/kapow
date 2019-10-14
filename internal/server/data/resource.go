package data

import (
	"io"
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
)

func getRequestBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	n, err := io.Copy(w, h.Request.Body)
	if err != nil {
		if n == 0 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// Only way to abort current connection as of go 1.13
			// https://github.com/golang/go/issues/16542
			panic("Truncated body")
		}
	}
}
