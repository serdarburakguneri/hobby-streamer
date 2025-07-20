package auth

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/token"
)

type AuthRepository interface {
	Login(ctx context.Context, req *LoginRequest) (*token.Token, error)
	ValidateToken(ctx context.Context, req *TokenValidationRequest) (*TokenValidationResponse, error)
	RefreshToken(ctx context.Context, req *TokenRefreshRequest) (*token.Token, error)
}

type Service struct {
	repository AuthRepository
}

func NewService(repository AuthRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*token.Token, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repository.Login(ctx, req)
}

func (s *Service) ValidateToken(ctx context.Context, req *TokenValidationRequest) (*TokenValidationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repository.ValidateToken(ctx, req)
}

func (s *Service) RefreshToken(ctx context.Context, req *TokenRefreshRequest) (*token.Token, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repository.RefreshToken(ctx, req)
}
