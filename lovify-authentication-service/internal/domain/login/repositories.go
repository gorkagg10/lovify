package login

import "context"

type UserRepository interface {
	EmailExists(context context.Context, email string) (bool, error)
	GetUser(context context.Context, email string) (*User, error)
	CreateUser(context context.Context, user *User) error
	ConnectProfile(email string, profileID string) error
}

type SecurityRepository interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword string, password string) (bool, error)
}

type TokenRepository interface {
	GenerateToken(tokenType string, email string) (*Token, error)
	ValidateToken(token string) error
}
