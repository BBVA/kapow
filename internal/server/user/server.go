package user

import (
	"net/http"

	"github.com/BBVA/kapow/internal/server/user/mux"
)

var Server = http.Server{
	Handler: mux.New(),
}

func Run() {
}
