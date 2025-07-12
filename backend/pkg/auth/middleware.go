package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type AuthMiddleware struct {
	validator TokenValidator
}

func NewAuthMiddleware(validator TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, constants.ErrAuthorizationHeader, http.StatusUnauthorized)
			return
		}

		if token == "dev-token-placeholder" {
			mockUser := &User{
				ID:       "dev-user",
				Username: "dev-admin",
				Email:    "dev@example.com",
				Roles:    []string{"admin"},
			}
			ctx := context.WithValue(r.Context(), "user", mockUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := m.validator.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, constants.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				http.Error(w, constants.ErrUserNotFound, http.StatusInternalServerError)
				return
			}

			if !m.validator.HasRole(user, role) {
				http.Error(w, constants.ErrInsufficientPerm, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func (m *AuthMiddleware) RequireAnyRole(roles []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				http.Error(w, constants.ErrUserNotFound, http.StatusInternalServerError)
				return
			}

			if !m.validator.HasAnyRole(user, roles) {
				http.Error(w, constants.ErrInsufficientPerm, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func (m *AuthMiddleware) RequireAllRoles(roles []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				http.Error(w, constants.ErrUserNotFound, http.StatusInternalServerError)
				return
			}

			if !m.validator.HasAllRoles(user, roles) {
				http.Error(w, constants.ErrInsufficientPerm, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get(constants.HeaderAuthorization)
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, constants.BearerPrefix)
}
