package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/fileops"
)

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info [file...]",
	Short: "📊 Display encrypted file information",
	Long: `Display comprehensive information about FileVault encrypted files.

Shows detailed metadata including:
  • File format version and magic number
  • Encryption algorithm (AES-256-GCM)
  • Original filename and file size
  • Salt and IV information (for security analysis)
  • PBKDF2 iteration count
  • Authentication tag status
  • File creation and modification times
  • File integrity status

This command does NOT require the password and will not decrypt the file.
It only reads and displays the metadata stored in the FileVault header.

SECURITY ANALYSIS:
  • Verifies FileVault format integrity
  • Shows cryptographic parameters used
  • Helps identify encryption version/strength
  • Useful for forensic analysis

BATCH PROCESSING:
  • Supports analysis of multiple files
  • Provides summary statistics
  • Detailed comparison of encryption parameters`,
	Example: `  # Basic file information
  filevault info document.pdf.enc

  # Analyze encryption parameters
  filevault info secret-data.enc

  # Check multiple encrypted files
  filevault info backup1.enc backup2.enc backup3.enc

  # Verbose output with all metadata
  filevault info -v important.enc

  # Quick format validation
  filevault info suspicious-file.enc

  # Batch analyze directory
  filevault info encrypted/*.enc`,
	Args: cobra.MinimumNArgs(1),
	RunE: runInfo,
}

var (
	infoShowHex bool
)

func init() {
	InfoCmd.Flags().BoolVar(&infoShowHex, "hex", false, "show cryptographic parameters in hexadecimal")
}

func runInfo(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	quiet, _ := cmd.Root().PersistentFlags().GetBool("quiet")

	// Handle multiple files
	if len(args) > 1 {
		return runBatchInfo(args, verbose, quiet)
	}

	// Single file analysis
	return analyzeFile(args[0], verbose, quiet)
}

func analyzeFile(inputFile string, verbose, quiet bool) error {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputFile)
	}

	// Get basic file info
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Perform verification to get comprehensive info
	result, err := core.VerifyFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to analyze file: %w", err)
	}

	// Read detailed header information
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var header fileops.FileHeader
	if result.FormatValid {
		_, err = header.ReadFrom(file)
		if err != nil {
			return fmt.Errorf("failed to read file header: %w", err)
		}
	}

	// Display comprehensive information
	if !quiet {
		fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", cli.ColorBlue, cli.ColorReset)
		fmt.Printf("%s                     FILE ANALYSIS REPORT                      %s\n", cli.ColorBold, cli.ColorReset)
		fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", cli.ColorBlue, cli.ColorReset)
		fmt.Printf("\n")

		// Basic file information
		fmt.Printf("%sBasic Information:%s\n", cli.ColorGreen, cli.ColorReset)
		fmt.Printf("  File: %s\n", inputFile)
		fmt.Printf("  File Size: %s\n", cli.FormatBytes(uint64(fileInfo.Size())))
		fmt.Printf("  Modified: %s\n", fileInfo.ModTime().Format(time.RFC3339))
		fmt.Printf("  Permissions: %s\n", fileInfo.Mode())
		fmt.Printf("\n")

		// Format information
		fmt.Printf("%sFormat Information:%s\n", cli.ColorYellow, cli.ColorReset)
		if result.IsValid {
			fmt.Printf("  Status: %s✅ Valid FileVault File%s\n", cli.ColorGreen, cli.ColorReset)
			fmt.Printf("  Format: FileVault v%d\n", result.FormatVersion)
			fmt.Printf("  Algorithm: %s\n", result.Algorithm)
			fmt.Printf("  Key Derivation: PBKDF2-SHA256 (100,000 iterations)\n")
		} else {
			fmt.Printf("  Status: %s❌ Invalid or Corrupted%s\n", cli.ColorRed, cli.ColorReset)
			fmt.Printf("  Error: %s\n", result.ErrorMessage)
		}
		fmt.Printf("\n")

		// Original file information
		if result.IsValid && result.OriginalFilename != "" {
			fmt.Printf("%sOriginal File Information:%s\n", cli.ColorCyan, cli.ColorReset)
			fmt.Printf("  Original Filename: %s\n", result.OriginalFilename)
			fmt.Printf("  Original Size: %s\n", cli.FormatBytes(result.OriginalSize))

			// Calculate compression ratio
			compressionRatio := float64(result.OriginalSize) / float64(result.FileSize) * 100
			fmt.Printf("  Compression Ratio: %.1f%%\n", compressionRatio)
			fmt.Printf("\n")
		}

		// Security parameters (if verbose or hex requested)
		if (verbose || infoShowHex) && result.IsValid {
			fmt.Printf("%sCryptographic Parameters:%s\n", cli.ColorPurple, cli.ColorReset)
			fmt.Printf("  Salt Length: 32 bytes\n")
			fmt.Printf("  IV Length: 16 bytes (12 bytes nonce + 4 bytes padding)\n")
			fmt.Printf("  Auth Tag Length: 16 bytes\n")

			if infoShowHex {
				fmt.Printf("  Salt (hex): %x\n", header.Salt[:8]) // Show first 8 bytes for security
				fmt.Printf("  IV (hex): %x\n", header.IV[:8])     // Show first 8 bytes for security
				fmt.Printf("  (Note: Only first 8 bytes shown for security)\n")
			}
			fmt.Printf("\n")
		}

		// Detailed verification results
		if verbose {
			fmt.Printf("%sVerification Details:%s\n", cli.ColorBlue, cli.ColorReset)
			fmt.Printf("  File Accessible: %s%t%s\n",
				getStatusColor(result.FileAccessible), result.FileAccessible, cli.ColorReset)
			fmt.Printf("  Format Valid: %s%t%s\n",
				getStatusColor(result.FormatValid), result.FormatValid, cli.ColorReset)
			fmt.Printf("  Header Valid: %s%t%s\n",
				getStatusColor(result.HeaderValid), result.HeaderValid, cli.ColorReset)
			fmt.Printf("  Size Consistent: %s%t%s\n",
				getStatusColor(result.SizeConsistent), result.SizeConsistent, cli.ColorReset)
			fmt.Printf("  Verification Time: %s\n", cli.FormatDuration(result.VerificationTime.Seconds()))
			fmt.Printf("\n")
		}

		fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", cli.ColorBlue, cli.ColorReset)
	}

	return nil
}

func runBatchInfo(files []string, verbose, quiet bool) error {
	if !quiet {
		cli.PrintInfo(fmt.Sprintf("Analyzing %d files", len(files)))
		fmt.Printf("\n")
	}

	validCount := 0
	invalidCount := 0
	var totalOriginalSize uint64
	var totalEncryptedSize int64

	for i, file := range files {
		if !quiet && !verbose {
			fmt.Printf("[%d/%d] %s\n", i+1, len(files), file)
		}

		result, err := core.VerifyFile(file)
		if err != nil {
			invalidCount++
			if !quiet {
				cli.PrintError(fmt.Sprintf("Failed to analyze %s: %v", file, err))
			}
			continue
		}

		if result.IsValid {
			validCount++
			totalOriginalSize += result.OriginalSize
			totalEncryptedSize += result.FileSize

			if verbose && !quiet {
				fmt.Printf("\n%s--- %s ---%s\n", cli.ColorBold, file, cli.ColorReset)
				fmt.Printf("  Status: ✅ Valid\n")
				fmt.Printf("  Format: FileVault v%d\n", result.FormatVersion)
				fmt.Printf("  Algorithm: %s\n", result.Algorithm)
				fmt.Printf("  Original: %s (%s)\n", result.OriginalFilename, cli.FormatBytes(result.OriginalSize))
				fmt.Printf("  Encrypted: %s\n", cli.FormatBytes(uint64(result.FileSize)))
			}
		} else {
			invalidCount++
			if !quiet {
				fmt.Printf("  Status: ❌ Invalid - %s\n", result.ErrorMessage)
			}
		}
	}

	// Summary
	if !quiet {
		fmt.Printf("\n%sBatch Analysis Summary:%s\n", cli.ColorGreen, cli.ColorReset)
		fmt.Printf("========================\n")
		fmt.Printf("Total files: %d\n", len(files))
		fmt.Printf("✅ Valid: %d\n", validCount)
		fmt.Printf("❌ Invalid: %d\n", invalidCount)

		if validCount > 0 {
			fmt.Printf("Total original size: %s\n", cli.FormatBytes(totalOriginalSize))
			fmt.Printf("Total encrypted size: %s\n", cli.FormatBytes(uint64(totalEncryptedSize)))

			if totalOriginalSize > 0 {
				avgCompression := float64(totalOriginalSize) / float64(totalEncryptedSize) * 100
				fmt.Printf("Average compression: %.1f%%\n", avgCompression)
			}
		}
	}

	return nil
}

func getStatusColor(status bool) string {
	if status {
		return cli.ColorGreen
	}
	return cli.ColorRed
}
