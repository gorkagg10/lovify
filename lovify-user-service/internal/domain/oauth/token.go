package oauth

import "time"

type Token struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func NewToken(accessToken, refreshToken string, expiresAt time.Time) *Token {
	return &Token{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
	}
}

func (t *Token) AccessToken() string {
	return t.accessToken
}

func (t *Token) RefreshToken() string {
	return t.refreshToken
}

func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}
