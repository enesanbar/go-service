package errors

import (
	"errors"
	"testing"
)

// BenchmarkNewError benchmarks creating a new error with no underlying error
func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewError("TEST001", "test message", "TestOp", nil)
	}
}

// BenchmarkNewErrorWithWrapping benchmarks creating an error with wrapping
func BenchmarkNewErrorWithWrapping(b *testing.B) {
	underlyingErr := errors.New("underlying error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewError("TEST001", "test message", "TestOp", underlyingErr)
	}
}

// BenchmarkTypedConstructor benchmarks using a typed constructor
func BenchmarkTypedConstructor(b *testing.B) {
	underlyingErr := errors.New("database error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewInternalError("Database.Query", "query failed", underlyingErr)
	}
}

// BenchmarkFluentAPI benchmarks building errors with method chaining
func BenchmarkFluentAPI(b *testing.B) {
	underlyingErr := errors.New("test error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewError("", "", "", nil).
			WithCode(EINVALID).
			WithMessage("validation failed").
			WithOperation("ValidateUser").
			WrapErr(underlyingErr).
			SetData(map[string]string{"field": "email"})
	}
}

// BenchmarkErrorMethod benchmarks calling the Error() method
func BenchmarkErrorMethod(b *testing.B) {
	err := NewInternalError("TestOp", "test message", errors.New("underlying"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
