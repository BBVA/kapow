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
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/BBVA/kapow/internal/server/httperror"
	"github.com/gorilla/mux"
)

type routeSpec struct {
	route  string
	method string
	rh     resourceHandler
}

func configRouter(rs []routeSpec) (r *mux.Router) {
	r = mux.NewRouter()
	for _, s := range rs {
		r.HandleFunc(s.route, checkHandler(s.rh)).Methods(s.method)
	}
	r.HandleFunc(
		"/handlers/{handlerID}/{resource:.*}",
		func(w http.ResponseWriter, r *http.Request) {
			httperror.ErrorJSON(w, "Invalid Resource Path", http.StatusBadRequest)
		})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
		httperror.ErrorJSON(w, "Data server: Not found", http.StatusNotFound)
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
		httperror.ErrorJSON(w, "Data server: Method not allowed", http.StatusMethodNotAllowed)
	})

	return r
}

func Run(bindAddr string, wg *sync.WaitGroup) {
	rs := []routeSpec{
		// request
		{"/handlers/{handlerID}/request/method", "GET", getRequestMethod},
		{"/handlers/{handlerID}/request/host", "GET", getRequestHost},
		{"/handlers/{handlerID}/request/version", "GET", getRequestVersion},
		{"/handlers/{handlerID}/request/path", "GET", getRequestPath},
		{"/handlers/{handlerID}/request/remote", "GET", getRequestRemote},
		{"/handlers/{handlerID}/request/matches/{name}", "GET", getRequestMatches},
		{"/handlers/{handlerID}/request/params/{name}", "GET", getRequestParams},
		{"/handlers/{handlerID}/request/headers/{name}", "GET", getRequestHeaders},
		{"/handlers/{handlerID}/request/cookies/{name}", "GET", getRequestCookies},
		{"/handlers/{handlerID}/request/form/{name}", "GET", getRequestForm},
		{"/handlers/{handlerID}/request/files/{name}/filename", "GET", getRequestFileName},
		{"/handlers/{handlerID}/request/files/{name}/content", "GET", getRequestFileContent},
		{"/handlers/{handlerID}/request/body", "GET", getRequestBody},

		// route
		{"/handlers/{handlerID}/route/id", "GET", getRouteId},

		// SSL stuff
		{"/handlers/{handlerID}/ssl/client/i/dn", "GET", getSSLClietnDN},

		// response
		{"/handlers/{handlerID}/response/status", "PUT", lockResponseWriter(setResponseStatus)},
		{"/handlers/{handlerID}/response/headers/{name}", "PUT", lockResponseWriter(setResponseHeaders)},
		{"/handlers/{handlerID}/response/cookies/{name}", "PUT", lockResponseWriter(setResponseCookies)},
		{"/handlers/{handlerID}/response/body", "PUT", lockResponseWriter(setResponseBody)},
		{"/handlers/{handlerID}/response/stream", "PUT", lockResponseWriter(setResponseBody)},
	}

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Signal startup
	log.Printf("DataServer listening at %s\n", bindAddr)
	wg.Done()

	log.Fatal(http.Serve(listener, configRouter(rs)))
}
