package data

import (
	"errors"
	"net/http"
)

func getRequestMethod(req *http.Request) (string, error) { return req.Method, nil }

func getRequestHost(req *http.Request) (string, error) { return req.Host, nil }

func getRequestPath(req *http.Request) (string, error) { return req.URL.EscapedPath(), nil }

func getRequestHeader(req *http.Request, name string) (string, error) {

	if val, ok := req.Header[name]; ok {
		return val[0], nil
	}
	return "", errors.New("Header not found")
}

func getRequestCookie(req *http.Request, name string) (string, error) {

	if val, err := req.Cookie(name); err != nil {
		return "", err
	} else {
		return val.Value, nil
	}
}

func setResponseStatus(res http.ResponseWriter, value int) { res.WriteHeader(value) }

func setResponseHeader(res http.ResponseWriter, name string, value string) {

	res.Header().Add(name, value)
}

func setResponseCookie(res http.ResponseWriter, name string, value string) {

	http.SetCookie(res, &http.Cookie{Name: name, Value: value})
}
