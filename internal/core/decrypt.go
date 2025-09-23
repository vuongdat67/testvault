package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/crypto"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/fileops"
)

// DecryptFile decrypts a FileVault encrypted file
func DecryptFile(inputPath, outputPath, password string) error {
	return DecryptFileWithProgress(inputPath, outputPath, password, nil)
}

// DecryptFileWithProgress decrypts a file with progress reporting
func DecryptFileWithProgress(inputPath, outputPath, password string, progressCallback ProgressCallback) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Report initial progress
	if progressCallback != nil {
		progressCallback(0, 100, "Reading file header")
	}

	// Read and validate header
	var header fileops.FileHeader
	_, err = header.ReadFrom(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if err := header.IsValid(); err != nil {
		return fmt.Errorf("invalid file format: %w", err)
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(10, 100, "Validating file format")
	}

	// Determine output path if not specified
	if outputPath == "" {
		outputPath = header.GetBaseFileName()
		if outputPath == "" {
			// Fallback: remove .enc extension
			baseName := filepath.Base(inputPath)
			if filepath.Ext(baseName) == ".enc" {
				outputPath = baseName[:len(baseName)-4]
			} else {
				outputPath = baseName + ".decrypted"
			}
		}
	}

	// Create AES cipher from password and salt
	cipher, err := crypto.NewAESCipherFromPassword(password, header.Salt)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(20, 100, "Deriving decryption key")
	}

	// Calculate encrypted data size (total - header - auth tag)
	inputInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get input file info: %w", err)
	}

	encryptedDataSize := int64(inputInfo.Size()) - int64(header.GetTotalSize()) - fileops.AuthTagSize

	// Report progress
	if progressCallback != nil {
		progressCallback(30, 100, "Reading encrypted data")
	}

	// Read encrypted data
	encryptedData := make([]byte, encryptedDataSize)
	_, err = io.ReadFull(inputFile, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to read encrypted data: %w", err)
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(50, 100, "Reading authentication tag")
	}

	// Read authentication tag
	authTag := make([]byte, fileops.AuthTagSize)
	_, err = io.ReadFull(inputFile, authTag)
	if err != nil {
		return fmt.Errorf("failed to read auth tag: %w", err)
	}

	// Create encrypted data structure
	cryptoData := &crypto.EncryptedData{
		Nonce:      header.IV[:12], // Use first 12 bytes of IV as nonce
		Ciphertext: encryptedData,
		Tag:        authTag,
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(70, 100, "Decrypting data")
	}

	// Decrypt
	plaintext, err := cipher.Decrypt(cryptoData)
	if err != nil {
		return fmt.Errorf("decryption failed (wrong password or corrupted file): %w", err)
	}

	// Verify original size
	if uint64(len(plaintext)) != header.OriginalSize {
		return fmt.Errorf("decrypted size mismatch: expected %d, got %d", header.OriginalSize, len(plaintext))
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(90, 100, "Writing decrypted file")
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write decrypted data
	_, err = outputFile.Write(plaintext)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data: %w", err)
	}

	// Report completion
	if progressCallback != nil {
		progressCallback(100, 100, "Decryption completed")
	}

	// Secure cleanup
	crypto.SecureZero(plaintext)
	crypto.SecureZero(encryptedData)

	return nil
}
