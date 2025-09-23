package unit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCLIFileOperations(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "This is test content for FileVault testing"

	// Test file creation
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test file reading
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Content mismatch: expected %q, got %q", testContent, string(content))
	}

	// Test file info
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Size() != int64(len(testContent)) {
		t.Errorf("File size mismatch: expected %d, got %d", len(testContent), info.Size())
	}
}

func TestBatchFileProcessing(t *testing.T) {
	tempDir := t.TempDir()

	testFiles := []string{"test1.txt", "test2.txt", "test3.txt"}

	// Create multiple test files
	for i, file := range testFiles {
		path := filepath.Join(tempDir, file)
		content := []byte("Test content for file " + string(rune('1'+i)))
		err := os.WriteFile(path, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Verify all files were created
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(entries) != len(testFiles) {
		t.Errorf("Expected %d files, found %d", len(testFiles), len(entries))
	}

	// Test batch file processing simulation
	processedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".txt" {
			processedCount++
		}
	}

	if processedCount != len(testFiles) {
		t.Errorf("Expected to process %d files, processed %d", len(testFiles), processedCount)
	}
}
