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
	
	if unwrapped != underlyingErr {
		t.Errorf("Expected Unwrap() to return the underlying error, got '%v'", unwrapped)
	}
}