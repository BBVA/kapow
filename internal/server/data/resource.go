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
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strconv"

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
			// Only way to abort current connection as of go 1.13
			// https://github.com/golang/go/issues/16542
			panic("Truncated body")
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

func getRequestPath(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	w.Header().Add("Content-Type", "application/octet-stream")
	// TODO: Discuss a how to obtain URL.EscapedPath() instead
	_, _ = w.Write([]byte(h.Request.URL.Path))
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
	if values, ok := h.Request.Header[textproto.CanonicalMIMEHeaderKey(name)]; ok {
		_, _ = w.Write([]byte(values[0]))
	} else {
		httperror.ErrorJSON(w, ResourceItemNotFound, http.StatusNotFound)
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

// FIXME: Allow any  HTTP status code. Now we are limited by WriteHeader
// capabilities
func setResponseStatus(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	sb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if si, err := strconv.Atoi(string(sb)); err != nil {
		httperror.ErrorJSON(w, NonIntegerValue, http.StatusUnprocessableEntity)
	} else if http.StatusText(si) == "" {
		httperror.ErrorJSON(w, InvalidStatusCode, http.StatusBadRequest)
	} else {
		h.Writer.WriteHeader(int(si))
	}
}

func setResponseHeaders(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	name := mux.Vars(r)["name"]
	vb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hds := h.Writer.Header()
	hds[name] = append(hds[name], string(vb))
}

func setResponseCookies(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	name := mux.Vars(r)["name"]
	vb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{Name: name, Value: string(vb)}
	http.SetCookie(h.Writer, c)
}

func setResponseBody(w http.ResponseWriter, r *http.Request, h *model.Handler) {
	if n, err := io.Copy(h.Writer, r.Body); err != nil {
		if n > 0 {
			panic("Truncated body")
		}
		httperror.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
