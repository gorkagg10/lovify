package base64

import (
	"crypto/rand"
	"encoding/base64"
	autherrors "github.com/gorkagg10/lovify-authentication-service/errors"
	"github.com/gorkagg10/lovify-authentication-service/internal/domain/login"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const tokenLength = 32

type SecurityRepository struct {
}

func NewSecurityRepository() *SecurityRepository {
	return &SecurityRepository{}
}

func (t *SecurityRepository) HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error("generating hashed password", slog.String("error", err.Error()))
		return "", autherrors.ErrHashedPasswordGenerationFailed
	}
	return string(passwordBytes), nil
}

func (t *SecurityRepository) GenerateToken(tokenType string) (*login.Token, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		slog.Error("generating token", slog.String("error", err.Error()))
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	return login.NewToken(token, tokenType, time.Now().Add(time.Hour*24)), nil
}

func (t *SecurityRepository) CheckPassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("checking password", slog.String("error", err.Error()))
		return false, autherrors.ErrIncorrectPassword
	}
	return true, nil
}
