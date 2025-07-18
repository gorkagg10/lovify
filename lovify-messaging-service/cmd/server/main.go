package server

import (
	"context"
	"fmt"
	"github.com/gorkagg10/lovify/lovify-messaging-service/config"
	"github.com/gorkagg10/lovify/lovify-messaging-service/database"
	service "github.com/gorkagg10/lovify/lovify-messaging-service/grpc/messaging-service"
	"github.com/gorkagg10/lovify/lovify-messaging-service/internal/infra/server"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"syscall"
)

func main() {
	port := 8084

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

	matchCollection := dbClient.Database("matchingService").Collection("matches")
	messagesCollection := dbClient.Database("messagingService").Collection("messages")
	matchingServer := server.NewMessagingServer(
		matchCollection,
		messagesCollection,
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

func SetupGrpcServer(messagingServer *server.MessagingServer) *grpc.Server {
	grpcServer := grpc.NewServer()
	service.RegisterMessagingServiceServer(grpcServer, messagingServer)
	return grpcServer
}
