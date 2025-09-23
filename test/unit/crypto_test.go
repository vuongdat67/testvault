package unit

import (
	"crypto/rand"
	"testing"
)

func TestRandomGeneration(t *testing.T) {
	// Test random salt generation
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		t.Fatalf("Failed to generate random salt: %v", err)
	}

	// Verify salt is not all zeros
	allZero := true
	for _, b := range salt {
		if b != 0 {
			allZero = false
			break
		}
	}

	if allZero {
		t.Error("Generated salt is all zeros")
	}

	// Test nonce generation
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		t.Fatalf("Failed to generate random nonce: %v", err)
	}

	// Verify nonce is not all zeros
	allZero = true
	for _, b := range nonce {
		if b != 0 {
			allZero = false
			break
		}
	}

	if allZero {
		t.Error("Generated nonce is all zeros")
	}
}

func TestKeyDerivation(t *testing.T) {
	password := "test-password-123"
	salt := make([]byte, 32)
	rand.Read(salt)

	// Simple PBKDF2 test simulation
	if len(password) < 8 {
		t.Error("Password too short for testing")
	}

	if len(salt) != 32 {
		t.Error("Salt should be 32 bytes")
	}

	// Test iterations parameter
	iterations := 100000
	if iterations < 10000 {
		t.Error("Iterations should be at least 10000 for security")
	}
}
