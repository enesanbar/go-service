package errors

import (
	"bytes"
	"fmt"
)

type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable message.
	Message string

	// Logical operation and nested error.
	// Should be supplied by every layer to construct
	// a call stack that has led to this error
	Op  string
	Err error

	// Data returns an arbitrary data related to error, e.g. validation error
	Data interface{}
}

// Error returns detailed error message for developer to debug
func (e Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

func (e Error) Unwrap() error {
	return e.Err
}
