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
	"os"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/client"
	"github.com/BBVA/kapow/internal/logger"
)

// GetCmd is the command line interface for get kapow data operation
var GetCmd = &cobra.Command{
	Use:     "get [flags] resource",
	Short:   "Retrieve a Kapow! resource",
	Long:    "Retrieve a Kapow! resource for the current request",
	Args:    cobra.ExactArgs(1),
	PreRunE: handlerIDRequired,
	Run: func(cmd *cobra.Command, args []string) {
		dataURL, _ := cmd.Flags().GetString("data-url")
		handler, _ := cmd.Flags().GetString("handler")

		err := client.GetData(dataURL, handler, args[0], os.Stdout)
		if err != nil {
			logger.L.Fatal(err)
		}
	},
}

func init() {
	GetCmd.Flags().String("data-url", getEnv("KAPOW_DATA_URL", "http://localhost:8082"), "Kapow! data interface URL")
	GetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler ID")
}
