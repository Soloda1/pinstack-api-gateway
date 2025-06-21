package notification_client

import (
	"context"
	"pinstack-api-gateway/internal/models"
)

type NotificationClient interface {
	SendNotification(ctx context.Context, userID int64, notificationType string, payload []byte) (int64, error)
	GetNotificationDetails(ctx context.Context, notificationID int64) (*models.Notification, error)
	GetUserNotificationFeed(ctx context.Context, userID int64, page, limit int32) ([]*models.Notification, int32, error)
	ReadNotification(ctx context.Context, notificationID int64) error
	ReadAllUserNotifications(ctx context.Context, userID int64) error
	RemoveNotification(ctx context.Context, notificationID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int32, error)
}
