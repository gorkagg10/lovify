package main

import (
	userServiceGrpc "github.com/gorkagg10/lovify-user-service/grpc/user-service"
	"log/slog"
	"os"

	"github.com/gorkagg10/lovify-api/config"
	"github.com/gorkagg10/lovify-api/grpc"
	transportHttp "github.com/gorkagg10/lovify-api/internal/transport/http"
	authServiceGrpc "github.com/gorkagg10/lovify-authentication-service/grpc/auth-service"
)

func Run(
	config *config.Config,
	authServiceClient authServiceGrpc.AuthServiceClient,
	userServiceClient userServiceGrpc.UserServiceClient) error {
	handler := transportHttp.NewHandler(config, authServiceClient, userServiceClient)
	if err := handler.Serve(); err != nil {
		slog.Error("gracefully serving the application")
		return err
	}
	return nil
}

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		slog.Error("loading configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	authServiceClient, err := grpc.NewClient(conf.AuthServiceEndpoint, authServiceGrpc.NewAuthServiceClient)
	if err != nil {
		slog.Error("connecting to auth service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	userServiceClient, err := grpc.NewClient(conf.UserServiceEndpoint, userServiceGrpc.NewUserServiceClient)
	if err = Run(conf, authServiceClient, userServiceClient); err != nil {
		slog.Error("running the application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
