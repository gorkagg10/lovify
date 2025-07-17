package config

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	Port                = "PORT"
	AuthServiceEndpoint = "AUTH_SERVICE_ENDPOINT"
	UserServiceEndpoint = "USER_SERVICE_ENDPOINT"
	FrontEndHost        = "FRONT_END_HOST"
)

type Config struct {
	Port                string
	AuthServiceEndpoint string
	UserServiceEndpoint string
	FrontEndHost        string
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
	frontEndHost := os.Getenv(FrontEndHost)
	if frontEndHost == "" {
		frontEndHost = "http://localhost:3000"
	}
	return &Config{
		Port:                port,
		AuthServiceEndpoint: authServiceEndpoint,
		UserServiceEndpoint: userServiceEndpoint,
		FrontEndHost:        frontEndHost,
	}, nil
}
