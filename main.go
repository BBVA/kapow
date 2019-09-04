package main

import (
	"fmt"
	b "github.com/BBVA/kapow/pkg/banner"
	"net/http"
	"os/exec"
)

func main() {
	ban := b.Banner("0.1.0")
	fmt.Println(ban)

	go func() {
		http.ListenAndServe(":8080", &userServerHandler{})
	}()
	http.ListenAndServe(":8081", &controlAPIHandler{})

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

type controlAPIHandler struct {
}

func (m *controlAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the control API!"))
}
