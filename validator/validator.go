package validator

import (
	"fmt"
	"net/http"
)

type Validator struct {
}

func IsRequestValid(r *http.Request, method, contentType string) (int, error) {
	if r.Method != method {
		return http.StatusMethodNotAllowed, fmt.Errorf("not valid Method")
	}

	if r.Header.Get("Content-Type") != contentType {
		return http.StatusBadRequest, fmt.Errorf("not valid Content-Type")
	}
	return http.StatusOK, nil
}
