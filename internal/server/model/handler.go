package model

import (
	"net/http"
	"sync"
)

// Handler represents an open HTTP connection in the User Server.
//
// This struct contains the connection Writer and Request to be managed
// by endpoints of the Data Server.
type Handler struct {
	// ID is unique identifier of the request.
	ID string

	// Route is a reference to the original route that matched this
	// request.
	Route *Route

	// Writing is a mutex that prevents two goroutines from writing at
	// the same time in the response.
	Writing sync.Mutex

	// Request is a pointer to the in-progress request.
	Request *http.Request

	// Writer is the original http.ResponseWriter of the request.
	Writer http.ResponseWriter
}
