package relation_client

import (
	"context"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/models"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/relation/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
				errMsg := st.Message()
				if errMsg == custom_errors.ErrSelfFollow.Error() {
					return custom_errors.ErrSelfFollow
				}
				return custom_errors.ErrValidationFailed
			case codes.AlreadyExists:
				return custom_errors.ErrAlreadyFollowing
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			case codes.Internal:
				return custom_errors.ErrExternalServiceError
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
				errMsg := st.Message()
				if errMsg == custom_errors.ErrSelfUnfollow.Error() {
					return custom_errors.ErrSelfUnfollow
				}
				return custom_errors.ErrValidationFailed
			case codes.NotFound:
				errMsg := st.Message()
				if errMsg == custom_errors.ErrFollowRelationNotFound.Error() {
					return custom_errors.ErrFollowRelationNotFound
				}
				return custom_errors.ErrUserNotFound
			case codes.Internal:
				return custom_errors.ErrExternalServiceError
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}

func (c *relationClient) GetFollowers(ctx context.Context, followeeID int64, limit, page int32) ([]*models.RelationUser, int64, error) {
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
				return nil, 0, custom_errors.ErrValidationFailed
			case codes.NotFound:
				return nil, 0, custom_errors.ErrUserNotFound
			case codes.Internal:
				errMsg := st.Message()
				if errMsg == custom_errors.ErrDatabaseQuery.Error() {
					return nil, 0, custom_errors.ErrDatabaseQuery
				}
				return nil, 0, custom_errors.ErrExternalServiceError
			default:
				return nil, 0, custom_errors.ErrExternalServiceError
			}
		}
		return nil, 0, custom_errors.ErrExternalServiceError
	}

	followers := make([]*models.RelationUser, len(resp.Followers))
	for i, user := range resp.Followers {
		followers[i] = models.RelationUserFromProto(user)
	}

	return followers, resp.Total, nil
}

func (c *relationClient) GetFollowees(ctx context.Context, followerID int64, limit, page int32) ([]*models.RelationUser, int64, error) {
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
				return nil, 0, custom_errors.ErrValidationFailed
			case codes.NotFound:
				return nil, 0, custom_errors.ErrUserNotFound
			case codes.Internal:
				errMsg := st.Message()
				if errMsg == custom_errors.ErrDatabaseQuery.Error() {
					return nil, 0, custom_errors.ErrDatabaseQuery
				}
				return nil, 0, custom_errors.ErrExternalServiceError
			default:
				return nil, 0, custom_errors.ErrExternalServiceError
			}
		}
		return nil, 0, custom_errors.ErrExternalServiceError
	}

	followees := make([]*models.RelationUser, len(resp.Followees))
	for i, user := range resp.Followees {
		followees[i] = models.RelationUserFromProto(user)
	}

	return followees, resp.Total, nil
}
