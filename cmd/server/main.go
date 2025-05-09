package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	auth_client "pinstack-api-gateway/internal/clients/auth"
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
	defer userConn.Close()

	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Services.Auth.Address, cfg.Services.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to Auth Service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer authConn.Close()

	userClient := user_client.NewUserClient(userConn, log)
	authClient := auth_client.NewAuthClient(authConn, log)

	server := api.NewAPIServer(
		fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		log,
		userClient,
		authClient,
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
