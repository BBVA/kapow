package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// ServerCmd is the command line interface for kapow server
var ServerCmd = &cobra.Command{
	Use:   "server [optional flags] [optional pow file(s)]",
	Short: "Start a kapow server",
	Long: `Start a Kapow server with, by default with client interface, data interface
	and admin interface`,
	PreRunE: validateServerCommandArguments,
	Run: func(cmd *cobra.Command, args []string) {
		cert, _ := cmd.Flags().GetString("certfile")
		key, _ := cmd.Flags().GetString("keyfile")
		fmt.Println("waka server feliz :)", cert, key)
	},
}

func init() {
	ServerCmd.Flags().String("certfile", "", "Cert file to serve thru https")
	ServerCmd.Flags().String("keyfile", "", "Key file to serve thru https")
	ServerCmd.Flags().String("bind", "", "IP address and port to listen to")
	ServerCmd.Flags().BoolP("interactive", "i", false, "Boot an empty kapow server with a shell")
}

func validateServerCommandArguments(cmd *cobra.Command, args []string) error {
	cert, _ := cmd.Flags().GetString("certfile")
	key, _ := cmd.Flags().GetString("keyfile")
	if (cert == "") != (key == "") {
		return errors.New("expected both or neither (certfile and keyfile)")
	}
	return nil
}
