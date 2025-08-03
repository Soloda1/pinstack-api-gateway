package relation_client

import (
	"context"
	"pinstack-api-gateway/internal/models"
)

type RelationClient interface {
	Follow(ctx context.Context, followerID, followeeID int64) error
	Unfollow(ctx context.Context, followerID, followeeID int64) error
	GetFollowers(ctx context.Context, followeeID int64, limit, page int32) ([]*models.RelationUser, int64, error)
	GetFollowees(ctx context.Context, followerID int64, limit, page int32) ([]*models.RelationUser, int64, error)
}
