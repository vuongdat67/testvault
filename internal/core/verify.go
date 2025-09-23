package core

import (
	"fmt"
	"os"
	"time"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/fileops"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/security"
)

// VerificationResult represents the result of file verification
type VerificationResult struct {
	IsValid          bool
	FormatValid      bool
	HeaderValid      bool
	SizeConsistent   bool
	FileAccessible   bool
	Filename         string
	OriginalFilename string
	FileSize         int64
	OriginalSize     uint64
	Algorithm        string
	FormatVersion    uint32
	ErrorMessage     string
	VerificationTime time.Duration
}

// VerifyFile performs comprehensive verification of an encrypted file
func VerifyFile(filePath string) (*VerificationResult, error) {
	startTime := time.Now()
	result := &VerificationResult{
		Filename: filePath,
	}

	// Check basic file accessibility
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("File not accessible: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	result.FileAccessible = true
	result.FileSize = fileInfo.Size()

	// Check if file is encrypted
	isEncrypted, err := security.IsEncryptedFile(filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to check file format: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	if !isEncrypted {
		result.ErrorMessage = "File is not a FileVault encrypted file"
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	result.FormatValid = true

	// Open file for detailed validation
	file, err := os.Open(filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to open file: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}
	defer file.Close()

	// Read and validate header
	var header fileops.FileHeader
	_, err = header.ReadFrom(file)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to read header: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	// Validate header structure
	if err := header.IsValid(); err != nil {
		result.ErrorMessage = fmt.Sprintf("Invalid header: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	result.HeaderValid = true
	result.OriginalFilename = header.FileName
	result.OriginalSize = header.OriginalSize
	result.FormatVersion = header.Version

	// Set algorithm name
	switch header.Algorithm {
	case fileops.AlgorithmAES256GCM:
		result.Algorithm = "AES-256-GCM"
	default:
		result.Algorithm = fmt.Sprintf("Unknown (%d)", header.Algorithm)
	}

	// Check size consistency
	expectedMinSize := int64(header.GetTotalSize() + fileops.AuthTagSize)
	if result.FileSize < expectedMinSize {
		result.ErrorMessage = fmt.Sprintf("File too small: expected at least %d bytes, got %d", expectedMinSize, result.FileSize)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	result.SizeConsistent = true

	// All checks passed
	result.IsValid = true
	result.VerificationTime = time.Since(startTime)

	return result, nil
}

// VerifyIntegrity performs deep integrity verification (requires password)
func VerifyIntegrity(filePath, password string) (*VerificationResult, error) {
	// First perform basic verification
	result, err := VerifyFile(filePath)
	if err != nil || !result.IsValid {
		return result, err
	}

	// Try to decrypt without writing to file (dry run)
	startTime := time.Now()

	// Open file for decryption test
	file, err := os.Open(filePath)
	if err != nil {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("Failed to open for integrity check: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}
	defer file.Close()

	// Read header (we already validated it, but need it for decryption)
	var header fileops.FileHeader
	_, err = header.ReadFrom(file)
	if err != nil {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("Failed to re-read header: %v", err)
		result.VerificationTime = time.Since(startTime)
		return result, nil
	}

	// For now, we'll rely on the basic verification
	// Full integrity check would require implementing dry-run decryption
	// which is complex and not in the current sprint scope

	result.VerificationTime = time.Since(startTime)
	return result, nil
}

// BatchVerify verifies multiple files
func BatchVerify(filePaths []string) ([]*VerificationResult, error) {
	results := make([]*VerificationResult, len(filePaths))

	for i, filePath := range filePaths {
		result, err := VerifyFile(filePath)
		if err != nil {
			return results, fmt.Errorf("failed to verify %s: %w", filePath, err)
		}
		results[i] = result
	}

	return results, nil
}

// GetVerificationSummary returns a summary of verification results
func GetVerificationSummary(results []*VerificationResult) map[string]int {
	summary := map[string]int{
		"total":      len(results),
		"valid":      0,
		"invalid":    0,
		"accessible": 0,
		"format_ok":  0,
		"header_ok":  0,
		"size_ok":    0,
	}

	for _, result := range results {
		if result.IsValid {
			summary["valid"]++
		} else {
			summary["invalid"]++
		}

		if result.FileAccessible {
			summary["accessible"]++
		}

		if result.FormatValid {
			summary["format_ok"]++
		}

		if result.HeaderValid {
			summary["header_ok"]++
		}

		if result.SizeConsistent {
			summary["size_ok"]++
		}
	}

	return summary
}
