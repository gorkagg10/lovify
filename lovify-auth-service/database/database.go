package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gorkagg10/lovify-authentication-service/config"
)

const migrationsPath = "database/migrations"

func Migrate(databaseConfig *config.DatabaseConfig) error {
	migration, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			databaseConfig.Username, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port,
			databaseConfig.Database, databaseConfig.SSLMode),
	)
	if err != nil {
		return fmt.Errorf("loading migration files: %w", err)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	sourceError, dbError := migration.Close()
	if sourceError != nil {
		return sourceError
	}
	if dbError != nil {
		return dbError
	}
	return nil
}

func NewDatabaseClient(ctx context.Context, databaseConfig *config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.Username, databaseConfig.Password,
		databaseConfig.Database, databaseConfig.SSLMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging database connection: %w", err)
	}
	return db, nil
}
