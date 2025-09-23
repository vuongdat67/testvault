package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/cli/commands"
	"github.com/vuongdat67/NT140.Q11.ANTT-Group15/internal/errors"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "filevault",
	Short: "🔒 Secure file encryption tool using AES-256-GCM",
	Long: cli.ColorCyan + `
 ███████ ██ ██      ███████ ██    ██  █████  ██    ██ ██   ████████ 
 ██      ██ ██      ██      ██    ██ ██   ██ ██    ██ ██      ██    
 █████   ██ ██      █████   ██    ██ ███████ ██    ██ ██      ██    
 ██      ██ ██      ██       ██  ██  ██   ██ ██    ██ ██      ██    
 ██      ██ ███████ ███████   ████   ██   ██  ██████  ███████ ██    
` + cli.ColorReset + `
FileVault is a modern, secure file encryption tool that uses industry-standard
cryptography to protect your sensitive files. It provides authenticated encryption
using AES-256-GCM with PBKDF2 key derivation for maximum security.

` + cli.ColorGreen + "🔐 SECURITY FEATURES:" + cli.ColorReset + `
  • AES-256-GCM authenticated encryption
  • PBKDF2 key derivation with 100,000 iterations  
  • 32-byte random salt for each file
  • Secure memory handling and cleanup
  • File integrity verification

` + cli.ColorBlue + "💡 COMMON USAGE:" + cli.ColorReset + `
  filevault encrypt document.pdf           # Encrypt a file
  filevault decrypt document.pdf.enc       # Decrypt a file  
  filevault info document.pdf.enc          # View file information
  filevault verify document.pdf.enc        # Verify file integrity

` + cli.ColorYellow + "🚀 ADVANCED OPTIONS:" + cli.ColorReset + `
  Use -v/--verbose for detailed output
  Use -q/--quiet for minimal output
  Use -f/--force to overwrite existing files
  Use -o/--output to specify custom output paths

For detailed help on any command, use: filevault <command> --help`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Show banner for main commands (not for help/version)
		if cmd.Use != "help" && cmd.Use != "version" && !cmd.Flags().Changed("help") {
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose && len(args) > 0 {
				cli.PrintBanner()
			}
		}
	},
	SilenceErrors: true, // We'll handle errors ourselves
	SilenceUsage:  true,
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "📋 Print version information",
	Long: `Print detailed version information including build details.

This shows the current FileVault version, git commit hash, and build date.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
		if verbose {
			cli.PrintBanner()
		}
		cli.PrintVersion(version, commit, date)
	},
}

// helpCmd provides enhanced help
var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "❓ Help about any command",
	Long: `Help provides detailed information about FileVault commands and usage.

Get help for a specific command:
  filevault help encrypt
  filevault help decrypt
  filevault help info

Or use the --help flag with any command:
  filevault encrypt --help`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			rootCmd.Help()
		} else {
			if subcmd, _, err := rootCmd.Find(args); err == nil {
				subcmd.Help()
			} else {
				fmt.Printf("Unknown command: %s\n", args[0])
				fmt.Println("\nAvailable commands:")
				for _, c := range rootCmd.Commands() {
					if !c.Hidden {
						fmt.Printf("  %-12s %s\n", c.Name(), c.Short)
					}
				}
			}
		}
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(commands.EncryptCmd)
	rootCmd.AddCommand(commands.DecryptCmd)
	rootCmd.AddCommand(commands.InfoCmd)
	rootCmd.AddCommand(commands.VerifyCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helpCmd)

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output with detailed information")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet output (errors only)")
	
	// Add usage examples
	rootCmd.SetUsageTemplate(getUsageTemplate())
	
	// Configure help
	rootCmd.SetHelpTemplate(getHelpTemplate())
	
	// Customize flag usage
	rootCmd.SetFlagErrorFunc(flagErrorFunc)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		exitCode := errors.HandleError(err, false)
		os.Exit(exitCode)
	}
}

func getUsageTemplate() string {
	return cli.ColorBold + `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}` + cli.ColorReset + `{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

` + cli.ColorGreen + `Available Commands:` + cli.ColorReset + `{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

` + cli.ColorYellow + `Flags:` + cli.ColorReset + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

` + cli.ColorCyan + `Global Flags:` + cli.ColorReset + `
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func getHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

func flagErrorFunc(cmd *cobra.Command, err error) error {
	cli.PrintError(fmt.Sprintf("Invalid flag usage: %v", err))
	fmt.Println()
	cmd.Help()
	return nil
}
