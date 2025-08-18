package decorator

import (
	"context"
	relation_client "pinstack-api-gateway/internal/clients/relation"
	"pinstack-api-gateway/internal/metrics"
	"pinstack-api-gateway/internal/models"
	"strconv"
	"time"
)

// RelationClientWithMetrics decorates RelationClient with metrics collection
type RelationClientWithMetrics struct {
	client          relation_client.RelationClient
	metricsProvider metrics.MetricsProvider
}

func NewRelationClientWithMetrics(client relation_client.RelationClient, metricsProvider metrics.MetricsProvider) relation_client.RelationClient {
	return &RelationClientWithMetrics{
		client:          client,
		metricsProvider: metricsProvider,
	}
}

func (c *RelationClientWithMetrics) Follow(ctx context.Context, followerID, followeeID int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("relation-service", "Follow", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("relation-service", "Follow", duration)
		c.metricsProvider.IncProxyRequestsTotal("relation-service", "/relation/follow", status)
		c.metricsProvider.ObserveProxyRequestDuration("relation-service", "/relation/follow", duration)
	}()

	return c.client.Follow(ctx, followerID, followeeID)
}

func (c *RelationClientWithMetrics) Unfollow(ctx context.Context, followerID, followeeID int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("relation-service", "Unfollow", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("relation-service", "Unfollow", duration)
		c.metricsProvider.IncProxyRequestsTotal("relation-service", "/relation/unfollow", status)
		c.metricsProvider.ObserveProxyRequestDuration("relation-service", "/relation/unfollow", duration)
	}()

	return c.client.Unfollow(ctx, followerID, followeeID)
}

func (c *RelationClientWithMetrics) GetFollowers(ctx context.Context, followeeID int64, limit, page int32) (users []*models.RelationUser, total int64, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("relation-service", "GetFollowers", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("relation-service", "GetFollowers", duration)
		c.metricsProvider.IncProxyRequestsTotal("relation-service", "/relation/"+strconv.FormatInt(followeeID, 10)+"/followers", status)
		c.metricsProvider.ObserveProxyRequestDuration("relation-service", "/relation/"+strconv.FormatInt(followeeID, 10)+"/followers", duration)
	}()

	return c.client.GetFollowers(ctx, followeeID, limit, page)
}

func (c *RelationClientWithMetrics) GetFollowees(ctx context.Context, followerID int64, limit, page int32) (users []*models.RelationUser, total int64, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("relation-service", "GetFollowees", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("relation-service", "GetFollowees", duration)
		c.metricsProvider.IncProxyRequestsTotal("relation-service", "/relation/"+strconv.FormatInt(followerID, 10)+"/followees", status)
		c.metricsProvider.ObserveProxyRequestDuration("relation-service", "/relation/"+strconv.FormatInt(followerID, 10)+"/followees", duration)
	}()

	return c.client.GetFollowees(ctx, followerID, limit, page)
}
