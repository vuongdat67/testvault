package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
)

func TestBasicDecryption(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "filevault_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test data
	testFile := filepath.Join(tempDir, "test.txt")
	encryptedFile := filepath.Join(tempDir, "test.txt.enc")
	decryptedFile := filepath.Join(tempDir, "decrypted.txt")
	password := "testpassword123"
	testData := []byte("Hello FileVault Test!")

	// Create test file
	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Encrypt file
	if err := core.EncryptFile(testFile, encryptedFile, password); err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}

	// Decrypt file
	if err := core.DecryptFile(encryptedFile, decryptedFile, password); err != nil {
		t.Fatalf("Failed to decrypt file: %v", err)
	}

	// Verify decrypted content
	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}

	if string(decryptedData) != string(testData) {
		t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s",
			string(testData), string(decryptedData))
	}
}
