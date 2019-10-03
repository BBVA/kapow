package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/client"
)

// GetCmd is the command line interface for get kapow data operation
var GetCmd = &cobra.Command{
	Use:     "get [flags] resource",
	Short:   "Retrive a Kapow! resource",
	Long:    "Retrive a Kapow! resource for the current request",
	Args:    cobra.ExactArgs(1),
	PreRunE: handlerIDRequired,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		handler, _ := cmd.Flags().GetString("handler")

		err := client.GetData(url, handler, args[0], os.Stdout)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
			os.Exit(1)
		}
	},
}

func init() {
	GetCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")
	GetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler id")
}
