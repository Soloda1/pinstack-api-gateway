package api

import (
	"net/http"
	"pinstack-api-gateway/config"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	user_client "pinstack-api-gateway/internal/clients/user"
	auth_handler "pinstack-api-gateway/internal/handlers/auth"
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
	authClient auth_client.AuthClient
}

func NewRouter(log *logger.Logger, userClient user_client.UserClient, authClient auth_client.AuthClient) *Router {
	return &Router{
		router:     chi.NewRouter(),
		log:        log,
		userClient: userClient,
		authClient: authClient,
	}
}

func (r *Router) Setup(cfg *config.Config) {
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Recoverer)
	r.router.Use(middlewares.RequestLoggerMiddleware(r.log))
	r.router.Use(middleware.Timeout(time.Duration(cfg.HTTPServer.Timeout) * time.Second))

	jwtMiddleware := middlewares.JWTValidationMiddleware(cfg.JWT.Secret, r.log)

	r.router.Route("/api/v1", func(v1 chi.Router) {
		v1.Mount("/users", r.setupUserRoutes(jwtMiddleware))
		v1.Mount("/auth", r.setupAuthRoutes(jwtMiddleware))
	})
}

func (r *Router) setupUserRoutes(jwtMiddleware func(next http.Handler) http.Handler) http.Handler {
	userHandler := user_handler.NewUserHandler(r.userClient, r.log)
	router := chi.NewRouter()

	router.Get("/{id}", userHandler.GetUser)
	router.Get("/username/{username}", userHandler.GetUserByUsername)
	router.Get("/email/{email}", userHandler.GetUserByEmail)
	router.Get("/search", userHandler.SearchUsers)

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/", userHandler.CreateUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Put("/{id}/avatar", userHandler.UpdateAvatar)
	})

	return router
}

func (r *Router) setupAuthRoutes(jwtMiddleware func(next http.Handler) http.Handler) http.Handler {
	authHandler := auth_handler.NewAuthHandler(r.authClient, r.log)
	router := chi.NewRouter()

	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.Refresh)
	router.Post("/logout", authHandler.Logout)

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/update-password", authHandler.UpdatePassword)
	})

	return router
}

func (r *Router) GetRouter() *chi.Mux {
	return r.router
}
