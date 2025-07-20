package http

import (
	"github.com/gorilla/mux"
	appauth "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/application/auth"
)

func NewRouter(authService *appauth.Service) *mux.Router {
	r := mux.NewRouter()

	authHandler := NewHandler(authService)

	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/auth/validate", authHandler.ValidateToken).Methods("POST")
	r.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	r.HandleFunc("/health", authHandler.Health).Methods("GET")

	return r
}
