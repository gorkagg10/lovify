package http

import (
	"context"
	"encoding/json"
	"github.com/gorkagg10/lovify/lovify-api/config"
	authServiceGrpc "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
	matchingServiceGrpc "github.com/gorkagg10/lovify/lovify-matching-service/grpc/matching-service"
	messagingServiceGrpc "github.com/gorkagg10/lovify/lovify-messaging-service"
	userServiceGrpc "github.com/gorkagg10/lovify/lovify-user-service/grpc/user-service"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Handler
type Handler struct {
	config                 *config.Config
	Router                 *mux.Router
	Server                 *http.Server
	AuthServiceClient      authServiceGrpc.AuthServiceClient
	UserServiceClient      userServiceGrpc.UserServiceClient
	MatchingServiceClient  matchingServiceGrpc.MatchingServiceClient
	MessagingServiceClient messagingServiceGrpc.MessagingServiceClient
}

// Response object
type Response struct {
	Message string `json:"message"`
}

// NewHandler - returns a pointer to a Handler
func NewHandler(
	config *config.Config,
	authServiceClient authServiceGrpc.AuthServiceClient,
	userServiceClient userServiceGrpc.UserServiceClient,
	matchingServiceClient matchingServiceGrpc.MatchingServiceClient,
) *Handler {
	slog.Info("setting up our handler")
	h := &Handler{
		config:                config,
		AuthServiceClient:     authServiceClient,
		UserServiceClient:     userServiceClient,
		MatchingServiceClient: matchingServiceClient,
	}

	h.Router = mux.NewRouter()
	// Sets up our middleware functions
	h.Router.Use(JSONMiddleware, TimeoutMiddleware)
	// Sets up CSRF token middleware for the /auth/ prefix
	//h.Router.PathPrefix("/users/").Subrouter().Use(AuthMiddleware)
	// setup the routes
	h.mapRoutes()

	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://lovify-app:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})

	handler := c.Handler(h.Router)
	h.Server = &http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      handler,
	}
	return h
}

// mapRoutes - sets up all the routes for our application
func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", h.AliveCheck).Methods("GET")
	h.Router.HandleFunc("/auth/register", h.Register).Methods("POST")
	h.Router.HandleFunc("/auth/login", h.Login).Methods("POST")
	h.Router.HandleFunc("/users/{user_id}/photos", h.StorePhotos).Methods("POST")
	h.Router.HandleFunc("/users/{user_id}/login/spotify", h.LoginSpotify)
	h.Router.HandleFunc("/callback/spotify", h.RegisterSpotify)
	h.Router.HandleFunc("/users", h.CreateUser).Methods("POST")
	h.Router.HandleFunc("/users/{user_id}/recommendations", h.GetRecommendations).Methods("GET")
	h.Router.HandleFunc("/users/{user_id}/matches", h.GetMatches).Methods("GET")
	h.Router.HandleFunc("/users/{from_id}/likes/{to_id}", h.HandleLike)
	h.Router.HandleFunc("/users/{user_id}/messages/{match_id}", h.SendMessage).Methods("POST")
	//h.Router.HandleFunc("/users/{user_id}/messages/{match_id}", h.ListMessages).Methods("GET")
	//h.Router.HandleFunc("/auth/protected", h.Protected).Methods("POST")
}

func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		Response{
			Message: "I am Alive!",
		}); err != nil {
		slog.Error("encoding response", slog.String("error", err.Error()))
		return
	}
}

/*
func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.ReadyCheck(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		Response{
			Message: "I am Ready!",
		}); err != nil {
		slog.Error("encoding response", slog.String("error", err.Error()))
		return
	}
}
*/

// Serve - gracefully serves our newly setup handler function
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			slog.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := h.Server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutting down the server", slog.String("error", err.Error()))
		return err
	}

	slog.Info("shutting down the server gracefully")
	return nil
}
