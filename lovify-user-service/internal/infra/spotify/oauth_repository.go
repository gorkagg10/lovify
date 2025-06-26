package spotify

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/gorkagg10/lovify-user-service/internal/domain/oauth"
)

type OAuthRepository struct {
	OAuthServiceConfig *oauth2.Config
}

func NewOAuthRepository(oauthServiceConfig *oauth2.Config) *OAuthRepository {
	return &OAuthRepository{
		OAuthServiceConfig: oauthServiceConfig,
	}
}

func (o *OAuthRepository) RequestAuthorization(state string) string {
	return o.OAuthServiceConfig.AuthCodeURL(state)
}

func (o *OAuthRepository) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	token, err := o.OAuthServiceConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return oauth.NewToken(token.AccessToken, token.RefreshToken, token.Expiry), nil
}
