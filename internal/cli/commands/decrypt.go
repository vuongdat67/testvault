package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/security"
)

// DecryptCmd represents the decrypt command
var DecryptCmd = &cobra.Command{
	Use:   "decrypt [file...]",
	Short: "🔓 Decrypt FileVault encrypted files",
	Long: `Decrypt FileVault encrypted files (.enc) using the original password.

The decryption process:
  1. Validates the FileVault format and magic number
  2. Prompts for the encryption password
  3. Derives the decryption key using stored salt
  4. Verifies authentication tag for integrity
  5. Decrypts and restores the original file

SECURITY VERIFICATION:
  • Validates FileVault format signature
  • Checks file integrity with authentication tags
  • Verifies HMAC to detect tampering
  • Secure memory handling during decryption

PERFORMANCE:
  • Progress tracking for large files
  • Optimized streaming decryption
  • Batch processing for multiple files`,
	Example: `  # Basic decryption
  filevault decrypt document.pdf.enc

  # Decrypt to specific output file
  filevault decrypt encrypted.enc -o recovered.pdf

  # Decrypt multiple files
  filevault decrypt *.enc

  # Decrypt to directory
  filevault decrypt file1.enc file2.enc -o decrypted/

  # Force overwrite existing files
  filevault decrypt backup.enc -o original.txt --force

  # Batch decrypt all .enc files in directory
  filevault decrypt encrypted/*.enc -o restored/`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDecrypt,
}

var (
	decryptOutput string
	decryptForce  bool
)

func init() {
	DecryptCmd.Flags().StringVarP(&decryptOutput, "output", "o", "", "output file or directory")
	DecryptCmd.Flags().BoolVarP(&decryptForce, "force", "f", false, "overwrite existing files")
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	quiet, _ := cmd.Root().PersistentFlags().GetBool("quiet")

	// Enhanced batch processing
	if len(args) > 1 {
		return processBatchDecrypt(args, verbose, quiet)
	}

	// Single file processing
	return decryptSingleFile(args[0], verbose, quiet)
}

// processBatchDecrypt handles multiple file decryption
func processBatchDecrypt(files []string, verbose, quiet bool) error {
	if !quiet {
		cli.PrintInfo(fmt.Sprintf("Starting batch decryption of %d files", len(files)))
	}

	// Get password once for all files
	password, err := security.PromptPassword("Enter password for batch decryption: ")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	successCount := 0
	failCount := 0

	for i, inputFile := range files {
		if verbose {
			cli.PrintProgress(fmt.Sprintf("Processing file %d/%d: %s", i+1, len(files), inputFile))
		}

		if err := decryptSingleFileWithPassword(inputFile, password, verbose, quiet); err != nil {
			if !quiet {
				cli.PrintError(fmt.Sprintf("Failed to decrypt %s: %v", inputFile, err))
			}
			failCount++
		} else {
			successCount++
		}
	}

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Batch decryption completed: %d success, %d failed", successCount, failCount))
	}

	if failCount > 0 {
		return fmt.Errorf("batch decryption had %d failures", failCount)
	}

	return nil
}

func decryptSingleFile(inputFile string, verbose, quiet bool) error {
	// Validate input file
	if err := security.ValidateInputFile(inputFile); err != nil {
		return err
	}

	// Check if it's actually an encrypted file
	isEncrypted, err := security.IsEncryptedFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	if !isEncrypted {
		cli.PrintWarning("File doesn't appear to be a FileVault encrypted file")
		if !cli.ConfirmAction("Continue anyway?") {
			return fmt.Errorf("decryption cancelled")
		}
	}

	// Get file info for progress tracking
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine output file
	outputFile := decryptOutput
	if outputFile == "" {
		// Auto-determine output filename
		baseName := filepath.Base(inputFile)
		if strings.HasSuffix(baseName, ".enc") {
			outputFile = strings.TrimSuffix(inputFile, ".enc")
		} else {
			outputFile = inputFile + ".decrypted"
		}
	} else if info, err := os.Stat(outputFile); err == nil && info.IsDir() {
		baseName := filepath.Base(inputFile)
		if strings.HasSuffix(baseName, ".enc") {
			baseName = strings.TrimSuffix(baseName, ".enc")
		}
		outputFile = filepath.Join(outputFile, baseName)
	}

	// Validate output file
	if err := security.ValidateOutputFile(outputFile, decryptForce); err != nil {
		return err
	}

	// Get password from user
	if verbose && !quiet {
		cli.PrintInfo("Getting password for decryption...")
	}

	password, err := security.PromptPassword("Enter password for decryption: ")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Show progress
	if verbose && !quiet {
		cli.PrintInfo(fmt.Sprintf("Decrypting %s -> %s", inputFile, outputFile))
		cli.PrintInfo(fmt.Sprintf("File size: %s", cli.FormatBytes(uint64(fileInfo.Size()))))
	}

	// Create progress bar for larger files
	var progress *cli.ProgressBar
	if fileInfo.Size() > 1024*1024 && !quiet { // Show progress for files > 1MB
		progress = cli.NewProgressBar(fileInfo.Size(), "Decrypting")
	}

	// Perform decryption
	startTime := time.Now()
	if progress != nil {
		// Use progress callback
		err = core.DecryptFileWithProgress(inputFile, outputFile, password, func(current, total int64, operation string) {
			// Convert percentage-based progress to file-size based
			actualProgress := (current * fileInfo.Size()) / total
			progress.Update(actualProgress)
		})
	} else {
		err = core.DecryptFile(inputFile, outputFile, password)
	}

	if err != nil {
		if progress != nil {
			progress.Finish()
		}
		if strings.Contains(err.Error(), "authentication failed") || strings.Contains(err.Error(), "decryption failed") {
			cli.PrintError("Decryption failed - wrong password or corrupted file")
		}
		return fmt.Errorf("decryption failed: %w", err)
	}

	if progress != nil {
		progress.Update(fileInfo.Size())
		progress.Finish()
	}

	elapsed := time.Since(startTime)

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Decrypted: %s -> %s", inputFile, outputFile))
		if verbose {
			cli.PrintInfo(fmt.Sprintf("Decryption completed in %s", cli.FormatDuration(elapsed.Seconds())))
		}
	}

	return nil
}

// decryptSingleFileWithPassword decrypts a file with pre-provided password
func decryptSingleFileWithPassword(inputFile, password string, verbose, quiet bool) error {
	// Validate input file
	if err := security.ValidateInputFile(inputFile); err != nil {
		return err
	}

	// Check if it's actually an encrypted file
	isEncrypted, err := security.IsEncryptedFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	if !isEncrypted {
		return fmt.Errorf("file doesn't appear to be a FileVault encrypted file")
	}

	// Get file info for progress tracking
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine output file
	outputFile := decryptOutput
	if outputFile == "" {
		// Auto-determine output filename
		baseName := filepath.Base(inputFile)
		if strings.HasSuffix(baseName, ".enc") {
			outputFile = strings.TrimSuffix(inputFile, ".enc")
		} else {
			outputFile = inputFile + ".decrypted"
		}
	} else if info, err := os.Stat(outputFile); err == nil && info.IsDir() {
		baseName := filepath.Base(inputFile)
		if strings.HasSuffix(baseName, ".enc") {
			baseName = strings.TrimSuffix(baseName, ".enc")
		}
		outputFile = filepath.Join(outputFile, baseName)
	}

	// Validate output file
	if err := security.ValidateOutputFile(outputFile, decryptForce); err != nil {
		return err
	}

	// Show progress
	if verbose && !quiet {
		cli.PrintInfo(fmt.Sprintf("Decrypting %s -> %s", inputFile, outputFile))
	}

	// Create progress bar for larger files
	var progress *cli.ProgressBar
	if fileInfo.Size() > 1024*1024 && !quiet { // Show progress for files > 1MB
		progress = cli.NewProgressBar(fileInfo.Size(), "Decrypting")
	}

	// Perform decryption
	startTime := time.Now()
	if progress != nil {
		// Use progress callback
		err = core.DecryptFileWithProgress(inputFile, outputFile, password, func(current, total int64, operation string) {
			// Convert percentage-based progress to file-size based
			actualProgress := (current * fileInfo.Size()) / total
			progress.Update(actualProgress)
		})
	} else {
		err = core.DecryptFile(inputFile, outputFile, password)
	}

	if err != nil {
		if progress != nil {
			progress.Finish()
		}
		return fmt.Errorf("decryption failed: %w", err)
	}

	if progress != nil {
		progress.Update(fileInfo.Size())
		progress.Finish()
	}

	elapsed := time.Since(startTime)

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Decrypted: %s -> %s", inputFile, outputFile))
		if verbose {
			cli.PrintInfo(fmt.Sprintf("Decryption completed in %s", cli.FormatDuration(elapsed.Seconds())))
		}
	}

	return nil
}
