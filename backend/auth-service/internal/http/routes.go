package http

import (
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/auth"
)

func NewRouter(authService auth.AuthService) *mux.Router {
	r := mux.NewRouter()

	authHandler := &auth.AuthHandler{Service: authService}

	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/auth/validate", authHandler.ValidateToken).Methods("POST")
	r.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	r.HandleFunc("/health", authHandler.Health).Methods("GET")

	return r
}
