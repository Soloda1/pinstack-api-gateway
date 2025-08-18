package decorator

import (
	"context"
	notification_client "pinstack-api-gateway/internal/clients/notification"
	"pinstack-api-gateway/internal/metrics"
	"pinstack-api-gateway/internal/models"
	"strconv"
	"time"
)

// NotificationClientWithMetrics decorates NotificationClient with metrics collection
type NotificationClientWithMetrics struct {
	client          notification_client.NotificationClient
	metricsProvider metrics.MetricsProvider
}

func NewNotificationClientWithMetrics(client notification_client.NotificationClient, metricsProvider metrics.MetricsProvider) notification_client.NotificationClient {
	return &NotificationClientWithMetrics{
		client:          client,
		metricsProvider: metricsProvider,
	}
}

func (c *NotificationClientWithMetrics) SendNotification(ctx context.Context, userID int64, notificationType string, payload []byte) (result int64, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "SendNotification", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "SendNotification", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/send", status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/send", duration)
	}()

	return c.client.SendNotification(ctx, userID, notificationType, payload)
}

func (c *NotificationClientWithMetrics) GetNotificationDetails(ctx context.Context, notificationID int64) (result *models.Notification, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "GetNotificationDetails", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "GetNotificationDetails", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10), duration)
	}()

	return c.client.GetNotificationDetails(ctx, notificationID)
}

func (c *NotificationClientWithMetrics) GetUserNotificationFeed(ctx context.Context, userID int64, page, limit int32) (notifications []*models.Notification, total int32, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "GetUserNotificationFeed", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "GetUserNotificationFeed", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/feed", status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/feed", duration)
	}()

	return c.client.GetUserNotificationFeed(ctx, userID, page, limit)
}

func (c *NotificationClientWithMetrics) ReadNotification(ctx context.Context, notificationID int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "ReadNotification", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "ReadNotification", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10)+"/read", status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10)+"/read", duration)
	}()

	return c.client.ReadNotification(ctx, notificationID)
}

func (c *NotificationClientWithMetrics) ReadAllUserNotifications(ctx context.Context, userID int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "ReadAllUserNotifications", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "ReadAllUserNotifications", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/read-all", status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/read-all", duration)
	}()

	return c.client.ReadAllUserNotifications(ctx, userID)
}

func (c *NotificationClientWithMetrics) RemoveNotification(ctx context.Context, notificationID int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "RemoveNotification", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "RemoveNotification", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/"+strconv.FormatInt(notificationID, 10), duration)
	}()

	return c.client.RemoveNotification(ctx, notificationID)
}

func (c *NotificationClientWithMetrics) GetUnreadCount(ctx context.Context, userID int64) (result int32, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("notification-service", "GetUnreadCount", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("notification-service", "GetUnreadCount", duration)
		c.metricsProvider.IncProxyRequestsTotal("notification-service", "/notification/unread-count", status)
		c.metricsProvider.ObserveProxyRequestDuration("notification-service", "/notification/unread-count", duration)
	}()

	return c.client.GetUnreadCount(ctx, userID)
}
