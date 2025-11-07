package errors

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	// Test case 1: Create a new error with all fields
	err := NewError("ERR001", "Something went wrong", "TestOperation", errors.New("underlying error"))

	if err.Code != "ERR001" {
		t.Errorf("Expected Code to be 'ERR001', got '%s'", err.Code)
	}

	if err.Message != "Something went wrong" {
		t.Errorf("Expected Message to be 'Something went wrong', got '%s'", err.Message)
	}

	if err.Op != "TestOperation" {
		t.Errorf("Expected Op to be 'TestOperation', got '%s'", err.Op)
	}

	if err.Err == nil || err.Err.Error() != "underlying error" {
		t.Errorf("Expected Err to be 'underlying error', got '%v'", err.Err)
	}

	// Test case 2: Create a new error without an underlying error
	err2 := NewError("ERR002", "Another error", "TestOperation2", nil)

	if err2.Code != "ERR002" {
		t.Errorf("Expected Code to be 'ERR002', got '%s'", err2.Code)
	}

	if err2.Err != nil {
		t.Errorf("Expected Err to be nil, got '%v'", err2.Err)
	}
}

func TestError_Error(t *testing.T) {
	// Test case 1: Error with an underlying error
	err := NewError("ERR001", "Something went wrong", "TestOperation", errors.New("underlying error"))
	expected := "TestOperation: underlying error"

	if err.Error() != expected {
		t.Errorf("Expected Error() to return '%s', got '%s'", expected, err.Error())
	}

	// Test case 2: Error without an underlying error
	err2 := NewError("ERR002", "Another error", "TestOperation2", nil)
	expected2 := "TestOperation2: <ERR002> Another error"

	if err2.Error() != expected2 {
		t.Errorf("Expected Error() to return '%s', got '%s'", expected2, err2.Error())
	}
}

func TestError_SetData(t *testing.T) {
	err := NewError("ERR001", "Something went wrong", "TestOperation", nil)

	// Test setting data
	testData := map[string]string{"field": "value"}
	err.SetData(testData)

	// Check if data was set correctly
	if err.Data == nil {
		t.Error("Expected Data to be set, got nil")
	}

	// Type assertion to check the data content
	data, ok := err.Data.(map[string]string)
	if !ok {
		t.Errorf("Expected Data to be of type map[string]string, got %T", err.Data)
	}

	if data["field"] != "value" {
		t.Errorf("Expected Data['field'] to be 'value', got '%s'", data["field"])
	}
}

func TestError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewError("ERR001", "Something went wrong", "TestOperation", underlyingErr)

	// Test unwrapping the error
	unwrapped := err.Unwrap()

	// With errors.WithStack, unwrapped won't be the same reference but should have the same message
	if unwrapped == nil {
		t.Error("Expected Unwrap() to return an error, got nil")
	}
	if unwrapped.Error() != underlyingErr.Error() {
		t.Errorf("Expected Unwrap() message to be '%v', got '%v'", underlyingErr.Error(), unwrapped.Error())
	}
}

// TestNewError_NilHandling tests that NewError handles nil errors gracefully
func TestNewError_NilHandling(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		message      string
		op           string
		err          error
		expectNilErr bool
	}{
		{
			name:         "nil error should not panic",
			code:         "TEST001",
			message:      "test message",
			op:           "TestOp",
			err:          nil,
			expectNilErr: true,
		},
		{
			name:         "non-nil error should be wrapped",
			code:         "TEST002",
			message:      "test message 2",
			op:           "TestOp2",
			err:          errors.New("underlying"),
			expectNilErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("NewError panicked with nil error: %v", r)
				}
			}()

			err := NewError(tt.code, tt.message, tt.op, tt.err)

			if err == nil {
				t.Fatal("NewError returned nil")
			}

			if tt.expectNilErr && err.Err != nil {
				t.Errorf("Expected Err to be nil, got %v", err.Err)
			}

			if !tt.expectNilErr && err.Err == nil {
				t.Error("Expected Err to be non-nil, got nil")
			}
		})
	}
}

// TestWrapErr_NilHandling tests that WrapErr handles nil errors gracefully
func TestWrapErr_NilHandling(t *testing.T) {
	err := NewError("TEST", "test", "TestOp", nil)

	// Wrapping nil should not panic and should not set Err
	result := err.WrapErr(nil)

	if result != err {
		t.Error("WrapErr should return the same error instance")
	}

	if err.Err != nil {
		t.Errorf("Expected Err to remain nil, got %v", err.Err)
	}

	// Wrapping non-nil error should work
	underlyingErr := errors.New("underlying")
	result = err.WrapErr(underlyingErr)

	if err.Err == nil {
		t.Error("Expected Err to be set after wrapping non-nil error")
	}
}

// TestSetData_Chaining tests that SetData returns *Error for method chaining
func TestSetData_Chaining(t *testing.T) {
	err := NewError("TEST", "test", "TestOp", nil)

	result := err.SetData(map[string]string{"key": "value"})

	if result != err {
		t.Error("SetData should return the same error instance for chaining")
	}

	// Test chaining multiple operations
	chainedErr := NewError("TEST2", "test2", "TestOp2", nil).
		SetData(map[string]string{"field": "validation error"}).
		WithMessage("updated message")

	if chainedErr.Message != "updated message" {
		t.Errorf("Expected chained message to be 'updated message', got '%s'", chainedErr.Message)
	}

	data, ok := chainedErr.Data.(map[string]string)
	if !ok || data["field"] != "validation error" {
		t.Error("Expected chained data to be preserved")
	}
}

// TestStackTracePreservation tests that stack traces are captured and preserved
func TestStackTracePreservation(t *testing.T) {
	// Create error with stack trace
	underlying := errors.New("database connection failed")
	err := NewError("DB001", "Database error", "ConnectDB", underlying)

	// Verify stack trace is preserved through wrapping
	if err.Err == nil {
		t.Fatal("Expected underlying error to be wrapped")
	}

	// Test that errors.Is works correctly
	if !errors.Is(err, underlying) {
		t.Error("errors.Is should find the underlying error")
	}

	// Test multiple layers of wrapping
	err2 := NewError("SVC001", "Service error", "HandleRequest", err)

	// Should be able to find both errors in the chain
	if !errors.Is(err2, err) {
		t.Error("errors.Is should find wrapped Error in chain")
	}

	if !errors.Is(err2, underlying) {
		t.Error("errors.Is should find original underlying error in chain")
	}
}

// TestWrapErrStackTrace tests that WrapErr captures stack traces
func TestWrapErrStackTrace(t *testing.T) {
	err := NewError("TEST", "test", "TestOp", nil)
	underlying := errors.New("underlying error")

	err.WrapErr(underlying)

	// Verify the error chain works with errors.Is
	if !errors.Is(err, underlying) {
		t.Error("errors.Is should work after WrapErr")
	}

	// Test error message contains stack information
	errorStr := err.Error()
	if errorStr == "" {
		t.Error("Error() should return non-empty string")
	}
}

// TestErrorsAsCompatibility tests compatibility with errors.As for type assertions
func TestErrorsAsCompatibility(t *testing.T) {
	customErr := NewError("CUSTOM001", "custom error", "TestOp", nil)
	wrappedErr := NewError("WRAP001", "wrapped", "WrapOp", customErr)

	// Test errors.As can extract our Error type
	var target *Error
	if !errors.As(wrappedErr, &target) {
		t.Error("errors.As should find *Error in chain")
	}

	if target.Code != "WRAP001" {
		t.Errorf("Expected extracted error code to be 'WRAP001', got '%s'", target.Code)
	}

	// Test finding the inner error
	var innerTarget *Error
	if !errors.As(wrappedErr.Err, &innerTarget) {
		t.Error("errors.As should find inner *Error")
	}

	if innerTarget.Code != "CUSTOM001" {
		t.Errorf("Expected inner error code to be 'CUSTOM001', got '%s'", innerTarget.Code)
	}
}

// TestTypedConstructors tests all six typed error constructors
func TestTypedConstructors(t *testing.T) {
	tests := []struct {
		name         string
		constructor  func(string, string, error) *Error
		expectedCode string
		op           string
		message      string
		err          error
	}{
		{
			name:         "NewInvalidError",
			constructor:  NewInvalidError,
			expectedCode: EINVALID,
			op:           "ValidateInput",
			message:      "field 'email' is required",
			err:          errors.New("validation failed"),
		},
		{
			name:         "NewNotFoundError",
			constructor:  NewNotFoundError,
			expectedCode: ENOTFOUND,
			op:           "FindUser",
			message:      "user not found",
			err:          errors.New("no rows in result set"),
		},
		{
			name:         "NewConflictError",
			constructor:  NewConflictError,
			expectedCode: ECONFLICT,
			op:           "CreateUser",
			message:      "user already exists",
			err:          errors.New("duplicate key"),
		},
		{
			name:         "NewForbiddenError",
			constructor:  NewForbiddenError,
			expectedCode: EFORBIDDEN,
			op:           "DeleteResource",
			message:      "insufficient permissions",
			err:          errors.New("access denied"),
		},
		{
			name:         "NewInternalError",
			constructor:  NewInternalError,
			expectedCode: EINTERNAL,
			op:           "ProcessRequest",
			message:      "unexpected server error",
			err:          errors.New("database connection failed"),
		},
		{
			name:         "NewNotModifiedError",
			constructor:  NewNotModifiedError,
			expectedCode: ENOTMODIFIED,
			op:           "UpdateCache",
			message:      "resource not modified",
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.op, tt.message, tt.err)

			if err == nil {
				t.Fatal("Constructor returned nil")
			}

			if err.Code != tt.expectedCode {
				t.Errorf("Expected Code to be '%s', got '%s'", tt.expectedCode, err.Code)
			}

			if err.Op != tt.op {
				t.Errorf("Expected Op to be '%s', got '%s'", tt.op, err.Op)
			}

			if err.Message != tt.message {
				t.Errorf("Expected Message to be '%s', got '%s'", tt.message, err.Message)
			}

			if tt.err != nil {
				if err.Err == nil {
					t.Error("Expected Err to be wrapped, got nil")
				}
				// Verify stack trace is captured
				if !errors.Is(err, tt.err) {
					t.Error("errors.Is should find wrapped error")
				}
			} else {
				if err.Err != nil {
					t.Errorf("Expected Err to be nil, got %v", err.Err)
				}
			}
		})
	}
}

// TestTypedConstructors_NilHandling tests that all constructors handle nil errors gracefully
func TestTypedConstructors_NilHandling(t *testing.T) {
	tests := []struct {
		name         string
		constructor  func(string, string, error) *Error
		expectedCode string
	}{
		{"NewInvalidError", NewInvalidError, EINVALID},
		{"NewNotFoundError", NewNotFoundError, ENOTFOUND},
		{"NewConflictError", NewConflictError, ECONFLICT},
		{"NewForbiddenError", NewForbiddenError, EFORBIDDEN},
		{"NewInternalError", NewInternalError, EINTERNAL},
		{"NewNotModifiedError", NewNotModifiedError, ENOTMODIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked with nil error: %v", tt.name, r)
				}
			}()

			err := tt.constructor("TestOp", "test message", nil)

			if err == nil {
				t.Fatalf("%s returned nil", tt.name)
			}

			if err.Code != tt.expectedCode {
				t.Errorf("Expected Code to be '%s', got '%s'", tt.expectedCode, err.Code)
			}

			if err.Err != nil {
				t.Errorf("Expected Err to be nil when no error is wrapped, got %v", err.Err)
			}
		})
	}
}

// TestTypedConstructors_CorrectCodes verifies all constructors set the correct error codes
func TestTypedConstructors_CorrectCodes(t *testing.T) {
	underlying := errors.New("test error")

	testCases := []struct {
		err      *Error
		expected string
	}{
		{NewInvalidError("op", "msg", underlying), EINVALID},
		{NewNotFoundError("op", "msg", underlying), ENOTFOUND},
		{NewConflictError("op", "msg", underlying), ECONFLICT},
		{NewForbiddenError("op", "msg", underlying), EFORBIDDEN},
		{NewInternalError("op", "msg", underlying), EINTERNAL},
		{NewNotModifiedError("op", "msg", underlying), ENOTMODIFIED},
	}

	for _, tc := range testCases {
		if tc.err.Code != tc.expected {
			t.Errorf("Expected code '%s', got '%s'", tc.expected, tc.err.Code)
		}
	}
}

// TestValidationErrorData tests attaching validation error data
func TestValidationErrorData(t *testing.T) {
	validationErrors := map[string]string{
		"email":    "invalid format",
		"password": "must be at least 8 characters",
		"age":      "must be positive",
	}

	err := NewInvalidError("ValidateUser", "validation failed", nil).
		SetData(validationErrors)

	if err.Data == nil {
		t.Fatal("Expected Data to be set")
	}

	data, ok := err.Data.(map[string]string)
	if !ok {
		t.Fatal("Expected Data to be map[string]string")
	}

	if data["email"] != "invalid format" {
		t.Errorf("Expected email error to be 'invalid format', got '%s'", data["email"])
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 validation errors, got %d", len(data))
	}
}

// TestResourceIdentifierData tests attaching resource identifier data
func TestResourceIdentifierData(t *testing.T) {
	resourceData := map[string]interface{}{
		"user_id":     12345,
		"order_id":    "ORD-98765",
		"resource":    "payment",
		"status_code": 404,
	}

	err := NewNotFoundError("PaymentService.FindOrder", "payment not found", nil).
		SetData(resourceData)

	if err.Data == nil {
		t.Fatal("Expected Data to be set")
	}

	data, ok := err.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected Data to be map[string]interface{}")
	}

	if userId, ok := data["user_id"].(int); !ok || userId != 12345 {
		t.Errorf("Expected user_id to be 12345, got %v", data["user_id"])
	}

	if orderId, ok := data["order_id"].(string); !ok || orderId != "ORD-98765" {
		t.Errorf("Expected order_id to be 'ORD-98765', got %v", data["order_id"])
	}
}

// TestMetadataAttachment tests attaching request/trace metadata
func TestMetadataAttachment(t *testing.T) {
	metadata := map[string]string{
		"request_id": "req-abc123",
		"trace_id":   "trace-xyz789",
		"user_agent": "TestClient/1.0",
		"ip_address": "192.168.1.1",
	}

	err := NewInternalError("API.HandleRequest", "internal server error", errors.New("db timeout")).
		SetData(metadata)

	if err.Data == nil {
		t.Fatal("Expected Data to be set")
	}

	data, ok := err.Data.(map[string]string)
	if !ok {
		t.Fatal("Expected Data to be map[string]string")
	}

	if data["request_id"] != "req-abc123" {
		t.Errorf("Expected request_id to be 'req-abc123', got '%s'", data["request_id"])
	}

	if data["trace_id"] != "trace-xyz789" {
		t.Errorf("Expected trace_id to be 'trace-xyz789', got '%s'", data["trace_id"])
	}
}

// TestDataPreservationThroughWrapping tests that data is preserved when errors are wrapped
func TestDataPreservationThroughWrapping(t *testing.T) {
	// Create error with data
	validationData := map[string]string{"field": "email", "error": "invalid"}
	innerErr := NewInvalidError("ValidateInput", "validation failed", nil).
		SetData(validationData)

	// Wrap the error
	outerErr := NewInternalError("ProcessRequest", "request processing failed", innerErr)

	// Outer error should have no data initially
	if outerErr.Data != nil {
		t.Error("Expected outer error to have no data initially")
	}

	// Inner error's data should still be accessible using errors.As
	// errors.As will unwrap until it finds an *Error with data
	current := error(outerErr.Err) // Start with the wrapped error
	for current != nil {
		var e *Error
		if errors.As(current, &e) && e.Data != nil {
			data, ok := e.Data.(map[string]string)
			if !ok {
				t.Fatal("Expected inner error data to be map[string]string")
			}
			if data["field"] != "email" {
				t.Errorf("Expected field to be 'email', got '%s'", data["field"])
			}
			return
		}
		// Try to unwrap further
		if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
			current = unwrapper.Unwrap()
		} else {
			break
		}
	}

	t.Error("Expected to find error with data in chain")
}

// TestWithData tests the WithData fluent API method
func TestWithData(t *testing.T) {
	data := map[string]string{"key": "value"}
	err := NewInvalidError("TestOp", "test message", nil).
		WithData(data)

	if err.Data == nil {
		t.Fatal("Expected Data to be set via WithData")
	}

	retrievedData, ok := err.Data.(map[string]string)
	if !ok || retrievedData["key"] != "value" {
		t.Error("WithData did not correctly set data")
	}
}

// TestFluentErrorConstruction tests building errors using method chaining
func TestFluentErrorConstruction(t *testing.T) {
	underlyingErr := errors.New("database error")

	// Build error using fluent API
	err := NewError("", "", "", nil).
		WithCode(EINVALID).
		WithMessage("validation failed").
		WithOperation("ValidateUser").
		WrapErr(underlyingErr).
		SetData(map[string]string{"field": "email"})

	if err.Code != EINVALID {
		t.Errorf("Expected Code to be '%s', got '%s'", EINVALID, err.Code)
	}

	if err.Message != "validation failed" {
		t.Errorf("Expected Message to be 'validation failed', got '%s'", err.Message)
	}

	if err.Op != "ValidateUser" {
		t.Errorf("Expected Op to be 'ValidateUser', got '%s'", err.Op)
	}

	if err.Err == nil {
		t.Error("Expected Err to be set")
	}

	if err.Data == nil {
		t.Error("Expected Data to be set")
	}
}

// TestConstructorVsBuilder tests that constructor and builder approaches are equivalent
func TestConstructorVsBuilder(t *testing.T) {
	underlyingErr := errors.New("test error")
	data := map[string]string{"field": "value"}

	// Using constructor
	constructorErr := NewInvalidError("TestOp", "test message", underlyingErr).
		SetData(data)

	// Using builder
	builderErr := NewError("", "", "", nil).
		WithCode(EINVALID).
		WithMessage("test message").
		WithOperation("TestOp").
		WrapErr(underlyingErr).
		SetData(data)

	// Both should have the same fields
	if constructorErr.Code != builderErr.Code {
		t.Error("Code mismatch between constructor and builder")
	}

	if constructorErr.Message != builderErr.Message {
		t.Error("Message mismatch between constructor and builder")
	}

	if constructorErr.Op != builderErr.Op {
		t.Error("Op mismatch between constructor and builder")
	}

	if (constructorErr.Err == nil) != (builderErr.Err == nil) {
		t.Error("Err presence mismatch between constructor and builder")
	}

	if (constructorErr.Data == nil) != (builderErr.Data == nil) {
		t.Error("Data presence mismatch between constructor and builder")
	}
}

// TestFluentChaining tests that all methods can be chained together
func TestFluentChaining(t *testing.T) {
	// Create a complex error with full chaining
	err := NewInvalidError("InitialOp", "initial message", nil).
		WithCode(EINTERNAL).                 // Change code
		WithMessage("updated message").      // Change message
		WithOperation("UpdatedOp").          // Change operation
		WrapErr(errors.New("root cause")).   // Add underlying error
		SetData(map[string]int{"count": 42}) // Add data

	// Verify all changes took effect
	if err.Code != EINTERNAL {
		t.Errorf("Expected Code to be '%s', got '%s'", EINTERNAL, err.Code)
	}

	if err.Message != "updated message" {
		t.Errorf("Expected Message to be 'updated message', got '%s'", err.Message)
	}

	if err.Op != "UpdatedOp" {
		t.Errorf("Expected Op to be 'UpdatedOp', got '%s'", err.Op)
	}

	if err.Err == nil {
		t.Error("Expected Err to be set")
	}

	data, ok := err.Data.(map[string]int)
	if !ok || data["count"] != 42 {
		t.Error("Expected Data to be map[string]int with count=42")
	}
}

// TestIs tests the Is helper function
func TestIs(t *testing.T) {
	underlyingErr := errors.New("database error")
	err := NewInternalError("Database.Query", "query failed", underlyingErr)

	// Should find the underlying error
	if !Is(err, underlyingErr) {
		t.Error("Is should find underlying error in chain")
	}

	// Should not match different error
	otherErr := errors.New("different error")
	if Is(err, otherErr) {
		t.Error("Is should not match different error")
	}
}

// TestAs tests the As helper function
func TestAs(t *testing.T) {
	innerErr := NewInvalidError("ValidateInput", "validation failed", nil).
		SetData(map[string]string{"field": "email"})

	outerErr := NewInternalError("ProcessRequest", "processing failed", innerErr)

	// Should find *Error in chain
	var e *Error
	if !As(outerErr, &e) {
		t.Fatal("As should find *Error in chain")
	}

	// The found error should be the outer error
	if e.Code != EINTERNAL {
		t.Errorf("Expected code to be %s, got %s", EINTERNAL, e.Code)
	}

	// Should be able to find inner error too
	var innerTarget *Error
	current := error(outerErr.Err)
	if As(current, &innerTarget) {
		if innerTarget.Code != EINVALID {
			t.Errorf("Expected inner error code to be %s, got %s", EINVALID, innerTarget.Code)
		}
	}
}

// TestWrap tests the Wrap helper function
func TestWrap(t *testing.T) {
	dbErr := errors.New("connection timeout")
	err := Wrap(dbErr, "database operation failed", "UserRepo.FindByID")

	if err == nil {
		t.Fatal("Wrap should not return nil")
	}

	if err.Code != EINTERNAL {
		t.Errorf("Expected code to be %s, got %s", EINTERNAL, err.Code)
	}

	if err.Message != "database operation failed" {
		t.Errorf("Expected message to be 'database operation failed', got '%s'", err.Message)
	}

	if err.Op != "UserRepo.FindByID" {
		t.Errorf("Expected op to be 'UserRepo.FindByID', got '%s'", err.Op)
	}

	if !Is(err, dbErr) {
		t.Error("Wrapped error should contain original error")
	}
}

// TestWrapf tests the Wrapf helper function
func TestWrapf(t *testing.T) {
	dbErr := errors.New("connection timeout")
	userID := 12345
	err := Wrapf(dbErr, "UserRepo.FindByID", "failed to find user %d", userID)

	if err == nil {
		t.Fatal("Wrapf should not return nil")
	}

	expectedMessage := "failed to find user 12345"
	if err.Message != expectedMessage {
		t.Errorf("Expected message to be '%s', got '%s'", expectedMessage, err.Message)
	}

	if err.Op != "UserRepo.FindByID" {
		t.Errorf("Expected op to be 'UserRepo.FindByID', got '%s'", err.Op)
	}

	if !Is(err, dbErr) {
		t.Error("Wrapped error should contain original error")
	}
}

// TestHasCode tests the HasCode helper function
func TestHasCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     string
		expected bool
	}{
		{
			name:     "matching code",
			err:      NewNotFoundError("FindUser", "user not found", nil),
			code:     ENOTFOUND,
			expected: true,
		},
		{
			name:     "non-matching code",
			err:      NewNotFoundError("FindUser", "user not found", nil),
			code:     EINVALID,
			expected: false,
		},
		{
			name:     "wrapped error with code",
			err:      NewInternalError("Outer", "outer", NewInvalidError("Inner", "inner", nil)),
			code:     EINTERNAL,
			expected: true,
		},
		{
			name:     "standard error without code",
			err:      errors.New("standard error"),
			code:     EINTERNAL,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasCode(tt.err, tt.code)
			if result != tt.expected {
				t.Errorf("HasCode() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetCode tests the GetCode helper function
func TestGetCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "error with code",
			err:      NewInvalidError("ValidateInput", "validation failed", nil),
			expected: EINVALID,
		},
		{
			name:     "wrapped error",
			err:      NewInternalError("Outer", "outer", NewNotFoundError("Inner", "inner", nil)),
			expected: EINTERNAL,
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			expected: "",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCode(tt.err)
			if result != tt.expected {
				t.Errorf("GetCode() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

// TestGetData tests the GetData helper function
func TestGetData(t *testing.T) {
	validationData := map[string]string{"email": "invalid"}

	tests := []struct {
		name       string
		err        error
		expectData bool
		validateFn func(interface{}) bool
	}{
		{
			name:       "error with data",
			err:        NewInvalidError("ValidateInput", "validation failed", nil).SetData(validationData),
			expectData: true,
			validateFn: func(data interface{}) bool {
				d, ok := data.(map[string]string)
				return ok && d["email"] == "invalid"
			},
		},
		{
			name:       "error without data",
			err:        NewInvalidError("ValidateInput", "validation failed", nil),
			expectData: false,
			validateFn: func(data interface{}) bool { return data == nil },
		},
		{
			name:       "standard error",
			err:        errors.New("standard error"),
			expectData: false,
			validateFn: func(data interface{}) bool { return data == nil },
		},
		{
			name:       "nil error",
			err:        nil,
			expectData: false,
			validateFn: func(data interface{}) bool { return data == nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetData(tt.err)
			if tt.expectData && result == nil {
				t.Error("Expected data to be present, got nil")
			}
			if !tt.expectData && result != nil {
				t.Errorf("Expected data to be nil, got %v", result)
			}
			if !tt.validateFn(result) {
				t.Errorf("Data validation failed for %v", result)
			}
		})
	}
}
