package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	SpotifyOAuthClientID     = "SPOTIFY_OAUTH_CLIENT_ID"
	SpotifyOAuthClientSecret = "SPOTIFY_OAUTH_CLIENT_SECRET"
	SpotifyOAuthRedirectURL  = "SPOTIFY_OAUTH_REDIRECT_URL"
)

var (
	EmptyEnvVariableErrMsg = "missing %s environment variable"
)

type Config struct {
	SpotifyOAuthConfig *SpotifyOAuthConfig
}

func NewConfig() (*Config, error) {
	spotifyOAuthConfig, err := NewSpotifyOAuthConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		SpotifyOAuthConfig: spotifyOAuthConfig,
	}, nil
}

type SpotifyOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewSpotifyOAuthConfig() (*SpotifyOAuthConfig, error) {
	clientID := os.Getenv(SpotifyOAuthClientID)
	if clientID == "" {
		return nil, errors.New(fmt.Sprintf(EmptyEnvVariableErrMsg, SpotifyOAuthClientID))
	}
	clientSecret := os.Getenv(SpotifyOAuthClientSecret)
	if clientSecret == "" {
		return nil, errors.New(fmt.Sprintf(EmptyEnvVariableErrMsg, SpotifyOAuthClientSecret))
	}
	redirectURL := os.Getenv(SpotifyOAuthRedirectURL)
	if redirectURL == "" {
		return nil, errors.New(fmt.Sprintf(EmptyEnvVariableErrMsg, SpotifyOAuthRedirectURL))
	}
	return &SpotifyOAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}, nil
}
