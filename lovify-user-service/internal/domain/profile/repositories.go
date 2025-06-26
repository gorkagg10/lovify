package profile

import (
	"context"

	"github.com/gorkagg10/lovify-user-service/internal/domain/oauth"
)

type UserRepository interface {
	CreateUserProfile(context.Context, *UserProfile) (string, error)
	ConnectWithMusicProvider(ctx context.Context, userID string) error
	StoreMusicProviderToken(ctx context.Context, userID string, token *oauth.Token) error
	StoreMusicProviderData(ctx context.Context, userID string, musicProviderData *MusicProviderData) error
}

type SecurityRepository interface {
	EncryptToken(string) (string, error)
	DecryptToken(string) (string, error)
}

type MusicProviderRepository interface {
	GetTopTracks(token string) ([]Track, error)
	GetTopArtists(token string) ([]Artist, error)
}
