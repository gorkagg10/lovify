package http

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"

	autherrors "github.com/gorkagg10/lovify-authentication-service/errors"
	authServiceGrpc "github.com/gorkagg10/lovify-authentication-service/grpc/auth-service"
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
		HttpOnly: true,
	})

	// Set CSRF token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFToken,
		Value:    loginResponse.CsrfToken.GetToken(),
		Expires:  loginResponse.CsrfToken.ExpirationDate.AsTime(),
		HttpOnly: false, // Needs to be accessible to the client side
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful!")
}

func (h *Handler) Protected(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintln(w, "Protected route")
}
