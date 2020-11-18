/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mux

import (
	"bufio"
	"log"
	"net/http"
	"os"

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

		stdOutR, stdOutW, err := os.Pipe()
		defer stdOutW.Close()
		if err != nil {
			log.Println(err)
			return
		}
		stdErrR, stdErrW, err := os.Pipe()
		defer stdErrW.Close()
		if err != nil {
			log.Println(err)
			return
		}

		go logStream(h.ID, "stdout", stdOutR)
		go logStream(h.ID, "stderr", stdErrR)

		err = spawner(h, stdOutW, stdErrW)

		if err != nil {
			log.Println(err)
		}

	})
}

func logStream(handlerId string, streamName string, stream *os.File) {
	defer stream.Close()
	execLog := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		execLog.Printf("%s %s: %s", handlerId, streamName, scanner.Text())
	}
}
