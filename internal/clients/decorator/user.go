package decorator

import (
	"context"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/metrics"
	"pinstack-api-gateway/internal/models"
	"strconv"
	"time"
)

// UserClientWithMetrics decorates UserClient with metrics collection
type UserClientWithMetrics struct {
	client          user_client.UserClient
	metricsProvider metrics.MetricsProvider
}

func NewUserClientWithMetrics(client user_client.UserClient, metricsProvider metrics.MetricsProvider) user_client.UserClient {
	return &UserClientWithMetrics{
		client:          client,
		metricsProvider: metricsProvider,
	}
}

func (c *UserClientWithMetrics) GetUser(ctx context.Context, id int64) (user *models.User, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "GetUser", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "GetUser", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/"+strconv.FormatInt(id, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/"+strconv.FormatInt(id, 10), duration)
	}()

	return c.client.GetUser(ctx, id)
}

func (c *UserClientWithMetrics) CreateUser(ctx context.Context, user *models.User) (result *models.User, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "CreateUser", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "CreateUser", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users", status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users", duration)
	}()

	return c.client.CreateUser(ctx, user)
}

func (c *UserClientWithMetrics) UpdateUser(ctx context.Context, user *models.User) (result *models.User, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "UpdateUser", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "UpdateUser", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users", status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users", duration)
	}()

	return c.client.UpdateUser(ctx, user)
}

func (c *UserClientWithMetrics) DeleteUser(ctx context.Context, id int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "DeleteUser", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "DeleteUser", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/"+strconv.FormatInt(id, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/"+strconv.FormatInt(id, 10), duration)
	}()

	return c.client.DeleteUser(ctx, id)
}

func (c *UserClientWithMetrics) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "GetUserByUsername", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "GetUserByUsername", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/username/"+username, status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/username/"+username, duration)
	}()

	return c.client.GetUserByUsername(ctx, username)
}

func (c *UserClientWithMetrics) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "GetUserByEmail", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "GetUserByEmail", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/email/"+email, status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/email/"+email, duration)
	}()

	return c.client.GetUserByEmail(ctx, email)
}

func (c *UserClientWithMetrics) SearchUsers(ctx context.Context, query string, page, limit int) (users []*models.User, total int64, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "SearchUsers", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "SearchUsers", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/search", status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/search", duration)
	}()

	return c.client.SearchUsers(ctx, query, page, limit)
}

func (c *UserClientWithMetrics) UpdateAvatar(ctx context.Context, id int64, avatarURL string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("user-service", "UpdateAvatar", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("user-service", "UpdateAvatar", duration)
		c.metricsProvider.IncProxyRequestsTotal("user-service", "/users/avatar", status)
		c.metricsProvider.ObserveProxyRequestDuration("user-service", "/users/avatar", duration)
	}()

	return c.client.UpdateAvatar(ctx, id, avatarURL)
}
