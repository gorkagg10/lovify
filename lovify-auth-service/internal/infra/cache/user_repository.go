package cache

import (
	"context"
	"fmt"

	"github.com/gorkagg10/lovify-authentication-service/internal/domain/login"
)

type UserRepository struct {
	database map[string]User
}

func NewUserRepository(database map[string]User) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (u *UserRepository) EmailExists(_ context.Context, email string) (bool, error) {
	_, ok := u.database[email]
	return ok, nil
}

func (u *UserRepository) GetUser(_ context.Context, email string) (*login.User, error) {
	user, ok := u.database[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return login.NewUser(user.Email(), user.HashedPassword()), nil
}

func (u *UserRepository) CreateUser(_ context.Context, user *login.User) error {
	u.database[user.Email()] = *NewUser(user.Email(), user.HashedPassword())
	return nil
}
