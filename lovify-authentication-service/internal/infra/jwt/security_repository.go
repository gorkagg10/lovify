package base64

import (
	jwt "github.com/golang-jwt/jwt/v5"
	autherrors "github.com/gorkagg10/lovify/lovify-authentication-service/errors"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "mF3T&^g@j8!WqZ9#Ln2sR4tXkP%vYbE7"

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

func (t *SecurityRepository) GenerateToken(tokenType string, email string) (*login.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return login.NewToken(tokenString, tokenType, time.Now().Add(time.Hour*24)), nil
}

func (t *SecurityRepository) CheckPassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("checking password", slog.String("error", err.Error()))
		return false, autherrors.ErrIncorrectPassword
	}
	return true, nil
}
