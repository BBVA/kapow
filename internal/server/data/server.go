package data

import (
	"io/ioutil"
	"net/http"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

// Rutas a registrar:
// /handlers/{handlerId}/response/headers/{item} GET|PUT
// /handlers/{handlerId}/response/cookies/{item} GET|PUT
// /handlers/{handlerId}/request/headers/{item} GET|PUT
// /handlers/{handlerId}/request/cookies/{item} GET|PUT

var getHandlerId func(string) (*model.Handler, bool) = Handlers.Get

func updateResource(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	hID := vars["handlerId"]

	hnd, ok := getHandlerId(hID)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	resource := vars["resource"]
	if resource == "response/headers" || resource == "response/cookies" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Body.Close()
	value := string(bodyBytes)

	if hnd != nil {
		hnd.Writing.Lock()
		hnd.Writer.Header().Add("pepe", value)
		hnd.Writing.Unlock()
	}

	res.WriteHeader(http.StatusOK)
}
