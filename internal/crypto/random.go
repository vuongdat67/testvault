package crypto

import (
    "crypto/rand"
    "io"
)

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(size int) ([]byte, error) {
    bytes := make([]byte, size)
    if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
        return nil, err
    }
    return bytes, nil
}

// GenerateNonce generates a random nonce for GCM
func GenerateNonce() ([]byte, error) {
    return GenerateRandomBytes(NonceSize)
}

// GenerateSalt generates a random salt for key derivation
func GenerateSalt() ([]byte, error) {
    return GenerateRandomBytes(SaltSize)
}

// GenerateKey generates a random 256-bit encryption key
func GenerateKey() ([]byte, error) {
    return GenerateRandomBytes(KeySize)
}
