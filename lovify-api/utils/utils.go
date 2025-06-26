package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		slog.Error("generating token", slog.String("error", err.Error()))
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
