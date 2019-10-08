package control

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ControlServer struct {
	bindAddr     string
	ctrlMux      *mux.Router
	traceChannel chan string
	useTLS       bool
	certfile     string
	keyfile      string
}

var server *ControlServer

func NewControlServer(bindAddr string, listenPort int, certfile, keyfile string) (*ControlServer, error) {

	if server == nil {
		var useTLS bool

		if certfile == "" && keyfile != "" {
			return nil, errors.New("No keyfile provided")
		} else if certfile != "" && keyfile == "" {
			return nil, errors.New("No certfile provided")
		} else if certfile != "" && keyfile != "" {
			useTLS = true
		} else {
			useTLS = false
		}

		server = &ControlServer{
			bindAddr: fmt.Sprintf("%s:%d", bindAddr, listenPort),
			useTLS:   useTLS,
			certfile: certfile,
			keyfile:  keyfile,
			ctrlMux:  mux.NewRouter(),
		}

		server.ctrlMux.HandleFunc("/routes/{id}", server.removeRoute).Methods("DELETE")
		server.ctrlMux.HandleFunc("/routes", server.listRoutes).Methods("GET")
		server.ctrlMux.HandleFunc("/routes", server.addRoute).Methods("POST")
	}

	return server, nil
}

func (srv *ControlServer) Start(traceChannel chan string) {

	srv.traceChannel = traceChannel

	// Start the server
	var err error
	if srv.useTLS {
		err = http.ListenAndServeTLS(srv.bindAddr, srv.certfile, srv.keyfile, srv.ctrlMux)
	} else {
		err = http.ListenAndServe(srv.bindAddr, srv.ctrlMux)
	}

	srv.traceChannel <- err.Error()
}

func (srv *ControlServer) removeRoute(http.ResponseWriter, *http.Request) {
	var logRecord string
	defer func() { srv.traceChannel <- logRecord }()

}

func (srv *ControlServer) listRoutes(http.ResponseWriter, *http.Request) {
	var logRecord string
	defer func() { srv.traceChannel <- logRecord }()

}

func (srv *ControlServer) addRoute(http.ResponseWriter, *http.Request) {
	var logRecord string
	defer func() { srv.traceChannel <- logRecord }()

}
