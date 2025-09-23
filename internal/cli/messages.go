package cli

import (
	"fmt"
	"os"
)

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Message types
type MessageType int

const (
	Success MessageType = iota
	Error
	Warning
	Info
	Progress
)

// IsColorSupported checks if terminal supports colors
func IsColorSupported() bool {
	// Check if we're outputting to a terminal
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) == os.ModeCharDevice
	}
	return false
}

// PrintColored prints a colored message if supported
func PrintColored(msg string, color string) {
	if IsColorSupported() {
		fmt.Print(color + msg + ColorReset)
	} else {
		fmt.Print(msg)
	}
}

// PrintMessage prints a formatted message with appropriate icon and color
func PrintMessage(msgType MessageType, message string) {
	var icon, color string
	
	switch msgType {
	case Success:
		icon = "✅"
		color = ColorGreen
	case Error:
		icon = "❌"
		color = ColorRed
	case Warning:
		icon = "⚠️ "
		color = ColorYellow
	case Info:
		icon = "ℹ️ "
		color = ColorBlue
	case Progress:
		icon = "⏳"
		color = ColorCyan
	}

	if IsColorSupported() {
		fmt.Printf("%s %s%s%s\n", icon, color, message, ColorReset)
	} else {
		fmt.Printf("%s %s\n", icon, message)
	}
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	PrintMessage(Success, message)
}

// PrintError prints an error message
func PrintError(message string) {
	PrintMessage(Error, message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	PrintMessage(Warning, message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	PrintMessage(Info, message)
}

// PrintProgress prints a progress message
func PrintProgress(message string) {
	PrintMessage(Progress, message)
}

// ConfirmAction prompts user for confirmation
func ConfirmAction(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	var response string
	fmt.Scanln(&response)
	
	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}

// FormatBytes formats byte count as human-readable string
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats duration in a human-readable way
func FormatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.1fm", seconds/60)
	} else {
		return fmt.Sprintf("%.1fh", seconds/3600)
	}
}

// PrintBanner prints the application banner
func PrintBanner() {
	banner := `
 ███████ ██ ██      ███████ ██    ██  █████  ██    ██ ██   ████████ 
 ██      ██ ██      ██      ██    ██ ██   ██ ██    ██ ██      ██    
 █████   ██ ██      █████   ██    ██ ███████ ██    ██ ██      ██    
 ██      ██ ██      ██       ██  ██  ██   ██ ██    ██ ██      ██    
 ██      ██ ███████ ███████   ████   ██   ██  ██████  ███████ ██    
                                                                    
 Secure File Encryption Tool - AES-256-GCM with PBKDF2
`
	PrintColored(banner, ColorCyan)
}

// PrintVersion prints version information
func PrintVersion(version, commit, date string) {
	fmt.Printf("FileVault %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built: %s\n", date)
	fmt.Printf("Go version: %s\n", "go1.25.0")
}

// Usage examples for commands
var UsageExamples = map[string][]string{
	"encrypt": {
		"filevault encrypt document.pdf",
		"filevault encrypt document.pdf -o secure.enc",
		"filevault encrypt *.txt -o encrypted/",
		"filevault encrypt large-file.zip --iterations 200000",
	},
	"decrypt": {
		"filevault decrypt document.pdf.enc",
		"filevault decrypt secure.enc original.pdf",
		"filevault decrypt *.enc -o decrypted/",
	},
	"info": {
		"filevault info document.pdf.enc",
		"filevault info encrypted-files/*.enc",
	},
	"verify": {
		"filevault verify document.pdf.enc",
		"filevault verify *.enc",
	},
}

// PrintUsageExamples prints usage examples for a command
func PrintUsageExamples(command string) {
	if examples, exists := UsageExamples[command]; exists {
		fmt.Printf("\nExamples:\n")
		for _, example := range examples {
			fmt.Printf("  %s\n", example)
		}
	}
}
