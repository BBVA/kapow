package model

import (
	"net/http"
	"sync"
)

type Handler struct {
	Id      string
	Route   *Route
	Writing sync.Mutex
	Writer  http.ResponseWriter
	Request *http.Request
}
