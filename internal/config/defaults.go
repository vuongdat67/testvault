package config

// Default values for FileVault configuration
const (
	// Cryptographic defaults
	DefaultPBKDF2Iterations = 100000
	DefaultAlgorithmName    = "AES-256-GCM"
	DefaultSaltSize         = 32
	DefaultIVSize           = 16
	DefaultKeySize          = 32
	DefaultTagSize          = 16

	// File operation defaults
	DefaultBufferSize         = 64 * 1024   // 64KB
	DefaultSmallFileBuffer    = 8 * 1024    // 8KB
	DefaultLargeFileBuffer    = 256 * 1024  // 256KB
	DefaultStreamingThreshold = 1024 * 1024 // 1MB

	// Security defaults
	DefaultPasswordMinLength = 8
	DefaultPasswordMaxLength = 128
	MinSecurePasswordLength  = 12

	// File size limits
	DefaultMaxFileSize     = 10 * 1024 * 1024 * 1024  // 10GB
	MaxAllowedFileSize     = 100 * 1024 * 1024 * 1024 // 100GB
	MinProcessableFileSize = 1                        // 1 byte

	// Performance defaults
	DefaultBatchSize      = 100
	DefaultProgressUpdate = 100 * 1024 * 1024 // Update every 100MB
	DefaultTimeoutSeconds = 300               // 5 minutes

	// UI defaults
	DefaultProgressBarWidth = 50
	DefaultColorSupport     = true
	DefaultVerboseMode      = false

	// File extensions
	DefaultEncryptedExtension = ".enc"
	DefaultDecryptedExtension = ".dec"
	DefaultBackupExtension    = ".bak"

	// Temporary file settings
	TempFilePrefix = "filevault_"
	TempFileSuffix = ".tmp"

	// Configuration file settings
	ConfigFileName        = "config.json"
	ConfigDirName         = ".filevault"
	ConfigFilePermissions = 0644
	ConfigDirPermissions  = 0755

	// Exit codes
	ExitSuccess               = 0
	ExitGeneralError          = 1
	ExitFileNotFound          = 2
	ExitPermissionDenied      = 3
	ExitAuthenticationFailed  = 4
	ExitCorruptedFile         = 5
	ExitInsufficientResources = 6
	ExitInvalidArguments      = 7
)

// Default paths
var (
	// These will be set at runtime based on user's environment
	DefaultConfigDir string
	DefaultTempDir   string
	DefaultOutputDir string
	DefaultLogDir    string
)

// Validation limits
const (
	MinIterations = 1000
	MaxIterations = 10000000

	MinBufferSize = 1024             // 1KB
	MaxBufferSize = 10 * 1024 * 1024 // 10MB

	MinPasswordLength = 8
	MaxPasswordLength = 128

	MinProgressBarWidth = 10
	MaxProgressBarWidth = 100
)

// Supported algorithms
var SupportedAlgorithms = []string{
	"AES-256-GCM",
	// Future algorithms can be added here
}

// Default file patterns for batch operations
var DefaultFilePatterns = []string{
	"*.txt",
	"*.doc",
	"*.docx",
	"*.pdf",
	"*.jpg",
	"*.jpeg",
	"*.png",
	"*.zip",
	"*.tar",
	"*.gz",
}

// Default exclusion patterns (files to skip in batch operations)
var DefaultExclusionPatterns = []string{
	"*.enc", // Already encrypted files
	"*.tmp", // Temporary files
	"*.log", // Log files
	".*",    // Hidden files (can be overridden)
}

// Performance profiles
type PerformanceProfile struct {
	Name               string
	BufferSize         int
	StreamingThreshold int64
	BatchSize          int
	Description        string
}

var PerformanceProfiles = map[string]PerformanceProfile{
	"fast": {
		Name:               "Fast",
		BufferSize:         256 * 1024, // 256KB
		StreamingThreshold: 512 * 1024, // 512KB
		BatchSize:          50,
		Description:        "Optimized for speed, uses more memory",
	},
	"balanced": {
		Name:               "Balanced",
		BufferSize:         64 * 1024,   // 64KB
		StreamingThreshold: 1024 * 1024, // 1MB
		BatchSize:          25,
		Description:        "Balanced performance and memory usage (default)",
	},
	"memory": {
		Name:               "Memory Efficient",
		BufferSize:         8 * 1024,   // 8KB
		StreamingThreshold: 256 * 1024, // 256KB
		BatchSize:          10,
		Description:        "Minimizes memory usage, slower performance",
	},
}

// Security levels
type SecurityLevel struct {
	Name           string
	Iterations     int
	MinPasswordLen int
	RequireStrong  bool
	Description    string
}

var SecurityLevels = map[string]SecurityLevel{
	"standard": {
		Name:           "Standard",
		Iterations:     100000,
		MinPasswordLen: 8,
		RequireStrong:  false,
		Description:    "Standard security level (default)",
	},
	"high": {
		Name:           "High",
		Iterations:     200000,
		MinPasswordLen: 12,
		RequireStrong:  true,
		Description:    "High security with stronger password requirements",
	},
	"paranoid": {
		Name:           "Paranoid",
		Iterations:     500000,
		MinPasswordLen: 16,
		RequireStrong:  true,
		Description:    "Maximum security (slower but more secure)",
	},
}

// GetDefaultProfile returns the default performance profile
func GetDefaultProfile() PerformanceProfile {
	return PerformanceProfiles["balanced"]
}

// GetDefaultSecurityLevel returns the default security level
func GetDefaultSecurityLevel() SecurityLevel {
	return SecurityLevels["standard"]
}

// IsAlgorithmSupported checks if an algorithm is supported
func IsAlgorithmSupported(algorithm string) bool {
	for _, supported := range SupportedAlgorithms {
		if supported == algorithm {
			return true
		}
	}
	return false
}

// ValidateIterations validates PBKDF2 iteration count
func ValidateIterations(iterations int) bool {
	return iterations >= MinIterations && iterations <= MaxIterations
}

// ValidateBufferSize validates buffer size
func ValidateBufferSize(size int) bool {
	return size >= MinBufferSize && size <= MaxBufferSize
}

// ValidatePasswordLength validates password length
func ValidatePasswordLength(length int) bool {
	return length >= MinPasswordLength && length <= MaxPasswordLength
}
