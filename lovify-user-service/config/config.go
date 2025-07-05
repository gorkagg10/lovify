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
	DBHost                   = "DB_HOST"
	DBPort                   = "DB_PORT"
	NatsURL                  = "NATS_URL"

	DefaultDBHost = "localhost"
	DefaultDBPort = "27017"
)

var (
	EmptyEnvVariableErrMsg = "missing %s environment variable"
)

type Config struct {
	SpotifyOAuthConfig *SpotifyOAuthConfig
	DatabaseConfig     *DatabaseConfig
	NatsURL            string
}

func NewConfig() (*Config, error) {
	spotifyOAuthConfig, err := NewSpotifyOAuthConfig()
	if err != nil {
		return nil, err
	}
	databaseConfig, err := NewDatabaseConfig()
	if err != nil {
		return nil, err
	}
	natsURL := os.Getenv(NatsURL)
	if natsURL == "" {
		return nil, fmt.Errorf(EmptyEnvVariableErrMsg, NatsURL)
	}
	return &Config{
		SpotifyOAuthConfig: spotifyOAuthConfig,
		DatabaseConfig:     databaseConfig,
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

type DatabaseConfig struct {
	Host string
	Port string
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	hostname := os.Getenv(DBHost)
	if hostname == "" {
		hostname = DefaultDBHost
	}
	port := os.Getenv(DBPort)
	if port == "" {
		port = DefaultDBPort
	}

	return &DatabaseConfig{
		Host: hostname,
		Port: port,
	}, nil
}
