package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BBVA/kapow/command"
)

func main() {
	var kapowCmd = &cobra.Command{Use: "kapow [action]"}

	kapowCmd.AddCommand(command.ServerCmd)
	kapowCmd.AddCommand(command.GetCmd)
	kapowCmd.AddCommand(command.SetCmd)
	kapowCmd.AddCommand(command.RouteCmd)

	err := kapowCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
