package config

import "os"

const (
	DBHost = "DB_HOST"
	DBPort = "DB_PORT"

	DefaultDBHost = "localhost"
	DefaultDBPort = "27017"
)

type Config struct {
	DatabaseConfig *DatabaseConfig
}

func NewConfig() (*Config, error) {
	databaseConfig, err := NewDatabaseConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		DatabaseConfig: databaseConfig,
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
