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
	"io"
	"log"
	"os"
	"strings"

	"github.com/BBVA/kapow/internal/client"

	"github.com/spf13/cobra"
)

// SetCmd is the command line interface for set kapow data operation
var SetCmd = &cobra.Command{
	Use:     "set [flags] resource [value]",
	Short:   "Set a Kapow! resource value",
	Long:    "Set a Kapow! resource value for the current request",
	Args:    cobra.RangeArgs(1, 2),
	PreRunE: handlerIDRequired,
	Run: func(cmd *cobra.Command, args []string) {
		var r io.Reader
		dataURL, _ := cmd.Flags().GetString("data-url")
		handler, _ := cmd.Flags().GetString("handler")
		path, args := args[0], args[1:]

		if len(args) == 1 {
			r = strings.NewReader(args[0])
		} else {
			r = os.Stdin
		}

		if err := client.SetData(dataURL, handler, path, r); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	SetCmd.Flags().String("data-url", getEnv("KAPOW_DATA_URL", "http://localhost:8082"), "Kapow! data interface URL")
	SetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler ID")
}
