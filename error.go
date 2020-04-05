package main

type HttpError struct {
	error
	status int
}

func ErrorToHttpCode(err error, defaultStatus int) int {
	switch e := err.(type) {
	case HttpError:
		return e.status
	default:
		return defaultStatus
	}
}
