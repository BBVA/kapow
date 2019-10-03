package client

import (
	"fmt"
	"io"
	"net/http"
)

const (
	errMandatoryParam   = "Mandatory parameter %s missing"
	errInvalidURL       = "kapowURL, handlerId or path has invalid format"
	errNotFound         = "Resource Item Not Found"
	errNotValidResource = "Invalid Resource Path"
	serverURLTemplate   = "%s/%s%s"
)

func SetData(kapowURL, handlerId, path string, r io.Reader) error {

	req, err := http.NewRequest("PUT", fmt.Sprintf(serverURLTemplate, kapowURL, handlerId, path), r)
	if err != nil {
		return err
	}

	kpowClient := &http.Client{}
	if resp, err := kpowClient.Do(req); err != nil {
		return err
	} else if resp.StatusCode == http.StatusNoContent {
		return fmt.Errorf(errNotFound)
	} else if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf(errNotValidResource)
	} else if resp.StatusCode >= http.StatusNotFound {
		return fmt.Errorf(resp.Status[4:])
	}

	return nil
}
