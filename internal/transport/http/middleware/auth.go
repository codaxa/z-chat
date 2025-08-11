// Package middleware provides HTTP middleware functionalities for the chat application.
package middleware

import (
	"context"
	"net/http"
	appContext "z-chat/internal/context"
	"z-chat/internal/services"
)

// AuthMiddleware holds the auth service dependency
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new auth middleware with the provided auth service
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate checks the request for a valid JWT token and extracts user claims
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		claims, err := am.authService.ValidateToken(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), appContext.UserClaimsKey, claims))
		next.ServeHTTP(w, r)
	})
}
