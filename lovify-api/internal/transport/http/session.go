package http

import (
	"errors"
	"github.com/gorkagg10/lovify/lovify-api/internal/database"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, ok := database.Users[username]
	if !ok {
		return AuthError
	}

	// Get the Session Token from the cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return AuthError
	}

	// Get the CSRF token from the headers
	csrf, err := r.Cookie("csrf_token")
	if err != nil || csrf.Value != user.CSRFToken || csrf.Value == "" {
		return AuthError
	}
	return nil
}
