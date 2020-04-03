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
	Long: `Start a Kapow server with a client interface, a data interface	and an
	admin interface`,
	PreRunE: validateServerCommandArguments,
	Run: func(cmd *cobra.Command, args []string) {
		var sConf server.ServerConfig = server.ServerConfig{}
		sConf.UserBindAddr, _ = cmd.Flags().GetString("bind")
		sConf.ControlBindAddr, _ = cmd.Flags().GetString("control-bind")
		sConf.DataBindAddr, _ = cmd.Flags().GetString("data-bind")

		sConf.CertFile, _ = cmd.Flags().GetString("certfile")
		sConf.KeyFile, _ = cmd.Flags().GetString("keyfile")

		sConf.ClientAuth, _ = cmd.Flags().GetBool("clientauth")
		sConf.ClientCaFile, _ = cmd.Flags().GetString("clientcafile")

		// Set environment variables KAPOW_DATA_URL and KAPOW_CONTROL_URL only if they aren't set so we don't overwrite user's preferences
		if _, exist := os.LookupEnv("KAPOW_DATA_URL"); !exist {
			os.Setenv("KAPOW_DATA_URL", "http://"+sConf.DataBindAddr)
		}
		if _, exist := os.LookupEnv("KAPOW_CONTROL_URL"); !exist {
			os.Setenv("KAPOW_CONTROL_URL", "http://"+sConf.DataBindAddr)
		}

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
			kapowCMD.Env = os.Environ()

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

	ServerCmd.Flags().Bool("clientauth", false, "Activate client mutual tls authentication")
	ServerCmd.Flags().String("clientcafile", "", "Cert file to validate client certificates")
}

func validateServerCommandArguments(cmd *cobra.Command, args []string) error {
	cert, _ := cmd.Flags().GetString("certfile")
	key, _ := cmd.Flags().GetString("keyfile")
	cliAuth, _ := cmd.Flags().GetBool("clientauth")

	if (cert == "") != (key == "") {
		return errors.New("expected both or neither (certfile and keyfile)")
	}

	if cert == "" {
		// If we don't serve thru https client authentication can't be enabled
		if cliAuth {
			return errors.New("Client authentication can't be active in a non https server")
		}
	}

	return nil
}
