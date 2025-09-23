package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/crypto"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/fileops"
)

// ProgressCallback is a function type for progress updates
type ProgressCallback func(current, total int64, operation string)

// EncryptFile encrypts a file using AES-256-GCM with PBKDF2 key derivation
func EncryptFile(inputPath, outputPath, password string) error {
	return EncryptFileWithProgress(inputPath, outputPath, password, nil)
}

// EncryptFileWithProgress encrypts a file with progress reporting
func EncryptFileWithProgress(inputPath, outputPath, password string, progressCallback ProgressCallback) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Get input file info
	inputInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get input file info: %w", err)
	}

	// Generate cryptographic parameters
	salt, err := crypto.GenerateSalt32()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	iv, err := crypto.GenerateIV16()
	if err != nil {
		return fmt.Errorf("failed to generate IV: %w", err)
	}

	// Create file header
	originalFileName := filepath.Base(inputPath)
	header := fileops.NewFileHeader(uint64(inputInfo.Size()), originalFileName, salt, iv)

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write header
	_, err = header.WriteTo(outputFile)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Create AES cipher from password
	cipher, err := crypto.NewAESCipherFromPassword(password, salt)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// For small files, read all at once
	if inputInfo.Size() <= 64*1024*1024 { // 64MB threshold
		return encryptSmallFile(inputFile, outputFile, cipher, iv, inputInfo.Size(), progressCallback)
	}

	// For large files, use streaming encryption
	return encryptLargeFile(inputFile, outputFile, cipher, iv, inputInfo.Size(), progressCallback)
}

// encryptSmallFile encrypts smaller files in one go
func encryptSmallFile(inputFile, outputFile *os.File, cipher *crypto.AESCipher, iv [16]byte, fileSize int64, progressCallback ProgressCallback) error {
	// Report initial progress
	if progressCallback != nil {
		progressCallback(0, fileSize, "Reading file")
	}

	// Read entire file content
	plaintext, err := io.ReadAll(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Report progress after reading
	if progressCallback != nil {
		progressCallback(fileSize/2, fileSize, "Encrypting")
	}

	// Encrypt the content using the nonce from header
	nonce := iv[:12] // Use first 12 bytes of IV as nonce
	encryptedData, err := cipher.EncryptWithNonce(plaintext, nonce)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	// Report progress after encryption
	if progressCallback != nil {
		progressCallback(fileSize*3/4, fileSize, "Writing encrypted data")
	}

	// Write encrypted data
	_, err = outputFile.Write(encryptedData.Ciphertext)
	if err != nil {
		return fmt.Errorf("failed to write encrypted data: %w", err)
	}

	// Write authentication tag at the end
	_, err = outputFile.Write(encryptedData.Tag)
	if err != nil {
		return fmt.Errorf("failed to write auth tag: %w", err)
	}

	// Report completion
	if progressCallback != nil {
		progressCallback(fileSize, fileSize, "Encryption completed")
	}

	// Secure cleanup
	crypto.SecureZero(plaintext)

	return nil
}

// encryptLargeFile encrypts large files with streaming
func encryptLargeFile(inputFile, outputFile *os.File, cipher *crypto.AESCipher, iv [16]byte, fileSize int64, progressCallback ProgressCallback) error {
	const chunkSize = 64 * 1024 // 64KB chunks
	
	// For now, fallback to small file method
	// TODO: Implement proper streaming encryption in future sprints
	return encryptSmallFile(inputFile, outputFile, cipher, iv, fileSize, progressCallback)
}
