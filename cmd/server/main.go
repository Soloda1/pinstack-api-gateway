package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	user_client "pinstack-api-gateway/internal/clients/user"
	"syscall"
	"time"

	"pinstack-api-gateway/config"
	"pinstack-api-gateway/internal/handlers/user"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("Starting API Gateway")

	userConn, err := grpc.NewClient(
		cfg.Services.User.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to User Service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	userClient := user_client.NewUserClient(userConn, log)

	userHandler := user_handler.NewUserHandler(userClient, log)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RequestLoggerMiddleware(log))
	r.Use(middleware.Timeout(time.Duration(cfg.HTTPServer.Timeout) * time.Second))

	userHandler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.HTTPServer.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPServer.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPServer.IdleTimeout) * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		log.Info("Starting HTTP server", "port", cfg.HTTPServer.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server error", slog.String("error", err.Error()))
		}
		done <- true
	}()

	<-quit
	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	err = userConn.Close()
	if err != nil {
		log.Error("Failed to close User Service", slog.String("error", err.Error()))
	}

	<-done
	log.Info("Server exited")
}
