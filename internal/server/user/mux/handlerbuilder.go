package mux

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/BBVA/kapow/internal/server/data"
	"github.com/BBVA/kapow/internal/server/model"
	"github.com/BBVA/kapow/internal/server/user/spawn"
)

var spawner = spawn.Spawn
var idGenerator = uuid.NewUUID

func handlerBuilder(route model.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := idGenerator()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h := &model.Handler{
			ID:      id.String(),
			Route:   route,
			Request: r,
			Writer:  w,
		}

		data.Handlers.Add(h)
		defer data.Handlers.Remove(h.ID)

		err = spawner(h, nil)
		if err != nil {
			log.Println(err)
		}
	})
}
