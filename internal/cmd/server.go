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
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/certs"
	"github.com/BBVA/kapow/internal/logger"
	"github.com/BBVA/kapow/internal/server"
)

func banner() {
	fmt.Fprintln(os.Stderr, `
                                                              %%  %%%%
                                                            %%%   %%%
                                                %%         %%%    %%%
                                       %%%%%%% %%%    %%%  %%%    %%%
                            *%%     %%%%%%%%%%%%%%%  %%%% %%%     %%
                   %%   %%%%%%%%%. %%%     %%%% %%% %%%%%%%%
                 %%%%   %%%   %%% %%%       %%% %%%%%%  %%%%
   %%%   %%%    %%%%%%  %%% %%%%  %%%      %%%%  %%%%   %%%      %%%
   %%%  %%%     %%  %%% %%%%%    %%%%%    %%%%   %%%
   %%% %%%     %% %%%%%%%%%       %%%%%%%%%%
   %%%%%%     %%%    %%%%%%         %%%
   %%% %%%%%  %%     %%%%%%
   %%%   %%%%%%%
   %%%%
    %           If you can script it, you can HTTP it.

	`)
}

// ServerCmd is the command line interface for kapow server
var ServerCmd = &cobra.Command{
	Use:   "server [optional flags] [optional init program(s)]",
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
		sConf.Debug, _ = cmd.Flags().GetBool("debug")

		sConf.ControlServerCert = certs.GenCert("control_server", "localhost")
		sConf.ControlClientCert = certs.GenCert("control_client", "localhost")

		// Set environment variables KAPOW_DATA_URL and KAPOW_CONTROL_URL only if they aren't set so we don't overwrite user's preferences
		if _, exist := os.LookupEnv("KAPOW_DATA_URL"); !exist {
			os.Setenv("KAPOW_DATA_URL", "http://"+sConf.DataBindAddr)
		}
		if _, exist := os.LookupEnv("KAPOW_CONTROL_URL"); !exist {
			os.Setenv("KAPOW_CONTROL_URL", "https://"+sConf.ControlBindAddr)
		}
		banner()

		server.StartServer(sConf)

		controlServerCertPEM := new(bytes.Buffer)
		err := pem.Encode(controlServerCertPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: sConf.ControlServerCert.SignedCert,
		})
		if err != nil {
			logger.L.Fatal(err)
		}

		controlClientCertPEM := new(bytes.Buffer)
		err = pem.Encode(controlClientCertPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: sConf.ControlClientCert.SignedCert,
		})
		if err != nil {
			logger.L.Fatal(err)
		}

		controlClientCertPrivKeyPEM := new(bytes.Buffer)
		err = pem.Encode(controlClientCertPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(sConf.ControlClientCert.PrivKey.(*rsa.PrivateKey)),
		})
		if err != nil {
			logger.L.Fatal(err)
		}

		for _, path := range args {
			go Run(
				path,
				sConf.Debug,
				controlServerCertPEM.String(),
				controlClientCertPEM.String(),
				controlClientCertPrivKeyPEM.String(),
			)
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

	ServerCmd.Flags().Bool("debug", false, "Activate debug mode for script executions to standard output")
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

func Run(
	path string,
	debug bool,
	controlServerCertPEM,
	controlClientCertPEM,
	controlClientCertPrivKeyPEM string,
) {
	logger.L.Printf("Running init program %+q", path)
	cmd := BuildCmd(path)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("KAPOW_CONTROL_SERVER_CERT=%v", controlServerCertPEM))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KAPOW_CONTROL_CLIENT_CERT=%v", controlClientCertPEM))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KAPOW_CONTROL_CLIENT_KEY=%v", controlClientCertPrivKeyPEM))

	var wg sync.WaitGroup
	if debug {
		if stdout, err := cmd.StdoutPipe(); err == nil {
			wg.Add(1)
			go logPipe(path, "stdout", stdout, &wg)
		}
		if stderr, err := cmd.StderrPipe(); err == nil {
			wg.Add(1)
			go logPipe(path, "stderr", stderr, &wg)
		}
	}
	err := cmd.Start()
	if err != nil {
		logger.L.Fatalf("Unable to run init program %+q: %s", path, err)
	}

	wg.Wait()
	err = cmd.Wait()
	if err != nil {
		logger.L.Printf("Init program exited with error: %s", err)
	} else {
		logger.L.Printf("Init program %+q finished OK", path)
	}
}

func logPipe(path, name string, pipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	in := bufio.NewScanner(pipe)

	for in.Scan() {
		logger.L.Printf("%+q (%s): %s", path, name, in.Text())
	}
	if err := in.Err(); err != nil {
		logger.L.Printf("Error reading from %+q’s %s: %s", path, name, err)
	}
}
