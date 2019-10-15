package data

import (
	"log"
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

var DataServer http.Server

func RunServer() {
	router := mux.NewRouter()
	// /request
	router.HandleFunc("/handlers/{handlerId}/request/method", checkResource(getRequestMethod)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/host", checkResource(getRequestHost)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/path", checkResource(getRequestPath)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/matches/{name}", checkResource(getRequestMatches)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/params/{name}", checkResource(getRequestParams)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/headers/{name}", checkResource(getRequestHeaders)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/cookies/{name}", checkResource(getRequestCookies)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/form/{name}", checkResource(getRequestForm)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/files/{name}", checkResource(getRequestFiles)).Methods("GET")
	router.HandleFunc("/handlers/{handlerId}/request/body", checkResource(getRequestBody)).Methods("GET")

	// /response
	router.HandleFunc("/handlers/{handlerId}/response/status", checkResource(lockResponseWriter(setResponseStatus))).Methods("POST")
	router.HandleFunc("/handlers/{handlerId}/response/headers/{name}", checkResource(lockResponseWriter(setResponseHeaders))).Methods("POST")
	router.HandleFunc("/handlers/{handlerId}/response/cookies/{name}", checkResource(lockResponseWriter(setResponseCookies))).Methods("POST")
	router.HandleFunc("/handlers/{handlerId}/response/body", checkResource(lockResponseWriter(setResponseBody))).Methods("POST")

	// Invalid resource paths return 400
	router.HandleFunc("/handlers/{handlerId}/{resourcePath:.*}", invalidResourcePath)

	DataServer = http.Server{
		Addr:    "127.0.0.1:8082",
		Handler: router}
	if err := DataServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("DataServer failed: %s", err)
	}
}

func invalidResourcePath(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

// Check that the requested resource exists prior to call the given
// handler
func checkResource(fn func(http.ResponseWriter, *http.Request, *model.Handler)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		h, ok := Handlers.Get(params["handlerId"])
		if ok {
			fn(w, r, h)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func lockResponseWriter(fn func(http.ResponseWriter, *http.Request, *model.Handler)) func(http.ResponseWriter, *http.Request, *model.Handler) {
	return func(w http.ResponseWriter, r *http.Request, h *model.Handler) {
		h.Writing.Lock()
		defer h.Writing.Unlock()
		fn(w, r, h)
	}
}
