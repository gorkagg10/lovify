package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	spotifyoauth2 "golang.org/x/oauth2/spotify"
	"google.golang.org/grpc"

	"github.com/gorkagg10/lovify/lovify-user-service/config"
	"github.com/gorkagg10/lovify/lovify-user-service/database"
	service "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/oauth"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/profile"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/infra/aescgm"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/infra/mongodb"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/infra/server"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/infra/spotify"
)

func main() {
	port := 8082

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	conf, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", err)
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

	natsConnection, err := nats.Connect(conf.NatsURL)
	if err != nil {
		slog.Error("failed to connect to nats server: %v", err)
		os.Exit(1)
	}
	defer natsConnection.Drain()

	jetStream, err := jetstream.New(natsConnection)
	if err != nil {
		slog.Error("failed to connect to jetstream: %v", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("failed to listen", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("listening", slog.String("port", fmt.Sprintf(":%d", port)))

	userServer := setupUserServer(dbClient, conf.SpotifyOAuthConfig, jetStream, conf.UploadsDir)
	srv := SetupGrpcServer(userServer)

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

func SetupGrpcServer(userServer *server.UserServer) *grpc.Server {
	grpcServer := grpc.NewServer()
	service.RegisterUserServiceServer(grpcServer, userServer)
	return grpcServer
}

func NewSecureHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}
}

func setupUserServer(
	dbClient *mongo.Client,
	spotifyOAuthConfig *config.SpotifyOAuthConfig,
	jetStream jetstream.JetStream,
	uploadsDir string,
) *server.UserServer {
	userCollection := dbClient.Database("userService").Collection("profiles")
	musicProviderTokensCollection := dbClient.Database("userService").Collection("musicProviderTokens")
	musicProviderDataCollection := dbClient.Database("userService").Collection("musicProviderData")

	userRepository := mongodb.NewUserRepository(userCollection, musicProviderTokensCollection, musicProviderDataCollection)
	securityRepository := aescgm.NewSecurityRepository("Q1wMqfEdIkjKZfxNNlf3qVaFOSv5AiRKrPCpW3kBB4c=")
	musicProviderRepository := spotify.NewMusicProviderRepository(NewSecureHTTPClient())
	oAuthRepository := spotify.NewOAuthRepository(
		&oauth2.Config{
			ClientID:     spotifyOAuthConfig.ClientID,
			ClientSecret: spotifyOAuthConfig.ClientSecret,
			Endpoint:     spotifyoauth2.Endpoint,
			RedirectURL:  spotifyOAuthConfig.RedirectURL,
			Scopes:       []string{"user-read-email", "user-read-recently-played", "user-top-read"},
		},
	)

	profileManager := profile.NewManager(userRepository, securityRepository, musicProviderRepository, jetStream)
	oAuthService := oauth.NewService(oAuthRepository)
	return server.NewUserServer(profileManager, oAuthService, uploadsDir)
}
