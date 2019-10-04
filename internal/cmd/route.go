package cmd

import (
	"fmt"
	"os"

	"github.com/BBVA/kapow/internal/client"

	"github.com/spf13/cobra"
)

// RouteCmd is the command line interface for kapow route handling
var RouteCmd = &cobra.Command{
	Use: "route [action]",
}

func init() {
	var routeListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: "List the current Kapow! routes",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")

			if err := client.ListRoutes(url, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	routeListCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")

	// TODO: Manage args for url_pattern and command_file (2 exact args)
	var routeAddCmd = &cobra.Command{
		Use:   "add [flags] url_pattern [command_file]",
		Short: "Add a route",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			method, _ := cmd.Flags().GetString("method")
			command, _ := cmd.Flags().GetString("command")
			entrypoint, _ := cmd.Flags().GetString("entrypoint")
			urlPattern := args[0]

			// TODO: Read command from parameter, file or stdin

			if err := client.AddRoute(url, urlPattern, method, entrypoint, command, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	// TODO: Add default values for flags and remove path flag
	routeAddCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")
	routeAddCmd.Flags().StringP("method", "X", "get", "HTTP method to accept")
	routeAddCmd.Flags().StringP("entrypoint", "e", "/bin/sh -c", "Command to execute")
	routeAddCmd.Flags().StringP("command", "c", "", "Command to pass to the shell")

	var routeRemoveCmd = &cobra.Command{
		Use:   "remove [flags] route_id",
		Short: "Remove the given route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")

			if err := client.RemoveRoute(url, args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	routeRemoveCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8082"), "Kapow! data interface URL")

	RouteCmd.AddCommand(routeListCmd)
	RouteCmd.AddCommand(routeAddCmd)
	RouteCmd.AddCommand(routeRemoveCmd)
}
