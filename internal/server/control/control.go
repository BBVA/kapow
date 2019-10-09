package control

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/model"
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

// user.Routes.Remove() []model.Route
var funcRemove func(id string) error

func removeRoute(http.ResponseWriter, *http.Request) {

	if err := funcRemove(""); err != nil {
		fmt.Printf("Mostoles, we've had a problem")
	}
}

// user.Routes.Snapshot() []model.Route
var funcSnapshot func() []model.Route = user.Routes.Snapshot

func listRoutes(http.ResponseWriter, *http.Request) {

	funcSnapshot()

}

// user.Routes.Append(r model.Route)
var funcAdd func(model.Route) model.Route = user.Routes.Append

func addRoute(res http.ResponseWriter, req *http.Request) {

	funcAdd(model.Route{})
	res.WriteHeader(http.StatusCreated)
	_, _ = io.Copy(res, req.Body)
}
