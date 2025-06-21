package notification_client

import (
	"context"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/notification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/models"
)

type notificationClient struct {
	client pb.NotificationServiceClient
	log    *logger.Logger
}

func NewNotificationClient(conn *grpc.ClientConn, log *logger.Logger) NotificationClient {
	return &notificationClient{
		client: pb.NewNotificationServiceClient(conn),
		log:    log,
	}
}

func (c *notificationClient) SendNotification(ctx context.Context, userID int64, notificationType string, payload []byte) (int64, error) {
	c.log.Debug("Sending notification", slog.Int64("userID", userID), slog.String("type", notificationType))

	resp, err := c.client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:  userID,
		Type:    notificationType,
		Payload: payload,
	})

	if err != nil {
		c.log.Error("Failed to send notification", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return 0, custom_errors.ErrValidationFailed
			default:
				return 0, custom_errors.ErrExternalServiceError
			}
		}
		return 0, custom_errors.ErrExternalServiceError
	}

	return resp.NotificationId, nil
}

func (c *notificationClient) GetNotificationDetails(ctx context.Context, notificationID int64) (*models.Notification, error) {
	c.log.Debug("Getting notification details", slog.Int64("notificationID", notificationID))

	resp, err := c.client.GetNotificationDetails(ctx, &pb.GetNotificationDetailsRequest{
		NotificationId: notificationID,
	})

	if err != nil {
		c.log.Error("Failed to get notification details", slog.Int64("notificationID", notificationID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, custom_errors.ErrNotificationNotFound
			case codes.InvalidArgument:
				return nil, custom_errors.ErrValidationFailed
			default:
				return nil, custom_errors.ErrExternalServiceError
			}
		}
		return nil, custom_errors.ErrExternalServiceError
	}

	return models.NotificationFromProto(resp), nil
}

func (c *notificationClient) GetUserNotificationFeed(ctx context.Context, userID int64, page, limit int32) ([]*models.Notification, int32, error) {
	c.log.Debug("Getting user notification feed", slog.Int64("userID", userID), slog.Int("page", int(page)), slog.Int("limit", int(limit)))

	resp, err := c.client.GetUserNotificationFeed(ctx, &pb.GetUserNotificationFeedRequest{
		UserId: userID,
		Page:   page,
		Limit:  limit,
	})

	if err != nil {
		c.log.Error("Failed to get user notification feed", slog.Int64("userID", userID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, 0, custom_errors.ErrValidationFailed
			default:
				return nil, 0, custom_errors.ErrExternalServiceError
			}
		}
		return nil, 0, custom_errors.ErrExternalServiceError
	}

	notifications := make([]*models.Notification, 0, len(resp.Notifications))
	for _, n := range resp.Notifications {
		notifications = append(notifications, models.NotificationFromProto(n))
	}

	return notifications, resp.Total, nil
}

func (c *notificationClient) ReadNotification(ctx context.Context, notificationID int64) error {
	c.log.Debug("Reading notification", slog.Int64("notificationID", notificationID))

	_, err := c.client.ReadNotification(ctx, &pb.ReadNotificationRequest{
		NotificationId: notificationID,
	})

	if err != nil {
		c.log.Error("Failed to read notification", slog.Int64("notificationID", notificationID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrNotificationNotFound
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

func (c *notificationClient) ReadAllUserNotifications(ctx context.Context, userID int64) error {
	c.log.Debug("Reading all user notifications", slog.Int64("userID", userID))

	_, err := c.client.ReadAllUserNotifications(ctx, &pb.ReadAllUserNotificationsRequest{
		UserId: userID,
	})

	if err != nil {
		c.log.Error("Failed to read all user notifications", slog.Int64("userID", userID), slog.String("error", err.Error()))
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

func (c *notificationClient) RemoveNotification(ctx context.Context, notificationID int64) error {
	c.log.Debug("Removing notification", slog.Int64("notificationID", notificationID))

	_, err := c.client.RemoveNotification(ctx, &pb.RemoveNotificationRequest{
		NotificationId: notificationID,
	})

	if err != nil {
		c.log.Error("Failed to remove notification", slog.Int64("notificationID", notificationID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return custom_errors.ErrNotificationNotFound
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

func (c *notificationClient) GetUnreadCount(ctx context.Context, userID int64) (int32, error) {
	c.log.Debug("Getting unread count", slog.Int64("userID", userID))

	resp, err := c.client.GetUnreadCount(ctx, &pb.GetUnreadCountRequest{
		UserId: userID,
	})

	if err != nil {
		c.log.Error("Failed to get unread count", slog.Int64("userID", userID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return 0, custom_errors.ErrValidationFailed
			default:
				return 0, custom_errors.ErrExternalServiceError
			}
		}
		return 0, custom_errors.ErrExternalServiceError
	}

	return resp.Count, nil
}
