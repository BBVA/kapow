package data

import (
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

// Rutas a registrar:
// /handlers/{handler_id}/{resource_path}/request GET
// /handlers/{handler_id}/{resource_path}/response PUT
//func configRouter() *mux.Router {
//	r := mux.NewRouter()
//
//	r.HandleFunc("/handlers/{handler_id}/response/headers/", updateResource).Methods("PUT")
//	r.HandleFunc("/handlers/{handler_id}/response/headers/{key}", updateResource).Methods("PUT")
//	return r
//}
//
//func Run(bindAddr string) {
//	r := configRouter()
//
//	log.Fatal(http.ListenAndServe(bindAddr, r))
//}
//
//func readResource(res http.ResponseWriter, req *http.Request) {
//
//}

var getHandlerId func(string) (*model.Handler, bool) = Handlers.Get

func updateResource(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	hID := vars["handler_id"]

	if _, ok := getHandlerId(hID); !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if resource := vars["resource"]; resource == "response/headers" {
		res.WriteHeader(http.StatusBadRequest)
	}

	res.WriteHeader(http.StatusOK)
	//
	//if _, ok := vars["key"]; !ok {
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
}
