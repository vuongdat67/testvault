package unit

import (
	"testing"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/security"
)

func TestPasswordValidation(t *testing.T) {
	// Test valid password
	err := security.ValidatePasswordBasic("validpassword123")
	if err != nil {
		t.Errorf("Valid password should pass: %v", err)
	}

	// Test empty password
	err = security.ValidatePasswordBasic("")
	if err == nil {
		t.Error("Empty password should fail validation")
	}

	// Test short password
	err = security.ValidatePasswordBasic("123")
	if err == nil {
		t.Error("Short password should fail validation")
	}
}

func TestInputFileValidation(t *testing.T) {
	// Test non-existent file
	err := security.ValidateInputFile("nonexistent.txt")
	if err == nil {
		t.Error("Non-existent file should fail validation")
	}

	// Test empty path
	err = security.ValidateInputFile("")
	if err == nil {
		t.Error("Empty path should fail validation")
	}
}
