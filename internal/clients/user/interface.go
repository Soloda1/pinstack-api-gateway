package user_client

import (
	"context"
	"pinstack-api-gateway/internal/models"
)

type UserClient interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SearchUsers(ctx context.Context, query string, page, limit int) ([]*models.User, int64, error)
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}
