package cae

import (
	"encoding/json"
	"errors"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

// ParseQueryError400 is a method used to parse whether a 404 error message means the resources not found.
// For the CAE service, there are some known 404 error codes:
// + CAE.01500208: application or component does not found.
// + CAE.01500404: environment does not found.
func ParseQueryError400(err error, specErrors []string) error {
	var err400 golangsdk.ErrDefault400
	if errors.As(err, &err400) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(err400.Body, &apiError); jsonErr != nil {
			return err
		}

		errCode, searchErr := jmespath.Search("error_code", apiError)
		if searchErr != nil {
			return err
		}

		for _, v := range specErrors {
			if errCode == v {
				return golangsdk.ErrDefault404{}
			}
		}
	}
	return err
}
