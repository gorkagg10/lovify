package cache

import "time"

type User struct {
	email          string
	hashedPassword string
}

func NewUser(
	email string,
	hashedPassword string) *User {
	return &User{
		email:          email,
		hashedPassword: hashedPassword,
	}
}

func (u *User) Email() string {
	return u.email
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

type Token struct {
	token          string
	expirationDate time.Time
	email          string
}

func NewToken(token string, expirationDate time.Time, email string) *Token {
	return &Token{
		token:          token,
		expirationDate: expirationDate,
		email:          email,
	}
}
