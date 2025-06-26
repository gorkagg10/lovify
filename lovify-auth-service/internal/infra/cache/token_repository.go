package cache

import (
	"context"

	"github.com/gorkagg10/lovify-authentication-service/internal/domain/login"
)

type TokenRepository struct {
	database map[string]Token
}

func NewTokenRepository(database map[string]Token) *TokenRepository {
	return &TokenRepository{
		database: database,
	}
}

func (t *TokenRepository) StoreToken(_ context.Context, token *login.Token, username string) error {
	t.database[token.Token()] = *NewToken(token.Token(), token.ExpirationDate(), username)
	return nil
}
