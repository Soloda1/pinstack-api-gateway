package auth_client

import (
	"context"
	"pinstack-api-gateway/internal/models"
)

type AuthClient interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.TokenPair, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	UpdatePassword(ctx context.Context, req *models.UpdatePasswordRequest) error
}
