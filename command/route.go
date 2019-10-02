package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

//RouteCmd is the command line interface for kapow route manipulation
var RouteCmd = &cobra.Command{
	Use: "route [action]",
}

func init() {
	var routeListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: "List the current Kapow! routes",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			fmt.Println("niano: ", url)
		},
	}
	routeListCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")

	var routeAddCmd = &cobra.Command{
		Use:   "add [flags] url_pattern [command_file]",
		Short: "Add a route",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			fmt.Println("niano: ", url)
		},
	}
	routeAddCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")
	routeAddCmd.Flags().StringP("command", "c", "", "Command to pass to the shell")
	routeAddCmd.Flags().StringP("entrypoint", "e", "", "Command to execute")
	routeAddCmd.Flags().StringP("method", "X", "", "HTTP method to accept")

	var routeRemoveCmd = &cobra.Command{
		Use:   "remove [flags] route_id",
		Short: "Remove the given route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			fmt.Println("niano: ", url)
		},
	}
	routeRemoveCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")

	RouteCmd.AddCommand(routeListCmd)
	RouteCmd.AddCommand(routeAddCmd)
	RouteCmd.AddCommand(routeRemoveCmd)
}
