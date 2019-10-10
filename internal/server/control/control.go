package control

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user"
)

// Run must start the control server in a specific address
func Run(bindAddr string) {
	r := configRouter()

	log.Fatal(http.ListenAndServe(bindAddr, r))
}

func configRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/routes/{id}", removeRoute).
		Methods("DELETE")
	r.HandleFunc("/routes", listRoutes).
		Methods("GET")
	r.HandleFunc("/routes", addRoute).
		Methods("POST")
	return r
}

// user.Routes.Remove() []model.Route
var funcRemove func(id string) error

func removeRoute(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if err := funcRemove(id); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}

// user.Routes.List() []model.Route
var funcList func() []model.Route = user.Routes.List

func listRoutes(res http.ResponseWriter, req *http.Request) {

	list := funcList()

	listBytes, _ := json.Marshal(list)
	res.Header().Set("Content-Type", "application/json")
	_, _ = res.Write(listBytes)
}

// user.Routes.Append(r model.Route) model.Route
var funcAdd func(model.Route) model.Route = user.Routes.Append

func addRoute(res http.ResponseWriter, req *http.Request) {
	var route model.Route

	payload, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(payload, &route)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if route.Method == "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if route.Pattern == "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	created := funcAdd(route)
	createdBytes, _ := json.Marshal(created)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "application/json")
	_, _ = res.Write(createdBytes)
}
