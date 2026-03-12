package usecase

import (
	"testing"
)

func TestNewConvertUseCase(t *testing.T) {
	uc := NewConvertUseCase(nil, 1000)
	if uc == nil {
		t.Fatal("Expected NewConvertUseCase to return a non-nil object")
	}
	if uc.threshold != 1000 {
		t.Errorf("Expected threshold to be 1000, got %d", uc.threshold)
	}
}

// Note: Further tests requiring the artifact client or python execution
// should use interfaces and mocks for the client and a mocked command runner.
