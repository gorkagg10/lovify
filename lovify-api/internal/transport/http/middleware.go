package http

import (
	"context"
	authServiceGrpc "github.com/gorkagg10/lovify/lovify-authentication-service/grpc/auth-service"
	"net/http"
	"time"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(client authServiceGrpc.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			email := r.FormValue("email")

			// Get the Session Token from the cookie
			st, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Token requerido", http.StatusUnauthorized)
				return

			}

			// Get the CSRF token from the headers
			csrf, err := r.Cookie("csrf_token")
			if err != nil {
				http.Error(w, "Token requerido", http.StatusUnauthorized)
				return

			}

			_, err = client.Authorize(
				r.Context(),
				&authServiceGrpc.AuthorizationRequest{
					Email:        &email,
					CsrfToken:    &csrf.Value,
					SessionToken: &st.Value,
				},
			)
			if err != nil {
				http.Error(w, "Token requerido", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

/*
func CSRFTokenMiddleware(authKey []byte) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return csrf.Protect(
			authKey,
			csrf.Secure(false),
		)(next)
	}
}
*/

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
