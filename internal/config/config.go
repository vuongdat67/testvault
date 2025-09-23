package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the FileVault configuration
type Config struct {
	// Encryption settings
	DefaultIterations int    `json:"default_iterations"`
	DefaultAlgorithm  string `json:"default_algorithm"`
	BufferSize        int    `json:"buffer_size"`

	// Security settings
	PasswordMinLength     int  `json:"password_min_length"`
	RequireStrongPassword bool `json:"require_strong_password"`
	SecureMemory          bool `json:"secure_memory"`

	// UI settings
	UseColors     bool `json:"use_colors"`
	ShowProgress  bool `json:"show_progress"`
	VerboseOutput bool `json:"verbose_output"`

	// Performance settings
	MaxFileSize        int64 `json:"max_file_size"`
	StreamingThreshold int64 `json:"streaming_threshold"`

	// Paths
	ConfigDir        string `json:"config_dir"`
	TempDir          string `json:"temp_dir"`
	DefaultOutputDir string `json:"default_output_dir"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".filevault")

	return &Config{
		// Encryption defaults
		DefaultIterations: 100000,
		DefaultAlgorithm:  "AES-256-GCM",
		BufferSize:        64 * 1024, // 64KB

		// Security defaults
		PasswordMinLength:     12,
		RequireStrongPassword: true,
		SecureMemory:          true,

		// UI defaults
		UseColors:     true,
		ShowProgress:  true,
		VerboseOutput: false,

		// Performance defaults
		MaxFileSize:        10 * 1024 * 1024 * 1024, // 10GB
		StreamingThreshold: 1024 * 1024,             // 1MB

		// Path defaults
		ConfigDir:        configDir,
		TempDir:          os.TempDir(),
		DefaultOutputDir: "",
	}
}

// Load loads configuration from file
func Load() (*Config, error) {
	config := DefaultConfig()

	configFile := config.GetConfigFilePath()

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(c.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with proper formatting
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configFile := c.GetConfigFilePath()
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigFilePath returns the path to the configuration file
func (c *Config) GetConfigFilePath() string {
	return filepath.Join(c.ConfigDir, "config.json")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DefaultIterations < 1000 {
		return fmt.Errorf("default_iterations must be at least 1000")
	}

	if c.DefaultIterations > 10000000 {
		return fmt.Errorf("default_iterations must not exceed 10,000,000")
	}

	if c.PasswordMinLength < 8 {
		return fmt.Errorf("password_min_length must be at least 8")
	}

	if c.PasswordMinLength > 128 {
		return fmt.Errorf("password_min_length must not exceed 128")
	}

	if c.BufferSize < 1024 {
		return fmt.Errorf("buffer_size must be at least 1024 bytes")
	}

	if c.BufferSize > 10*1024*1024 {
		return fmt.Errorf("buffer_size must not exceed 10MB")
	}

	if c.MaxFileSize < 1024 {
		return fmt.Errorf("max_file_size must be at least 1024 bytes")
	}

	if c.StreamingThreshold < 1024 {
		return fmt.Errorf("streaming_threshold must be at least 1024 bytes")
	}

	return nil
}

// Reset resets configuration to defaults
func (c *Config) Reset() {
	defaultConfig := DefaultConfig()
	*c = *defaultConfig
}

// GetTempFile returns a temporary file path
func (c *Config) GetTempFile(prefix string) string {
	return filepath.Join(c.TempDir, fmt.Sprintf("%s_%d.tmp", prefix, os.Getpid()))
}

// IsColorSupported returns true if colored output is supported and enabled
func (c *Config) IsColorSupported() bool {
	if !c.UseColors {
		return false
	}

	// Check if we're on Windows (basic check)
	if runtime.GOOS == "windows" {
		// Windows Terminal supports colors, but older cmd.exe might not
		return os.Getenv("WT_SESSION") != "" || os.Getenv("ConEmuANSI") != ""
	}

	// Unix-like systems generally support colors
	return true
}

// ShouldShowProgress returns true if progress should be shown
func (c *Config) ShouldShowProgress(fileSize int64) bool {
	return c.ShowProgress && fileSize >= c.StreamingThreshold
}

// GetEffectiveBufferSize returns the buffer size to use for a given file size
func (c *Config) GetEffectiveBufferSize(fileSize int64) int {
	// Use smaller buffer for small files
	if fileSize < 1024*1024 { // < 1MB
		return min(c.BufferSize, 8*1024) // Max 8KB for small files
	}

	// Use larger buffer for very large files
	if fileSize > 100*1024*1024 { // > 100MB
		return max(c.BufferSize, 256*1024) // Min 256KB for large files
	}

	return c.BufferSize
}

// UpdateFromFlags updates config from command line flags
func (c *Config) UpdateFromFlags(iterations int, verbose bool, quiet bool, useColors bool) {
	if iterations > 0 {
		c.DefaultIterations = iterations
	}

	if verbose {
		c.VerboseOutput = true
	}

	if quiet {
		c.VerboseOutput = false
		c.ShowProgress = false
	}

	// Override color setting if explicitly specified
	c.UseColors = useColors && c.IsColorSupported()
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
