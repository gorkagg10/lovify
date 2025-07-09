package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	autherrors "github.com/gorkagg10/lovify/lovify-authentication-service/errors"
	"log/slog"
	"time"

	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
)

type TokenRepository struct {
	pgClient *sql.DB
}

func NewTokenRepository(pgClient *sql.DB) *TokenRepository {
	return &TokenRepository{
		pgClient: pgClient,
	}
}

func (t *TokenRepository) StoreToken(ctx context.Context, token *login.Token, email string) error {
	var tokenID int64
	if err := t.pgClient.QueryRowContext(
		ctx,
		`INSERT INTO tokens (token, type, expiration_date, email) 
				VALUES ($1, $2, $3, $4)
				RETURNING id;`, token.Token(), token.TokenType(), token.ExpirationDate().String(), email).Scan(&tokenID); err != nil {
		slog.Error("inserting token into database", slog.String("error", err.Error()))
		return autherrors.ErrDatabaseQueryFailed
	}
	return nil
}

func (t *TokenRepository) GetToken(ctx context.Context, token string, tokenType string, email string) (*login.Token, error) {
	var dbToken Token
	if err := t.pgClient.QueryRowContext(
		ctx,
		`SELECT token, token_type, expiration_date, email from tokens WHERE token = $1, token_type = $2, email = $3;`,
		token, tokenType, email,
	).Scan(&dbToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("token not found")
		}
		return nil, fmt.Errorf("fetching token: %w", err)
	}
	expirationDate, err := time.Parse(time.RFC3339, dbToken.ExpirationDate)
	if err != nil {
		return nil, fmt.Errorf("parsing token expiration date: %w", err)
	}
	return login.NewToken(dbToken.Token, dbToken.TokenType, expirationDate), nil
}
