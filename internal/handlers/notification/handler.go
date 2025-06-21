package notification_handler

import (
	notification_client "pinstack-api-gateway/internal/clients/notification"
	"pinstack-api-gateway/internal/logger"
)

type NotificationHandler struct {
	notificationClient notification_client.NotificationClient
	log                *logger.Logger
}

func NewNotificationHandler(notificationClient notification_client.NotificationClient, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationClient: notificationClient,
		log:                log,
	}
}
