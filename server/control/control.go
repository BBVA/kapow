package control

import (
	"fmt"
	"net/http"
)

type ControlServer struct {
	bindAddr     string
	mux          *http.ServeMux
	traceChannel chan string
	certfile     string
	keyfile      string
}

var server *ControlServer

func NewControlServer(bindAddr string, listenPort int, certfile, keyfile string) *ControlServer {

	if server == nil {
		server = &ControlServer{bindAddr: fmt.Sprintf("%s:%d", bindAddr, listenPort),
			certfile: certfile,
			keyfile:  keyfile}
	}

	return server
}

func (*ControlServer) Start(traceChannel chan string) {

}
