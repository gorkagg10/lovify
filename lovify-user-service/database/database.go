package database

import (
	"context"
	"fmt"
	"github.com/gorkagg10/lovify/lovify-user-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, databaseConfig *config.DatabaseConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", databaseConfig.Host, databaseConfig.Port)))
	if err != nil {
		return nil, err
	}
	return client, nil
}
