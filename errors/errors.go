package errors

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func NotFoundError(err error) bool {
	notFoundErr := "(?i)not Found"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}

func ExtractErrorCode(err error) (int, error) {
	// error returned from machinebox/graphql is of graphql type
	// errorNon200Template = "graphql: server returned a non-200 status code: 503"
	rootError := err
	if strings.Contains(err.Error(), "graphql") {
		errorStringArray := strings.Split(err.Error(), ":")
		if len(errorStringArray) == 3 {
			errCodeString := strings.TrimSpace(errorStringArray[2])
			errCode, err := strconv.ParseUint(errCodeString, 10, 32)
			if err != nil {
				return 0, rootError
			}
			return int(errCode), nil
		}
	}
	return 0, rootError
}

func BuildErrorMessage(err error) error {
	// if it's a Not Found error, we return the actual graphql error.
	if NotFoundError(err) {
		return err
	}
	errCode, err := ExtractErrorCode(err)
	// if we fail to decode the error code, just return the error directly
	if http.StatusText(errCode) == "" {
		return err
	}
	var errString string
	return errors.New(errString)
}
