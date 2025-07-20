package http

import (
	"encoding/json"
	"net/http"
	"time"

	appauth "github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/application/auth"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Handler struct {
	authService *appauth.Service
	logger      *logger.Logger
}

func NewHandler(authService *appauth.Service) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger.Get().WithService("auth-handler"),
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	loginReq := appauth.NewLoginRequest(req.Username, req.Password, req.ClientID)
	log.Debug("Processing login request", "username", req.Username, "client_id", req.ClientID)

	domainToken, err := h.authService.Login(r.Context(), loginReq)
	if err != nil {
		h.handleError(w, err, "Login failed")
		return
	}

	response := TokenResponse{
		AccessToken:  domainToken.AccessToken().Value(),
		TokenType:    domainToken.TokenType().Value(),
		ExpiresIn:    domainToken.ExpiresIn().Value(),
		RefreshToken: domainToken.RefreshToken().Value(),
		ExpiresAt:    domainToken.ExpiresAt().Value(),
	}

	log.Info("Login successful", "username", req.Username, "client_id", req.ClientID)
	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
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

	validationReq := appauth.NewTokenValidationRequest(req.Token)
	log.Debug("Validating token")

	validation, err := h.authService.ValidateToken(r.Context(), validationReq)
	if err != nil {
		h.handleError(w, err, "Token validation failed")
		return
	}

	var responseUser *UserResponse
	if validation.User() != nil {
		// Convert domain roles to string slice for JSON response
		var roleStrings []string
		for _, role := range validation.User().Roles().Values() {
			roleStrings = append(roleStrings, role.Value())
		}

		responseUser = &UserResponse{
			ID:       validation.User().ID().Value(),
			Username: validation.User().Username().Value(),
			Email:    validation.User().Email().Value(),
			Roles:    roleStrings,
			Enabled:  validation.User().Enabled(),
		}
	}

	response := TokenValidationResponse{
		Valid:   validation.Valid(),
		User:    responseUser,
		Message: validation.Message(),
		Roles:   validation.Roles(),
	}

	log.Info("Token validation successful")
	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
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

	refreshReq := appauth.NewTokenRefreshRequest(req.RefreshToken)
	log.Debug("Refreshing token")

	domainToken, err := h.authService.RefreshToken(r.Context(), refreshReq)
	if err != nil {
		h.handleError(w, err, "Token refresh failed")
		return
	}

	response := TokenResponse{
		AccessToken:  domainToken.AccessToken().Value(),
		TokenType:    domainToken.TokenType().Value(),
		ExpiresIn:    domainToken.ExpiresIn().Value(),
		RefreshToken: domainToken.RefreshToken().Value(),
		ExpiresAt:    domainToken.ExpiresAt().Value(),
	}

	log.Info("Token refresh successful")
	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	log := h.logger.WithContext(r.Context())
	log.Debug("Health check requested")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "auth-service"}); err != nil {
		log.WithError(err).Error("Failed to encode health check response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
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

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
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
