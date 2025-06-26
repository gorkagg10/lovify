package aescgm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type SecurityRepository struct {
	SecretKey string
}

func NewSecurityRepository(secretKey string) *SecurityRepository {
	return &SecurityRepository{
		SecretKey: secretKey,
	}
}

func (s *SecurityRepository) EncryptToken(token string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(s.SecretKey)
	if err != nil {
		return "", fmt.Errorf("decoding secret key: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("creating cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("creating GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generating nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *SecurityRepository) DecryptToken(token string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(s.SecretKey)
	if err != nil {
		return "", fmt.Errorf("decoding secret key: %w", err)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("decoding token: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("creating cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("creating GCM: %w", err)
	}
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypting token: %w", err)
	}
	return string(plaintext), nil
}
