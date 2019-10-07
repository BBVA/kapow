package model

// Route contains the data needed to represent a Kapow! user route.
type Route struct {
	// ID is the unique identifier of the Route.
	ID string `json:"id"`

	// Method is the HTTP method that will match this Route.
	Method string `json:"method"`

	// Pattern is the gorilla/mux path pattern that will match this
	// Route.
	Pattern string `json:"url_pattern"`

	// Entrypoint is the string that will be executed when the Route
	// match.
	//
	// This string will be split according to the shell parsing rules to
	// be passed as a list to exec.Command.
	Entrypoint string `json:"entrypoint"`

	// Command is the last argument to be passed to exec.Command when
	// executing the Entrypoint
	Command string `json:"command"`
}
