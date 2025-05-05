package user_client

import (
	"context"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/models"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userClient struct {
	client pb.UserServiceClient
	log    *logger.Logger
}

func NewUserClient(conn *grpc.ClientConn, log *logger.Logger) UserClient {
	return &userClient{
		client: pb.NewUserServiceClient(conn),
		log:    log,
	}
}

func (c *userClient) GetUser(ctx context.Context, id int64) (*models.User, error) {
	c.log.Info("Getting user by ID", "id", id)
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		c.log.Error("Failed to get user", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, custom_errors.ErrUserNotFound
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully got user", "id", id)
	return models.UserFromProto(resp), nil
}

func (c *userClient) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	c.log.Info("Creating new user", "username", user.Username, "email", user.Email)
	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		c.log.Error("Failed to create user", "username", user.Username, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				switch st.Message() {
				case "username already exists":
					return nil, custom_errors.ErrUsernameExists
				case "email already exists":
					return nil, custom_errors.ErrEmailExists
				}
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid username":
					return nil, custom_errors.ErrInvalidUsername
				case "invalid email":
					return nil, custom_errors.ErrInvalidEmail
				case "invalid password":
					return nil, custom_errors.ErrInvalidPassword
				}
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully created user", "id", resp.Id)
	return models.UserFromProto(resp), nil
}

func (c *userClient) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	c.log.Info("Updating user", "id", user.ID)
	resp, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       user.ID,
		Username: &user.Username,
		Email:    &user.Email,
		FullName: user.FullName,
		Bio:      user.Bio,
	})
	if err != nil {
		c.log.Error("Failed to update user", "id", user.ID, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, custom_errors.ErrUserNotFound
			case codes.AlreadyExists:
				switch st.Message() {
				case "username already exists":
					return nil, custom_errors.ErrUsernameExists
				case "email already exists":
					return nil, custom_errors.ErrEmailExists
				}
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid username":
					return nil, custom_errors.ErrInvalidUsername
				case "invalid email":
					return nil, custom_errors.ErrInvalidEmail
				}
			case codes.PermissionDenied:
				return nil, custom_errors.ErrOperationNotAllowed
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully updated user", "id", user.ID)
	return models.UserFromProto(resp), nil
}

func (c *userClient) DeleteUser(ctx context.Context, id int64) error {
	c.log.Info("Deleting user", "id", id)
	_, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		c.log.Error("Failed to delete user", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			case codes.PermissionDenied:
				return custom_errors.ErrOperationNotAllowed
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully deleted user", "id", id)
	return nil
}

func (c *userClient) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	c.log.Info("Getting user by username", "username", username)
	resp, err := c.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{Username: username})
	if err != nil {
		c.log.Error("Failed to get user by username", "username", username, "error", err)
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, custom_errors.ErrUserNotFound
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully got user by username", "username", username)
	return models.UserFromProto(resp), nil
}

func (c *userClient) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	c.log.Info("Getting user by email", "email", email)
	resp, err := c.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: email})
	if err != nil {
		c.log.Error("Failed to get user by email", "email", email, "error", err)
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, custom_errors.ErrUserNotFound
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully got user by email", "email", email)
	return models.UserFromProto(resp), nil
}

func (c *userClient) SearchUsers(ctx context.Context, query string, page, limit int) ([]*models.User, int64, error) {
	c.log.Info("Searching users", "query", query, "page", page, "limit", limit)
	resp, err := c.client.SearchUsers(ctx, &pb.SearchUsersRequest{
		Query:  query,
		Offset: int32(page),
		Limit:  int32(limit),
	})
	if err != nil {
		c.log.Error("Failed to search users", "query", query, "error", err)
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				return nil, 0, custom_errors.ErrInvalidSearchQuery
			}
		}
		return nil, 0, custom_errors.ErrExternalServiceError
	}

	users := make([]*models.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, models.UserFromProto(u))
	}
	c.log.Info("Successfully searched users", "query", query, "total", resp.Total)
	return users, resp.Total, nil
}

func (c *userClient) UpdatePassword(ctx context.Context, id int64, password string) error {
	c.log.Info("Updating user password", "id", id)
	_, err := c.client.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		Id:       id,
		Password: password,
	})
	if err != nil {
		c.log.Error("Failed to update user password", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			case codes.InvalidArgument:
				if st.Message() == "invalid password" {
					return custom_errors.ErrInvalidPassword
				}
			case codes.PermissionDenied:
				return custom_errors.ErrOperationNotAllowed
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully updated user password", "id", id)
	return nil
}

func (c *userClient) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	c.log.Info("Updating user avatar", "id", id)
	_, err := c.client.UpdateAvatar(ctx, &pb.UpdateAvatarRequest{
		Id:        id,
		AvatarUrl: avatarURL,
	})
	if err != nil {
		c.log.Error("Failed to update user avatar", "id", id, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			case codes.InvalidArgument:
				if st.Message() == "invalid avatar format" {
					return custom_errors.ErrInvalidAvatarFormat
				}
			case codes.PermissionDenied:
				return custom_errors.ErrOperationNotAllowed
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully updated user avatar", "id", id)
	return nil
}
