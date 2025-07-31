package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/models"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/service"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.Get().WithService("auth-handler"),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		ClientID string `json:"client_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid login request body")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" || req.ClientID == "" {
		log.Warn("Missing required login fields", "username_provided", req.Username != "", "password_provided", req.Password != "", "client_id_provided", req.ClientID != "")
		h.writeError(w, http.StatusBadRequest, "Username, password, and client_id are required")
		return
	}

	loginReq, err := models.NewLoginRequest(req.Username, req.Password, req.ClientID)
	if err != nil {
		log.WithError(err).Warn("Invalid login request fields")
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Debug("Processing login request", "username", req.Username, "client_id", req.ClientID)

	token, err := h.authService.Login(r.Context(), loginReq)
	if err != nil {
		h.handleError(w, err, "Login failed")
		return
	}

	response := TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	log.Info("Login successful", "username", req.Username, "client_id", req.ClientID)
	h.writeJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid validate token request body")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		h.writeError(w, http.StatusBadRequest, "Token is required")
		return
	}

	validationReq, err := models.NewTokenValidationRequest(req.Token)
	if err != nil {
		log.WithError(err).Warn("Invalid token value")
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Debug("Validating token")

	validation, err := h.authService.ValidateToken(r.Context(), validationReq)
	if err != nil {
		h.handleError(w, err, "Token validation failed")
		return
	}

	var responseUser *UserResponse
	if validation.User != nil {
		responseUser = &UserResponse{
			ID:       validation.User.ID,
			Username: validation.User.Username,
			Email:    validation.User.Email,
			Roles:    validation.User.Roles,
			Enabled:  validation.User.Enabled,
		}
	}

	response := TokenValidationResponse{
		Valid:   validation.IsValid,
		User:    responseUser,
		Message: validation.Message,
		Roles:   validation.User.Roles,
	}

	log.Info("Token validation successful")
	h.writeJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid refresh token request body")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		h.writeError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	refreshReq, err := models.NewTokenRefreshRequest(req.RefreshToken)
	if err != nil {
		log.WithError(err).Warn("Invalid refresh token value")
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Debug("Refreshing token")

	token, err := h.authService.RefreshToken(r.Context(), refreshReq)
	if err != nil {
		h.handleError(w, err, "Token refresh failed")
		return
	}

	response := TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	log.Info("Token refresh successful")
	h.writeJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())
	log.Debug("Health check requested")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "auth-service"}); err != nil {
		log.WithError(err).Error("Failed to encode health check response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
	if apperrors.IsAppError(err) {
		appErr := err.(*apperrors.AppError)
		h.logger.WithError(err).Error("Application error", "error_type", appErr.Type, "context", appErr.Context)

		status := appErr.HTTPStatus()
		message := appErr.Message
		if appErr.Type == apperrors.ErrorTypeCircuitBreaker {
			message = "Service temporarily unavailable"
		}

		h.writeError(w, status, message)
		return
	}

	h.logger.WithError(err).Error("Unexpected error")
	h.writeError(w, http.StatusInternalServerError, defaultMessage)
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{
		"error": message,
	})
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Enabled  bool     `json:"enabled"`
}

type TokenValidationResponse struct {
	Valid   bool          `json:"valid"`
	User    *UserResponse `json:"user,omitempty"`
	Message string        `json:"message,omitempty"`
	Roles   []string      `json:"roles,omitempty"`
}
