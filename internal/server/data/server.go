package data

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Rutas a registrar:
// /handlers/{handler_id}/{resource_path}/request GET
// /handlers/{handler_id}/{resource_path}/response PUT

func configRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/handlers/{handler_id}/{root}/{resource:.*$}", readResource).Methods("GET")
	r.HandleFunc("/handlers/{handler_id}/{root}/{resource:.*$}", updateResource).Methods("PUT")
	return r
}

func Run(bindAddr string) {
	r := configRouter()

	log.Fatal(http.ListenAndServe(bindAddr, r))
}

func readResource(res http.ResponseWriter, req *http.Request) {

}

func updateResource(res http.ResponseWriter, req *http.Request) {

}
