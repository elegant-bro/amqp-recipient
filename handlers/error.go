package handlers

import "fmt"

type internalError struct {
	origin error
	msg    string
}

func newInternalError(origin error, msg string) *internalError {
	return &internalError{origin: origin, msg: msg}
}

func (i internalError) Error() string {
	return fmt.Sprintf("%s: %v", i.msg, i.origin)
}
