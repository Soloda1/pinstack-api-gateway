package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	"pinstack-api-gateway/internal/clients/decorator"
	notification_client "pinstack-api-gateway/internal/clients/notification"
	post_client "pinstack-api-gateway/internal/clients/post"
	relation_client "pinstack-api-gateway/internal/clients/relation"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/metrics/prometheus"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"pinstack-api-gateway/config"
	"pinstack-api-gateway/internal/api"
	"pinstack-api-gateway/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	log.Info("Starting API Gateway")

	// Initialize Prometheus metrics provider
	metricsProvider := prometheus.NewPrometheusMetrics()
	log.Info("Prometheus metrics provider initialized")

	userConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.User.Address, cfg.Services.User.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to User Service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := userConn.Close(); err != nil {
			log.Error("Failed to close User Service connection", slog.String("error", err.Error()))
		}
	}()

	// Record connection metrics
	metricsProvider.IncGRPCClientConnectionsTotal("user-service", "success")

	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Auth.Address, cfg.Services.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Auth Service", slog.String("error", err.Error()))
		metricsProvider.IncGRPCClientConnectionsTotal("auth-service", "error")
		os.Exit(1)
	}
	defer func() {
		if err := authConn.Close(); err != nil {
			log.Error("Failed to close Auth Service connection", slog.String("error", err.Error()))
		}
	}()
	metricsProvider.IncGRPCClientConnectionsTotal("auth-service", "success")

	postConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Post.Address, cfg.Services.Post.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Post Service", slog.String("error", err.Error()))
		metricsProvider.IncGRPCClientConnectionsTotal("post-service", "error")
		os.Exit(1)
	}
	metricsProvider.IncGRPCClientConnectionsTotal("post-service", "success")

	relationConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Relation.Address, cfg.Services.Relation.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Relation Service", slog.String("error", err.Error()))
		metricsProvider.IncGRPCClientConnectionsTotal("relation-service", "error")
		os.Exit(1)
	}
	metricsProvider.IncGRPCClientConnectionsTotal("relation-service", "success")

	notificationConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Notification.Address, cfg.Services.Notification.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Notification Service", slog.String("error", err.Error()))
		metricsProvider.IncGRPCClientConnectionsTotal("notification-service", "error")
		os.Exit(1)
	}
	metricsProvider.IncGRPCClientConnectionsTotal("notification-service", "success")

	// Create base clients
	baseUserClient := user_client.NewUserClient(userConn, log)
	baseAuthClient := auth_client.NewAuthClient(authConn, log)
	basePostClient := post_client.NewPostClient(postConn, log)
	baseRelationClient := relation_client.NewRelationClient(relationConn, log)
	baseNotificationClient := notification_client.NewNotificationClient(notificationConn, log)

	// Wrap clients with metrics decorators
	userClient := decorator.NewUserClientWithMetrics(baseUserClient, metricsProvider)
	authClient := decorator.NewAuthClientWithMetrics(baseAuthClient, metricsProvider)
	postClient := decorator.NewPostClientWithMetrics(basePostClient, metricsProvider)
	relationClient := decorator.NewRelationClientWithMetrics(baseRelationClient, metricsProvider)
	notificationClient := decorator.NewNotificationClientWithMetrics(baseNotificationClient, metricsProvider)

	server := api.NewAPIServer(
		fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		log,
		userClient,
		authClient,
		postClient,
		relationClient,
		notificationClient,
		metricsProvider,
	)

	metricsAddr := fmt.Sprintf("%s:%d", cfg.Prometheus.Address, cfg.Prometheus.Port)
	metricsServer := &http.Server{
		Addr:    metricsAddr,
		Handler: nil,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool, 1)
	metricsDone := make(chan bool, 1)

	go func() {
		if err := server.Run(cfg); err != nil {
			log.Error("Server error", slog.String("error", err.Error()))
		}
		done <- true
	}()

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Info("Starting Prometheus metrics server", slog.String("address", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Prometheus metrics server error", slog.String("error", err.Error()))
		}
		metricsDone <- true
	}()

	<-quit
	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Metrics server shutdown error", slog.String("error", err.Error()))
	}

	<-done
	<-metricsDone

	log.Info("Server exited")
}
