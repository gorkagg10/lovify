package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/gorkagg10/lovify/lovify-matching-service/config"
	"github.com/gorkagg10/lovify/lovify-matching-service/database"
	service "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/infra/server"
)

func main() {
	port := 8083

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	conf, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	dbClient, err := database.Connect(ctx, conf.DatabaseConfig)
	if err != nil {
		slog.Error("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			slog.Error("failed to disconnect from database: %v", err)
			os.Exit(1)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("failed to listen", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("listening", slog.String("port", fmt.Sprintf(":%d", port)))

	userCollection := dbClient.Database("userService").Collection("profiles")
	musicProviderDataCollection := dbClient.Database("userService").Collection("musicProviderData")
	likeCollection := dbClient.Database("matchingService").Collection("likes")
	matchCollection := dbClient.Database("matchingService").Collection("matches")
	matchingServer := server.NewMatchingServer(
		userCollection,
		musicProviderDataCollection,
		likeCollection,
		matchCollection,
	)
	srv := SetupGrpcServer(matchingServer)

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

func SetupGrpcServer(matchingServer *server.MatchingServer) *grpc.Server {
	grpcServer := grpc.NewServer()
	service.RegisterMatchingServiceServer(grpcServer, matchingServer)
	return grpcServer
}
