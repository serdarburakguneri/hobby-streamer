package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type contextKey string

const (
	userContextKey        contextKey = "user"
	serviceUserContextKey contextKey = "service_user"
)

type AuthMiddleware struct {
	validator   TokenValidator
	userAuth    bool
	serviceAuth bool
}

func NewAuthMiddleware(validator TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{
		validator:   validator,
		userAuth:    false,
		serviceAuth: false,
	}
}

func (m *AuthMiddleware) RequireUserAuth() *AuthMiddleware {
	m.userAuth = true
	return m
}

func (m *AuthMiddleware) RequireServiceAuth() *AuthMiddleware {
	m.serviceAuth = true
	return m
}

func (m *AuthMiddleware) Build() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			token := extractToken(r)
			if token == "" {
				http.Error(w, constants.ErrAuthorizationHeader, http.StatusUnauthorized)
				return
			}

			if m.userAuth && m.serviceAuth {
				user, err := m.validator.ValidateToken(r.Context(), token)
				if err == nil {
					ctx := context.WithValue(r.Context(), userContextKey, user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				keycloakValidator, ok := m.validator.(*KeycloakValidator)
				if !ok {
					http.Error(w, constants.ErrInvalidToken, http.StatusUnauthorized)
					return
				}

				serviceValidator := NewServiceTokenValidator(keycloakValidator.keycloakURL, keycloakValidator.realm, keycloakValidator.clientID)
				serviceUser, serviceErr := serviceValidator.ValidateServiceToken(r.Context(), token)
				if serviceErr != nil {
					http.Error(w, constants.ErrInvalidToken, http.StatusUnauthorized)
					return
				}

				if !serviceValidator.IsServiceToken(serviceUser) {
					http.Error(w, "Invalid service token", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), serviceUserContextKey, serviceUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if m.userAuth {
				user, err := m.validator.ValidateToken(r.Context(), token)
				if err != nil {
					http.Error(w, constants.ErrInvalidToken, http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), userContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if m.serviceAuth {
				keycloakValidator, ok := m.validator.(*KeycloakValidator)
				if !ok {
					http.Error(w, "Service auth not supported with this validator", http.StatusInternalServerError)
					return
				}

				serviceValidator := NewServiceTokenValidator(keycloakValidator.keycloakURL, keycloakValidator.realm, keycloakValidator.clientID)
				serviceUser, err := serviceValidator.ValidateServiceToken(r.Context(), token)
				if err != nil {
					http.Error(w, constants.ErrInvalidToken, http.StatusUnauthorized)
					return
				}

				if !serviceValidator.IsServiceToken(serviceUser) {
					http.Error(w, "Invalid service token", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), serviceUserContextKey, serviceUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "No auth requirements configured", http.StatusInternalServerError)
		}
	}
}

func (m *AuthMiddleware) RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userContextKey).(*User)
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
			user, ok := r.Context().Value(userContextKey).(*User)
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
			user, ok := r.Context().Value(userContextKey).(*User)
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
