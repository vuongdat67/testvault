// Package filevault provides a simple client interface for file encryption and decryption operations.
// This package wraps the internal FileVault functionality in a clean, easy-to-use API.
package filevault

import (
	"fmt"
	"path/filepath"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/security"
)

// Client represents the main FileVault client for encryption/decryption operations
type Client struct {
	// Configuration options for the client
	verbose bool
}

// ClientOption represents configuration options for the FileVault client
type ClientOption func(*Client)

// NewClient creates a new FileVault client with optional configuration
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		verbose: false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithVerbose enables verbose logging for operations
func WithVerbose(verbose bool) ClientOption {
	return func(c *Client) {
		c.verbose = verbose
	}
}

// EncryptFile encrypts a file using AES-256-GCM with the provided password
func (c *Client) EncryptFile(inputPath, password string) error {
	return c.EncryptFileWithOutput(inputPath, "", password)
}

// EncryptFileWithOutput encrypts a file with a custom output path
func (c *Client) EncryptFileWithOutput(inputPath, outputPath, password string) error {
	// Validate password strength
	if err := security.ValidatePasswordBasic(password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Generate default output path if not provided
	if outputPath == "" {
		outputPath = inputPath + ".enc"
	}

	// Validate input file
	if err := security.ValidateInputFile(inputPath); err != nil {
		return fmt.Errorf("input file validation failed: %w", err)
	}

	// Check if output file already exists
	if err := security.ValidateOutputFile(outputPath, false); err != nil {
		return fmt.Errorf("output validation failed: %w", err)
	}

	// Perform encryption
	if c.verbose {
		fmt.Printf("Encrypting: %s -> %s\n", inputPath, outputPath)
	}

	return core.EncryptFile(inputPath, outputPath, password)
}

// DecryptFile decrypts a FileVault encrypted file using the provided password
func (c *Client) DecryptFile(encryptedPath, password string) error {
	return c.DecryptFileWithOutput(encryptedPath, "", password)
}

// DecryptFileWithOutput decrypts a file with a custom output path
func (c *Client) DecryptFileWithOutput(encryptedPath, outputPath, password string) error {
	// Validate input file
	if err := security.ValidateEncryptedFile(encryptedPath); err != nil {
		return fmt.Errorf("encrypted file validation failed: %w", err)
	}

	// Generate default output path if not provided
	if outputPath == "" {
		outputPath = c.getOriginalFilename(encryptedPath)
	}

	// Check if output file already exists
	if err := security.ValidateOutputFile(outputPath, false); err != nil {
		return fmt.Errorf("output validation failed: %w", err)
	}

	// Perform decryption
	if c.verbose {
		fmt.Printf("Decrypting: %s -> %s\n", encryptedPath, outputPath)
	}

	return core.DecryptFile(encryptedPath, outputPath, password)
}

// VerifyFile checks the integrity and format of an encrypted file
func (c *Client) VerifyFile(encryptedPath string) (*VerificationResult, error) {
	if err := security.ValidateEncryptedFile(encryptedPath); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	if c.verbose {
		fmt.Printf("Verifying file: %s\n", encryptedPath)
	}

	coreResult, err := core.VerifyFile(encryptedPath)
	if err != nil {
		return nil, err
	}

	// Convert from core.VerificationResult to our VerificationResult
	result := &VerificationResult{
		Valid:            coreResult.IsValid,
		FormatValid:      coreResult.FormatValid,
		HeaderValid:      coreResult.HeaderValid,
		SizeConsistent:   coreResult.SizeConsistent,
		FileAccessible:   coreResult.FileAccessible,
		Filename:         coreResult.Filename,
		OriginalFilename: coreResult.OriginalFilename,
		FileSize:         coreResult.FileSize,
		OriginalSize:     coreResult.OriginalSize,
		Algorithm:        coreResult.Algorithm,
		FormatVersion:    coreResult.FormatVersion,
		ErrorMessage:     coreResult.ErrorMessage,
	}

	return result, nil
}

// VerificationResult contains the result of file verification
type VerificationResult struct {
	Valid            bool   `json:"valid"`
	FormatValid      bool   `json:"format_valid"`
	HeaderValid      bool   `json:"header_valid"`
	SizeConsistent   bool   `json:"size_consistent"`
	FileAccessible   bool   `json:"file_accessible"`
	Filename         string `json:"filename"`
	OriginalFilename string `json:"original_filename"`
	FileSize         int64  `json:"file_size"`
	OriginalSize     uint64 `json:"original_size"`
	Algorithm        string `json:"algorithm"`
	FormatVersion    uint32 `json:"format_version"`
	ErrorMessage     string `json:"error_message"`
}

// IsValid returns true if the file passed all verification checks
func (vr *VerificationResult) IsValid() bool {
	return vr.Valid && vr.FormatValid && vr.HeaderValid && vr.SizeConsistent && vr.FileAccessible
}

// GetErrorMessage returns the verification error message if any
func (vr *VerificationResult) GetErrorMessage() string {
	return vr.ErrorMessage
}

// getOriginalFilename attempts to determine the original filename from an encrypted file
func (c *Client) getOriginalFilename(encryptedPath string) string {
	// Fallback: remove .enc extension if present
	if filepath.Ext(encryptedPath) == ".enc" {
		return encryptedPath[:len(encryptedPath)-4]
	}

	// Final fallback: add .dec suffix
	return encryptedPath + ".dec"
}
