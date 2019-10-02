package command

import (
	"fmt"

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
		url, _ := cmd.Flags().GetString("url")
		handler, _ := cmd.Flags().GetString("handler")
		fmt.Println("niano: ", url, handler)
	},
}

func init() {
	SetCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")
	SetCmd.Flags().String("handler", getEnv("KAPOW_HANDLER_ID", ""), "Kapow! handler id")
}
