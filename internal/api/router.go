package api

import (
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

type Router struct {
	router     *chi.Mux
	log        *logger.Logger
	userClient user_client.UserClient
}

func NewRouter(log *logger.Logger, userClient user_client.UserClient) *Router {
	return &Router{
		router:     chi.NewRouter(),
		log:        log,
		userClient: userClient,
	}
}

func (r *Router) Setup(cfg *config.Config) {
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Recoverer)
	r.router.Use(middlewares.RequestLoggerMiddleware(r.log))
	r.router.Use(middleware.Timeout(time.Duration(cfg.HTTPServer.Timeout) * time.Second))

	r.router.Route("/api/v1", func(v1 chi.Router) {
		v1.Mount("/users", r.setupUserRoutes())
	})
}

func (r *Router) setupUserRoutes() http.Handler {
	userHandler := user_handler.NewUserHandler(r.userClient, r.log)
	router := chi.NewRouter()

	router.Get("/{id}", userHandler.GetUser)
	router.Post("/", userHandler.CreateUser)
	router.Put("/{id}", userHandler.UpdateUser)
	router.Delete("/{id}", userHandler.DeleteUser)
	router.Get("/username/{username}", userHandler.GetUserByUsername)
	router.Get("/email/{email}", userHandler.GetUserByEmail)
	router.Get("/search", userHandler.SearchUsers)
	router.Put("/{id}/avatar", userHandler.UpdateAvatar)

	return router
}

func (r *Router) GetRouter() *chi.Mux {
	return r.router
}
