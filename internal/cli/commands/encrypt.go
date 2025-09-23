package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/security"
)

// EncryptCmd represents the encrypt command
var EncryptCmd = &cobra.Command{
	Use:   "encrypt [file...]",
	Short: "🔐 Encrypt files using AES-256-GCM",
	Long: `Encrypt one or more files using secure AES-256-GCM authenticated encryption.

The encryption process:
  1. Prompts for a password (with strength checking)
  2. Derives encryption key using PBKDF2 (100,000 iterations)
  3. Generates random salt and IV for each file
  4. Encrypts with AES-256-GCM (provides authentication)
  5. Creates .enc file with custom FileVault format

SECURITY FEATURES:
  • Each file gets unique salt and IV
  • Password strength validation
  • Memory is securely cleaned after use
  • File integrity protection with authentication tags

PERFORMANCE:
  • Progress bars for files > 1MB
  • Optimized streaming for large files
  • Multi-file batch processing support`,
	Example: `  # Basic encryption
  filevault encrypt document.pdf

  # Encrypt to specific output file
  filevault encrypt document.pdf -o secure.enc

  # Encrypt multiple files
  filevault encrypt *.txt *.pdf

  # Encrypt to directory
  filevault encrypt file1.txt file2.pdf -o encrypted/

  # Keep original files after encryption
  filevault encrypt important.doc --keep

  # Custom PBKDF2 iterations for extra security
  filevault encrypt secret.txt --iterations 200000

  # Force overwrite existing files
  filevault encrypt data.xlsx -o backup.enc --force`,
	Args: cobra.MinimumNArgs(1),
	RunE: runEncrypt,
}

var (
	encryptOutput     string
	encryptForce      bool
	encryptKeep       bool
	encryptIterations int
)

func init() {
	EncryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "output file or directory")
	EncryptCmd.Flags().BoolVarP(&encryptForce, "force", "f", false, "overwrite existing files")
	EncryptCmd.Flags().BoolVarP(&encryptKeep, "keep", "k", false, "keep original file after encryption")
	EncryptCmd.Flags().IntVar(&encryptIterations, "iterations", 100000, "PBKDF2 iterations")
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	quiet, _ := cmd.Root().PersistentFlags().GetBool("quiet")

	// Enhanced batch processing
	if len(args) > 1 {
		return processBatchEncrypt(args, verbose, quiet)
	}

	// Single file processing
	return encryptSingleFile(args[0], verbose, quiet)
}

// processBatchEncrypt handles multiple file encryption
func processBatchEncrypt(files []string, verbose, quiet bool) error {
	if !quiet {
		cli.PrintInfo(fmt.Sprintf("Starting batch encryption of %d files", len(files)))
	}

	// Get password once for all files
	password, err := security.PromptPassword("Enter password for batch encryption: ")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Confirm password
	confirmPassword, err := security.PromptPassword("Confirm password: ")
	if err != nil {
		return fmt.Errorf("failed to get password confirmation: %w", err)
	}

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	successCount := 0
	failCount := 0

	for i, inputFile := range files {
		if verbose {
			cli.PrintProgress(fmt.Sprintf("Processing file %d/%d: %s", i+1, len(files), inputFile))
		}

		if err := encryptSingleFileWithPassword(inputFile, password, verbose, quiet); err != nil {
			if !quiet {
				cli.PrintError(fmt.Sprintf("Failed to encrypt %s: %v", inputFile, err))
			}
			failCount++
		} else {
			successCount++
		}
	}

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Batch encryption completed: %d success, %d failed", successCount, failCount))
	}

	if failCount > 0 {
		return fmt.Errorf("batch encryption had %d failures", failCount)
	}

	return nil
}

func encryptSingleFile(inputFile string, verbose, quiet bool) error {
	// Validate input file
	if err := security.ValidateInputFile(inputFile); err != nil {
		return err
	}

	// Get file info for progress tracking
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine output file
	outputFile := encryptOutput
	if outputFile == "" {
		outputFile = inputFile + ".enc"
	} else if info, err := os.Stat(outputFile); err == nil && info.IsDir() {
		outputFile = filepath.Join(outputFile, filepath.Base(inputFile)+".enc")
	}

	// Validate output file
	if err := security.ValidateOutputFile(outputFile, encryptForce); err != nil {
		return err
	}

	// Get password from user
	if verbose && !quiet {
		cli.PrintInfo("Getting password for encryption...")
	}

	password, err := security.PromptPassword("Enter password for encryption: ")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Confirm password
	confirmPassword, err := security.PromptPassword("Confirm password: ")
	if err != nil {
		return fmt.Errorf("failed to get password confirmation: %w", err)
	}

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Check password strength
	strength := security.CheckPasswordStrength(password)
	if strength == security.Weak && !encryptForce {
		if !quiet {
			cli.PrintWarning(fmt.Sprintf("Password strength is %s", strength))
			if !cli.ConfirmAction("Continue with weak password?") {
				return fmt.Errorf("encryption cancelled due to weak password")
			}
		}
	} else if verbose {
		cli.PrintInfo(fmt.Sprintf("Password strength: %s", strength))
	}

	// Show progress
	if verbose && !quiet {
		cli.PrintInfo(fmt.Sprintf("Encrypting %s -> %s", inputFile, outputFile))
		cli.PrintInfo(fmt.Sprintf("File size: %s", cli.FormatBytes(uint64(fileInfo.Size()))))
		cli.PrintInfo(fmt.Sprintf("Using PBKDF2 with %d iterations", encryptIterations))
	}

	// Create progress bar for larger files
	var progress *cli.ProgressBar
	if fileInfo.Size() > 1024*1024 && !quiet { // Show progress for files > 1MB
		progress = cli.NewProgressBar(fileInfo.Size(), "Encrypting")
	}

	// Perform encryption
	startTime := time.Now()
	if progress != nil {
		// Use progress callback
		err = core.EncryptFileWithProgress(inputFile, outputFile, password, func(current, total int64, operation string) {
			progress.Update(current)
		})
	} else {
		err = core.EncryptFile(inputFile, outputFile, password)
	}

	if err != nil {
		if progress != nil {
			progress.Finish()
		}
		return fmt.Errorf("encryption failed: %w", err)
	}

	if progress != nil {
		progress.Update(fileInfo.Size())
		progress.Finish()
	}

	elapsed := time.Since(startTime)

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Encrypted: %s -> %s", inputFile, outputFile))
		if verbose {
			cli.PrintInfo(fmt.Sprintf("Encryption completed in %s", cli.FormatDuration(elapsed.Seconds())))
		}
	}

	// Remove original file if not keeping
	if !encryptKeep {
		if err := os.Remove(inputFile); err != nil {
			if !quiet {
				cli.PrintWarning(fmt.Sprintf("Could not remove original file: %v", err))
			}
		} else if verbose {
			cli.PrintInfo("Original file removed")
		}
	}

	return nil
}

// encryptSingleFileWithPassword encrypts a file with pre-provided password
func encryptSingleFileWithPassword(inputFile, password string, verbose, quiet bool) error {
	// Validate input file
	if err := security.ValidateInputFile(inputFile); err != nil {
		return err
	}

	// Get file info for progress tracking
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine output file
	outputFile := encryptOutput
	if outputFile == "" {
		outputFile = inputFile + ".enc"
	} else if info, err := os.Stat(outputFile); err == nil && info.IsDir() {
		outputFile = filepath.Join(outputFile, filepath.Base(inputFile)+".enc")
	}

	// Validate output file
	if err := security.ValidateOutputFile(outputFile, encryptForce); err != nil {
		return err
	}

	// Show progress
	if verbose && !quiet {
		cli.PrintInfo(fmt.Sprintf("Encrypting %s -> %s", inputFile, outputFile))
	}

	// Create progress bar for larger files
	var progress *cli.ProgressBar
	if fileInfo.Size() > 1024*1024 && !quiet { // Show progress for files > 1MB
		progress = cli.NewProgressBar(fileInfo.Size(), "Encrypting")
	}

	// Perform encryption
	startTime := time.Now()
	if progress != nil {
		// Use progress callback
		err = core.EncryptFileWithProgress(inputFile, outputFile, password, func(current, total int64, operation string) {
			progress.Update(current)
		})
	} else {
		err = core.EncryptFile(inputFile, outputFile, password)
	}

	if err != nil {
		if progress != nil {
			progress.Finish()
		}
		return fmt.Errorf("encryption failed: %w", err)
	}

	if progress != nil {
		progress.Update(fileInfo.Size())
		progress.Finish()
	}

	elapsed := time.Since(startTime)

	if !quiet {
		cli.PrintSuccess(fmt.Sprintf("Encrypted: %s -> %s", inputFile, outputFile))
		if verbose {
			cli.PrintInfo(fmt.Sprintf("Encryption completed in %s", cli.FormatDuration(elapsed.Seconds())))
		}
	}

	// Remove original file if not keeping
	if !encryptKeep {
		if err := os.Remove(inputFile); err != nil {
			if !quiet {
				cli.PrintWarning(fmt.Sprintf("Could not remove original file: %v", err))
			}
		} else if verbose {
			cli.PrintInfo("Original file removed")
		}
	}

	return nil
}
