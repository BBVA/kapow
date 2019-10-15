package data

import (
	"io"
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
)

func getRequestMethod(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestHost(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestPath(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestMatches(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestParams(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestForm(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func getRequestFiles(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

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

func setResponseStatus(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func setResponseHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func setResponseCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}

func setResponseBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	// FIXME
}
