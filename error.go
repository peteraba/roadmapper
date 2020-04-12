package main

// HttpError represents an error that indicates an HTTP status
type HttpError struct {
	error
	status int
}

// ErrorToHttpCode
func ErrorToHttpCode(err error, defaultStatus int) int {
	switch e := err.(type) {
	case HttpError:
		return e.status
	default:
		return defaultStatus
	}
}
