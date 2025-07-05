package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	autherrors "github.com/gorkagg10/lovify/lovify-authentication-service/errors"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
)

type UserRepository struct {
	pgClient *sql.DB
}

func NewUserRepository(pgClient *sql.DB) *UserRepository {
	return &UserRepository{pgClient}
}

func (u *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var userId int64
	if err := u.pgClient.QueryRowContext(
		ctx,
		`SELECT id from users where email = $1;`, email).Scan(&userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		slog.Error("checking if email exists", slog.String("errors", err.Error()))
		return false, autherrors.ErrDatabaseQueryFailed
	}
	return true, autherrors.ErrUserAlreadyExists
}

func (u *UserRepository) GetUser(ctx context.Context, email string) (*login.User, error) {
	var user User
	if err := u.pgClient.QueryRowContext(
		ctx,
		`SELECT email, password from users where email = $1;`, email).Scan(&user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("getting user from database: %w", err)
	}
	return login.NewUser(user.Email, user.Password), nil
}

func (u *UserRepository) CreateUser(ctx context.Context, user *login.User) error {
	var userID int64
	if err := u.pgClient.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password)
				VALUES($1, $2)
				RETURNING id;`, user.Email(), user.HashedPassword()).Scan(&userID); err != nil {
		slog.Error("inserting user in database", slog.String("error", err.Error()))
		return autherrors.ErrDatabaseQueryFailed
	}
	return nil
}
