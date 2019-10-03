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

	err := AddRoute("http://localhost", "/hello", "GET", "", "echo Hello World | kapow set /response/body")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if gock.IsDone() == false {
		t.Error("Expected endpoint call not made")
	}
}
