package token

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, token *Token) error
	FindByAccessToken(ctx context.Context, accessToken AccessToken) (*Token, error)
	FindByRefreshToken(ctx context.Context, refreshToken RefreshToken) (*Token, error)
	FindExpiredTokens(ctx context.Context) ([]*Token, error)
	Delete(ctx context.Context, accessToken AccessToken) error
	DeleteExpiredTokens(ctx context.Context) error
}
