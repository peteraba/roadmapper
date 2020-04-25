package herr

import "fmt"

// HttpError represents an error that indicates an HTTP status
type HttpError struct {
	error
	status int
}

func NewHttpError(err error, status int) HttpError {
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

// ToHttpCode returns a status code indicated by an HttpError or a default status if the error is not an HttpError instance
func ToHttpCode(err error, defaultStatus int) int {
	switch e := err.(type) {
	case HttpError:
		return e.status
	default:
		return defaultStatus
	}
}
