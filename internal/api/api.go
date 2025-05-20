package api

import (
	"context"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/config"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	post_client "pinstack-api-gateway/internal/clients/post"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/logger"
	"time"
)

type APIServer struct {
	address    string
	log        *logger.Logger
	router     *Router
	server     *http.Server
	userClient user_client.UserClient
	authClient auth_client.AuthClient
	postClient post_client.PostClient
}

func NewAPIServer(address string, log *logger.Logger, userClient user_client.UserClient, authClient auth_client.AuthClient, postClient post_client.PostClient) *APIServer {
	return &APIServer{
		address:    address,
		log:        log,
		userClient: userClient,
		authClient: authClient,
		postClient: postClient,
	}
}

func (s *APIServer) Run(cfg *config.Config) error {
	s.router = NewRouter(s.log, s.userClient, s.authClient)
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
