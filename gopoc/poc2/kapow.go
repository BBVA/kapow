package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

const defaultMessage = "Hello Mr. %v\n"
const commandErrorMessage = "Error executing command: %v\n"
const commandExecutionMessage = "Command executed: %v\n"

type serverSpec struct {
	id, listenAddr string
	routes         *http.ServeMux
}

var servers = []serverSpec{{id: "commandServer", listenAddr: ":8090"}, {id: "defaultServer", listenAddr: ":8080"}}

func startServer(spec serverSpec) {
	err := http.ListenAndServe(spec.listenAddr, spec.routes)
	if err != nil {
		log.Fatal("Error serving ", err)
	}
}

func main() {

	pipe := make(chan string)

	// Create default route handler
	servers[1].routes = http.NewServeMux()
	servers[1].routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/speech; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(defaultMessage, "Unknown")))

		pipe <- fmt.Sprintf("Received request on server %v", servers[1].id)
		servers[0].routes.HandleFunc("/testNewRoute", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/speech; charset=utf-8")

			cmd := exec.Command("ls", "-la", "../")
			output, err := cmd.CombinedOutput()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf(commandErrorMessage, err)))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf(commandExecutionMessage, string(output))))
			}

			pipe <- fmt.Sprintf("Received request \"/\" on server %v", servers[0].id)
		})
	})

	// Create commands route handler
	servers[0].routes = http.NewServeMux()
	servers[0].routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/speech; charset=utf-8")

		cmd := exec.Command("ls", "-la", "./")
		output, err := cmd.CombinedOutput()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(commandErrorMessage, err)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(commandExecutionMessage, string(output))))
		}

		pipe <- fmt.Sprintf("Received request \"/testNewRoute\" on server %v", servers[0].id)
	})

	for i := 0; i < len(servers); i++ {
		go startServer(servers[i])
	}

	for {
		fmt.Println("Tracing: ", <-pipe)
	}
}
