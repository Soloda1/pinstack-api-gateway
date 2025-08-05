package auth_client

import (
	"context"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/models"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authClient struct {
	client pb.AuthServiceClient
	log    *logger.Logger
}

func NewAuthClient(conn *grpc.ClientConn, log *logger.Logger) AuthClient {
	return &authClient{
		client: pb.NewAuthServiceClient(conn),
		log:    log,
	}
}

func (c *authClient) Register(ctx context.Context, req *models.RegisterRequest) (*models.TokenPair, error) {
	c.log.Info("Registering new user", "username", req.Username, "email", req.Email)
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarURL,
	})
	if err != nil {
		c.log.Error("Failed to register user", "username", req.Username, "error", err)
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
	c.log.Info("Successfully registered user", "username", req.Username)
	return &models.TokenPair{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (c *authClient) Login(ctx context.Context, req *models.LoginRequest) (*models.TokenPair, error) {
	c.log.Info("Logging in user", "login", req.Login)
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		c.log.Error("Failed to login user", "login", req.Login, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, custom_errors.ErrUserNotFound
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid credentials":
					return nil, custom_errors.ErrInvalidCredentials
				case "invalid password":
					return nil, custom_errors.ErrInvalidCredentials
				case "invalid login":
					return nil, custom_errors.ErrInvalidCredentials
				}
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully logged in user", "login", req.Login)
	return &models.TokenPair{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (c *authClient) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	c.log.Info("Refreshing tokens")
	resp, err := c.client.Refresh(ctx, &pb.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		c.log.Error("Failed to refresh tokens", "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid refresh token":
					return nil, custom_errors.ErrInvalidRefreshToken
				case "token expired":
					return nil, custom_errors.ErrTokenExpired
				case "invalid token":
					return nil, custom_errors.ErrInvalidRefreshToken
				}
			case codes.Unauthenticated:
				return nil, custom_errors.ErrUnauthenticated
			case codes.NotFound:
				return nil, custom_errors.ErrUserNotFound
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully refreshed tokens")
	return &models.TokenPair{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (c *authClient) Logout(ctx context.Context, refreshToken string) error {
	c.log.Info("Logging out user")
	_, err := c.client.Logout(ctx, &pb.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		c.log.Error("Failed to logout user", "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid refresh token":
					return custom_errors.ErrInvalidRefreshToken
				case "token expired":
					return custom_errors.ErrTokenExpired
				case "invalid token":
					return custom_errors.ErrInvalidRefreshToken
				}
			case codes.Unauthenticated:
				return custom_errors.ErrUnauthenticated
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully logged out user")
	return nil
}

func (c *authClient) UpdatePassword(ctx context.Context, req *models.UpdatePasswordRequest) error {
	c.log.Info("Updating user password", "id", req.ID)
	_, err := c.client.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		Id:          req.ID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		c.log.Error("Failed to update password", "id", req.ID, "error", err)
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrUserNotFound
			case codes.InvalidArgument:
				switch st.Message() {
				case "invalid password":
					return custom_errors.ErrInvalidCredentials
				case "invalid old password":
					return custom_errors.ErrInvalidCredentials
				case "invalid new password":
					return custom_errors.ErrInvalidPassword
				}
			case codes.PermissionDenied:
				return custom_errors.ErrOperationNotAllowed
			}
		}
		return custom_errors.ErrExternalServiceError
	}
	c.log.Info("Successfully updated password", "id", req.ID)
	return nil
}
