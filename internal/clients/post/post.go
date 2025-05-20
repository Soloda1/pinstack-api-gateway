package post_client

import (
	"context"
	"pinstack-api-gateway/internal/custom_errors"
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
	c.log.Info("Creating post", "title", post.Title)
	resp, err := c.client.CreatePost(ctx, models.CreatePostDTOToProto(post))
	if err != nil {
		c.log.Error("Failed to create post", "error", err)
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
	return models.PostDetailedFromProto(resp), nil
}

func (c *postClient) GetPostByID(ctx context.Context, id int64) (*models.PostDetailed, error) {
	c.log.Info("Getting post", "id", id)
	resp, err := c.client.GetPost(ctx, &pb.GetPostRequest{Id: id})
	if err != nil {
		c.log.Error("Failed to get post", "id", id, "error", err)
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

func (c *postClient) ListPosts(ctx context.Context, filters *models.PostFilters) ([]*models.PostDetailed, error) {
	c.log.Info("Listing posts", "filters", filters)
	req := models.PostFiltersToProto(filters)
	resp, err := c.client.ListPosts(ctx, req)
	if err != nil {
		c.log.Error("Failed to list posts", "error", err)
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				return nil, custom_errors.ErrPostValidation
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	posts := make([]*models.PostDetailed, 0, len(resp.Posts))
	for _, p := range resp.Posts {
		posts = append(posts, models.PostDetailedFromProto(p))
	}
	return posts, nil
}

func (c *postClient) UpdatePost(ctx context.Context, id int64, post *models.UpdatePostDTO) error {
	c.log.Info("Updating post", "id", id)
	_, err := c.client.UpdatePost(ctx, models.UpdatePostDTOToProto(id, post))
	if err != nil {
		c.log.Error("Failed to update post", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrPostNotFound
			case codes.InvalidArgument:
				return custom_errors.ErrPostValidation
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}

func (c *postClient) DeletePost(ctx context.Context, id int64) error {
	c.log.Info("Deleting post", "id", id)
	_, err := c.client.DeletePost(ctx, &pb.DeletePostRequest{Id: id})
	if err != nil {
		c.log.Error("Failed to delete post", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrPostNotFound
			case codes.InvalidArgument:
				return custom_errors.ErrPostValidation
			default:
				return custom_errors.ErrExternalServiceError
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	return nil
}
