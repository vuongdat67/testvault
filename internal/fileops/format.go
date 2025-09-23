package fileops

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"path/filepath"
)

// FileVault binary format constants
const (
	MagicBytes         = "FVLT"
	FormatVersion      = 1
	AlgorithmAES256GCM = 1

	MagicSize          = 4
	VersionSize        = 4
	AlgorithmSize      = 4
	SaltSize           = 32
	IVSize             = 16
	OriginalSizeSize   = 8
	FileNameLengthSize = 4
	ReservedSize       = 32
	ChecksumSize       = 16

	BaseHeaderSize = MagicSize + VersionSize + AlgorithmSize + SaltSize +
		IVSize + OriginalSizeSize + FileNameLengthSize +
		ReservedSize + ChecksumSize

	AuthTagSize = 16
)

// FileHeader represents the FileVault file header
type FileHeader struct {
	Magic          [4]byte
	Version        uint32
	Algorithm      uint32
	Salt           [32]byte
	IV             [16]byte
	OriginalSize   uint64
	FileNameLength uint32
	FileName       string
	Reserved       [32]byte
	Checksum       [16]byte
}

// NewFileHeader creates a new file header
func NewFileHeader(originalSize uint64, fileName string, salt [32]byte, iv [16]byte) *FileHeader {
	header := &FileHeader{
		Version:        FormatVersion,
		Algorithm:      AlgorithmAES256GCM,
		Salt:           salt,
		IV:             iv,
		OriginalSize:   originalSize,
		FileNameLength: uint32(len(fileName)),
		FileName:       fileName,
	}

	copy(header.Magic[:], []byte(MagicBytes))
	header.calculateChecksum()

	return header
}

// calculateChecksum calculates and sets the header checksum
func (h *FileHeader) calculateChecksum() {
	hasher := sha256.New()

	hasher.Write(h.Magic[:])
	binary.Write(hasher, binary.LittleEndian, h.Version)
	binary.Write(hasher, binary.LittleEndian, h.Algorithm)
	hasher.Write(h.Salt[:])
	hasher.Write(h.IV[:])
	binary.Write(hasher, binary.LittleEndian, h.OriginalSize)
	binary.Write(hasher, binary.LittleEndian, h.FileNameLength)
	hasher.Write([]byte(h.FileName))
	hasher.Write(h.Reserved[:])

	hash := hasher.Sum(nil)
	copy(h.Checksum[:], hash[:16])
}

// IsValid checks if the header is valid
func (h *FileHeader) IsValid() error {
	if string(h.Magic[:]) != MagicBytes {
		return fmt.Errorf("invalid magic number")
	}

	if h.Version != FormatVersion {
		return fmt.Errorf("unsupported version: %d", h.Version)
	}

	return nil
}

// GetTotalSize returns the total header size
func (h *FileHeader) GetTotalSize() int {
	return BaseHeaderSize + len(h.FileName)
}

// GetBaseFileName returns the base filename
func (h *FileHeader) GetBaseFileName() string {
	return filepath.Base(h.FileName)
}

// WriteTo writes the header to an io.Writer
func (h *FileHeader) WriteTo(w io.Writer) (int64, error) {
	bytesWritten := int64(0)

	if err := binary.Write(w, binary.LittleEndian, h.Magic); err != nil {
		return bytesWritten, fmt.Errorf("failed to write magic: %w", err)
	}
	bytesWritten += MagicSize

	if err := binary.Write(w, binary.LittleEndian, h.Version); err != nil {
		return bytesWritten, fmt.Errorf("failed to write version: %w", err)
	}
	bytesWritten += VersionSize

	if err := binary.Write(w, binary.LittleEndian, h.Algorithm); err != nil {
		return bytesWritten, fmt.Errorf("failed to write algorithm: %w", err)
	}
	bytesWritten += AlgorithmSize

	if err := binary.Write(w, binary.LittleEndian, h.Salt); err != nil {
		return bytesWritten, fmt.Errorf("failed to write salt: %w", err)
	}
	bytesWritten += SaltSize

	if err := binary.Write(w, binary.LittleEndian, h.IV); err != nil {
		return bytesWritten, fmt.Errorf("failed to write IV: %w", err)
	}
	bytesWritten += IVSize

	if err := binary.Write(w, binary.LittleEndian, h.OriginalSize); err != nil {
		return bytesWritten, fmt.Errorf("failed to write original size: %w", err)
	}
	bytesWritten += OriginalSizeSize

	if err := binary.Write(w, binary.LittleEndian, h.FileNameLength); err != nil {
		return bytesWritten, fmt.Errorf("failed to write filename length: %w", err)
	}
	bytesWritten += FileNameLengthSize

	if h.FileNameLength > 0 {
		n, err := w.Write([]byte(h.FileName))
		if err != nil {
			return bytesWritten, fmt.Errorf("failed to write filename: %w", err)
		}
		bytesWritten += int64(n)
	}

	if err := binary.Write(w, binary.LittleEndian, h.Reserved); err != nil {
		return bytesWritten, fmt.Errorf("failed to write reserved: %w", err)
	}
	bytesWritten += ReservedSize

	if err := binary.Write(w, binary.LittleEndian, h.Checksum); err != nil {
		return bytesWritten, fmt.Errorf("failed to write checksum: %w", err)
	}
	bytesWritten += ChecksumSize

	return bytesWritten, nil
}

// ReadFrom reads the header from an io.Reader
func (h *FileHeader) ReadFrom(r io.Reader) (int64, error) {
	bytesRead := int64(0)

	if err := binary.Read(r, binary.LittleEndian, &h.Magic); err != nil {
		return bytesRead, fmt.Errorf("failed to read magic: %w", err)
	}
	bytesRead += MagicSize

	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return bytesRead, fmt.Errorf("failed to read version: %w", err)
	}
	bytesRead += VersionSize

	if err := binary.Read(r, binary.LittleEndian, &h.Algorithm); err != nil {
		return bytesRead, fmt.Errorf("failed to read algorithm: %w", err)
	}
	bytesRead += AlgorithmSize

	if err := binary.Read(r, binary.LittleEndian, &h.Salt); err != nil {
		return bytesRead, fmt.Errorf("failed to read salt: %w", err)
	}
	bytesRead += SaltSize

	if err := binary.Read(r, binary.LittleEndian, &h.IV); err != nil {
		return bytesRead, fmt.Errorf("failed to read IV: %w", err)
	}
	bytesRead += IVSize

	if err := binary.Read(r, binary.LittleEndian, &h.OriginalSize); err != nil {
		return bytesRead, fmt.Errorf("failed to read original size: %w", err)
	}
	bytesRead += OriginalSizeSize

	if err := binary.Read(r, binary.LittleEndian, &h.FileNameLength); err != nil {
		return bytesRead, fmt.Errorf("failed to read filename length: %w", err)
	}
	bytesRead += FileNameLengthSize

	if h.FileNameLength > 0 {
		if h.FileNameLength > 4096 {
			return bytesRead, fmt.Errorf("filename length too large: %d", h.FileNameLength)
		}

		fileNameBytes := make([]byte, h.FileNameLength)
		n, err := io.ReadFull(r, fileNameBytes)
		if err != nil {
			return bytesRead, fmt.Errorf("failed to read filename: %w", err)
		}
		h.FileName = string(fileNameBytes)
		bytesRead += int64(n)
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Reserved); err != nil {
		return bytesRead, fmt.Errorf("failed to read reserved: %w", err)
	}
	bytesRead += ReservedSize

	if err := binary.Read(r, binary.LittleEndian, &h.Checksum); err != nil {
		return bytesRead, fmt.Errorf("failed to read checksum: %w", err)
	}
	bytesRead += ChecksumSize

	return bytesRead, nil
}
