package cmd

import "fmt"

type userError struct {
	s string
}

func (ue userError) Error() string {
	return ue.s
}

func newUserError(a ...interface{}) userError {
	return userError{s: fmt.Sprintln(a...)}
}

func isUserError(err error) bool {
	if _, ok := err.(userError); ok {
		return true
	}

	return false
}
