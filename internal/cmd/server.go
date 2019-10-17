package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/BBVA/kapow/internal/server/control"
	"github.com/BBVA/kapow/internal/server/user"
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
		//cert, _ := cmd.Flags().GetString("certfile")
		//key, _ := cmd.Flags().GetString("keyfile")
		userURL, _ := cmd.Flags().GetString("url")
		controlURL, _ := cmd.Flags().GetString("control-url")
		dataURL, _ := cmd.Flags().GetString("data-url")

		// FIXME: If is a hostport change the name
		fmt.Println("urls:")
		fmt.Println("\t" + userURL)
		fmt.Println("\t" + controlURL)
		fmt.Println("\t" + dataURL)

		// TODO: run data server
		go func() { log.Fatal(http.ListenAndServe(dataURL, nil)) }()

		// TODO: run control server
		go control.Run(controlURL)

		// TODO: run user server
		go user.Run(userURL)

		// start sub shell + ENV(KAPOW_CONTROL_URL)
		// TODO: process several files... tomorrow
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			log.Fatalf("%s does not exist", args[0])
		}
		kapowCMD := exec.Command("bash", args[0])
		kapowCMD.Stdout = os.Stdout
		kapowCMD.Stderr = os.Stderr
		kapowCMD.Env = append(os.Environ(), "KAPOW_URL=http://"+userURL)
		kapowCMD.Env = append(kapowCMD.Env, "KAPOW_CONTROL_URL=http://"+controlURL)
		kapowCMD.Env = append(kapowCMD.Env, "KAPOW_DATA_URL=http://"+dataURL)

		// run bash -c "[pow files contents]"
		err = kapowCMD.Run()
		if err != nil {
			fmt.Println(err)
		}

		select {}
	},
}

func init() {
	ServerCmd.Flags().String("certfile", "", "Cert file to serve thru https")
	ServerCmd.Flags().String("keyfile", "", "Key file to serve thru https")
	ServerCmd.Flags().String("bind", "", "IP address and port to listen to")
	ServerCmd.Flags().BoolP("interactive", "i", false, "Boot an empty kapow server with a shell")
	ServerCmd.Flags().String("url", getEnv("KAPOW_URL", "http://localhost:8080"), "Kapow! user interface URL")
	ServerCmd.Flags().String("control-url", getEnv("KAPOW_CONTROL_URL", "http://localhost:8081"), "Kapow! control interface URL")
	ServerCmd.Flags().String("data-url", getEnv("KAPOW_DATA_URL", "http://localhost:8082"), "Kapow! data interface URL")
}

func validateServerCommandArguments(cmd *cobra.Command, args []string) error {
	cert, _ := cmd.Flags().GetString("certfile")
	key, _ := cmd.Flags().GetString("keyfile")
	if (cert == "") != (key == "") {
		return errors.New("expected both or neither (certfile and keyfile)")
	}
	return nil
}
