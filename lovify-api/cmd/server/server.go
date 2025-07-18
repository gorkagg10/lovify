package main

import (
	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
	userServiceGrpc "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
	"log/slog"
	"os"

	"github.com/gorkagg10/lovify/lovify-api/config"
	"github.com/gorkagg10/lovify/lovify-api/grpc"
	transportHttp "github.com/gorkagg10/lovify/lovify-api/internal/transport/http"
	authServiceGrpc "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
)

func Run(
	config *config.Config,
	authServiceClient authServiceGrpc.AuthServiceClient,
	userServiceClient userServiceGrpc.UserServiceClient,
	matchingServiceClient matchingServiceGrpc.MatchingServiceClient,
) error {
	handler := transportHttp.NewHandler(config, authServiceClient, userServiceClient, matchingServiceClient)
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
	if err != nil {
		slog.Error("connecting to user service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	matchingServiceClient, err := grpc.NewClient(conf.MatchingServiceEndpoint, matchingServiceGrpc.NewMatchingServiceClient)
	if err != nil {
		slog.Error("connecting to matching service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if err = Run(conf, authServiceClient, userServiceClient, matchingServiceClient); err != nil {
		slog.Error("running the application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
