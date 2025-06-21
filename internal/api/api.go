// Package api provides HTTP API server implementation
// @title Pinstack API Gateway
// @version 1.0
// @description API Gateway for Pinstack social media platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package api

import (
	"context"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/config"
	_ "pinstack-api-gateway/docs"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	notification_client "pinstack-api-gateway/internal/clients/notification"
	post_client "pinstack-api-gateway/internal/clients/post"
	relation_client "pinstack-api-gateway/internal/clients/relation"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/logger"
	"time"
)

type APIServer struct {
	address            string
	log                *logger.Logger
	router             *Router
	server             *http.Server
	userClient         user_client.UserClient
	authClient         auth_client.AuthClient
	postClient         post_client.PostClient
	relationClient     relation_client.RelationClient
	notificationClient notification_client.NotificationClient
}

func NewAPIServer(address string,
	log *logger.Logger,
	userClient user_client.UserClient,
	authClient auth_client.AuthClient,
	postClient post_client.PostClient,
	relationClient relation_client.RelationClient,
	notificationClient notification_client.NotificationClient,
) *APIServer {
	return &APIServer{
		address:            address,
		log:                log,
		userClient:         userClient,
		authClient:         authClient,
		postClient:         postClient,
		relationClient:     relationClient,
		notificationClient: notificationClient,
	}
}

func (s *APIServer) Run(cfg *config.Config) error {
	s.router = NewRouter(s.log, s.userClient, s.authClient, s.postClient, s.relationClient, s.notificationClient)
	s.router.Setup(cfg)

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.router.GetRouter(),
		ReadTimeout:  time.Duration(cfg.HTTPServer.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPServer.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPServer.IdleTimeout) * time.Second,
	}

	s.log.Info("Starting server", slog.String("address", s.address))
	s.log.Debug("Debug logger enabled")

	return s.server.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
