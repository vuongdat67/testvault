package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/core"
)

// VerifyCmd represents the verify command
var VerifyCmd = &cobra.Command{
	Use:   "verify [file...]",
	Short: "🔍 Verify encrypted file integrity",
	Long: `Verify the integrity and format of FileVault encrypted files.

This command performs comprehensive validation WITHOUT requiring passwords:
  • Validates FileVault magic number ("FVLT")
  • Checks file format version compatibility  
  • Verifies header structure and fields
  • Validates salt and IV lengths
  • Checks authentication tag presence
  • Confirms file size consistency
  • Tests read accessibility

INTEGRITY CHECKS:
  • File format corruption detection
  • Header field validation
  • Cryptographic parameter verification
  • File system consistency checks

This is useful for:
  • Batch verification of backup files
  • Detecting file corruption or tampering
  • Validating FileVault format compliance
  • Pre-decryption integrity assessment

BATCH PROCESSING:
  • Supports multiple files in one command
  • Provides summary statistics for batch operations
  • Detailed per-file results available with -v flag`,
	Example: `  # Verify single encrypted file
  filevault verify document.pdf.enc

  # Verify multiple files
  filevault verify *.enc

  # Batch verify backup directory
  filevault verify backups/*.enc

  # Verify with detailed output
  filevault verify -v important.enc

  # Quick integrity check
  filevault verify -q suspicious.enc

  # Verify all files in directory
  filevault verify encrypted-data/*`,
	Args: cobra.MinimumNArgs(1),
	RunE: runVerify,
}

var (
	verifyDeep bool
)

func init() {
	VerifyCmd.Flags().BoolVar(&verifyDeep, "deep", false, "perform deep integrity verification (requires password)")
}

func runVerify(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	quiet, _ := cmd.Root().PersistentFlags().GetBool("quiet")

	// Handle batch verification
	if len(args) > 1 {
		return runBatchVerify(args, verbose, quiet)
	}

	// Single file verification
	return verifySingleFile(args[0], verbose, quiet)
}

func verifySingleFile(inputFile string, verbose, quiet bool) error {
	// Check if input file exists first
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputFile)
	}

	// Perform verification
	result, err := core.VerifyFile(inputFile)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	// Display results
	if result.IsValid {
		if !quiet {
			cli.PrintSuccess("File verification successful")
			if verbose {
				fmt.Printf("   Format: FileVault v%d\n", result.FormatVersion)
				fmt.Printf("   Algorithm: %s\n", result.Algorithm)
				fmt.Printf("   Original file: %s (%s)\n",
					result.OriginalFilename,
					cli.FormatBytes(result.OriginalSize))
				fmt.Printf("   Encrypted size: %s\n", cli.FormatBytes(uint64(result.FileSize)))
				fmt.Printf("   Verification time: %s\n", cli.FormatDuration(result.VerificationTime.Seconds()))
			}
		}
	} else {
		if !quiet {
			cli.PrintError(fmt.Sprintf("File verification failed: %s", result.ErrorMessage))
			if verbose {
				fmt.Printf("   File accessible: %t\n", result.FileAccessible)
				fmt.Printf("   Format valid: %t\n", result.FormatValid)
				fmt.Printf("   Header valid: %t\n", result.HeaderValid)
				fmt.Printf("   Size consistent: %t\n", result.SizeConsistent)
			}
		}
		return fmt.Errorf("verification failed: %s", result.ErrorMessage)
	}

	return nil
}

func runBatchVerify(files []string, verbose, quiet bool) error {
	if !quiet {
		cli.PrintInfo(fmt.Sprintf("Starting batch verification of %d files", len(files)))
	}

	// Perform batch verification
	results, err := core.BatchVerify(files)
	if err != nil {
		return fmt.Errorf("batch verification failed: %w", err)
	}

	// Calculate summary
	summary := core.GetVerificationSummary(results)

	// Display individual results if verbose
	if verbose && !quiet {
		fmt.Printf("\nDetailed Results:\n")
		fmt.Printf("================\n")
		for _, result := range results {
			status := "❌ FAILED"
			if result.IsValid {
				status = "✅ PASSED"
			}

			fmt.Printf("%-50s %s\n", result.Filename, status)
			if !result.IsValid && result.ErrorMessage != "" {
				fmt.Printf("    Error: %s\n", result.ErrorMessage)
			} else if result.IsValid {
				fmt.Printf("    Format: FileVault v%d, %s\n", result.FormatVersion, result.Algorithm)
				if result.OriginalFilename != "" {
					fmt.Printf("    Original: %s (%s)\n",
						result.OriginalFilename,
						cli.FormatBytes(result.OriginalSize))
				}
			}
		}
	}

	// Display summary
	if !quiet {
		fmt.Printf("\nVerification Summary:\n")
		fmt.Printf("====================\n")
		fmt.Printf("Total files: %d\n", summary["total"])
		fmt.Printf("✅ Valid: %d\n", summary["valid"])
		fmt.Printf("❌ Invalid: %d\n", summary["invalid"])

		if verbose {
			fmt.Printf("📁 Accessible: %d\n", summary["accessible"])
			fmt.Printf("📄 Format OK: %d\n", summary["format_ok"])
			fmt.Printf("📋 Header OK: %d\n", summary["header_ok"])
			fmt.Printf("📏 Size OK: %d\n", summary["size_ok"])
		}
	}

	// Return error if any files failed
	if summary["invalid"] > 0 {
		return fmt.Errorf("verification failed for %d out of %d files", summary["invalid"], summary["total"])
	}

	if !quiet {
		cli.PrintSuccess("All files verified successfully")
	}

	return nil
}
