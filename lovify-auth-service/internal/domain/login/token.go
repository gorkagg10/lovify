package login

import "time"

type Token struct {
	token          string
	tokenType      string
	expirationDate time.Time
}

func NewToken(token string, tokenType string, expirationDate time.Time) *Token {
	return &Token{
		token:          token,
		tokenType:      tokenType,
		expirationDate: expirationDate,
	}
}

func (t *Token) Token() string {
	return t.token
}

func (t *Token) TokenType() string {
	return t.tokenType
}

func (t *Token) ExpirationDate() time.Time {
	return t.expirationDate
}
