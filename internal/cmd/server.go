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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/server"
)

// ServerCmd is the command line interface for kapow server
var ServerCmd = &cobra.Command{
	Use:   "server [optional flags] [optional pow file(s)]",
	Short: "Start a kapow server",
	Long: `Start a Kapow server with, by default with client interface, data interface
	and admin interface`,
	PreRunE: validateServerCommandArguments,
	Run: func(cmd *cobra.Command, args []string) {
		var sConf server.ServerConfig = server.ServerConfig{}
		sConf.UserBindAddr, _ = cmd.Flags().GetString("bind")
		sConf.ControlBindAddr, _ = cmd.Flags().GetString("control-bind")
		sConf.DataBindAddr, _ = cmd.Flags().GetString("data-bind")

		sConf.CertFile, _ = cmd.Flags().GetString("certfile")
		sConf.KeyFile, _ = cmd.Flags().GetString("keyfile")

		go server.StartServer(sConf)

		// start sub shell + ENV(KAPOW_CONTROL_URL)
		if len(args) > 0 {
			powfile := args[0]
			_, err := os.Stat(powfile)
			if os.IsNotExist(err) {
				log.Fatalf("%s does not exist", powfile)
			}
			log.Printf("Running powfile: %q\n", powfile)
			kapowCMD := exec.Command("bash", powfile)
			kapowCMD.Stdout = os.Stdout
			kapowCMD.Stderr = os.Stderr
			kapowCMD.Env = append(os.Environ(), "KAPOW_CONTROL_URL=http://"+sConf.ControlBindAddr)

			err = kapowCMD.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println()
			log.Printf("Done running powfile: %q\n", powfile)
		}

		select {}
	},
}

func init() {
	ServerCmd.Flags().String("bind", "0.0.0.0:8080", "IP address and port to bind the user interface to")
	ServerCmd.Flags().String("control-bind", "localhost:8081", "IP address and port to bind the control interface to")
	ServerCmd.Flags().String("data-bind", "localhost:8082", "IP address and port to bind the data interface to")

	ServerCmd.Flags().String("certfile", "", "Cert file to serve thru https")
	ServerCmd.Flags().String("keyfile", "", "Key file to serve thru https")
}

func validateServerCommandArguments(cmd *cobra.Command, args []string) error {
	cert, _ := cmd.Flags().GetString("certfile")
	key, _ := cmd.Flags().GetString("keyfile")
	if (cert == "") != (key == "") {
		return errors.New("expected both or neither (certfile and keyfile)")
	}
	return nil
}
