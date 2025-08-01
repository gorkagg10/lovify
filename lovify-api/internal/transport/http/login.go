package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"

	autherrors "github.com/gorkagg10/lovify/lovify-authentication-service/errors"
	authServiceGrpc "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
)

const (
	SessionToken = "session_token"
	CSRFToken    = "csrf_token"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if len(email) < 8 || len(password) < 8 {
		http.Error(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	_, err := h.AuthServiceClient.RegisterUser(
		r.Context(),
		&authServiceGrpc.RegisterRequest{
			Email:    &email,
			Password: &password,
		},
	)
	if err != nil {
		var errorMessage string
		switch status.Code(err) {
		case codes.InvalidArgument:
			switch {
			case errors.Is(err, autherrors.ErrUserAlreadyExists):
				errorMessage = "User already exists"
			default:
				errorMessage = "Invalid email/password"
			}
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		case codes.Internal:
			errorMessage = "Internal server error"
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		default:
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User successfully registered")
}

type LoginResponse struct {
	IsProfileConnected bool   `json:"is_profile_connected"`
	ProfileID          string `json:"profile_id"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	loginResponse, err := h.AuthServiceClient.Login(
		r.Context(),
		&authServiceGrpc.LoginRequest{
			Email:    &email,
			Password: &password,
		},
	)
	if err != nil {
		var errorMessage string
		switch status.Code(err) {
		case codes.NotFound:
			errorMessage = "User not found"
			http.Error(w, errorMessage, http.StatusNotFound)
			return
		default:
			http.Error(w, "Error login user", http.StatusInternalServerError)
			return
		}
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionToken,
		Value:    loginResponse.SessionToken.GetToken(),
		Expires:  loginResponse.SessionToken.ExpirationDate.AsTime(),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})

	// Set CSRF token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFToken,
		Value:    loginResponse.CsrfToken.GetToken(),
		Expires:  loginResponse.CsrfToken.ExpirationDate.AsTime(),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false, // Needs to be accessible to the client side
	})

	loginAPIResponse := LoginResponse{
		IsProfileConnected: loginResponse.GetIsProfileConnected(),
		ProfileID:          loginResponse.GetProfileID(),
	}
	loginAPIResponseJSON, err := json.Marshal(loginAPIResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(loginAPIResponseJSON)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Fuerza expiración inmediata
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Eliminar cookie CSRF
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Devolver respuesta en JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
