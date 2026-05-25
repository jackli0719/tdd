package errors

import (
	"errors"
	"testing"
)

func TestErrorCodes(t *testing.T) {
	// Test that error codes match the contract
	if ErrCodeSuccess != 0 {
		t.Errorf("ErrCodeSuccess = %d, want 0", ErrCodeSuccess)
	}
	if ErrCodeParam != 400 {
		t.Errorf("ErrCodeParam = %d, want 400", ErrCodeParam)
	}
	if ErrCodeUnauthorized != 401 {
		t.Errorf("ErrCodeUnauthorized = %d, want 401", ErrCodeUnauthorized)
	}
	if ErrCodeForbidden != 403 {
		t.Errorf("ErrCodeForbidden = %d, want 403", ErrCodeForbidden)
	}
	if ErrCodeNotFound != 404 {
		t.Errorf("ErrCodeNotFound = %d, want 404", ErrCodeNotFound)
	}
	if ErrCodeInternal != 500 {
		t.Errorf("ErrCodeInternal = %d, want 500", ErrCodeInternal)
	}
}

func TestAppError(t *testing.T) {
	// Test Error() method
	e := New(400, "bad request")
	if e.Error() != "400: bad request" {
		t.Errorf("Error() = %s, want '400: bad request'", e.Error())
	}

	// Test Error() with wrapped error
	origErr := errors.New("underlying error")
	e2 := NewWithError(500, "server error", origErr)
	expected := "500: server error (underlying error)"
	if e2.Error() != expected {
		t.Errorf("Error() = %s, want '%s'", e2.Error(), expected)
	}

	// Test Unwrap
	if e2.Unwrap() != origErr {
		t.Errorf("Unwrap() = %v, want %v", e2.Unwrap(), origErr)
	}

	// Test errors.Is
	if !errors.Is(e2, origErr) {
		t.Error("errors.Is check failed")
	}
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		wantCode int
	}{
		{"Param", Param("invalid"), 400},
		{"Unauthorized", Unauthorized("no token"), 401},
		{"Forbidden", Forbidden("no access"), 403},
		{"NotFound", NotFound("user"), 404},
		{"Internal", Internal("oops"), 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %d, want %d", tt.err.Code, tt.wantCode)
			}
		})
	}

	// Test InternalWithError
	origErr := errors.New("db connection failed")
	innerErr := InternalWithError("database error", origErr)
	if innerErr.Code != 500 {
		t.Errorf("InternalWithError Code = %d, want 500", innerErr.Code)
	}
	if innerErr.Err != origErr {
		t.Errorf("InternalWithError Err = %v, want %v", innerErr.Err, origErr)
	}
}