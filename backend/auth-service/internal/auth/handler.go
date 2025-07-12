package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AuthHandler struct {
	Service AuthService
	logger  *logger.Logger
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		Service: service,
		logger:  logger.WithService("auth-handler"),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid login request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.ClientID == "" {
		log.Warn("Missing required login fields", "username_provided", req.Username != "", "password_provided", req.Password != "", "client_id_provided", req.ClientID != "")
		http.Error(w, "Username, password, and client_id are required", http.StatusBadRequest)
		return
	}

	log.Debug("Processing login request", "username", req.Username, "client_id", req.ClientID)
	token, err := h.Service.Login(r.Context(), &req)
	if err != nil {
		log.WithError(err).Warn("Login failed", "username", req.Username, "client_id", req.ClientID)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Info("Login successful", "username", req.Username, "client_id", req.ClientID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid token validation request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		log.Warn("Missing token in validation request")
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(req.Token, "Bearer ")

	log.Debug("Validating token")
	response, err := h.Service.ValidateToken(r.Context(), tokenString)
	if err != nil {
		log.WithError(err).Warn("Token validation failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Token validation completed", "valid", response.Valid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid token refresh request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		log.Warn("Missing refresh token in request")
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	log.Debug("Processing token refresh")
	token, err := h.Service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		log.WithError(err).Warn("Token refresh failed")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Info("Token refresh successful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())
	log.Debug("Health check requested")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "auth-service"})
}

// TokenRefreshRequest represents a token refresh request
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
