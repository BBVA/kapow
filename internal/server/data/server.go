package data

import (
	"errors"
	"net/http"
)

func getRequestMethod(req *http.Request) (string, error) { return req.Method, nil }

func getRequestHost(req *http.Request) (string, error) { return req.Host, nil }

func getRequestPath(req *http.Request) (string, error) { return req.URL.EscapedPath(), nil }

func setResponseStatus(res http.ResponseWriter, value interface{}) error {

	if val, ok := value.(int); ok {
		res.WriteHeader(val)
		return nil
	}

	return errors.New("Not an integer value")
}
