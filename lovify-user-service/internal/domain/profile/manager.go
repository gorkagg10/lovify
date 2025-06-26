package profile

import (
	"context"
	"encoding/base64"
	userServiceGrpc "github.com/gorkagg10/lovify-user-service/grpc/user-service"
	"github.com/gorkagg10/lovify-user-service/internal/domain/oauth"
)

type Manager struct {
	userRepository          UserRepository
	securityRepository      SecurityRepository
	musicProviderRepository MusicProviderRepository
}

func NewManager(
	userRepository UserRepository,
	securityRepository SecurityRepository,
	musicProviderRepository MusicProviderRepository,
) *Manager {
	return &Manager{
		userRepository:          userRepository,
		securityRepository:      securityRepository,
		musicProviderRepository: musicProviderRepository,
	}
}

func (m *Manager) CreateUserProfile(ctx context.Context, req *userServiceGrpc.CreateUserRequest) (string, error) {
	userProfile := NewUserProfile(
		req.GetEmail(),
		req.GetBirthday().AsTime(),
		req.GetName(),
		req.GetGender().String(),
		req.GetSexualOrientation().String(),
		req.GetDescription(),
	)
	return m.userRepository.CreateUserProfile(ctx, userProfile)
}

func (m *Manager) encryptToken(token *oauth.Token) (*oauth.Token, error) {
	encryptedAccessToken, err := m.securityRepository.EncryptToken(token.AccessToken())
	if err != nil {
		return nil, err
	}
	encryptedRefreshToken, err := m.securityRepository.EncryptToken(token.RefreshToken())
	if err != nil {
		return nil, err
	}
	return oauth.NewToken(
		encryptedAccessToken,
		encryptedRefreshToken,
		token.ExpiresAt(),
	), nil
}

func (m *Manager) getMusicProviderData(token *oauth.Token) (*MusicProviderData, error) {
	topTracks, err := m.musicProviderRepository.GetTopTracks(token.AccessToken())
	if err != nil {
		return nil, err
	}
	topArtists, err := m.musicProviderRepository.GetTopArtists(token.AccessToken())
	if err != nil {
		return nil, err
	}
	return NewMusicProviderData(topTracks, topArtists), nil
}

func (m *Manager) ConnectWithMusicProvider(ctx context.Context, state string, token *oauth.Token) error {
	userID, err := getUserID(state)
	if err != nil {
		return err
	}
	err = m.userRepository.ConnectWithMusicProvider(ctx, userID)
	if err != nil {
		return err
	}

	musicProviderData, err := m.getMusicProviderData(token)
	if err != nil {
		return err
	}
	err = m.userRepository.StoreMusicProviderData(ctx, userID, musicProviderData)
	if err != nil {
		return err
	}

	token, err = m.encryptToken(token)
	if err != nil {
		return err
	}
	err = m.userRepository.StoreMusicProviderToken(ctx, userID, token)
	if err != nil {
		return err
	}

	return nil
}

func getUserID(state string) (string, error) {
	stateBytes, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}
	return string(stateBytes), nil
}
