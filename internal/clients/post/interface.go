package post_client

import (
	"context"
	"pinstack-api-gateway/internal/models"
)

type PostClient interface {
	CreatePost(ctx context.Context, post *models.CreatePostDTO) (*models.PostDetailed, error)
	GetPostByID(ctx context.Context, id int64) (*models.PostDetailed, error)
	ListPosts(ctx context.Context, filters *models.PostFilters) ([]*models.PostDetailed, int64, error)
	UpdatePost(ctx context.Context, id int64, post *models.UpdatePostDTO) error
	DeletePost(ctx context.Context, userID int64, id int64) error
}
