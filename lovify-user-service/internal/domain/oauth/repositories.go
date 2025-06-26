package oauth

import "context"

type AuthRepository interface {
	RequestAuthorization(state string) string
	Exchange(ctx context.Context, code string) (*Token, error)
}
