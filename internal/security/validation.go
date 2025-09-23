package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/errors"
)

// ValidateInputFile validates an input file for encryption/decryption
func ValidateInputFile(filePath string) error {
	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return errors.NewFileNotFoundError(filePath)
	}
	if err != nil {
		return errors.NewError(errors.ErrFileReadError, "cannot access file", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return errors.NewError(errors.ErrInvalidInput, fmt.Sprintf("not a regular file: %s", filePath), nil)
	}

	// Check file permissions
	file, err := os.Open(filePath)
	if err != nil {
		return errors.NewPermissionDeniedError(filePath, err)
	}
	file.Close()

	// Check file size limits (prevent extremely large files)
	const maxFileSize = 10 * 1024 * 1024 * 1024 // 10GB limit
	if info.Size() > maxFileSize {
		return errors.NewError(errors.ErrFileTooLarge, 
			fmt.Sprintf("file too large: %d bytes (max %d)", info.Size(), maxFileSize), nil)
	}

	return nil
}

// ValidateOutputFile validates an output file path and checks for overwrite protection
func ValidateOutputFile(filePath string, force bool) error {
	// Validate filename security
	if err := ValidateFilename(filepath.Base(filePath)); err != nil {
		return err
	}

	// Check if output directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.NewError(errors.ErrFileNotFound, fmt.Sprintf("output directory does not exist: %s", dir), nil)
	}

	// Check if output file already exists
	if _, err := os.Stat(filePath); err == nil && !force {
		fvErr := errors.NewError(errors.ErrFileAlreadyExists, fmt.Sprintf("output file already exists: %s", filePath), nil)
		fvErr.Context["filename"] = filePath
		return fvErr
	}

	// Check if we can create the file
	testFile, err := os.Create(filePath)
	if err != nil {
		return errors.NewPermissionDeniedError(filePath, err)
	}
	testFile.Close()
	os.Remove(filePath) // Clean up test file

	return nil
}

// ValidatePasswordBasic performs basic password validation
func ValidatePasswordBasic(password string) error {
	if len(password) == 0 {
		return errors.NewError(errors.ErrInvalidPassword, "password cannot be empty", nil)
	}

	if len(password) < 8 {
		return errors.NewError(errors.ErrWeakPassword, "password must be at least 8 characters long", nil)
	}

	return nil
}

// ValidatePasswordStrict performs strict password validation
func ValidatePasswordStrict(password string, policy PasswordPolicy) error {
	// Basic validation first
	if err := ValidatePasswordBasic(password); err != nil {
		return err
	}

	// Apply policy validation
	if err := ValidatePassword(password, policy); err != nil {
		return errors.NewError(errors.ErrWeakPassword, err.Error(), nil)
	}

	return nil
}

// ValidateFilename validates a filename for security issues
func ValidateFilename(filename string) error {
	if filename == "" {
		return errors.NewError(errors.ErrInvalidInput, "filename cannot be empty", nil)
	}

	// Check for path traversal attempts
	if strings.Contains(filename, "..") {
		return errors.NewError(errors.ErrSecurityViolation, 
			fmt.Sprintf("filename contains path traversal sequence: %s", filename), nil)
	}

	// Check for absolute paths
	if filepath.IsAbs(filename) {
		return errors.NewError(errors.ErrSecurityViolation, 
			fmt.Sprintf("filename cannot be an absolute path: %s", filename), nil)
	}

	// Check for suspicious characters
	suspiciousChars := []string{"\x00", "\n", "\r"}
	for _, char := range suspiciousChars {
		if strings.Contains(filename, char) {
			return errors.NewError(errors.ErrSecurityViolation, 
				fmt.Sprintf("filename contains suspicious character: %s", filename), nil)
		}
	}

	// Check filename length (Windows has 255 char limit)
	if len(filename) > 255 {
		return errors.NewError(errors.ErrInvalidInput, "filename too long", nil)
	}

	return nil
}

// IsEncryptedFile checks if a file appears to be a FileVault encrypted file
func IsEncryptedFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, errors.NewPermissionDeniedError(filePath, err)
	}
	defer file.Close()

	// Read the first 4 bytes to check magic number
	magic := make([]byte, 4)
	n, err := file.Read(magic)
	if err != nil || n < 4 {
		return false, nil
	}

	return string(magic) == "FVLT", nil
}

// ValidateEncryptedFile validates that a file is a proper FileVault file
func ValidateEncryptedFile(filePath string) error {
	isEncrypted, err := IsEncryptedFile(filePath)
	if err != nil {
		return err
	}
	
	if !isEncrypted {
		return errors.NewInvalidFormatError(filePath)
	}
	
	return nil
}

// SanitizeFilename sanitizes a filename by removing or replacing dangerous characters
func SanitizeFilename(filename string) string {
	// Replace path separators and other dangerous characters
	dangerous := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := filename

	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	// Remove leading/trailing spaces and dots
	sanitized = strings.Trim(sanitized, " .")

	// Ensure filename is not empty
	if sanitized == "" {
		sanitized = "unnamed_file"
	}

	// Limit length
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	return sanitized
}

// ValidateFileIntegrity performs basic file integrity checks
func ValidateFileIntegrity(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return errors.NewFileNotFoundError(filePath)
	}

	// Check for zero-byte files
	if info.Size() == 0 {
		return errors.NewCorruptedFileError(filePath, fmt.Errorf("file is empty"))
	}

	// Check for minimum header size for encrypted files
	if strings.HasSuffix(filePath, ".enc") {
		const minHeaderSize = 100
		if info.Size() < minHeaderSize {
			return errors.NewCorruptedFileError(filePath, fmt.Errorf("file too small to be valid"))
		}
	}

	return nil
}
