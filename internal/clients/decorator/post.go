package decorator

import (
	"context"
	post_client "pinstack-api-gateway/internal/clients/post"
	"pinstack-api-gateway/internal/metrics"
	"pinstack-api-gateway/internal/models"
	"strconv"
	"time"
)

// PostClientWithMetrics decorates PostClient with metrics collection
type PostClientWithMetrics struct {
	client          post_client.PostClient
	metricsProvider metrics.MetricsProvider
}

func NewPostClientWithMetrics(client post_client.PostClient, metricsProvider metrics.MetricsProvider) post_client.PostClient {
	return &PostClientWithMetrics{
		client:          client,
		metricsProvider: metricsProvider,
	}
}

func (c *PostClientWithMetrics) CreatePost(ctx context.Context, post *models.CreatePostDTO) (result *models.PostDetailed, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("post-service", "CreatePost", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("post-service", "CreatePost", duration)
		c.metricsProvider.IncProxyRequestsTotal("post-service", "/posts", status)
		c.metricsProvider.ObserveProxyRequestDuration("post-service", "/posts", duration)
	}()

	return c.client.CreatePost(ctx, post)
}

func (c *PostClientWithMetrics) GetPostByID(ctx context.Context, id int64) (result *models.PostDetailed, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("post-service", "GetPostByID", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("post-service", "GetPostByID", duration)
		c.metricsProvider.IncProxyRequestsTotal("post-service", "/posts/"+strconv.FormatInt(id, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("post-service", "/posts/"+strconv.FormatInt(id, 10), duration)
	}()

	return c.client.GetPostByID(ctx, id)
}

func (c *PostClientWithMetrics) ListPosts(ctx context.Context, filters *models.PostFilters) (posts []*models.PostDetailed, total int64, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("post-service", "ListPosts", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("post-service", "ListPosts", duration)
		c.metricsProvider.IncProxyRequestsTotal("post-service", "/posts/list", status)
		c.metricsProvider.ObserveProxyRequestDuration("post-service", "/posts/list", duration)
	}()

	return c.client.ListPosts(ctx, filters)
}

func (c *PostClientWithMetrics) UpdatePost(ctx context.Context, id int64, post *models.UpdatePostDTO) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("post-service", "UpdatePost", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("post-service", "UpdatePost", duration)
		c.metricsProvider.IncProxyRequestsTotal("post-service", "/posts/"+strconv.FormatInt(id, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("post-service", "/posts/"+strconv.FormatInt(id, 10), duration)
	}()

	return c.client.UpdatePost(ctx, id, post)
}

func (c *PostClientWithMetrics) DeletePost(ctx context.Context, userID int64, id int64) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		status := "success"
		if err != nil {
			status = "error"
		}
		c.metricsProvider.IncGRPCClientRequestsTotal("post-service", "DeletePost", status)
		c.metricsProvider.ObserveGRPCClientRequestDuration("post-service", "DeletePost", duration)
		c.metricsProvider.IncProxyRequestsTotal("post-service", "/posts/"+strconv.FormatInt(id, 10), status)
		c.metricsProvider.ObserveProxyRequestDuration("post-service", "/posts/"+strconv.FormatInt(id, 10), duration)
	}()

	return c.client.DeletePost(ctx, userID, id)
}
