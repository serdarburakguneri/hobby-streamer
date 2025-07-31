package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(authHandler *AuthHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/auth/validate", authHandler.ValidateToken).Methods("POST")
	router.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	router.HandleFunc("/health", authHandler.Health).Methods("GET")

	return router
}
