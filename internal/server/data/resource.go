/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package data

import (
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/BBVA/kapow/internal/logger"
	"github.com/BBVA/kapow/internal/server/httperror"
	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

// Constants for error reasons
const (
	ResourceItemNotFound = "Resource Item Not Found"
	NonIntegerValue      = "Non Integer Value"
	InvalidStatusCode    = "Invalid Status Code"
)

func getRequestBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	n, err := io.Copy(w, h.Request.Body)
	if err != nil {
		if n == 0 {
			httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		} else {
			panic(http.ErrAbortHandler)
		}
	}
}

func getRequestMethod(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.Method))
}

func getRequestHost(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.Host))
}

func getRequestVersion(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.Proto))
}

func getRequestPath(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	// TODO: Discuss a how to obtain URL.EscapedPath() instead
	_, _ = w.Write([]byte(h.Request.URL.Path))
}

func getRequestRemote(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Request.RemoteAddr))
}

func getRequestMatches(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	vars := mux.Vars(h.Request)
	if value, ok := vars[name]; ok {
		_, _ = w.Write([]byte(value))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

func getRequestParams(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	if values, ok := h.Request.URL.Query()[name]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

func getRequestHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	name = textproto.CanonicalMIMEHeaderKey(name)
	if values, ok := h.Request.Header[name]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		if name == "Host" {
			_, _ = w.Write([]byte(h.Request.Host))
		} else {
			httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
		}
	}
}

func getRequestCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	if cookie, err := h.Request.Cookie(name); err == nil {
		_, _ = w.Write([]byte(cookie.Value))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

// NOTE: The current implementation doesn't allow us to decode
// form encoded data sent in a request with an arbitrary method. This is
// needed for Kapow! semantic so it MUST be changed in the future
// FIXME: Implement a ParseForm function that doesn't care about Method
// nor Content-Type
func getRequestForm(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	// FIXME: This SHOULD? return an error when Body is empty and IS NOT
	// We tried to exercise this execution path but didn't know how.
	err := h.Request.ParseForm()
	if err != nil {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	} else if values, ok := h.Request.Form[name]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

func getRequestFileName(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	_, header, err := h.Request.FormFile(name)
	if err == nil {
		_, _ = w.Write([]byte(header.Filename))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

func getRequestFileContent(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	name := mux.Vars(r)["name"]
	file, _, err := h.Request.FormFile(name)
	if err == nil {
		_, _ = io.Copy(w, file)
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	}
}

func getSSLClietnDN(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	if h.Request.TLS == nil {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	} else if h.Request.TLS.VerifiedChains == nil {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
	} else {
		w.Header().Add("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte(h.Request.TLS.VerifiedChains[0][0].Subject.String()))
	}
}

func getRouteId(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	_, _ = w.Write([]byte(h.Route.ID))

}

// FIXME: Allow any  HTTP status code. Now we are limited by WriteHeader
// capabilities
func setResponseStatus(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	sb, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if si, err := strconv.Atoi(string(sb)); err != nil {
		httperror.ErrorJSON(w, NonIntegerValue, http.StatusUnprocessableEntity)
	} else if http.StatusText(si) == "" {
		httperror.ErrorJSON(w, InvalidStatusCode, http.StatusBadRequest)
	} else {
		h.Status = si
	}
}

func setResponseHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	name := mux.Vars(r)["name"]
	vb, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hds := h.Writer.Header()
	hds[name] = append(hds[name], string(vb))
}

func setResponseCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	name := mux.Vars(r)["name"]
	vb, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{Name: name, Value: string(vb)}
	http.SetCookie(h.Writer, c)
}

func setResponseBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {

	if !h.BodyOut {
		if h.Status != 0 {
			h.Writer.WriteHeader(h.Status)
		}
		h.BodyOut = true
	}

	n, err := io.Copy(h.Writer, r.Body)
	if err != nil {
		if n > 0 {
			panic(http.ErrAbortHandler)
		}
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	h.SentBytes += n
}

func setServerLog(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	handlerId := mux.Vars(r)["handlerID"]
	if prefix := mux.Vars(r)["prefix"]; prefix == "" {
		logger.L.Printf("%s %s\n", escapeString(handlerId), msg)
	} else {
		logger.L.Printf("%s %s: %s\n", escapeString(handlerId), escapeString(prefix), msg)
	}
}

// function to scape strings in order to be printed in a Log
func escapeString(s string) string {
  s = strings.Replace(s, "\n", "", -1)
  s = strings.Replace(s, "\r", "", -1)
  s = strings.Replace(s, "\t", "", -1)
  s = strings.Replace(s, "\b", "", -1)

  return s
}
