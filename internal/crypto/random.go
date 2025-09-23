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

// GenerateSalt32 generates a 32-byte salt as array for the new format
func GenerateSalt32() ([32]byte, error) {
	var salt [32]byte
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return salt, err
	}
	copy(salt[:], bytes)
	return salt, nil
}

// GenerateIV16 generates a 16-byte IV for the new format (12 bytes used + 4 zeros)
func GenerateIV16() ([16]byte, error) {
	var iv [16]byte
	nonce, err := GenerateNonce()
	if err != nil {
		return iv, err
	}
	copy(iv[:12], nonce) // First 12 bytes are the actual nonce
	// Last 4 bytes remain zero-padded
	return iv, nil
}

// GenerateKey generates a random 256-bit encryption key
func GenerateKey() ([]byte, error) {
	return GenerateRandomBytes(KeySize)
}
