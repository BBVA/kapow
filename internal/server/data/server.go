package data

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BBVA/kapow/internal/server/model"
	"github.com/gorilla/mux"
)

var requestOperations map[string]func([]string, *http.Request, http.ResponseWriter) = make(map[string]func([]string, *http.Request, http.ResponseWriter))

var responseOperations map[string]func([]string, io.ReadCloser, http.ResponseWriter, *model.Handler) = make(map[string]func([]string, io.ReadCloser, http.ResponseWriter, *model.Handler))

func init() {
	requestOperations["method"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		val := getRequestMethod(targetReq)
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(val))
	}

	requestOperations["host"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		val := getRequestHost(targetReq)
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(val))
	}

	requestOperations["path"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		val := getRequestPath(targetReq)
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(val))
	}

	requestOperations["matches"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 2 {
			res.WriteHeader(http.StatusBadRequest)
		} else if val, err := getRequestMatch(targetReq, resourceComponents[1]); err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(val))
		}
	}

	requestOperations["params"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 2 {
			res.WriteHeader(http.StatusBadRequest)
		} else if val, err := getRequestParam(targetReq, resourceComponents[1]); err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(val))
		}
	}

	requestOperations["headers"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 2 {
			res.WriteHeader(http.StatusBadRequest)
		} else if val, err := getRequestHeader(targetReq, resourceComponents[1]); err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(val))
		}
	}

	requestOperations["cookies"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 2 {
			res.WriteHeader(http.StatusBadRequest)
		} else if val, err := getRequestCookie(targetReq, resourceComponents[1]); err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(val))
		}
	}

	requestOperations["form"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 2 {
			res.WriteHeader(http.StatusBadRequest)
		} else if val, err := getRequestForm(targetReq, resourceComponents[1]); err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(val))
		}
	}

	requestOperations["files"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		if len(resourceComponents) != 3 {
			res.WriteHeader(http.StatusBadRequest)
		} else if resourceComponents[2] == "filename" {
			if val, err := getRequestFileName(targetReq, resourceComponents[1]); err != nil {
				res.WriteHeader(http.StatusNotFound)
			} else {
				res.WriteHeader(http.StatusOK)
				_, _ = res.Write([]byte(val))
			}
		} else if resourceComponents[2] == "content" {
			if err := copyRequestFile(targetReq, resourceComponents[1], res); err != nil {
				res.WriteHeader(http.StatusNotFound)
			}
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	}

	requestOperations["body"] = func(resourceComponents []string, targetReq *http.Request, res http.ResponseWriter) {
		buf := new(bytes.Buffer)
		if err := copyFromRequestBody(targetReq, buf); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write(buf.Bytes())
		}
	}

	responseOperations["status"] = func(resourceComponents []string, reqBody io.ReadCloser, reqRes http.ResponseWriter, target *model.Handler) {
		if val, err := readValueFromBody(reqBody); err != nil {
			reqRes.WriteHeader(http.StatusInternalServerError)
		} else {
			if status, err := strconv.Atoi(val); err != nil {
				reqRes.WriteHeader(http.StatusBadRequest)
			} else {
				target.Writing.Lock()
				setResponseStatus(target.Writer, status)
				target.Writing.Unlock()
				reqRes.WriteHeader(http.StatusOK)
			}
		}
	}

	responseOperations["headers"] = func(resourceComponents []string, reqBody io.ReadCloser, reqRes http.ResponseWriter, target *model.Handler) {
		if len(resourceComponents) != 2 {
			reqRes.WriteHeader(http.StatusBadRequest)
		} else if val, err := readValueFromBody(reqBody); err != nil {
			reqRes.WriteHeader(http.StatusInternalServerError)
		} else {
			target.Writing.Lock()
			setResponseHeader(target.Writer, resourceComponents[1], val)
			target.Writing.Unlock()
			reqRes.WriteHeader(http.StatusOK)
		}
	}

	responseOperations["cookies"] = func(resourceComponents []string, reqBody io.ReadCloser, reqRes http.ResponseWriter, target *model.Handler) {
		if len(resourceComponents) != 2 {
			reqRes.WriteHeader(http.StatusBadRequest)
		} else if val, err := readValueFromBody(reqBody); err != nil {
			reqRes.WriteHeader(http.StatusInternalServerError)
		} else {
			target.Writing.Lock()
			setResponseCookie(target.Writer, resourceComponents[1], val)
			target.Writing.Unlock()
			reqRes.WriteHeader(http.StatusOK)
		}
	}

	responseOperations["body"] = func(resourceComponents []string, reqBody io.ReadCloser, reqRes http.ResponseWriter, target *model.Handler) {
		target.Writing.Lock()
		defer func() {
			target.Writing.Unlock()
			_ = reqBody.Close()
		}()
		if err := copyToResponseBody(target.Writer, reqBody); err != nil {
			reqRes.WriteHeader(http.StatusInternalServerError)
		} else {
			reqRes.WriteHeader(http.StatusOK)
		}
	}

	responseOperations["stream"] = func(resourceComponents []string, reqBody io.ReadCloser, reqRes http.ResponseWriter, target *model.Handler) {
		target.Writing.Lock()
		defer func() {
			target.Writing.Unlock()
			_ = reqBody.Close()
		}()
		if err := copyToResponseStream(target.Writer, reqBody); err != nil {
			reqRes.WriteHeader(http.StatusInternalServerError)
		} else {
			reqRes.WriteHeader(http.StatusOK)
		}
	}
}

// Run must start the data server in a specific address
func Run(bindAddr string) { log.Fatal(http.ListenAndServe(bindAddr, configRouter())) }

func configRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/handlers/{handler_id}/request/{resource_path:.*$}", readRequestResources).Methods(http.MethodGet)
	r.HandleFunc("/handlers/{handler_id}/response/{resource_path:.*$}", writeResponseResources).Methods(http.MethodPut)
	return r
}

var getHandler func(id string) (*model.Handler, bool) = Handlers.Get

func readRequestResources(res http.ResponseWriter, req *http.Request) {
	rVars := mux.Vars(req)
	handlerId := rVars["handler_id"]

	// Check if we have handler to work with
	handler, ok := getHandler(handlerId)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the resource is valid
	resourcePath := rVars["resource_path"]
	resComp := strings.Split(resourcePath, "/")

	if operation, ok := requestOperations[resComp[0]]; !ok {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		operation(resComp, handler.Request, res)
	}
}

func getRequestMethod(req *http.Request) string { return req.Method }

func getRequestHost(req *http.Request) string { return req.Host }

func getRequestPath(req *http.Request) string { return req.URL.EscapedPath() }

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

func getRequestParam(req *http.Request, name string) (string, error) {

	queryParams := req.URL.Query()
	if val, ok := queryParams[name]; ok {
		return val[0], nil
	} else {
		return "", errors.New("Query string parameter not found")
	}
}

func getRequestForm(req *http.Request, name string) (string, error) {

	// Must work with both POST form and multipart forms
	if val := req.PostFormValue(name); val != "" {
		return val, nil
	} else {
		return "", errors.New("Form field not found")
	}
}

func getRequestFileName(req *http.Request, name string) (string, error) {

	_, fileHeader, err := req.FormFile(name)
	if err != nil {
		return "", errors.New("File not found")
	}

	return fileHeader.Filename, nil
}

func copyRequestFile(req *http.Request, name string, w io.Writer) error {

	file, _, err := req.FormFile(name)
	if err != nil {
		return errors.New("File not found")
	}

	_, err = io.Copy(w, file)
	if err != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func getRequestMatch(req *http.Request, name string) (string, error) {
	vars := mux.Vars(req)

	if val, ok := vars[name]; ok {
		return val, nil
	} else {
		return "", errors.New("Match not found")
	}
}

func copyFromRequestBody(req *http.Request, w io.Writer) error {

	defer req.Body.Close()
	if _, err := io.Copy(w, req.Body); err != nil {
		return err
	}

	return nil
}

func readValueFromBody(body io.ReadCloser) (string, error) {

	if bodyBytes, err := ioutil.ReadAll(body); err != nil {
		return "", nil
	} else {
		_ = body.Close()
		return string(bodyBytes), nil
	}
}

func writeResponseResources(res http.ResponseWriter, req *http.Request) {
	rVars := mux.Vars(req)
	handlerId := rVars["handler_id"]

	// Check if we have handler to work with
	handler, ok := getHandler(handlerId)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the resource is valid
	resourcePath := rVars["resource_path"]
	resComp := strings.Split(resourcePath, "/")

	if operation, ok := responseOperations[resComp[0]]; !ok {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		operation(resComp, req.Body, res, handler)
	}
}

func setResponseStatus(res http.ResponseWriter, value int) { res.WriteHeader(value) }

func setResponseHeader(res http.ResponseWriter, name string, value string) {
	res.Header().Add(name, value)
}

func setResponseCookie(res http.ResponseWriter, name string, value string) {
	http.SetCookie(res, &http.Cookie{Name: name, Value: value})
}

func copyToResponseBody(res http.ResponseWriter, r io.Reader) error {

	if _, err := io.Copy(res, r); err != nil {
		return err
	}

	return nil
}

func copyToResponseStream(res http.ResponseWriter, r io.Reader) error {

	if _, err := io.Copy(res, r); err != nil {
		return err
	}

	return nil
}
