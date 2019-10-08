package control

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ControlServer struct {
	bindAddr string
	ctrlMux  *mux.Router
}

var Server *ControlServer

func (srv *ControlServer) Run(bindAddr string) {

	Server = &ControlServer{bindAddr, mux.NewRouter()}

	Server.ctrlMux.HandleFunc("/routes/{id}", Server.removeRoute).Methods("DELETE")
	Server.ctrlMux.HandleFunc("/routes", Server.listRoutes).Methods("GET")
	Server.ctrlMux.HandleFunc("/routes", Server.addRoute).Methods("POST")
}

func (srv *ControlServer) removeRoute(http.ResponseWriter, *http.Request) {

}

func (srv *ControlServer) listRoutes(http.ResponseWriter, *http.Request) {

}

func (srv *ControlServer) addRoute(http.ResponseWriter, *http.Request) {

}
