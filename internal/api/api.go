package api

import (
	"context"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/config"
	user_client "pinstack-api-gateway/internal/clients/user"
	user_handler "pinstack-api-gateway/internal/handlers/user"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/middlewares"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	address    string
	log        *logger.Logger
	router     *chi.Mux
	server     *http.Server
	userClient user_client.UserClient
}

func NewAPIServer(address string, log *logger.Logger, userClient user_client.UserClient) *APIServer {
	return &APIServer{
		address:    address,
		log:        log,
		userClient: userClient,
	}
}

func (s *APIServer) Run(cfg *config.Config) error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.RequestLoggerMiddleware(s.log))
	router.Use(middleware.Timeout(time.Duration(cfg.HTTPServer.Timeout) * time.Second))

	userHandler := user_handler.NewUserHandler(s.userClient, s.log)

	router.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/{id}", userHandler.GetUser)
		r.Post("/", userHandler.CreateUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Get("/username/{username}", userHandler.GetUserByUsername)
		r.Get("/email/{email}", userHandler.GetUserByEmail)
		r.Get("/search", userHandler.SearchUsers)
		r.Put("/{id}/password", userHandler.UpdatePassword)
		r.Put("/{id}/avatar", userHandler.UpdateAvatar)
	})

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      router,
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
