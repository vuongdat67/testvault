package fileops

import (
    "encoding/binary"
    "fmt"
    "io"
)

// FileVault binary format constants
const (
    // Magic number for FileVault files
    MagicNumber = 0x46565431 // "FVT1" in hex
    
    // Version
    FormatVersion = 1
    
    // Header sizes
    HeaderSize     = 64  // Total header size
    MagicSize      = 4   // Magic number size
    VersionSize    = 4   // Version size
    MetadataSize   = 56  // Remaining metadata size
    
    // Crypto metadata sizes (from crypto package)
    SaltSize  = 16 // PBKDF2 salt
    NonceSize = 12 // GCM nonce
    TagSize   = 16 // GCM auth tag
)

// FileHeader represents the FileVault file header
type FileHeader struct {
    Magic         uint32   // Magic number "FVT1"
    Version       uint32   // Format version
    Salt          [16]byte // PBKDF2 salt
    Nonce         [12]byte // GCM nonce
    AuthTag       [16]byte // GCM authentication tag
    OriginalSize  uint64   // Original file size
    Reserved      [8]byte  // Reserved for future use
}

// NewFileHeader creates a new file header
func NewFileHeader() *FileHeader {
    return &FileHeader{
        Magic:   MagicNumber,
        Version: FormatVersion,
    }
}

// IsValid checks if the header has valid magic number and version
func (h *FileHeader) IsValid() error {
    if h.Magic != MagicNumber {
        return fmt.Errorf("invalid magic number: expected 0x%08X, got 0x%08X", MagicNumber, h.Magic)
    }
    
    if h.Version != FormatVersion {
        return fmt.Errorf("unsupported version: expected %d, got %d", FormatVersion, h.Version)
    }
    
    return nil
}

// WriteTo writes the header to an io.Writer
func (h *FileHeader) WriteTo(w io.Writer) (int64, error) {
    // Write in binary format (little-endian)
    if err := binary.Write(w, binary.LittleEndian, h.Magic); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.Version); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.Salt); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.Nonce); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.AuthTag); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.OriginalSize); err != nil {
        return 0, err
    }
    if err := binary.Write(w, binary.LittleEndian, h.Reserved); err != nil {
        return 0, err
    }
    
    return int64(HeaderSize), nil
}

// ReadFrom reads the header from an io.Reader
func (h *FileHeader) ReadFrom(r io.Reader) (int64, error) {
    // Read in binary format (little-endian)
    if err := binary.Read(r, binary.LittleEndian, &h.Magic); err != nil {
        return 0, fmt.Errorf("failed to read magic: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
        return 0, fmt.Errorf("failed to read version: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Salt); err != nil {
        return 0, fmt.Errorf("failed to read salt: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Nonce); err != nil {
        return 0, fmt.Errorf("failed to read nonce: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.AuthTag); err != nil {
        return 0, fmt.Errorf("failed to read auth tag: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.OriginalSize); err != nil {
        return 0, fmt.Errorf("failed to read original size: %w", err)
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Reserved); err != nil {
        return 0, fmt.Errorf("failed to read reserved: %w", err)
    }
    
    if err := h.IsValid(); err != nil {
        return 0, err
    }
    
    return int64(HeaderSize), nil
}

// GetSize returns the header size in bytes
func (h *FileHeader) GetSize() int {
    return HeaderSize
}

// String returns a string representation of the header
func (h *FileHeader) String() string {
    return fmt.Sprintf("FileVault Header v%d (Magic: 0x%08X, Size: %d bytes)", 
        h.Version, h.Magic, h.OriginalSize)
}