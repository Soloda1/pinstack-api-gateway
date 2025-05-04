package user_handler

import (
	"github.com/go-chi/chi/v5"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/logger"
)

type UserHandler struct {
	userClient user_client.UserClient
	log        *logger.Logger
}

func NewUserHandler(userClient user_client.UserClient, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		log:        log,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/{id}", h.GetUser)
		r.Post("/", h.CreateUser)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
		r.Get("/username/{username}", h.GetUserByUsername)
		r.Get("/email/{email}", h.GetUserByEmail)
		r.Get("/search", h.SearchUsers)
		r.Put("/{id}/password", h.UpdatePassword)
		r.Put("/{id}/avatar", h.UpdateAvatar)
	})
}
