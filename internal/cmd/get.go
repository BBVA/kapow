package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/client"
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
			log.Fatal(err)
		}
	},
}

func init() {
	GetCmd.Flags().String("data-url", getEnv("KAPOW_DATA_URL", "http://localhost:8082"), "Kapow! data interface URL")
	GetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler ID")
}
