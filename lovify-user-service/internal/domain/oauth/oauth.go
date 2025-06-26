package oauth

import (
	"context"
	"encoding/base64"
)

type Service struct {
	OAuthRepository AuthRepository
}

func NewService(oauthRepository AuthRepository) *Service {
	return &Service{
		OAuthRepository: oauthRepository,
	}
}

func (s *Service) RequestAuthorization(userID string) string {
	state := base64.URLEncoding.EncodeToString([]byte(userID))
	return s.OAuthRepository.RequestAuthorization(state)
}

func (s *Service) Exchange(ctx context.Context, code string) (*Token, error) {
	return s.OAuthRepository.Exchange(ctx, code)
}
