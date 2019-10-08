package control

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/user"
)

func Run(bindAddr string) {

	r := mux.NewRouter()

	r.HandleFunc("/routes/{id}", removeRoute).
		Methods("DELETE")
	r.HandleFunc("/routes", listRoutes).
		Methods("GET")
	r.HandleFunc("/routes", addRoute).
		Methods("POST")

	log.Fatal(http.ListenAndServe(bindAddr, r))
}

func removeRoute(http.ResponseWriter, *http.Request) {

}

func listRoutes(http.ResponseWriter, *http.Request) {

	user.Routes.Snapshot()

}

func addRoute(http.ResponseWriter, *http.Request) {

	user.Routes.Append(routeSpec)
}
