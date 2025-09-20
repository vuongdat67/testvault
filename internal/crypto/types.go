package crypto

import "errors"

// Crypto errors
var (
    ErrInvalidKeySize    = errors.New("invalid key size: must be 32 bytes for AES-256")
    ErrInvalidNonceSize  = errors.New("invalid nonce size: must be 12 bytes for GCM")
    ErrCiphertextTooShort = errors.New("ciphertext too short")
    ErrDecryptionFailed  = errors.New("decryption failed: authentication failed")
)

// Constants for AES-256-GCM
const (
    KeySize      = 32 // AES-256 key size
    NonceSize    = 12 // GCM recommended nonce size
    TagSize      = 16 // GCM authentication tag size
    SaltSize     = 16 // Salt size for PBKDF2
    DefaultIterations = 100000 // PBKDF2 iterations
)

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
    Nonce      []byte `json:"nonce"`
    Salt       []byte `json:"salt"`
    Ciphertext []byte `json:"ciphertext"`
    Tag        []byte `json:"tag"`
}

// KeyDerivationParams holds parameters for key derivation
type KeyDerivationParams struct {
    Salt       []byte
    Iterations int
    KeyLength  int
}