package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
	"time"
)

const jwtSecret = "mF3T&^g@j8!WqZ9#Ln2sR4tXkP%vYbE7"

type TokenRepository struct {
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (t *TokenRepository) GenerateToken(tokenType string, email string) (*login.Token, error) {
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

func (t *TokenRepository) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_, ok = claims["email"].(string)
		if !ok {
			return errors.New("invalid email")
		}
	}
	return nil
}
