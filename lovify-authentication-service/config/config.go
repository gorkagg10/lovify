package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	MigrationsPath = "MIGRATIONS_PATH"
	DBHost         = "DB_HOST"
	DBPort         = "DB_PORT"
	DBUser         = "DB_USER"
	DBPassword     = "DB_PASSWORD"
	DBName         = "DB_NAME"
	SSLMode        = "DB_SSL_MODE"

	DefaultDBHost    = "localhost"
	DefaultDBPort    = "5432"
	DefaultDBUser    = "postgres"
	DefaultDBName    = "postgres"
	DefaultDBSSLMode = "disable"
)

var (
	emptyDBPassword     = errors.New("empty database password")
	emptyMigrationsPath = errors.New("empty database migrations path")
)

type Config struct {
	DatabaseConfig *DatabaseConfig
}

func NewConfig() (*Config, error) {
	databaseConfig, err := NewDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("loading database config: %w", err)
	}
	return &Config{
		DatabaseConfig: databaseConfig,
	}, nil
}

type DatabaseConfig struct {
	MigrationsPath string
	Host           string
	Port           string
	Username       string
	Password       string
	Database       string
	SSLMode        string
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	migrationsPath := os.Getenv(MigrationsPath)
	if migrationsPath == "" {
		return nil, emptyMigrationsPath
	}
	hostname := os.Getenv(DBHost)
	if hostname == "" {
		hostname = DefaultDBHost
	}
	port := os.Getenv(DBPort)
	if port == "" {
		port = DefaultDBPort
	}
	username := os.Getenv(DBUser)
	if username == "" {
		username = DefaultDBUser
	}
	password := os.Getenv(DBPassword)
	if password == "" {
		return nil, emptyDBPassword
	}
	databaseName := os.Getenv(DBName)
	if databaseName == "" {
		databaseName = DefaultDBName
	}
	sslMode := os.Getenv(SSLMode)
	if sslMode == "" {
		sslMode = DefaultDBSSLMode
	}

	return &DatabaseConfig{
		MigrationsPath: migrationsPath,
		Host:           hostname,
		Port:           port,
		Username:       username,
		Password:       password,
		Database:       databaseName,
		SSLMode:        sslMode,
	}, nil
}
