package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorkagg10/lovify/lovify-authentication-service/events"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/gorkagg10/lovify/lovify-authentication-service/config"
	"github.com/gorkagg10/lovify/lovify-authentication-service/database"
	service "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/domain/login"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/infra/base64"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/infra/postgres"
	"github.com/gorkagg10/lovify/lovify-authentication-service/internal/infra/server"
)

func main() {
	port := 8081

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.NewConfig()
	if err != nil {
		slog.Error("loading configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = database.Migrate(conf.DatabaseConfig)
	if err != nil {
		slog.Error("migrating database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pgClient, err := database.NewDatabaseClient(ctx, conf.DatabaseConfig)
	if err != nil {
		slog.Error("creating database client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err = pgClient.Close(); err != nil {
			slog.Error("closing database client", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	natsHandler, err := events.NewNatsHandler(conf.NatsEndpoint)
	if err != nil {
		slog.Error("creating nats handler", slog.String("error", err.Error()))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("failed to listen", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("listening", slog.String("port", fmt.Sprintf(":%d", port)))

	userRepository := postgres.NewUserRepository(pgClient)
	tokenRepository := postgres.NewTokenRepository(pgClient)
	securityRepository := base64.NewSecurityRepository()

	authenticationService := login.NewAuthorization(userRepository, securityRepository, tokenRepository)
	authServer := server.NewAuthServer(authenticationService)

	srv := setupGrpcServer(authServer)

	go func() {
		if err := events.Listen(ctx, natsHandler, []events.Consumer{
			events.NewConsumer(
				events.CreateProfileConsumer, events.Workqueue, authenticationService.ProcessProfileCreation,
			),
		}); err != nil {
		}
	}()

	go func() {
		if err = srv.Serve(lis); err != nil {
			slog.Error("failed to serve", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	srv.GracefulStop()
	slog.Info("shutdown completed", slog.String("port", fmt.Sprintf(":%d", port)))
}

func setupGrpcServer(authServer *server.AuthServer) *grpc.Server {
	grpcServer := grpc.NewServer()
	service.RegisterAuthServiceServer(grpcServer, authServer)
	return grpcServer
}

func setupAuthServer(pgClient *sql.DB) *server.AuthServer {
	userRepository := postgres.NewUserRepository(pgClient)
	tokenRepository := postgres.NewTokenRepository(pgClient)
	securityRepository := base64.NewSecurityRepository()

	authenticationService := login.NewAuthorization(userRepository, securityRepository, tokenRepository)
	return server.NewAuthServer(authenticationService)
}
