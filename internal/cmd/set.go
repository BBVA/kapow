package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/BBVA/kapow/internal/client"

	"github.com/spf13/cobra"
)

//SetCmd is the command line interface for set kapow data operation
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

		if len(args) >= 2 {
			// We have a command line value create a stringReader
			r = strings.NewReader(strings.Join(args, " "))
		} else {
			// Use stdin
			r = os.Stdin
		}

		if err := client.SetData(dataURL, handler, args[0], r); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	SetCmd.Flags().String("data-url", getEnv("KAPOW_DATA_URL", "http://localhost:8082"), "Kapow! data interface URL")
	SetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler ID")
}
