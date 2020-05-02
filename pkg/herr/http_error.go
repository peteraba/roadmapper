package herr

import (
	"errors"
	"fmt"
)

// HttpError represents an error that indicates an HTTP status
type HttpError struct {
	error
	status int
}

func (he *HttpError) Unwrap() error {
	return he.error
}

func NewFromError(err error, status int) HttpError {
	switch e := err.(type) {
	case nil:
		panic(fmt.Errorf("new HTTP error from nil"))
	case HttpError:
		return e
	default:
		return HttpError{
			error:  err,
			status: status,
		}
	}
}

func NewFromString(msg string, status int) HttpError {
	return HttpError{
		error:  fmt.Errorf(msg),
		status: status,
	}
}

// ToHttpCode returns a status code indicated by an HttpError or a default status if the error is not an HttpError instance
func ToHttpCode(err error, defaultStatus int) int {
	he := &HttpError{}
	if errors.As(err, he) {
		return he.status
	}

	return defaultStatus
}
