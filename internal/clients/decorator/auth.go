package decorator

import (
	"context"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	"pinstack-api-gateway/internal/metrics"
	"pinstack-api-gateway/internal/models"
	"time"
)

// AuthClientWithMetrics decorates AuthClient with metrics collection
type AuthClientWithMetrics struct {
	client          auth_client.AuthClient
	metricsProvider metrics.MetricsProvider
}

func NewAuthClientWithMetrics(client auth_client.AuthClient, metricsProvider metrics.MetricsProvider) auth_client.AuthClient {
	return &AuthClientWithMetrics{
		client:          client,
		metricsProvider: metricsProvider,
	}
}

func (c *AuthClientWithMetrics) Register(ctx context.Context, req *models.RegisterRequest) (result *models.TokenPair, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("auth-service", "Register", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("auth-service", "Register", duration)
		c.metricsProvider.IncProxyRequestsTotal("auth-service", "/auth/register", status)
		c.metricsProvider.ObserveProxyRequestDuration("auth-service", "/auth/register", duration)
		c.metricsProvider.IncAuthenticationTotal(status)
	}()

	return c.client.Register(ctx, req)
}

func (c *AuthClientWithMetrics) Login(ctx context.Context, req *models.LoginRequest) (result *models.TokenPair, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("auth-service", "Login", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("auth-service", "Login", duration)
		c.metricsProvider.IncProxyRequestsTotal("auth-service", "/auth/login", status)
		c.metricsProvider.ObserveProxyRequestDuration("auth-service", "/auth/login", duration)
		c.metricsProvider.IncAuthenticationTotal(status)
	}()

	return c.client.Login(ctx, req)
}

func (c *AuthClientWithMetrics) Refresh(ctx context.Context, refreshToken string) (result *models.TokenPair, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("auth-service", "Refresh", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("auth-service", "Refresh", duration)
		c.metricsProvider.IncProxyRequestsTotal("auth-service", "/auth/refresh", status)
		c.metricsProvider.ObserveProxyRequestDuration("auth-service", "/auth/refresh", duration)
		c.metricsProvider.IncAuthenticationTotal(status)
	}()

	return c.client.Refresh(ctx, refreshToken)
}

func (c *AuthClientWithMetrics) Logout(ctx context.Context, refreshToken string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("auth-service", "Logout", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("auth-service", "Logout", duration)
		c.metricsProvider.IncProxyRequestsTotal("auth-service", "/auth/logout", status)
		c.metricsProvider.ObserveProxyRequestDuration("auth-service", "/auth/logout", duration)
		c.metricsProvider.IncAuthenticationTotal(status)
	}()

	return c.client.Logout(ctx, refreshToken)
}

func (c *AuthClientWithMetrics) UpdatePassword(ctx context.Context, req *models.UpdatePasswordRequest) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("auth-service", "UpdatePassword", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("auth-service", "UpdatePassword", duration)
		c.metricsProvider.IncProxyRequestsTotal("auth-service", "/auth/update-password", status)
		c.metricsProvider.ObserveProxyRequestDuration("auth-service", "/auth/update-password", duration)
	}()

	return c.client.UpdatePassword(ctx, req)
}
