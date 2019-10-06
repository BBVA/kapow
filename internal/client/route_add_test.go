package client

import (
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestSuccessOnCorrectRoute(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost").
		Put("/routes").
		MatchType("json").
		JSON(map[string]string{
			"method":      "GET",
			"url_pattern": "/hello",
			"entrypoint":  "",
			"command":     "echo Hello World | kapow set /response/body",
		}).
		Reply(http.StatusCreated).
		JSON(map[string]string{})

	// TODO: As per the spec¹, the call should return a JSON body with the info
	// of the newly created route.  Should we consider this in the mocked server
	// and in the test?
	// ¹: https://github.com/BBVA/kapow/tree/master/spec#insert-a-route

	err := AddRoute(
		"http://localhost",
		"/hello", "GET", "", "echo Hello World | kapow set /response/body", nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !gock.IsDone() {
		t.Error("Expected endpoint call not made")
	}
}
