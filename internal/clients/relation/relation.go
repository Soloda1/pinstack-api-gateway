package relation_client

import (
	"context"
	"log/slog"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/logger"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/relation/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// relationClient реализует интерфейс PostClient
// (название RelationClient, если потребуется, можно поменять в interface.go)
type relationClient struct {
	client pb.RelationServiceClient
	log    *logger.Logger
}

func NewRelationClient(conn *grpc.ClientConn, log *logger.Logger) RelationClient {
	return &relationClient{
		client: pb.NewRelationServiceClient(conn),
		log:    log,
	}
}

func (c *relationClient) Follow(ctx context.Context, followerID, followeeID int64) error {
	c.log.Info("Following user", slog.Int64("follower_id", followerID), slog.Int64("followee_id", followeeID))
	_, err := c.client.Follow(ctx, &pb.FollowRequest{
		FollowerId: followerID,
		FolloweeId: followeeID,
	})
	if err != nil {
		c.log.Error("Failed to follow", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return custom_errors.ErrValidationFailed
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}

func (c *relationClient) Unfollow(ctx context.Context, followerID, followeeID int64) error {
	c.log.Info("Unfollowing user", slog.Int64("follower_id", followerID), slog.Int64("followee_id", followeeID))
	_, err := c.client.Unfollow(ctx, &pb.UnfollowRequest{
		FollowerId: followerID,
		FolloweeId: followeeID,
	})
	if err != nil {
		c.log.Error("Failed to unfollow", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return custom_errors.ErrValidationFailed
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}

func (c *relationClient) GetFollowers(ctx context.Context, followeeID int64, limit, page int32) ([]int64, error) {
	c.log.Info("Getting followers", slog.Int64("followee_id", followeeID), slog.Int("limit", int(limit)), slog.Int("page", int(page)))
	resp, err := c.client.GetFollowers(ctx, &pb.GetFollowersRequest{
		FolloweeId: followeeID,
		Limit:      limit,
		Page:       page,
	})
	if err != nil {
		c.log.Error("Failed to get followers", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, custom_errors.ErrValidationFailed
			default:
				return nil, custom_errors.ErrExternalServiceError
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	return resp.FollowerIds, nil
}

func (c *relationClient) GetFollowees(ctx context.Context, followerID int64, limit, page int32) ([]int64, error) {
	c.log.Info("Getting followees", slog.Int64("follower_id", followerID), slog.Int("limit", int(limit)), slog.Int("page", int(page)))
	resp, err := c.client.GetFollowees(ctx, &pb.GetFolloweesRequest{
		FollowerId: followerID,
		Limit:      limit,
		Page:       page,
	})
	if err != nil {
		c.log.Error("Failed to get followees", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, custom_errors.ErrValidationFailed
			default:
				return nil, custom_errors.ErrExternalServiceError
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	return resp.FolloweeIds, nil
}
