package config

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	Port                    = "PORT"
	AuthServiceEndpoint     = "AUTH_SERVICE_ENDPOINT"
	UserServiceEndpoint     = "USER_SERVICE_ENDPOINT"
	MatchingServiceEndpoint = "MATCHING_SERVICE_ENDPOINT"
	MessageServiceEndpoint  = "MESSAGE_SERVICE_ENDPOINT"
	FrontEndHost            = "FRONT_END_HOST"
)

type Config struct {
	Port                     string
	AuthServiceEndpoint      string
	UserServiceEndpoint      string
	MatchingServiceEndpoint  string
	MessagingServiceEndpoint string
	FrontEndHost             string
}

func NewConfig() (*Config, error) {
	port := os.Getenv(Port)
	if port == "" {
		slog.Debug("PORT not set, defaulting to 8080")
		port = "8080"
	}
	authServiceEndpoint := os.Getenv(AuthServiceEndpoint)
	if authServiceEndpoint == "" {
		return nil, fmt.Errorf("%s must be set", AuthServiceEndpoint)
	}
	userServiceEndpoint := os.Getenv(UserServiceEndpoint)
	if userServiceEndpoint == "" {
		return nil, fmt.Errorf("%s must be set", UserServiceEndpoint)
	}
	matchingServiceEndpoint := os.Getenv(MatchingServiceEndpoint)
	if matchingServiceEndpoint == "" {
		return nil, fmt.Errorf("%s must be set", MatchingServiceEndpoint)
	}
	messageServiceEndpoint := os.Getenv(MessageServiceEndpoint)
	if messageServiceEndpoint == "" {
		return nil, fmt.Errorf("%s must be set", MessageServiceEndpoint)
	}
	frontEndHost := os.Getenv(FrontEndHost)
	if frontEndHost == "" {
		frontEndHost = "http://localhost:3000"
	}
	return &Config{
		Port:                     port,
		AuthServiceEndpoint:      authServiceEndpoint,
		UserServiceEndpoint:      userServiceEndpoint,
		MatchingServiceEndpoint:  matchingServiceEndpoint,
		MessagingServiceEndpoint: messageServiceEndpoint,
		FrontEndHost:             frontEndHost,
	}, nil
}
