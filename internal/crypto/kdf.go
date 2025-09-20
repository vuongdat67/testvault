package crypto

import (
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
)

// DeriveKey derives encryption key from password using PBKDF2
func DeriveKey(password string, salt []byte, iterations int) []byte {
    if iterations <= 0 {
        iterations = DefaultIterations
    }
    return pbkdf2.Key([]byte(password), salt, iterations, KeySize, sha256.New)
}

// DeriveKeyWithParams derives key using predefined parameters
func DeriveKeyWithParams(password string, params KeyDerivationParams) []byte {
    return DeriveKey(password, params.Salt, params.Iterations)
}

// CreateKeyDerivationParams creates new key derivation parameters
func CreateKeyDerivationParams() (*KeyDerivationParams, error) {
    salt, err := GenerateSalt()
    if err != nil {
        return nil, err
    }
    
    return &KeyDerivationParams{
        Salt:       salt,
        Iterations: DefaultIterations,
        KeyLength:  KeySize,
    }, nil
}