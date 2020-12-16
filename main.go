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

package main

import (
	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/cmd"
	"github.com/BBVA/kapow/internal/logger"
)

func main() {
	var kapowCmd = &cobra.Command{Use: "kapow [action]"}

	kapowCmd.AddCommand(cmd.ServerCmd)
	kapowCmd.AddCommand(cmd.GetCmd)
	kapowCmd.AddCommand(cmd.SetCmd)
	kapowCmd.AddCommand(cmd.RouteCmd)

	err := kapowCmd.Execute()
	if err != nil {
		logger.L.Fatalln(err)
	}
}
