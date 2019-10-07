package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func handlerIDRequired(cmd *cobra.Command, args []string) error {
	handler, _ := cmd.Flags().GetString("handler")
	if handler == "" {
		return errors.New("--handler or KAPOW_HANDLER_ID is mandatory")
	}
	return nil
}
