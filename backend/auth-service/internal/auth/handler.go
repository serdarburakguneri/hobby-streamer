package auth

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
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
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" || req.ClientID == "" {
		log.Warn("Missing required login fields", "username_provided", req.Username != "", "password_provided", req.Password != "", "client_id_provided", req.ClientID != "")
		h.writeError(w, http.StatusBadRequest, "Username, password, and client_id are required")
		return
	}

	log.Debug("Processing login request", "username", req.Username, "client_id", req.ClientID)
	token, err := h.Service.Login(r.Context(), &req)
	if err != nil {
		h.handleError(w, err, "Login failed")
		return
	}

	log.Info("Login successful", "username", req.Username, "client_id", req.ClientID)
	h.writeJSON(w, http.StatusOK, token)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid validate token request body")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		h.writeError(w, http.StatusBadRequest, "Token is required")
		return
	}

	log.Debug("Validating token")
	validation, err := h.Service.ValidateToken(r.Context(), req.Token)
	if err != nil {
		h.handleError(w, err, "Token validation failed")
		return
	}

	log.Info("Token validation successful")
	h.writeJSON(w, http.StatusOK, validation)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())

	var req TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid refresh token request body")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		h.writeError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	log.Debug("Refreshing token")
	token, err := h.Service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(w, err, "Token refresh failed")
		return
	}

	log.Info("Token refresh successful")
	h.writeJSON(w, http.StatusOK, token)
}

func (h *AuthHandler) Health(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())
	log.Debug("Health check requested")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "auth-service"})
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
	if apperrors.IsAppError(err) {
		appErr := err.(*apperrors.AppError)
		h.logger.WithError(err).Error("Application error", "error_type", appErr.Type, "context", appErr.Context)

		switch appErr.Type {
		case apperrors.ErrorTypeValidation:
			h.writeError(w, http.StatusBadRequest, appErr.Message)
			return
		case apperrors.ErrorTypeUnauthorized:
			h.writeError(w, http.StatusUnauthorized, appErr.Message)
			return
		case apperrors.ErrorTypeForbidden:
			h.writeError(w, http.StatusForbidden, appErr.Message)
			return
		case apperrors.ErrorTypeNotFound:
			h.writeError(w, http.StatusNotFound, appErr.Message)
			return
		case apperrors.ErrorTypeConflict:
			h.writeError(w, http.StatusConflict, appErr.Message)
			return
		case apperrors.ErrorTypeTransient:
			h.writeError(w, http.StatusServiceUnavailable, appErr.Message)
			return
		case apperrors.ErrorTypeTimeout:
			h.writeError(w, http.StatusGatewayTimeout, appErr.Message)
			return
		case apperrors.ErrorTypeExternal:
			h.writeError(w, http.StatusBadGateway, appErr.Message)
			return
		default:
			h.writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
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

// TokenRefreshRequest represents a token refresh request
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
