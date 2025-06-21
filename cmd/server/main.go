package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	notification_client "pinstack-api-gateway/internal/clients/notification"
	post_client "pinstack-api-gateway/internal/clients/post"
	relation_client "pinstack-api-gateway/internal/clients/relation"
	user_client "pinstack-api-gateway/internal/clients/user"
	"syscall"
	"time"

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

	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Auth.Address, cfg.Services.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Auth Service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := authConn.Close(); err != nil {
			log.Error("Failed to close Auth Service connection", slog.String("error", err.Error()))
		}
	}()

	postConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Post.Address, cfg.Services.Post.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Post Service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	relationConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Relation.Address, cfg.Services.Relation.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Relation Service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	notificationConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Notification.Address, cfg.Services.Notification.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Notification Service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	userClient := user_client.NewUserClient(userConn, log)
	authClient := auth_client.NewAuthClient(authConn, log)
	postClient := post_client.NewPostClient(postConn, log)
	relationClient := relation_client.NewRelationClient(relationConn, log)
	notificationClient := notification_client.NewNotificationClient(notificationConn, log)

	server := api.NewAPIServer(
		fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		log,
		userClient,
		authClient,
		postClient,
		relationClient,
		notificationClient,
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		if err := server.Run(cfg); err != nil {
			log.Error("Server error", slog.String("error", err.Error()))
		}
		done <- true
	}()

	<-quit
	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	<-done
	log.Info("Server exited")
}
