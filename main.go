package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/internal/cmd"
)

func main() {
	var kapowCmd = &cobra.Command{Use: "kapow [action]"}

	kapowCmd.AddCommand(cmd.ServerCmd)
	kapowCmd.AddCommand(cmd.GetCmd)
	kapowCmd.AddCommand(cmd.SetCmd)
	kapowCmd.AddCommand(cmd.RouteCmd)

	err := kapowCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
