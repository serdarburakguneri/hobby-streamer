package config

import (
	"net/http"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
)

type AuthConfig struct {
	Middleware func(http.Handler) http.Handler
}

func NewAuthConfig(configManager *config.Manager) (*AuthConfig, error) {
	dynamicCfg := configManager.GetDynamicConfig()

	keycloakURL := dynamicCfg.GetStringFromComponent("keycloak", "url")
	keycloakRealm := dynamicCfg.GetStringFromComponent("keycloak", "realm")
	keycloakClientID := dynamicCfg.GetStringFromComponent("keycloak", "client_id")

	keycloakValidator := auth.NewKeycloakValidator(keycloakURL, keycloakRealm, keycloakClientID)
	authMiddleware := auth.NewAuthMiddleware(keycloakValidator)

	authHandlerFunc := authMiddleware.RequireUserAuth().RequireServiceAuth().Build()

	return &AuthConfig{
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHandlerFunc(next.ServeHTTP)(w, r)
			})
		},
	}, nil
}
