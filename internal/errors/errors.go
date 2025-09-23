package errors

import (
	"fmt"
)

// ErrorCode represents different types of errors
type ErrorCode int

const (
	// General errors
	ErrUnknown ErrorCode = iota
	ErrInvalidArguments
	ErrInvalidConfig
	
	// File operation errors
	ErrFileNotFound
	ErrFileAlreadyExists
	ErrFilePermissionDenied
	ErrFileReadError
	ErrFileWriteError
	ErrFileCorrupted
	ErrFileTooLarge
	
	// Cryptographic errors
	ErrInvalidPassword
	ErrWeakPassword
	ErrAuthenticationFailed
	ErrKeyDerivationFailed
	ErrEncryptionFailed
	ErrDecryptionFailed
	ErrInvalidFormat
	ErrUnsupportedVersion
	
	// Security errors
	ErrInvalidInput
	ErrSecurityViolation
	ErrMemoryError
)

// FileVaultError represents a structured error with context
type FileVaultError struct {
	Code        ErrorCode
	Message     string
	Cause       error
	Context     map[string]interface{}
	Suggestions []string
}

// Error implements the error interface
func (e *FileVaultError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *FileVaultError) Unwrap() error {
	return e.Cause
}

// GetExitCode returns the appropriate exit code for the error
func (e *FileVaultError) GetExitCode() int {
	switch e.Code {
	case ErrFileNotFound:
		return 2
	case ErrFilePermissionDenied:
		return 3
	case ErrInvalidPassword, ErrAuthenticationFailed:
		return 4
	case ErrFileCorrupted, ErrInvalidFormat:
		return 5
	case ErrFileTooLarge, ErrMemoryError:
		return 6
	case ErrInvalidArguments, ErrInvalidInput:
		return 7
	default:
		return 1
	}
}

// GetUserFriendlyMessage returns a user-friendly error message
func (e *FileVaultError) GetUserFriendlyMessage() string {
	switch e.Code {
	case ErrFileNotFound:
		return "File not found. Please check the file path."
	case ErrFileAlreadyExists:
		return "Output file already exists. Use --force to overwrite."
	case ErrFilePermissionDenied:
		return "Permission denied. Please check file permissions."
	case ErrInvalidPassword:
		return "Authentication failed. Please check your password."
	case ErrWeakPassword:
		return "Password is too weak. Please use a stronger password."
	case ErrFileCorrupted:
		return "File appears to be corrupted or damaged."
	case ErrInvalidFormat:
		return "File is not a valid FileVault encrypted file."
	case ErrUnsupportedVersion:
		return "Unsupported file version. Please update FileVault."
	default:
		return e.Message
	}
}

// GetSuggestions returns helpful suggestions for resolving the error
func (e *FileVaultError) GetSuggestions() []string {
	if len(e.Suggestions) > 0 {
		return e.Suggestions
	}
	
	switch e.Code {
	case ErrFileNotFound:
		return []string{
			"Verify the file path is correct",
			"Check if the file exists in the current directory",
			"Use absolute path if relative path doesn't work",
		}
	case ErrFileAlreadyExists:
		return []string{
			"Use --force flag to overwrite existing file",
			"Choose a different output filename",
			"Remove the existing file manually",
		}
	case ErrInvalidPassword:
		return []string{
			"Make sure you're using the correct password",
			"Check for typos in the password",
			"Verify the file hasn't been corrupted",
		}
	case ErrWeakPassword:
		return []string{
			"Use at least 12 characters",
			"Include uppercase and lowercase letters",
			"Add numbers and special characters",
			"Use --force to proceed with weak password (not recommended)",
		}
	case ErrInvalidFormat:
		return []string{
			"Verify the file is encrypted with FileVault",
			"Check if the file has been corrupted",
			"Make sure you're using the correct file",
		}
	default:
		return []string{}
	}
}

// NewError creates a new FileVaultError
func NewError(code ErrorCode, message string, cause error) *FileVaultError {
	return &FileVaultError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewFileNotFoundError creates a file not found error
func NewFileNotFoundError(filename string) *FileVaultError {
	err := NewError(ErrFileNotFound, fmt.Sprintf("file not found: %s", filename), nil)
	err.Context["filename"] = filename
	return err
}

// NewPermissionDeniedError creates a permission denied error
func NewPermissionDeniedError(filename string, cause error) *FileVaultError {
	err := NewError(ErrFilePermissionDenied, fmt.Sprintf("permission denied: %s", filename), cause)
	err.Context["filename"] = filename
	return err
}

// NewInvalidPasswordError creates an invalid password error
func NewInvalidPasswordError(cause error) *FileVaultError {
	return NewError(ErrInvalidPassword, "authentication failed: invalid password", cause)
}

// NewWeakPasswordError creates a weak password error
func NewWeakPasswordError(strength string) *FileVaultError {
	err := NewError(ErrWeakPassword, fmt.Sprintf("password is too weak: %s", strength), nil)
	err.Context["strength"] = strength
	return err
}

// NewCorruptedFileError creates a corrupted file error
func NewCorruptedFileError(filename string, cause error) *FileVaultError {
	err := NewError(ErrFileCorrupted, fmt.Sprintf("file appears to be corrupted: %s", filename), cause)
	err.Context["filename"] = filename
	return err
}

// NewInvalidFormatError creates an invalid format error
func NewInvalidFormatError(filename string) *FileVaultError {
	err := NewError(ErrInvalidFormat, fmt.Sprintf("invalid file format: %s", filename), nil)
	err.Context["filename"] = filename
	return err
}

// WrapError wraps an existing error with FileVaultError context
func WrapError(code ErrorCode, message string, cause error) *FileVaultError {
	return NewError(code, message, cause)
}

// IsFileVaultError checks if an error is a FileVaultError
func IsFileVaultError(err error) bool {
	_, ok := err.(*FileVaultError)
	return ok
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if fvErr, ok := err.(*FileVaultError); ok {
		return fvErr.Code
	}
	return ErrUnknown
}

// HandleError provides centralized error handling and user feedback
func HandleError(err error, quiet bool) int {
	if err == nil {
		return 0
	}

	if fvErr, ok := err.(*FileVaultError); ok {
		if !quiet {
			fmt.Printf("❌ %s\n", fvErr.GetUserFriendlyMessage())
			
			suggestions := fvErr.GetSuggestions()
			if len(suggestions) > 0 {
				fmt.Println("\nSuggestions:")
				for _, suggestion := range suggestions {
					fmt.Printf("  • %s\n", suggestion)
				}
			}
		}
		return fvErr.GetExitCode()
	}

	// Handle regular errors
	if !quiet {
		fmt.Printf("❌ Error: %v\n", err)
	}
	return 1
}