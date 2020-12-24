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

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/BBVA/kapow/internal/client"
	"github.com/BBVA/kapow/internal/logger"

	"github.com/spf13/cobra"
)

// RouteCmd is the command line interface for kapow route handling
var RouteCmd = &cobra.Command{
	Use: "route [action]",
}

func init() {
	var routeListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: "List the current Kapow! routes",
		Run: func(cmd *cobra.Command, args []string) {
			controlURL, _ := cmd.Flags().GetString("control-url")

			if err := client.ListRoutes(controlURL, os.Stdout); err != nil {
				logger.L.Fatal(err)
			}
		},
	}
	routeListCmd.Flags().String("control-url", getEnv("KAPOW_CONTROL_URL", "http://localhost:8081"), "Kapow! control interface URL")

	// TODO: Manage args for url_pattern and command_file (2 exact args)
	var routeAddCmd = &cobra.Command{
		Use:   "add [flags] url_pattern [command_file]",
		Short: "Add a route",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			controlURL, _ := cmd.Flags().GetString("control-url")
			method, _ := cmd.Flags().GetString("method")
			command, _ := cmd.Flags().GetString("command")
			entrypoint, _ := cmd.Flags().GetString("entrypoint")
			urlPattern := args[0]

			if len(args) > 1 && command == "" {
				commandFile := args[1]
				var buf []byte
				var err error
				if commandFile == "-" {
					buf, err = ioutil.ReadAll(os.Stdin)
				} else {
					buf, err = ioutil.ReadFile(commandFile)
				}
				if err != nil {
					logger.L.Fatal(err)
				}
				command = string(buf)
			}

			if err := client.AddRoute(controlURL, urlPattern, method, entrypoint, command, os.Stdout); err != nil {
				logger.L.Fatal(err)
			}
		},
	}
	// TODO: Add default values for flags and remove path flag
	routeAddCmd.Flags().String("control-url", getEnv("KAPOW_CONTROL_URL", "http://localhost:8081"), "Kapow! control interface URL")
	routeAddCmd.Flags().StringP("method", "X", "GET", "HTTP method to accept")
	routeAddCmd.Flags().StringP("entrypoint", "e", "", "Command to execute")
	routeAddCmd.Flags().StringP("command", "c", "", "Command to pass to the shell")

	var routeRemoveCmd = &cobra.Command{
		Use:   "remove [flags] route_id",
		Short: "Remove the given route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			controlURL, _ := cmd.Flags().GetString("control-url")

			if err := client.RemoveRoute(controlURL, args[0]); err != nil {
				logger.L.Fatal(err)
			}
		},
	}
	routeRemoveCmd.Flags().String("control-url", getEnv("KAPOW_CONTROL_URL", "http://localhost:8081"), "Kapow! control interface URL")

	RouteCmd.AddCommand(routeListCmd)
	RouteCmd.AddCommand(routeAddCmd)
	RouteCmd.AddCommand(routeRemoveCmd)
}
