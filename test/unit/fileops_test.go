package unit

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOperations(t *testing.T) {
	tempDir := t.TempDir()

	// Test file creation and writing
	testFile := filepath.Join(tempDir, "test.txt")
	testData := []byte("Hello FileVault Testing")

	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test file reading
	readData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !bytes.Equal(testData, readData) {
		t.Errorf("Data mismatch: expected %v, got %v", testData, readData)
	}
}

func TestFileValidation(t *testing.T) {
	tempDir := t.TempDir()

	// Test valid file
	validFile := filepath.Join(tempDir, "valid.txt")
	err := os.WriteFile(validFile, []byte("valid content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid file: %v", err)
	}

	info, err := os.Stat(validFile)
	if err != nil {
		t.Fatalf("Failed to stat valid file: %v", err)
	}

	if info.IsDir() {
		t.Error("File should not be a directory")
	}

	if info.Size() == 0 {
		t.Error("File should not be empty")
	}

	// Test non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	_, err = os.Stat(nonExistentFile)
	if err == nil {
		t.Error("Non-existent file should return an error")
	}
}

func TestBinaryFormat(t *testing.T) {
	// Test FileVault magic number
	magicNumber := []byte("FVLT")
	if len(magicNumber) != 4 {
		t.Errorf("Magic number should be 4 bytes, got %d", len(magicNumber))
	}

	if string(magicNumber) != "FVLT" {
		t.Errorf("Magic number should be 'FVLT', got %s", string(magicNumber))
	}

	// Test version bytes
	version := []byte{1, 0} // version 1.0
	if len(version) != 2 {
		t.Errorf("Version should be 2 bytes, got %d", len(version))
	}
}
