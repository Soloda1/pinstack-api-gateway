package post_client

import (
	"context"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/models"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/post/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postClient struct {
	client pb.PostServiceClient
	log    *logger.Logger
}

func NewPostClient(conn *grpc.ClientConn, log *logger.Logger) PostClient {
	return &postClient{
		client: pb.NewPostServiceClient(conn),
		log:    log,
	}
}

func (c *postClient) CreatePost(ctx context.Context, post *models.CreatePostDTO) (*models.PostDetailed, error) {
	c.log.Debug("Creating post", slog.String("title", post.Title))
	resp, err := c.client.CreatePost(ctx, models.CreatePostDTOToProto(post))
	if err != nil {
		c.log.Error("Failed to create post", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, custom_errors.ErrPostValidation
			default:
				return nil, custom_errors.ErrExternalServiceError
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Debug("response client post", slog.Any("post", post))
	return models.PostDetailedFromProto(resp), nil
}

func (c *postClient) GetPostByID(ctx context.Context, id int64) (*models.PostDetailed, error) {
	c.log.Info("Getting post", slog.Int64("id", id))
	resp, err := c.client.GetPost(ctx, &pb.GetPostRequest{Id: id})
	if err != nil {
		c.log.Error("Failed to get post", slog.Int64("id", id), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, custom_errors.ErrPostNotFound
			case codes.InvalidArgument:
				return nil, custom_errors.ErrPostValidation
			default:
				return nil, custom_errors.ErrExternalServiceError
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	return models.PostDetailedFromProto(resp), nil
}

func (c *postClient) ListPosts(ctx context.Context, filters *models.PostFilters) ([]*models.PostDetailed, int64, error) {
	c.log.Info("Listing posts", slog.Any("filters", filters))
	req := models.PostFiltersToProto(filters)
	resp, err := c.client.ListPosts(ctx, req)
	if err != nil {
		c.log.Error("Failed to list posts", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				return nil, 0, custom_errors.ErrPostValidation
			}
		}
		return nil, 0, custom_errors.ErrExternalServiceError
	}
	posts := make([]*models.PostDetailed, 0, len(resp.Posts))
	for _, p := range resp.Posts {
		posts = append(posts, models.PostDetailedFromProto(p))
	}
	return posts, resp.Total, nil
}

func (c *postClient) UpdatePost(ctx context.Context, id int64, post *models.UpdatePostDTO) error {
	c.log.Info("Updating post", slog.Int64("id", id))
	_, err := c.client.UpdatePost(ctx, models.UpdatePostDTOToProto(id, post))
	if err != nil {
		c.log.Error("Failed to update post", slog.Int64("id", id), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrPostNotFound
			case codes.InvalidArgument:
				return custom_errors.ErrPostValidation
			case codes.PermissionDenied:
				return custom_errors.ErrForbidden
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}

func (c *postClient) DeletePost(ctx context.Context, userID int64, id int64) error {
	c.log.Info("Deleting post", slog.Int64("id", id))
	_, err := c.client.DeletePost(ctx, &pb.DeletePostRequest{UserId: userID, Id: id})
	if err != nil {
		c.log.Error("Failed to delete post", slog.Int64("id", id), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				c.log.Debug("Post not found", slog.Int64("id", id), slog.String("error", err.Error()))
				return custom_errors.ErrPostNotFound
			case codes.InvalidArgument:
				c.log.Debug("Post validation error", slog.Int64("id", id), slog.String("error", err.Error()))
				return custom_errors.ErrPostValidation
			case codes.PermissionDenied:
				c.log.Debug("Permission denied", slog.Int64("id", id), slog.String("error", err.Error()))
				return custom_errors.ErrForbidden
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}
