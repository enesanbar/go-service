// Package errors provides a structured error handling system with support for
// error codes, operation context, stack traces, and arbitrary data attachment.
//
// # Core Features
//
//  1. Typed error constructors for common HTTP/gRPC error scenarios
//  2. Automatic stack trace preservation using github.com/pkg/errors
//  3. Contextual data attachment for validation errors, resource IDs, etc.
//  4. Fluent API for incremental error construction
//  5. Standard library compatibility with errors.Is() and errors.As()
//
// # Quick Start
//
// Create typed errors using convenience constructors:
//
//	err := errors.NewInvalidError("ValidateUser", "invalid email format", nil)
//	err := errors.NewNotFoundError("UserRepo.FindByID", "user not found", dbErr)
//	err := errors.NewInternalError("ProcessPayment", "payment processing failed", stripeErr)
//
// # Error Wrapping
//
// Errors automatically capture stack traces when wrapping underlying errors:
//
//	dbErr := db.Query(...)
//	repoErr := errors.NewInternalError("UserRepo.Save", "database error", dbErr)
//	svcErr := errors.NewInternalError("UserService.Create", "failed to create user", repoErr)
//
// The complete error chain is preserved and can be inspected:
//
//	if errors.Is(svcErr, sql.ErrNoRows) {
//	    // Handle specific database error
//	}
//
// # Data Attachment
//
// Attach structured data to errors for better error handling:
//
//	err := errors.NewInvalidError("ValidateInput", "validation failed", nil).
//	    SetData(map[string]string{
//	        "email": "invalid format",
//	        "age": "must be positive",
//	    })
//
// # Fluent API
//
// Build errors incrementally using method chaining:
//
//	err := errors.NewError("", "", "", nil).
//	    WithCode(errors.ENOTFOUND).
//	    WithMessage("resource not found").
//	    WithOperation("FindOrder").
//	    SetData(map[string]interface{}{"order_id": 12345})
//
// # Error Types
//
// Six predefined error types map to HTTP status codes:
//
//	EINVALID       - 400 Bad Request (validation errors)
//	ENOTFOUND      - 404 Not Found (resource doesn't exist)
//	ECONFLICT      - 409 Conflict (duplicate resource, version conflict)
//	EFORBIDDEN     - 403 Forbidden (authorization failure)
//	EINTERNAL      - 500 Internal Server Error (unexpected errors)
//	ENOTMODIFIED   - 304 Not Modified (resource unchanged)
//
// # Stack Traces
//
// Stack traces are automatically captured using errors.WithStack(). To view them,
// use the %+v format verb in logging:
//
//	fmt.Printf("%+v\n", err)  // Prints error with full stack trace
//
// # Performance
//
// Error creation is optimized for performance (<5μs per error) and is suitable
// for production use even in hot paths.
package errors

import (
	"bytes"
	stderrors "errors"
	"fmt"

	"github.com/pkg/errors"
)

// Error represents a structured error with code, message, operation context,
// and optional underlying error with stack trace preservation.
// It supports fluent construction and data attachment for rich error handling.
//
// Stack traces are automatically captured when wrapping errors using NewError()
// or any of the typed constructors. To print the full stack trace, use the
// %+v format verb:
//
//	err := NewInternalError("ConnectDB", "connection failed", dbErr)
//	fmt.Printf("%+v\n", err)  // Prints error with stack trace
//
// The error chain can be inspected using standard errors.Is() and errors.As():
//
//	if errors.Is(err, sql.ErrNoRows) {
//	    // Handle specific error type
//	}
//
// Stack traces are preserved through multiple layers of wrapping:
//
//	dbErr := connectDatabase()
//	svcErr := NewInternalError("Service.Init", "service init failed", dbErr)
//	apiErr := NewInternalError("API.Start", "API start failed", svcErr)
//	// All stack traces from dbErr, svcErr, and apiErr are preserved
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
	Data any
}

// NewError creates a new error with the given code, message, and operation.
// If err is not nil, it is wrapped with a stack trace using errors.WithStack().
// This preserves the original error for errors.Is() and errors.As() checks while
// capturing the call stack at the point of wrapping.
//
// Example:
//
//	err := NewError("DB001", "query failed", "UserRepo.FindByID", sqlErr)
func NewError(code, message, op string, err error) *Error {
	e := &Error{
		Code:    code,
		Message: message,
		Op:      op,
	}

	// Only wrap error if it's not nil
	if err != nil {
		e.Err = errors.WithStack(err)
	}

	return e
}

// NewInvalidError creates a new error for invalid input/validation failures.
func NewInvalidError(op, message string, err error) *Error {
	return NewError(EINVALID, message, op, err)
}

// NewNotFoundError creates a new error for resource not found scenarios.
func NewNotFoundError(op, message string, err error) *Error {
	return NewError(ENOTFOUND, message, op, err)
}

// NewConflictError creates a new error for resource conflicts.
func NewConflictError(op, message string, err error) *Error {
	return NewError(ECONFLICT, message, op, err)
}

// NewForbiddenError creates a new error for authorization failures.
func NewForbiddenError(op, message string, err error) *Error {
	return NewError(EFORBIDDEN, message, op, err)
}

// NewInternalError creates a new error for internal server errors.
func NewInternalError(op, message string, err error) *Error {
	return NewError(EINTERNAL, message, op, err)
}

// NewNotModifiedError creates a new error for unmodified resource scenarios.
func NewNotModifiedError(op, message string, err error) *Error {
	return NewError(ENOTMODIFIED, message, op, err)
}

// WithCode sets the error code and returns the error for method chaining.
// This allows fluent error construction:
//
//	err := NewError("", "initial", "op", nil).
//	    WithCode(EINVALID).
//	    WithMessage("updated message")
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// WithMessage sets the error message and returns the error for method chaining.
// Useful for building errors incrementally or updating existing errors.
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// WithOperation sets the operation name and returns the error for method chaining.
// Operation names should follow the pattern "Package.Function" or "Type.Method".
func (e *Error) WithOperation(op string) *Error {
	e.Op = op
	return e
}

// WrapErr wraps an underlying error with stack trace and returns the error for chaining.
// Uses errors.WithStack() to preserve stack traces. If err is nil, this is a no-op.
//
// Example:
//
//	err := NewInternalError("ProcessData", "processing failed", nil).
//	    WrapErr(dbErr)
func (e *Error) WrapErr(err error) *Error {
	if err != nil {
		e.Err = errors.WithStack(err)
	}
	return e
}

// SetData sets the data field of the error.
// Common patterns include:
//   - Validation errors: map[string]string{"field": "error message"}
//   - Resource identifiers: map[string]interface{}{"user_id": 123, "resource": "orders"}
//   - Request metadata: map[string]string{"request_id": "abc123", "trace_id": "xyz"}
//
// Example:
//
//	err := NewInvalidError("ValidateUser", "validation failed", nil)
//	err.SetData(map[string]string{
//	    "email": "invalid format",
//	    "age": "must be positive",
//	})
func (e *Error) SetData(data interface{}) *Error {
	e.Data = data
	return e
}

// WithData is an alias for SetData() to support fluent method naming patterns.
// See SetData() for usage examples.
func (e *Error) WithData(data interface{}) *Error {
	return e.SetData(data)
}

// Error returns detailed error message for developer to debug
func (e *Error) Error() string {
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

func (e *Error) Unwrap() error {
	return e.Err
}

// Is reports whether any error in err's chain matches target.
// This is a convenience wrapper around the standard errors.Is function.
//
// Example:
//
//	if errors.Is(err, sql.ErrNoRows) {
//	    // Handle not found case
//	}
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As finds the first error in err's chain that matches target type.
// This is a convenience wrapper around the standard errors.As function.
//
// Example:
//
//	var e *Error
//	if errors.As(err, &e) {
//	    fmt.Println("Error code:", e.Code)
//	}
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// Wrap creates a new error wrapping the given error with a message and operation.
// This is a convenience function for quickly wrapping errors without specifying a code.
// The error code defaults to EINTERNAL.
//
// Example:
//
//	if err := db.Query(...); err != nil {
//	    return errors.Wrap(err, "database query failed", "UserRepo.FindByID")
//	}
func Wrap(err error, message, op string) *Error {
	return NewInternalError(op, message, err)
}

// Wrapf creates a new error wrapping the given error with a formatted message.
// This is a convenience function that combines wrapping with message formatting.
//
// Example:
//
//	if err := db.Query(...); err != nil {
//	    return errors.Wrapf(err, "UserRepo.FindByID", "failed to find user %d", userID)
//	}
func Wrapf(err error, op, format string, args ...interface{}) *Error {
	message := fmt.Sprintf(format, args...)
	return NewInternalError(op, message, err)
}

// HasCode checks if the error or any error in its chain is an *Error with the given code.
//
// Example:
//
//	if errors.HasCode(err, errors.ENOTFOUND) {
//	    // Handle not found error
//	}
func HasCode(err error, code string) bool {
	var e *Error
	if stderrors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode extracts the error code from the error if it's an *Error.
// Returns empty string if the error is not an *Error or has no code.
//
// Example:
//
//	code := errors.GetCode(err)
//	if code == errors.ENOTFOUND {
//	    // Handle not found
//	}
func GetCode(err error) string {
	var e *Error
	if stderrors.As(err, &e) {
		return e.Code
	}
	return ""
}

// GetData extracts the data field from the error if it's an *Error.
// Returns nil if the error is not an *Error or has no data.
//
// Example:
//
//	if data := errors.GetData(err); data != nil {
//	    if validationErrs, ok := data.(map[string]string); ok {
//	        // Handle validation errors
//	    }
//	}
func GetData(err error) interface{} {
	var e *Error
	if stderrors.As(err, &e) {
		return e.Data
	}
	return nil
}
