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
	GenerateToken(tokenType string) (*Token, error)
}

type TokenRepository interface {
	StoreToken(ctx context.Context, token *Token, email string) error
	GetToken(ctx context.Context, token string, tokenType string, email string) (*Token, error)
}
