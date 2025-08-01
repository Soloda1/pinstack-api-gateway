package api

import (
	"net/http"
	"pinstack-api-gateway/config"
	auth_client "pinstack-api-gateway/internal/clients/auth"
	notification_client "pinstack-api-gateway/internal/clients/notification"
	post_client "pinstack-api-gateway/internal/clients/post"
	relation_client "pinstack-api-gateway/internal/clients/relation"
	user_client "pinstack-api-gateway/internal/clients/user"
	auth_handler "pinstack-api-gateway/internal/handlers/auth"
	notification_handler "pinstack-api-gateway/internal/handlers/notification"
	post_handler "pinstack-api-gateway/internal/handlers/post"
	relation_handler "pinstack-api-gateway/internal/handlers/relation"
	user_handler "pinstack-api-gateway/internal/handlers/user"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/middlewares"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	router             *chi.Mux
	log                *logger.Logger
	userClient         user_client.UserClient
	authClient         auth_client.AuthClient
	postClient         post_client.PostClient
	relationClient     relation_client.RelationClient
	notificationClient notification_client.NotificationClient
}

func NewRouter(log *logger.Logger, userClient user_client.UserClient, authClient auth_client.AuthClient, postClient post_client.PostClient, relationClient relation_client.RelationClient, notificationClient notification_client.NotificationClient) *Router {
	return &Router{
		router:             chi.NewRouter(),
		log:                log,
		userClient:         userClient,
		authClient:         authClient,
		postClient:         postClient,
		relationClient:     relationClient,
		notificationClient: notificationClient,
	}
}

func (r *Router) Setup(cfg *config.Config) {
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Recoverer)
	r.router.Use(middlewares.RequestLoggerMiddleware(r.log))
	r.router.Use(middleware.Timeout(time.Duration(cfg.HTTPServer.Timeout) * time.Second))

	jwtMiddleware := middlewares.JWTValidationMiddleware(cfg.JWT.Secret, r.log)

	r.router.Get("/swagger/*", httpSwagger.WrapHandler)

	r.router.Route("/api/v1", func(v1 chi.Router) {
		v1.Mount("/users", r.setupUserRoutes(jwtMiddleware))
		v1.Mount("/auth", r.setupAuthRoutes(jwtMiddleware))
		v1.Mount("/posts", r.setupPostRoutes(jwtMiddleware))
		v1.Mount("/relation", r.setupRelationRoutes(jwtMiddleware))
		v1.Mount("/notification", r.setupNotificationRoutes(jwtMiddleware))
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
		r.Put("/", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Put("/avatar", userHandler.UpdateAvatar)
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

func (r *Router) setupPostRoutes(jwtMiddleware func(next http.Handler) http.Handler) http.Handler {
	postHandler := post_handler.NewPostHandler(r.postClient, r.userClient, r.log)
	router := chi.NewRouter()

	router.Get("/list", postHandler.List)
	router.Get("/{id}", postHandler.Get)
	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/", postHandler.Create)
		r.Put("/{id}", postHandler.Update)
		r.Delete("/{id}", postHandler.Delete)
	})

	return router
}

func (r *Router) setupRelationRoutes(jwtMiddleware func(next http.Handler) http.Handler) http.Handler {
	relationHandler := relation_handler.NewRelationHandler(r.relationClient, r.log)
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/follow", relationHandler.Follow)
		r.Post("/unfollow", relationHandler.Unfollow)
	})

	return router
}

func (r *Router) setupNotificationRoutes(jwtMiddleware func(next http.Handler) http.Handler) http.Handler {
	notificationHandler := notification_handler.NewNotificationHandler(r.notificationClient, r.log)
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Get("/{notification_id}", notificationHandler.GetNotificationDetails)
		r.Get("/feed", notificationHandler.GetUserNotificationFeed)
		r.Get("/unread-count", notificationHandler.GetUnreadCount)
		r.Put("/{notification_id}/read", notificationHandler.ReadNotification)
		r.Put("/read-all", notificationHandler.ReadAllUserNotifications)
		r.Delete("/{notification_id}", notificationHandler.RemoveNotification)
		r.Post("/send", notificationHandler.SendNotification)
	})

	return router
}

func (r *Router) GetRouter() *chi.Mux {
	return r.router
}
