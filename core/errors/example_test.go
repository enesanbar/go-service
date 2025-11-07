package errors_test

import (
	stderrors "errors"
	"fmt"

	"github.com/enesanbar/go-service/core/errors"
)

// Example_typedConstructors demonstrates using typed error constructors
func Example_typedConstructors() {
	// Create a validation error
	err := errors.NewInvalidError("ValidateUser", "validation failed", nil).
		SetData(map[string]string{
			"email": "invalid format",
			"age":   "must be positive",
		})

	fmt.Println(err.Code)
	fmt.Println(err.Op)
	// Output:
	// invalid
	// ValidateUser
}

// Example_errorWrapping demonstrates wrapping errors across multiple layers
func Example_errorWrapping() {
	// Simulate a database error
	dbErr := stderrors.New("connection timeout")

	// Repository layer wraps it
	repoErr := errors.NewInternalError("UserRepo.FindByID", "database query failed", dbErr)

	// Service layer wraps the repository error
	svcErr := errors.NewInternalError("UserService.GetUser", "failed to fetch user", repoErr)

	// Print the error chain
	fmt.Println(svcErr.Error())

	// Check if the original error is in the chain
	if stderrors.Is(svcErr, dbErr) {
		fmt.Println("Contains database error")
	}

	// Output:
	// UserService.GetUser: UserRepo.FindByID: connection timeout
	// Contains database error
}

// Example_dataAttachment demonstrates attaching contextual data to errors
func Example_dataAttachment() {
	// Attach validation errors
	err := errors.NewInvalidError("CreateUser", "validation failed", nil).
		SetData(map[string]string{
			"username": "already exists",
			"email":    "invalid format",
		})

	// Extract data for error handling
	if validationData, ok := err.Data.(map[string]string); ok {
		// Print in deterministic order for testing
		if msg, ok := validationData["email"]; ok {
			fmt.Printf("email: %s\n", msg)
		}
		if msg, ok := validationData["username"]; ok {
			fmt.Printf("username: %s\n", msg)
		}
	}

	// Output:
	// email: invalid format
	// username: already exists
}

// Example_fluentAPI demonstrates building errors with method chaining
func Example_fluentAPI() {
	// Start with an empty error and build it up
	err := errors.NewError("", "", "", nil).
		WithCode(errors.ENOTFOUND).
		WithMessage("resource not found").
		WithOperation("OrderService.FindOrder").
		SetData(map[string]interface{}{
			"order_id": "ORD-12345",
			"user_id":  42,
		})

	fmt.Println(err.Code)
	fmt.Println(err.Message)
	fmt.Println(err.Op)

	// Output:
	// not_found
	// resource not found
	// OrderService.FindOrder
}

// Example_constructorVsBuilder demonstrates two equivalent approaches
func Example_constructorVsBuilder() {
	underlying := stderrors.New("network error")

	// Approach 1: Constructor with chaining
	err1 := errors.NewInternalError("API.Call", "request failed", underlying).
		SetData(map[string]string{"endpoint": "/api/users"})

	// Approach 2: Builder pattern
	err2 := errors.NewError("", "", "", nil).
		WithCode(errors.EINTERNAL).
		WithMessage("request failed").
		WithOperation("API.Call").
		WrapErr(underlying).
		SetData(map[string]string{"endpoint": "/api/users"})

	// Both produce the same result
	fmt.Println(err1.Code == err2.Code)
	fmt.Println(err1.Message == err2.Message)

	// Output:
	// true
	// true
}

// Example_stackTraces demonstrates stack trace preservation
func Example_stackTraces() {
	// Create an error with stack trace
	dbErr := stderrors.New("deadlock detected")
	err := errors.NewInternalError("Database.Query", "query failed", dbErr)

	// Normal error message (without stack trace)
	fmt.Printf("%v\n", err)

	// To see stack traces, use %+v in production logging
	// fmt.Printf("%+v\n", err)  // Would show full stack trace

	// Output:
	// Database.Query: deadlock detected
}

// Example_helperFunctions demonstrates using helper functions
func Example_helperFunctions() {
	// Wrap a standard error quickly
	dbErr := stderrors.New("connection timeout")
	err := errors.Wrap(dbErr, "database query failed", "UserRepo.FindByID")

	// Check error type
	if errors.HasCode(err, errors.EINTERNAL) {
		fmt.Println("Internal error detected")
	}

	// Get error code
	code := errors.GetCode(err)
	fmt.Println("Error code:", code)

	// Check if specific error is in chain
	if errors.Is(err, dbErr) {
		fmt.Println("Contains database error")
	}

	// Output:
	// Internal error detected
	// Error code: internal
	// Contains database error
}

// Example_wrapf demonstrates formatted error wrapping
func Example_wrapf() {
	userID := 12345
	dbErr := stderrors.New("connection timeout")

	// Wrap with formatted message
	err := errors.Wrapf(dbErr, "UserService.GetUser", "failed to fetch user %d", userID)

	fmt.Println(err.Message)
	fmt.Println(err.Op)

	// Output:
	// failed to fetch user 12345
	// UserService.GetUser
}

// Example_getData demonstrates extracting data from errors
func Example_getData() {
	// Create error with validation data
	err := errors.NewInvalidError("ValidateUser", "validation failed", nil).
		SetData(map[string]string{
			"email": "invalid format",
			"age":   "must be positive",
		})

	// Extract data using helper
	if data := errors.GetData(err); data != nil {
		if validationErrs, ok := data.(map[string]string); ok {
			fmt.Println("Validation errors found:")
			if msg, ok := validationErrs["email"]; ok {
				fmt.Printf("  email: %s\n", msg)
			}
		}
	}

	// Output:
	// Validation errors found:
	//   email: invalid format
}
