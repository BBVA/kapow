package user

import (
	"log"
	"net/http"

	"github.com/BBVA/kapow/internal/server/user/mux"
)

var Server = http.Server{
	Handler: mux.New(),
}

func Run(bindAddr string) {
	Server = http.Server{
		Addr:    bindAddr,
		Handler: mux.New(),
	}
	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("UserServer failed: %s", err)
	}
}
