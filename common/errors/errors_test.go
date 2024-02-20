package errors

import (
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("test error message", SystemErrorType)

	if err.GetErrorMessage() != "test error message" {
		t.Errorf("expected error message 'test error message', got '%s'", err.GetErrorMessage())
	}

	if err.GetErrorType() != SystemErrorType {
		t.Errorf("expected error type '%s', got '%s'", SystemErrorType, err.GetErrorType())
	}
}

func TestNewPhoneAlreadyUsedErrorMessage(t *testing.T) {
	expected := "phone number +6282324842834 already used."
	result := NewPhoneAlreadyUsedErrorMessage("+6282324842834")

	if result != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, result)
	}
}
