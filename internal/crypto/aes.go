package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "fmt"
)

// AESCipher handles AES-256-GCM encryption/decryption
type AESCipher struct {
    key []byte
}

// NewAESCipher creates a new AES cipher with the given key
func NewAESCipher(key []byte) (*AESCipher, error) {
    if len(key) != KeySize {
        return nil, fmt.Errorf("%w: got %d bytes", ErrInvalidKeySize, len(key))
    }
    
    // Validate key by creating cipher
    _, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("invalid AES key: %w", err)
    }
    
    return &AESCipher{key: key}, nil
}

// NewAESCipherFromPassword creates cipher from password using PBKDF2
func NewAESCipherFromPassword(password string, salt []byte) (*AESCipher, error) {
    key := DeriveKey(password, salt, DefaultIterations)
    return NewAESCipher(key)
}

// Encrypt encrypts plaintext using AES-256-GCM
func (c *AESCipher) Encrypt(plaintext []byte) (*EncryptedData, error) {
    // Create AES cipher
    block, err := aes.NewCipher(c.key)
    if err != nil {
        return nil, err
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Generate random nonce
    nonce, err := GenerateNonce()
    if err != nil {
        return nil, err
    }
    
    // Encrypt and authenticate
    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    
    // Split ciphertext and tag (GCM appends tag to ciphertext)
    tagStart := len(ciphertext) - TagSize
    actualCiphertext := ciphertext[:tagStart]
    tag := ciphertext[tagStart:]
    
    return &EncryptedData{
        Nonce:      nonce,
        Ciphertext: actualCiphertext,
        Tag:        tag,
    }, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (c *AESCipher) Decrypt(data *EncryptedData) ([]byte, error) {
    // Validate inputs
    if len(data.Nonce) != NonceSize {
        return nil, fmt.Errorf("%w: got %d bytes", ErrInvalidNonceSize, len(data.Nonce))
    }
    
    if len(data.Tag) != TagSize {
        return nil, fmt.Errorf("invalid tag size: expected %d, got %d", TagSize, len(data.Tag))
    }
    
    // Create AES cipher
    block, err := aes.NewCipher(c.key)
    if err != nil {
        return nil, err
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Reconstruct full ciphertext with tag
    fullCiphertext := append(data.Ciphertext, data.Tag...)
    
    // Decrypt and verify
    plaintext, err := gcm.Open(nil, data.Nonce, fullCiphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
    }
    
    return plaintext, nil
}

// EncryptWithPassword is a convenience function for password-based encryption
func EncryptWithPassword(plaintext []byte, password string) (*EncryptedData, []byte, error) {
    // Generate salt
    salt, err := GenerateSalt()
    if err != nil {
        return nil, nil, err
    }
    
    // Create cipher from password
    cipher, err := NewAESCipherFromPassword(password, salt)
    if err != nil {
        return nil, nil, err
    }
    
    // Encrypt
    encryptedData, err := cipher.Encrypt(plaintext)
    if err != nil {
        return nil, nil, err
    }
    
    // Add salt to encrypted data
    encryptedData.Salt = salt
    
    return encryptedData, salt, nil
}

// DecryptWithPassword is a convenience function for password-based decryption
func DecryptWithPassword(data *EncryptedData, password string) ([]byte, error) {
    if data.Salt == nil {
        return nil, fmt.Errorf("salt is required for password-based decryption")
    }
    
    // Create cipher from password and salt
    cipher, err := NewAESCipherFromPassword(password, data.Salt)
    if err != nil {
        return nil, err
    }
    
    return cipher.Decrypt(data)
}

// SecureZero securely zeros out sensitive data in memory
func SecureZero(data []byte) {
    for i := range data {
        data[i] = 0
    }
}