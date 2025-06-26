package mongodb

import (
	"time"
)

type MusicProviderToken struct {
	UserID       string    `bson:"user_id"`
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at"`
}

func NewMusicProviderToken(
	UserProfileID string,
	AccessToken string,
	RefreshToken string,
	ExpiresAt time.Time) *MusicProviderToken {
	return &MusicProviderToken{
		UserID:       UserProfileID,
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		ExpiresAt:    ExpiresAt,
	}
}
