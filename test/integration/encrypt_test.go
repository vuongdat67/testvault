package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
)

func TestBasicEncryption(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "filevault_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test data
	testFile := filepath.Join(tempDir, "test.txt")
	encryptedFile := filepath.Join(tempDir, "test.txt.enc")
	password := "testpassword123"
	testData := []byte("Hello FileVault Encryption Test!")

	// Create test file
	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Encrypt file
	if err := core.EncryptFile(testFile, encryptedFile, password); err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}

	// Verify encrypted file exists
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatalf("Encrypted file was not created")
	}

	// Verify encrypted file is different from original
	encryptedData, err := os.ReadFile(encryptedFile)
	if err != nil {
		t.Fatalf("Failed to read encrypted file: %v", err)
	}

	// Should not contain original data (basic check)
	if string(encryptedData) == string(testData) {
		t.Error("Encrypted file appears to contain unencrypted data")
	}

	// Verify file size is larger (due to header + auth tag)
	if len(encryptedData) <= len(testData) {
		t.Error("Encrypted file should be larger than original due to metadata")
	}
}
