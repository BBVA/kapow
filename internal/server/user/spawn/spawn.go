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

package spawn

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/google/shlex"

	"github.com/BBVA/kapow/internal/server/model"
)

func Spawn(h *model.Handler, out io.Writer) error {
	if h.Route.Entrypoint == "" {
		return errors.New("Entrypoint cannot be empty")
	}
	args, err := shlex.Split(h.Route.Entrypoint)
	if err != nil {
		return err
	}

	if h.Route.Command != "" {
		args = append(args, h.Route.Command)
	}

	cmd := exec.Command(args[0], args[1:]...)
	if out != nil {
		cmd.Stdout = out
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "KAPOW_HANDLER_ID="+h.ID)

	err = cmd.Run()

	return err
}
