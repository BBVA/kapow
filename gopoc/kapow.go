package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {
	go func() {
		fmt.Println("Listening on port 8080")
		http.ListenAndServe(":8080", &userServerHandler{})
	}()

	http.ListenAndServe(":8081", &controlApiHandler{})

}

type userServerHandler struct {
}

func (m *userServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("date").Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(out)
	}
}

type controlApiHandler struct {
}

func (m *controlApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the control API!"))
}
